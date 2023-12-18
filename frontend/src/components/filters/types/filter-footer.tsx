import { Component, Fragment, h, Prop, State } from '@stencil/core';
import { ANALYTICS_CUSTOM_EVENTS, SEARCH_QUERY_EVERYTHING } from 'config';
import { SearchFacet } from 'config/search';
import services from 'services';
import stores from 'stores';
import { catchRetryableAction } from 'utils/error';

@Component({
  tag: 'mex-filter-footer',
  styleUrl: 'filter-footer.css',
})
export class FilterFooterComponent {
  @Prop() busy: boolean;
  @Prop() facet: SearchFacet;
  @Prop() hasItems: boolean;
  @Prop() isLoadMoreDisabled: boolean;

  @State() isLoadingMore = false;

  get axis() {
    return this.facet.axis;
  }

  get showPlaceholder() {
    return this.isLoadingMore || (this.busy && !this.hasItems);
  }

  async loadMore() {
    const { axis } = this;

    this.isLoadingMore = true;
    const response = await services.search.fetchResults({
      query: stores.search.query || SEARCH_QUERY_EVERYTHING,
      offset: 0,
      limit: 0,
      filters: stores.filters.all,
      facets: [this.facet],
      facetsOffset: stores.search.getFacetOffset(axis.name),
    });

    stores.search.addFacetBuckets(axis.name, response?.facets.find((facet) => facet.axis === axis.name)?.buckets);

    this.isLoadingMore = false;

    services.analytics.trackEvent(...ANALYTICS_CUSTOM_EVENTS.SEARCH_NAVIGATION, 'Filter extended');
  }

  render() {
    const { axis, hasItems: hasBuckets, isLoadMoreDisabled } = this;

    return (
      <Fragment>
        {this.showPlaceholder && <mex-placeholder lines={2} />}

        {hasBuckets && stores.search.hasMoreBuckets(axis.name) && (
          <button
            class="filter-footer__more"
            disabled={isLoadMoreDisabled}
            onClick={() => catchRetryableAction(async () => await this.loadMore())}
            data-test="filter:more"
            data-test-context={axis}
          >
            {stores.i18n.t('filters.more')}
          </button>
        )}
      </Fragment>
    );
  }
}
