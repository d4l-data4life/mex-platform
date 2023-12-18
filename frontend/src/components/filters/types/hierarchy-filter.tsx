import { Component, Event, EventEmitter, Fragment, h, Prop, State } from '@stencil/core';

import { ANALYTICS_CUSTOM_EVENTS } from 'config';
import { BrowseItemConfigType } from 'config/browse';
import { SearchFacet } from 'config/search';
import { BrowseItem } from 'models/browse-item';
import services from 'services';
import stores from 'stores';
import { Hierarchy, SearchResultsFacet } from 'stores/search';
import { CheckStatusEnum } from 'utils/filters';
import { getHierarchyNodesFirstLevel } from 'utils/search';

@Component({
  tag: 'mex-hierarchy-filter',
  styleUrl: 'hierarchy-filter.css',
})
export class HierarchyFilterComponent {
  @Prop() busy: boolean;
  @Prop() expanded: boolean;
  @Prop() facet: SearchFacet;
  @Prop() facetResponse?: SearchResultsFacet;
  @Prop() handleChange?: () => void;

  @State() hierarchy?: Hierarchy;
  @State() expandedNodeStates: { [key: string]: boolean } = {};

  @Event() updateSearch: EventEmitter;

  get axis() {
    return this.facet.axis;
  }

  get canToggle() {
    return this.expanded && !this.busy;
  }

  get hasItems() {
    return !!this.items?.length;
  }

  get isEmpty() {
    return !this.busy && (!this.hasItems || this.items.every(({ count }) => !count));
  }

  get items(): BrowseItem[] {
    const { facet, facetResponse, hierarchy } = this;
    if (!facetResponse) {
      return [];
    }

    const firstLevel = getHierarchyNodesFirstLevel(hierarchy?.nodes, facet.minLevel);

    const config = {
      ...facet,
      key: facet.axis.uiField,
      type: BrowseItemConfigType.hierarchy,
      enableSingleNodeVersion: true,
    };

    return (
      hierarchy?.nodes
        .filter(({ depth }) => depth === firstLevel)
        .map((node) => new BrowseItem({ config, facet: facetResponse, node, hierarchy })) ?? []
    );
  }

  get state() {
    return stores.filters.get(this.axis.name);
  }

  isExpandedNode(item: BrowseItem) {
    return (
      this.expandedNodeStates[item.value] ??
      (this.isFullyChecked(item) || item.descendants.some((descendant) => this.isFullyChecked(descendant)))
    );
  }

  isFullyChecked(item: BrowseItem, toggledItem?: BrowseItem, toggledItemChecked?: boolean) {
    const { state } = this;
    const { value, parents, children } = item;

    if (value === toggledItem?.value || parents.find((parent) => parent.value === toggledItem?.value)) {
      return toggledItemChecked;
    }

    const isChecked = state.includes(value) || parents.some((parent) => state.includes(parent.value));
    if (!children.length) {
      return isChecked;
    }

    const fullyCheckedChildrenCount = children.filter((child) =>
      this.isFullyChecked(child, toggledItem, toggledItemChecked)
    ).length;
    return fullyCheckedChildrenCount === children.length;
  }

  regenerateFilterValues(items: BrowseItem[] = this.items, toggledItem?: BrowseItem, toggledItemChecked?: boolean) {
    return items.reduce((values, item) => {
      if (this.isFullyChecked(item, toggledItem, toggledItemChecked)) {
        return values.concat([item.value]);
      }

      return values.concat(this.regenerateFilterValues(item.children, toggledItem, toggledItemChecked));
    }, []);
  }

  toggle(item: BrowseItem, previouslyChecked: CheckStatusEnum) {
    const {
      axis: { name },
    } = this;

    const isInputChecked = [CheckStatusEnum.SEMI, CheckStatusEnum.UNCHECKED].includes(previouslyChecked);
    stores.filters.set(name, this.regenerateFilterValues(this.items, item, isInputChecked));

    this.updateSearch.emit();

    services.analytics.trackEvent(
      ...ANALYTICS_CUSTOM_EVENTS.SEARCH_FILTER,
      `Switched ${isInputChecked ? 'on' : 'off'}: ${name}`
    );
  }

  async loadHierarchy() {
    const { axis, entityType, linkField, displayField } = this.facet;
    this.hierarchy = await services.search.fetchHierarchy(axis.uiField, entityType, linkField, displayField);
  }

  componentWillLoad() {
    this.loadHierarchy();
  }

  renderCheckbox(item: BrowseItem, isChecked = false, isParentExpanded = true) {
    const { axis, state } = this;
    const { text, value, count } = item;
    const isSemiChecked = !isChecked && item.descendants.some(({ value }) => this.state.includes(value));

    return (
      <mex-checkbox
        key={`${axis}--${value}`}
        classes="hierarchy-filter__item"
        label={text}
        secondaryText={`(${count})`}
        handleChange={(previouslyChecked) => this.toggle(item, previouslyChecked)}
        checked={isChecked ? CheckStatusEnum.CHECKED : isSemiChecked ? CheckStatusEnum.SEMI : CheckStatusEnum.UNCHECKED}
        disabled={!this.canToggle || !isParentExpanded || (!count && !state.includes(value))}
        data-test="filter:item"
        data-test-context={axis.name}
        data-test-active={String(state.includes(value))}
        testAttr="filter:item:hierarchy"
      />
    );
  }

  renderExpandableParent(item: BrowseItem, isChecked = false, isExpanded = true, isParentExpanded = true) {
    const isVisibleAndExpanded = isParentExpanded && isExpanded;

    return (
      <div class="hierarchy-filter__parent">
        <div class="hierarchy-filter__parent-label">{this.renderCheckbox(item, isChecked, isParentExpanded)}</div>
        <button
          class="hierarchy-filter__parent-button"
          onClick={() => {
            this.expandedNodeStates = {
              ...this.expandedNodeStates,
              [item.value]: !this.isExpandedNode(item),
            };

            services.analytics.trackEvent(
              ...ANALYTICS_CUSTOM_EVENTS.SEARCH_NAVIGATION,
              this.isExpandedNode(item) ? 'Hierarchy filter subtree expanded' : 'Hierarchy filter subtree collapsed'
            );
          }}
          disabled={!isParentExpanded || !this.canToggle}
          aria-expanded={String(isVisibleAndExpanded)}
          data-test="hierarchy-filter:toggle"
          data-test-active={String(isVisibleAndExpanded)}
        >
          <mex-icon-chevron
            classes={`hierarchy-filter__button-chevron icon--inline icon--mirrorable ${
              isExpanded ? 'icon--mirrored-vertical' : ''
            }`}
          />
        </button>
      </div>
    );
  }

  renderHierarchy(item: BrowseItem, isParentChecked = false, isParentExpanded = true) {
    if (!item.bucket) {
      return;
    }

    const canDescend = item.canDescend && item.children.some((childItem) => childItem?.bucket);
    const isChecked = isParentChecked || this.state.includes(item.value);
    const isExpanded = this.isExpandedNode(item);

    return (
      <li>
        {!canDescend && this.renderCheckbox(item, isChecked, isParentExpanded)}

        {canDescend && (
          <Fragment>
            {this.renderExpandableParent(item, isChecked, isExpanded, isParentExpanded)}
            <mex-accordion expanded={isExpanded}>
              <div class="hierarchy-filter__children-container">
                <ul class="hierarchy-filter__children">
                  {item.children.map((childItem) => this.renderHierarchy(childItem, isChecked, isExpanded))}
                </ul>
              </div>
            </mex-accordion>
          </Fragment>
        )}
      </li>
    );
  }

  render() {
    const { busy, hasItems, isEmpty, items } = this;
    return (
      <div class={{ 'hierarchy-filter': true, 'hierarchy-filter--empty': isEmpty }}>
        {busy && !this.facetResponse && <mex-placeholder lines={2} />}
        {isEmpty && stores.i18n.t('filters.empty')}
        {hasItems && <ul>{items.map((item) => this.renderHierarchy(item))}</ul>}
      </div>
    );
  }
}
