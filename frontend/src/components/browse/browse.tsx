import { Component, Fragment, h, Host, State } from '@stencil/core';
import {
  ANALYTICS_CUSTOM_EVENTS,
  BROWSE_CONFIG,
  BROWSE_FACETS_LIMIT,
  SEARCH_CONFIG,
  SEARCH_QUERY_EVERYTHING,
} from 'config';
import { BrowseItemConfig, BrowseItemConfigType } from 'config/browse';
import { BrowseItem } from 'models/browse-item';
import services from 'services';
import stores from 'stores';
import { Hierarchy, SearchResultsFacet } from 'stores/search';
import { catchRetryableAction } from 'utils/error';
import { translateFieldName } from 'utils/field';
import { getHierarchyNodesFirstLevel } from 'utils/search';

@Component({
  tag: 'mex-browse',
  styleUrl: 'browse.css',
})
export class BrowseComponent {
  @State() selectedTab: BrowseItemConfig = BROWSE_CONFIG.TABS[0];
  @State() facets?: SearchResultsFacet[];
  @State() hierarchies?: Hierarchy[];
  @State() selectedItem?: BrowseItem;

  get isBusy() {
    return !this.facets || (this.selectedTab?.type === BrowseItemConfigType.hierarchy && !this.hierarchies);
  }

  get items(): BrowseItem[] {
    const { selectedTab: tab, selectedItem } = this;

    const facet = this.facets?.find((facet) => facet.axis === tab.axis.name);
    const hierarchy = this.hierarchies?.find((hierarchy) => hierarchy.key === tab.key);

    if (tab.type === BrowseItemConfigType.hierarchy) {
      const firstLevel = getHierarchyNodesFirstLevel(hierarchy?.nodes, tab.minLevel);

      return (
        selectedItem?.children ??
        hierarchy?.nodes
          .filter(({ depth }) => depth === firstLevel)
          .map((node) => new BrowseItem({ config: tab, facet, node, hierarchy })) ??
        []
      );
    }

    return (
      this.facets
        ?.find((facet) => facet.axis === tab.axis.name)
        ?.buckets?.map((bucket) => new BrowseItem({ config: tab, bucket, hierarchy })) ?? []
    );
  }

  get breadcrumbs() {
    const { selectedItem } = this;
    if (!selectedItem) {
      return null;
    }

    const base: { label: string; handleClick?: () => void; icon?: string } = {
      label: stores.i18n.t('browse.hierarchy.all'),
      handleClick: () => (this.selectedItem = null),
      icon: 'back',
    };

    return [base]
      .concat(selectedItem.parents.map((item) => ({ label: item.text, handleClick: () => (this.selectedItem = item) })))
      .concat([{ label: selectedItem.text }])
      .map((breadcrumb, index) =>
        Object.assign(breadcrumb, { testAttr: 'browse:breadcrumbs:item', testKeyAttr: index })
      );
  }

  renderTile(item: BrowseItem, index?: number, customText?: string) {
    const { canDescend, count } = item;
    const disabled = !count;

    return (
      <mex-tile
        key={index}
        class="browse__item"
        text={customText ?? item.text}
        hint={stores.i18n.t(count === 1 ? 'browse.sourceCount' : 'browse.sourcesCount', {
          count,
        })}
        url={disabled ? null : item.url}
        icon={disabled ? null : canDescend ? 'drilldown' : 'arrow'}
        disabled={disabled}
        handleClick={() => {
          if (canDescend) {
            this.selectedItem = item;
          } else {
            services.analytics.trackEvent(
              ...ANALYTICS_CUSTOM_EVENTS.SEARCH_BROWSING,
              `Initiated: ${this.selectedTab.key.name}`
            );
          }
        }}
        testAttr="browse:grid:item"
        testKeyAttr={index}
      />
    );
  }

  componentDidLoad() {
    const facets = BROWSE_CONFIG.TABS.map(({ axis }) => axis)
      .map((axis) => SEARCH_CONFIG.FACETS.find((facet) => facet.axis.name === axis.name))
      .filter(Boolean);

    catchRetryableAction(async () => {
      this.facets = (
        await services.search.fetchResults({
          query: SEARCH_QUERY_EVERYTHING,
          limit: 0,
          facetsLimit: BROWSE_FACETS_LIMIT,
          fields: [],
          highlightFields: [],
          facets,
        })
      )?.facets;

      this.hierarchies = await Promise.all(
        BROWSE_CONFIG.TABS.filter(({ type }) => type === BrowseItemConfigType.hierarchy)
          .filter(({ axis }) => !!this.facets?.find((facet) => facet.axis === axis.name)?.bucketNo)
          .map((tab) => services.search.fetchHierarchy(tab.key, tab.entityType, tab.linkField, tab.displayField))
      );
    });
  }

  render() {
    const { selectedTab, selectedItem, breadcrumbs, isBusy } = this;
    const { t } = stores.i18n;

    return (
      !!selectedTab && (
        <Host class="browse view__wrapper" data-test="browse">
          <nav class="browse__nav" role="tablist" data-test="browse:nav">
            <h2 class="browse__title" data-test="browse:nav:title">
              {t('browse.title')}
            </h2>
            <mex-tabs
              class="browse__tabs"
              items={BROWSE_CONFIG.TABS.map((tab) => ({
                label: translateFieldName(tab.key),
                isActive: selectedTab === tab,
                data: tab,
              }))}
              handleClick={(item) => {
                this.selectedItem = null;
                this.selectedTab = item.data as BrowseItemConfig;
              }}
              testAttr="browse:nav:toggle"
            />
          </nav>

          {breadcrumbs && (
            <mex-breadcrumbs class="browse__breadcrumbs" items={breadcrumbs} data-test="browse:breadcrumbs" />
          )}

          <div class="browse__grid" data-test="browse:grid">
            {!isBusy && (
              <Fragment>
                {selectedItem &&
                  this.items.length > 1 &&
                  this.renderTile(
                    selectedItem.clone({ allowDescent: false }),
                    undefined,
                    t('browse.hierarchy.allOf', { text: selectedItem.text })
                  )}

                {this.items
                  .sort((a, b) => a.text.localeCompare(b.text))
                  .map((item, index) => this.renderTile(item, index))}

                {!this.items.length && (
                  <mex-not-found-state
                    class="browse__empty-state"
                    testAttr="browse:grid:empty"
                    caption={t('browse.empty.title')}
                    text={t('browse.empty.text', { name: translateFieldName(this.selectedTab.key) })}
                  />
                )}
              </Fragment>
            )}
            {isBusy &&
              new Array(SEARCH_CONFIG.FACETS_LIMIT)
                .fill(null)
                .map((_, index) => <mex-tile key={index} class="browse__item" isBusy />)}
          </div>
        </Host>
      )
    );
  }
}
