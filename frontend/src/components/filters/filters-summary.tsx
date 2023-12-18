import { Component, Event, EventEmitter, h, Host, Prop } from '@stencil/core';
import { ANALYTICS_CUSTOM_EVENTS, FIELDS, SEARCH_CONFIG, SEARCH_INVISIBLE_FACETS } from 'config';
import { BrowseItemConfigType } from 'config/browse';
import { SearchFacet } from 'config/search';
import { BrowseItem } from 'models/browse-item';
import { Field } from 'models/field';
import services from 'services';
import stores from 'stores';
import { translateFieldName } from 'utils/field';

@Component({
  tag: 'mex-filters-summary',
  styleUrl: 'filters-summary.css',
})
export class FiltersSummaryComponent {
  @Prop() label: string;
  @Prop() testAttr?: string;

  @Event() updateSearch: EventEmitter;

  get hierarchies() {
    return stores.search.hierarchies;
  }

  getBrowseItems(facet: SearchFacet, field: Field) {
    const config = {
      ...facet,
      key: field,
      type: facet.type === 'hierarchy' ? BrowseItemConfigType.hierarchy : BrowseItemConfigType.facet,
      enableSingleNodeVersion: facet.type === 'hierarchy',
    };
    const facetResponse = stores.search.facets.find((item) => item.axis === facet.axis.name);
    const buckets = facet && facetResponse?.buckets;
    const hierarchy = this.hierarchies[field.name];

    return (
      (hierarchy
        ? hierarchy.nodes.map((node) => new BrowseItem({ config, facet: facetResponse, node, hierarchy }))
        : buckets.map((bucket) => new BrowseItem({ config, facet: facetResponse, bucket }))) ?? []
    );
  }

  get items() {
    const { all } = stores.filters;
    const facets = SEARCH_CONFIG.FACETS.concat(SEARCH_INVISIBLE_FACETS());

    return all.reduce((items, [filterName, values]) => {
      const facet = facets.find(({ axis }) => axis.name === filterName);
      const field = facet.axis.uiField ?? FIELDS[facet.axis.name];
      const browseItems = this.getBrowseItems(facet, field)
        .flatMap((browseItem) => [browseItem, browseItem.singleNodeVersion])
        .filter(Boolean);

      return field
        ? items.concat(
            values.map((value) => {
              const browseItem = browseItems.find((browseItem) => browseItem.value === value);

              return {
                filterName,
                value,
                text: `${translateFieldName(field)}: ${browseItem?.text ?? value}`,
              };
            })
          )
        : items;
    }, []);
  }

  removeFilter(fieldName: string, value: string) {
    stores.filters.remove(fieldName, value);
    this.updateSearch.emit();

    services.analytics.trackEvent(...ANALYTICS_CUSTOM_EVENTS.SEARCH_FILTER, `Removed: ${fieldName}`);
  }

  render() {
    const { label, testAttr } = this;

    return (
      <Host class="filters-summary">
        {label}
        {this.items.map((item, index) => (
          <mex-tag
            key={item.text}
            text={item.text}
            handleClose={() => this.removeFilter(item.filterName, item.value)}
            closeTitle={stores.i18n.t('filters.remove', { name: item.text })}
            data-test={testAttr}
            data-test-key={index}
            testAttr={testAttr}
          />
        ))}
      </Host>
    );
  }
}
