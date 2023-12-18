import {
  FEATURE_FLAGS,
  ROUTES,
  SEARCH_FACETS_SINGLE_NODE_LABEL_KEY,
  SEARCH_FACETS_SINGLE_NODE_VALUE_PREFIX,
} from 'config';
import { BrowseItemConfig, BrowseItemConfigType } from 'config/browse';
import {
  Hierarchy,
  HierarchyNode,
  HierarchyNodeDisplayValue,
  SearchResultsFacet,
  SearchResultsFacetBucket,
} from 'stores/search';
import { getFilterQueryParamName, getQueryStringFromParams } from 'utils/search';
import { formatValue, normalizeConceptId, sortAndFilterValues } from 'utils/field';
import stores from 'stores';

export interface BrowseItemConstructorArgs {
  config: BrowseItemConfig;
  facet?: SearchResultsFacet;
  bucket?: SearchResultsFacetBucket;
  node?: HierarchyNode;
  hierarchy?: Hierarchy;
  allowDescent?: boolean;
  isSingleNodeVersion?: boolean;
}

interface BrowseItemCloneArgs extends Omit<BrowseItemConstructorArgs, 'config'> {
  config?: BrowseItemConfig;
}

export class BrowseItem {
  #config: BrowseItemConfig;
  #allowDescent: boolean;
  #facet?: SearchResultsFacet;
  #bucket?: SearchResultsFacetBucket;
  #node?: HierarchyNode;
  #hierarchy?: Hierarchy;
  #isSingleNodeVersion: boolean;

  constructor({
    config,
    facet,
    bucket,
    node,
    hierarchy,
    allowDescent = true,
    isSingleNodeVersion = false,
  }: BrowseItemConstructorArgs) {
    this.#config = config;
    this.#facet = facet;
    this.#bucket = bucket;
    this.#node = node;
    this.#hierarchy = hierarchy;
    this.#allowDescent = !!allowDescent;
    this.#isSingleNodeVersion = isSingleNodeVersion;
  }

  get level() {
    return this.#node?.depth ?? 0;
  }

  get minLevel() {
    return this.#config.minLevel ?? 0;
  }

  get maxLevel() {
    return Math.max(this.#config.maxLevel ?? 0, this.minLevel);
  }

  get parentNode() {
    if (this.isSingleNodeVersion) {
      return this.#node;
    }

    const parentNodeId = this.#node?.parentNodeId;
    return parentNodeId && this.#hierarchy?.nodes.find(({ nodeId }) => nodeId === parentNodeId);
  }

  get parent() {
    const { parentNode } = this;
    return parentNode && this.clone({ node: parentNode });
  }

  get parents() {
    const { parent } = this;
    if (!parent || parent.level < this.minLevel) {
      return [];
    }

    return [parent, ...parent.parents].reverse();
  }

  get childNodes() {
    const { nodeId } = this;
    return nodeId ? this.#hierarchy?.nodes.filter(({ parentNodeId }) => parentNodeId === nodeId) ?? [] : [];
  }

  get singleNodeVersion(): BrowseItem | null {
    if (!this.#config.enableSingleNodeVersion || !this.canDescend || !this.childNodes.length) {
      return null;
    }

    return this.clone({ isSingleNodeVersion: true });
  }

  get children() {
    const { childNodes, singleNodeVersion, isSingleNodeVersion } = this;
    if (isSingleNodeVersion) {
      return [];
    }

    const children = childNodes.map((node) => this.clone({ node }));

    if (!children.some(({ count }) => count) || !singleNodeVersion?.count) {
      return children;
    }

    return [singleNodeVersion].concat(children);
  }

  get hasChildren() {
    return !!this.childNodes.length;
  }

  get descendants() {
    return this.children.reduce(
      (descendants, child) => descendants.concat(child.hasChildren ? child.descendants : [], [child]),
      []
    );
  }

  get childlessDescendants() {
    return this.children.reduce(
      (childlessDescendants, child) =>
        childlessDescendants.concat(child.hasChildren ? child.childlessDescendants : [child]),
      []
    );
  }

  get canDescend() {
    return (
      this.#allowDescent &&
      this.#config.type === BrowseItemConfigType.hierarchy &&
      this.level < this.maxLevel &&
      this.hasChildren
    );
  }

  get isSingleNodeVersion() {
    return this.#isSingleNodeVersion;
  }

  get url() {
    if (this.canDescend) {
      return null;
    }

    const params = new URLSearchParams();

    if (!this.hasHierarchyInfo && this.hasChildren) {
      this.childlessDescendants.forEach((descendant) =>
        params.append(getFilterQueryParamName(this.#config.axis.name), normalizeConceptId(descendant.value))
      );
    } else {
      params.append(getFilterQueryParamName(this.#config.axis.name), normalizeConceptId(this.value));
    }

    return `${ROUTES.SEARCH}?${getQueryStringFromParams(params)}`;
  }

  get bucket() {
    const node = this.#node;
    return this.#bucket ?? this.#facet?.buckets.find(({ value }) => value === node.nodeId);
  }

  get hasHierarchyInfo() {
    return this.#facet?.buckets.some(({ hierarchyInfo }) => !!hierarchyInfo);
  }

  get #displayValues() {
    return (<HierarchyNodeDisplayValue[]>(
      sortAndFilterValues(this.#node?.display, FEATURE_FLAGS.BROWSE_HIERARCHY_TRANSLATIONS)
    ))?.map(({ display }) => display);
  }

  get text() {
    const text = formatValue(this.#displayValues ?? [this.bucket?.value ?? this.nodeId], this.#config.key);
    return this.isSingleNodeVersion ? stores.i18n.t(SEARCH_FACETS_SINGLE_NODE_LABEL_KEY, { text }) : text;
  }

  get value() {
    const prefix = this.isSingleNodeVersion ? SEARCH_FACETS_SINGLE_NODE_VALUE_PREFIX : '';
    return prefix + (this.bucket?.value ?? this.nodeId);
  }

  get count() {
    const count = this.hasHierarchyInfo
      ? this.bucket?.count ?? 0
      : this.bucket?.count ?? this.children.reduce((count, child) => count + child.count, 0);

    if (!this.isSingleNodeVersion) {
      return count;
    }

    const children = this.childNodes.map((node) => this.clone({ node }));
    return count - children.reduce((sum, { count }) => sum + count, 0);
  }

  get nodeId() {
    return this.#node?.nodeId;
  }

  clone(args: BrowseItemCloneArgs) {
    return new BrowseItem({
      config: this.#config,
      facet: this.#facet,
      bucket: this.#bucket,
      node: this.#node,
      hierarchy: this.#hierarchy,
      allowDescent: this.#allowDescent,
      ...args,
    });
  }
}
