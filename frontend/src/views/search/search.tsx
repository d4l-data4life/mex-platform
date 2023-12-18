import { Component, h, Host, Listen, Prop, Watch } from '@stencil/core';
import {
  SEARCH_PAGINATION_START,
  SEARCH_QUERY_EVERYTHING,
  SEARCH_PARAM_SORTING_AXIS,
  SEARCH_PARAM_SORTING_ORDER,
  SEARCH_PARAM_PAGE,
  SEARCH_PARAM_FOCUS,
  SEARCH_CONFIG,
  ROUTES,
} from 'config';
import services from 'services';
import { match } from 'stencil-router-v2';
import stores from 'stores';
import { SearchResults } from 'stores/search';
import { catchRetryableAction } from 'utils/error';
import { getSearchUrl } from 'utils/search';

@Component({
  tag: 'mex-view-search',
})
export class SearchView {
  @Prop() url: URL;

  @Listen('updateSearch')
  updateSearchListener() {
    this.updateHistory();
  }

  @Watch('url')
  watchMatchHandler(newUrl: URL, oldUrl?: URL) {
    const reset = newUrl.pathname !== oldUrl?.pathname;
    const queryMatch = match(ROUTES.SEARCH_QUERY)(newUrl.pathname);
    const queryParams = new URLSearchParams(document.location.search);
    const page = parseInt(queryParams.get(SEARCH_PARAM_PAGE) ?? `${SEARCH_PAGINATION_START}`, 10);
    const query = decodeURIComponent(queryMatch?.query ?? '');
    const offset = (page - SEARCH_PAGINATION_START) * stores.search.limit;
    const sortingAxis = queryParams.get(SEARCH_PARAM_SORTING_AXIS);
    const sortingOrder = queryParams.get(SEARCH_PARAM_SORTING_ORDER);
    const focus = queryParams.get(SEARCH_PARAM_FOCUS);

    stores.search.sorting = SEARCH_CONFIG.SORTING_OPTIONS.find(
      ({ axis, order }) => ((!axis && !sortingAxis) || axis?.name === sortingAxis) && order === sortingOrder
    );
    stores.search.query = query;
    stores.search.focus = focus;

    stores.filters.fromQueryParams(queryParams);
    this.performSearch(query, reset, offset);
  }

  async performSearch(query: string, reset: boolean, offset = 0) {
    catchRetryableAction(
      async () => await services.search.updateResults({ query: query || SEARCH_QUERY_EVERYTHING, offset, reset }),
      true,
      (fromRetry: boolean) => {
        if (!fromRetry) {
          stores.search.response = { numFound: 0 } as SearchResults;
          stores.search.isBusy = false;
        }
      }
    );
  }

  async search(value: string, field = null) {
    stores.search.query = value;
    stores.search.focus = field;

    this.updateHistory(true);
  }

  async paginate(page: number) {
    stores.search.page !== page && this.updateHistory(false, page);
  }

  updateHistory(resetFilters = false, page = SEARCH_PAGINATION_START) {
    stores.router.push(getSearchUrl(resetFilters, page));
  }

  componentWillLoad() {
    this.watchMatchHandler(this.url);
  }

  render() {
    const { query } = stores.search;
    return (
      <Host class="view">
        <mex-search
          value={query}
          searchFocus={stores.search.focus}
          handleSearch={(value, field) => this.search(value, field)}
          handleReset={() => query && this.search('')}
          autofocus
        />
        <div class="view__wrapper view__wrapper--flex">
          <mex-results handlePagination={(page) => this.paginate(page)} />
        </div>
      </Host>
    );
  }
}
