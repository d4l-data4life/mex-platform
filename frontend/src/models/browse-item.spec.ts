jest.mock('stencil-router-v2');

import * as config from 'config';
import { BrowseItemConfig, BrowseItemConfigType } from 'config/browse';
import stores from 'stores';
import { Hierarchy, HierarchyNode, SearchResultsFacet, SearchResultsFacetBucket } from 'stores/search';
import { BrowseItem } from './browse-item';
import { Field, FieldImportance, FieldRenderer } from './field';

const fooField = new Field({
  name: 'fooField',
  renderer: FieldRenderer.plain,
  importance: FieldImportance.mandatory,
  isEnumerable: false,
  isVirtual: false,
});

const barField = new Field({
  name: 'barField',
  renderer: FieldRenderer.plain,
  importance: FieldImportance.recommended,
  isEnumerable: false,
  isVirtual: false,
});

const TAB_HIERARCHY: BrowseItemConfig = {
  key: fooField,
  axis: { name: fooField.name, uiField: fooField },
  type: BrowseItemConfigType.hierarchy,
  entityType: 'OrganizationalUnit',
  linkField: barField,
  displayField: Field.createPending('bazField'),
  minLevel: 1,
  maxLevel: 2,
};

const TAB_FACET: BrowseItemConfig = {
  key: barField,
  axis: { name: barField.name, uiField: barField },
  type: BrowseItemConfigType.facet,
};

const BUCKET_1: SearchResultsFacetBucket = {
  count: 1,
  value: 'bar-id',
};

const BUCKET_2: SearchResultsFacetBucket = {
  count: 4,
  value: 'boo-id',
};

const BUCKET_3: SearchResultsFacetBucket = {
  count: 0,
  value: `${config.FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX}far`,
};

const BUCKET_4: SearchResultsFacetBucket = {
  count: 0,
  value: 'Something else',
};

const BUCKET_5: SearchResultsFacetBucket = {
  count: 6,
  value: 'foo-id',
};

const BUCKET_6: SearchResultsFacetBucket = {
  count: 3,
  value: 'baz-id',
};

const FACET: SearchResultsFacet = {
  type: 'exact',
  axis: 'fooField',
  bucketNo: 4,
  buckets: [BUCKET_1, BUCKET_2, BUCKET_3, BUCKET_4],
};

const createDisplay = (display: string, language?: string, place?: number) => {
  return { display, language, place };
};

const NODE_1: HierarchyNode = {
  depth: 0,
  nodeId: 'root-id',
  display: [createDisplay('Root Node')],
};

const NODE_2: HierarchyNode = {
  depth: 1,
  nodeId: 'foo-id',
  display: [createDisplay('Foo Node')],
  parentNodeId: 'root-id',
};

const NODE_3: HierarchyNode = {
  depth: 2,
  nodeId: 'bar-id',
  display: [createDisplay('Bar Node', 'de')],
  parentNodeId: 'foo-id',
};

const NODE_4: HierarchyNode = {
  depth: 2,
  nodeId: 'baz-id',
  parentNodeId: 'foo-id',
};

const NODE_5: HierarchyNode = {
  depth: 3,
  nodeId: 'boo-id',
  display: [
    createDisplay('Boo Node', null, 1),
    createDisplay('Boo Node DE', 'de', 0),
    createDisplay('Boo Node EN', 'en', 0),
  ],
  parentNodeId: 'baz-id',
};

const HIERARCHY: Hierarchy = {
  key: null,
  nodes: [NODE_1, NODE_2, NODE_3, NODE_4, NODE_5],
};

const FACET_WITH_HIERARCHY: SearchResultsFacet = {
  ...FACET,
  buckets: [BUCKET_1, BUCKET_2, BUCKET_3, BUCKET_4, BUCKET_5, BUCKET_6].map((bucket) => {
    const node = HIERARCHY.nodes.find(({ nodeId }) => nodeId === bucket.value);

    return {
      ...bucket,
      hierarchyInfo: {
        '@type': 'type.googleapis.com/mex.v0.HierarchyInfo',
        parentValue: node?.parentNodeId ?? '',
        display: node?.display?.[0]?.display ?? '',
        depth: node?.depth ?? 0,
      },
    };
  }),
};

const BROWSE_ITEM_1 = new BrowseItem({
  config: TAB_HIERARCHY,
  facet: FACET,
  node: NODE_1,
  hierarchy: HIERARCHY,
});

const BROWSE_ITEM_2 = new BrowseItem({
  config: TAB_HIERARCHY,
  facet: FACET,
  node: NODE_2,
  hierarchy: HIERARCHY,
});

const BROWSE_ITEM_2_WITH_SINGLE_NODE_CHILD = BROWSE_ITEM_2.clone({
  config: { ...TAB_HIERARCHY, enableSingleNodeVersion: true },
  facet: FACET_WITH_HIERARCHY,
});

const BROWSE_ITEM_2_SINGLE_NODE = BROWSE_ITEM_2.clone({
  isSingleNodeVersion: true,
  facet: FACET_WITH_HIERARCHY,
});

const BROWSE_ITEM_3 = new BrowseItem({
  config: TAB_HIERARCHY,
  facet: FACET,
  node: NODE_3,
  hierarchy: HIERARCHY,
});

const BROWSE_ITEM_4 = new BrowseItem({
  config: TAB_HIERARCHY,
  facet: FACET,
  node: NODE_4,
  hierarchy: HIERARCHY,
});

const BROWSE_ITEM_5 = new BrowseItem({
  config: TAB_HIERARCHY,
  facet: FACET,
  node: NODE_5,
  hierarchy: HIERARCHY,
});

const BROWSE_ITEM_FACET_BUCKET = new BrowseItem({
  config: TAB_FACET,
  bucket: BUCKET_2,
});

const BROWSE_ITEM_FACET_BUCKET_2 = new BrowseItem({
  config: TAB_FACET,
  bucket: BUCKET_3,
});

const BROWSE_ITEM_FACET_BUCKET_3 = new BrowseItem({
  config: TAB_FACET,
  bucket: BUCKET_4,
});

describe('BrowseItem model', () => {
  it("returns the node's level (default: 0)", () => {
    expect(BROWSE_ITEM_1.level).toBe(0);
    expect(BROWSE_ITEM_2.level).toBe(1);
    expect(BROWSE_ITEM_3.level).toBe(2);
  });

  it('returns the min level of the tab (default: 0)', () => {
    expect(BROWSE_ITEM_1.minLevel).toBe(1);
    expect(BROWSE_ITEM_2.minLevel).toBe(1);
    expect(BROWSE_ITEM_FACET_BUCKET.minLevel).toBe(0);
  });

  it('returns the max level of the tab (default: 0)', () => {
    expect(BROWSE_ITEM_1.maxLevel).toBe(2);
    expect(BROWSE_ITEM_2.maxLevel).toBe(2);
    expect(BROWSE_ITEM_FACET_BUCKET.maxLevel).toBe(0);
  });

  it('returns the parent node if exists', () => {
    expect(BROWSE_ITEM_1.parentNode).toBe(undefined);
    expect(BROWSE_ITEM_2.parentNode).toBe(NODE_1);
    expect(BROWSE_ITEM_3.parentNode).toBe(NODE_2);
    expect(BROWSE_ITEM_4.parentNode).toBe(NODE_2);
    expect(BROWSE_ITEM_5.parentNode).toBe(NODE_4);
    expect(BROWSE_ITEM_FACET_BUCKET.parentNode).toBe(undefined);
  });

  it('returns its own node as parent node if instance is a single node version', () => {
    expect(BROWSE_ITEM_2_SINGLE_NODE.parentNode).toBe(NODE_2);
  });

  it('creates a parent browse item from the parent node', () => {
    expect(BROWSE_ITEM_1.parent).toBe(undefined);
    expect(BROWSE_ITEM_2.parent).toStrictEqual(BROWSE_ITEM_1);
    expect(BROWSE_ITEM_3.parent).toStrictEqual(BROWSE_ITEM_2);
  });

  it('creates an unmodified version of itself as parent browse item if instance is a single node version', () => {
    expect(BROWSE_ITEM_2_SINGLE_NODE.value).not.toBe(BROWSE_ITEM_2.value);
    expect(BROWSE_ITEM_2_SINGLE_NODE.parent.value).toBe(BROWSE_ITEM_2.value);
  });

  it('returns all parents of an item in reversed order, respecting the min level', () => {
    expect(BROWSE_ITEM_1.parents).toEqual([]);
    expect(BROWSE_ITEM_2.parents).toStrictEqual([]);
    expect(BROWSE_ITEM_3.parents).toStrictEqual([BROWSE_ITEM_2]);
    expect(BROWSE_ITEM_4.parents).toStrictEqual([BROWSE_ITEM_2]);
    expect(BROWSE_ITEM_5.parents).toStrictEqual([BROWSE_ITEM_2, BROWSE_ITEM_4]);
  });

  it('returns the child nodes', () => {
    expect(BROWSE_ITEM_1.childNodes).toStrictEqual([NODE_2]);
    expect(BROWSE_ITEM_2.childNodes).toStrictEqual([NODE_3, NODE_4]);
    expect(BROWSE_ITEM_3.childNodes).toEqual([]);
    expect(BROWSE_ITEM_4.childNodes).toStrictEqual([NODE_5]);
    expect(BROWSE_ITEM_5.childNodes).toEqual([]);
    expect(BROWSE_ITEM_FACET_BUCKET.childNodes).toEqual([]);
  });

  it('returns a modified version of itself as the single node if applicable (enabled in config, can descend, has child nodes; fallback: null)', () => {
    expect(BROWSE_ITEM_2.singleNodeVersion).toBe(null);
    expect(BROWSE_ITEM_2_WITH_SINGLE_NODE_CHILD.singleNodeVersion).not.toBe(null);
    expect(BROWSE_ITEM_2_WITH_SINGLE_NODE_CHILD.singleNodeVersion.value).toBe(BROWSE_ITEM_2_SINGLE_NODE.value);
    expect(
      BROWSE_ITEM_2.clone({
        config: { ...TAB_HIERARCHY, maxLevel: 0, enableSingleNodeVersion: true },
      }).singleNodeVersion
    ).toBe(null); // can not descend
    expect(
      BROWSE_ITEM_3.clone({
        config: { ...TAB_HIERARCHY, enableSingleNodeVersion: true },
      }).singleNodeVersion
    ).toBe(null); // has no children
  });

  it('creates child browse items from the child nodes', () => {
    expect(BROWSE_ITEM_1.children).toStrictEqual([BROWSE_ITEM_2]);
    expect(BROWSE_ITEM_2.children).toStrictEqual([BROWSE_ITEM_3, BROWSE_ITEM_4]);
    expect(BROWSE_ITEM_3.children).toEqual([]);
  });

  it('creates child browse items plus a single node version of itself if applicable (enabled in config, can descend, has count, has at least one child node with count)', () => {
    expect(BROWSE_ITEM_2.children).toStrictEqual([BROWSE_ITEM_3, BROWSE_ITEM_4]);
    expect(BROWSE_ITEM_2_WITH_SINGLE_NODE_CHILD.children).toStrictEqual([
      BROWSE_ITEM_2_SINGLE_NODE,
      BROWSE_ITEM_3,
      BROWSE_ITEM_4,
    ]);
  });

  it('returns if there are children', () => {
    expect(BROWSE_ITEM_1.hasChildren).toBe(true);
    expect(BROWSE_ITEM_3.hasChildren).toBe(false);
    expect(BROWSE_ITEM_FACET_BUCKET.hasChildren).toBe(false);
  });

  it('returns the childless descendants', () => {
    expect(BROWSE_ITEM_1.childlessDescendants).toStrictEqual([BROWSE_ITEM_3, BROWSE_ITEM_5]);
    expect(BROWSE_ITEM_2.childlessDescendants).toStrictEqual([BROWSE_ITEM_3, BROWSE_ITEM_5]);
    expect(BROWSE_ITEM_3.childlessDescendants).toStrictEqual([]);
    expect(BROWSE_ITEM_4.childlessDescendants).toStrictEqual([BROWSE_ITEM_5]);
    expect(BROWSE_ITEM_5.childlessDescendants).toStrictEqual([]);
    expect(BROWSE_ITEM_FACET_BUCKET.childlessDescendants).toStrictEqual([]);
  });

  describe('canDescend getter', () => {
    it('returns false if allowDescent flag is set to false', () => {
      expect(BROWSE_ITEM_1.canDescend).toBe(true);
      expect(BROWSE_ITEM_1.clone({ allowDescent: false }).canDescend).toBe(false);
    });

    it('returns false if tab type is not hierarchy', () => {
      expect(BROWSE_ITEM_1.canDescend).toBe(true);
      expect(BROWSE_ITEM_1.clone({ config: TAB_FACET }).canDescend).toBe(false);
    });

    it("returns false if tab's max level is lower than child level", () => {
      expect(BROWSE_ITEM_2.canDescend).toBe(true);
      expect(BROWSE_ITEM_2.clone({ config: { ...TAB_HIERARCHY, maxLevel: 1 } }).canDescend).toBe(false);

      expect(BROWSE_ITEM_4.canDescend).toBe(false);
      expect(BROWSE_ITEM_4.clone({ config: { ...TAB_HIERARCHY, maxLevel: 3 } }).canDescend).toBe(true);
    });

    it('returns false if the item has no children', () => {
      expect(BROWSE_ITEM_2.canDescend).toBe(true);
      expect(BROWSE_ITEM_3.clone({ config: { ...TAB_HIERARCHY, maxLevel: 3 } }).canDescend).toBe(false);
    });
  });

  it('returns whether or not the instance is a single node version', () => {
    expect(BROWSE_ITEM_2.isSingleNodeVersion).toBe(false);
    expect(BROWSE_ITEM_2_WITH_SINGLE_NODE_CHILD.isSingleNodeVersion).toBe(false);
    expect(BROWSE_ITEM_2_SINGLE_NODE.isSingleNodeVersion).toBe(true);
    expect(BROWSE_ITEM_2_WITH_SINGLE_NODE_CHILD.singleNodeVersion?.isSingleNodeVersion).toBe(true);
  });

  describe('url getter', () => {
    beforeAll(() => {
      const originalRoutes = config.ROUTES;
      jest.spyOn(config, 'ROUTES', 'get').mockReturnValue({ ...originalRoutes, SEARCH: '/search-route' });
      jest.spyOn(config, 'SEARCH_PARAM_FILTER_PREFIX', 'get').mockReturnValue('');
    });

    it('returns null if canDescend is true', () => {
      expect(BROWSE_ITEM_1.url).toBe(null);
      expect(BROWSE_ITEM_1.clone({ allowDescent: false }).url).not.toBe(null);
    });

    it('returns search route with axis constraint or constraints for each childless descendant', () => {
      expect(BROWSE_ITEM_3.url).toBe('/search-route?fooField[]=bar-id');
      expect(BROWSE_ITEM_4.url).toBe('/search-route?fooField[]=boo-id');
      expect(BROWSE_ITEM_1.clone({ allowDescent: false }).url).toBe(
        '/search-route?fooField[]=bar-id&fooField[]=boo-id'
      );
      expect(BROWSE_ITEM_FACET_BUCKET.url).toBe('/search-route?barField[]=boo-id');
      expect(BROWSE_ITEM_FACET_BUCKET_2.url).toBe('/search-route?barField[]=far');
      expect(BROWSE_ITEM_FACET_BUCKET_3.url).toBe('/search-route?barField[]=Something+else');
      expect(
        BROWSE_ITEM_4.clone({
          hierarchy: {
            key: null,
            nodes: [NODE_1, NODE_2, NODE_3, NODE_4, { ...NODE_5, display: null }],
          },
        }).url
      ).toBe('/search-route?fooField[]=boo-id');
    });

    it('returns search route with only the bucket as axis constraint if hierarchy info facet', () => {
      expect(BROWSE_ITEM_1.clone({ allowDescent: false, facet: FACET_WITH_HIERARCHY }).url).toBe(
        '/search-route?fooField[]=root-id'
      );
    });
  });

  it('returns the facet bucket attached or matching for the item', () => {
    expect(BROWSE_ITEM_1.bucket).toBe(undefined);
    expect(BROWSE_ITEM_2.bucket).toBe(undefined);
    expect(BROWSE_ITEM_3.bucket).toBe(BUCKET_1);
    expect(BROWSE_ITEM_4.bucket).toBe(undefined);
    expect(BROWSE_ITEM_5.bucket).toBe(BUCKET_2);
  });

  it('returns if at least one of the facet buckets has a hierarchy info attached', () => {
    expect(BROWSE_ITEM_1.hasHierarchyInfo).toBe(false);
    expect(BROWSE_ITEM_1.clone({ facet: FACET_WITH_HIERARCHY }).hasHierarchyInfo).toBe(true);
  });

  it('returns the formatted text to be displayed in the tile as well as the value to be used as axis constraint (from node or bucket)', () => {
    config.FEATURE_FLAGS.BROWSE_HIERARCHY_TRANSLATIONS = false;

    expect(BROWSE_ITEM_1.text).toBe('Root Node');
    expect(BROWSE_ITEM_1.value).toBe('root-id');
    expect(BROWSE_ITEM_2.text).toBe('Foo Node');
    expect(BROWSE_ITEM_2.value).toBe('foo-id');
    expect(BROWSE_ITEM_3.text).toBe('Bar Node');
    expect(BROWSE_ITEM_3.value).toBe('bar-id');
    expect(BROWSE_ITEM_4.text).toBe('baz-id');
    expect(BROWSE_ITEM_4.value).toBe('baz-id');
    expect(BROWSE_ITEM_5.text).toBe('Boo Node DE, Boo Node EN, Boo Node');
    expect(BROWSE_ITEM_5.value).toBe('boo-id');

    config.FEATURE_FLAGS.BROWSE_HIERARCHY_TRANSLATIONS = true;
    stores.i18n.language = 'de';
    expect(BROWSE_ITEM_5.text).toBe('Boo Node DE, Boo Node');
    stores.i18n.language = 'en';
    expect(BROWSE_ITEM_5.text).toBe('Boo Node EN, Boo Node');
    stores.i18n.language = 'fr';
    expect(BROWSE_ITEM_5.text).toBe('Boo Node');

    expect(BROWSE_ITEM_FACET_BUCKET.text).toBe('boo-id');
    expect(BROWSE_ITEM_FACET_BUCKET.value).toBe('boo-id');
    expect(BROWSE_ITEM_FACET_BUCKET_2.text).toBe('far');
    expect(BROWSE_ITEM_FACET_BUCKET_2.value).toBe(`${config.FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX}far`);
    expect(BROWSE_ITEM_FACET_BUCKET_3.text).toBe('Something else');
    expect(BROWSE_ITEM_FACET_BUCKET_3.value).toBe('Something else');
  });

  it('returns a special text and value indicating that a single node version will be selected', () => {
    jest.spyOn(stores.i18n, 't').mockImplementation((key, options: any) => {
      if (key === config.SEARCH_FACETS_SINGLE_NODE_LABEL_KEY) {
        return `${options.text} (direct)`;
      }

      return key;
    });
    jest.spyOn(config, 'SEARCH_FACETS_SINGLE_NODE_VALUE_PREFIX', 'get').mockReturnValue('direct:');

    expect(BROWSE_ITEM_2_SINGLE_NODE.text).toBe('Foo Node (direct)');
    expect(BROWSE_ITEM_2_SINGLE_NODE.value).toBe('direct:foo-id');
  });

  it("returns the count of matches including descendants'", () => {
    expect(BROWSE_ITEM_1.count).toBe(5);
    expect(BROWSE_ITEM_2.count).toBe(5);
    expect(BROWSE_ITEM_3.count).toBe(1);
    expect(BROWSE_ITEM_4.count).toBe(4);
    expect(BROWSE_ITEM_5.count).toBe(4);
    expect(BROWSE_ITEM_FACET_BUCKET.count).toBe(4);
  });

  it("returns directly the bucket's counts if buckets have hierarchy info", () => {
    const HIERARCHY_BROWSE_ITEM_1 = BROWSE_ITEM_1.clone({ facet: FACET_WITH_HIERARCHY });
    const HIERARCHY_BROWSE_ITEM_2 = BROWSE_ITEM_2.clone({ facet: FACET_WITH_HIERARCHY });
    const HIERARCHY_BROWSE_ITEM_3 = BROWSE_ITEM_3.clone({ facet: FACET_WITH_HIERARCHY });
    const HIERARCHY_BROWSE_ITEM_4 = BROWSE_ITEM_4.clone({ facet: FACET_WITH_HIERARCHY });
    const HIERARCHY_BROWSE_ITEM_5 = BROWSE_ITEM_5.clone({ facet: FACET_WITH_HIERARCHY });

    expect(HIERARCHY_BROWSE_ITEM_1.count).toBe(0);
    expect(HIERARCHY_BROWSE_ITEM_2.count).toBe(6);
    expect(HIERARCHY_BROWSE_ITEM_3.count).toBe(1);
    expect(HIERARCHY_BROWSE_ITEM_4.count).toBe(3);
    expect(HIERARCHY_BROWSE_ITEM_5.count).toBe(4);
  });

  it('returns the single node count if applicable', () => {
    expect(BROWSE_ITEM_2_SINGLE_NODE.count).toBe(2); // 6 - (3 + 1)
    expect(BROWSE_ITEM_2_SINGLE_NODE.count).not.toBe(BROWSE_ITEM_2_WITH_SINGLE_NODE_CHILD.count);
  });

  it('returns the node ID if exists', () => {
    expect(BROWSE_ITEM_1.nodeId).toBe('root-id');
    expect(BROWSE_ITEM_2.nodeId).toBe('foo-id');
    expect(BROWSE_ITEM_3.nodeId).toBe('bar-id');
    expect(BROWSE_ITEM_4.nodeId).toBe('baz-id');
    expect(BROWSE_ITEM_5.nodeId).toBe('boo-id');
    expect(BROWSE_ITEM_FACET_BUCKET.nodeId).toBe(undefined);
  });

  it('allows to clone the browse item, overriding individual properties', () => {
    const clonedItem1 = BROWSE_ITEM_1.clone({});
    expect(clonedItem1).not.toBe(BROWSE_ITEM_1);
    expect(clonedItem1).toEqual(BROWSE_ITEM_1);

    const clonedItem2 = BROWSE_ITEM_1.clone({ node: { depth: 9, nodeId: 'far-id' } });
    expect(clonedItem2.bucket).toBe(BROWSE_ITEM_1.bucket);
    expect(clonedItem2.value).toBe('far-id');
  });
});
