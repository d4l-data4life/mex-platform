import { Component, Event, EventEmitter, h, Host, Prop } from '@stencil/core';
import { href } from 'stencil-router-v2';
import { DropdownOption } from 'components/dropdown/dropdown';
import { ANALYTICS_CUSTOM_EVENTS, SEARCH_CONFIG, SEARCH_PAGINATION_START } from 'config';
import services from 'services';
import stores from 'stores';
import { getSearchUrl } from 'utils/search';

@Component({
  tag: 'mex-results',
  styleUrl: 'results.css',
})
export class ResultsComponent {
  #sortingOptions: DropdownOption[];

  @Prop() handlePagination?: (page: number) => void;

  @Event() updateSearch: EventEmitter;

  constructor() {
    this.#sortingOptions = SEARCH_CONFIG.SORTING_OPTIONS.map((option) => ({
      value: option,
      label: `sorting.${option.order}.${option.axis?.name ?? null}`,
    }));
  }

  get summaryText() {
    const { t } = stores.i18n;
    const { count } = stores.search;

    if (count === 0) {
      return t('results.summary.none');
    }

    if (count === 1) {
      return t('results.summary.one');
    }

    return t('results.summary.multiple', { count: count ?? stores.search.limit });
  }

  get hasFacets() {
    return !!stores.search.facets.length;
  }

  get isEmpty() {
    const { isBusy, response } = stores.search;
    return !isBusy && !response.items?.length;
  }

  changeSorting = (value: any) => {
    stores.search.sorting = this.#sortingOptions.find((option) => option.value === value)?.value;
    this.updateSearch.emit();

    const { axis, order } = stores.search.sorting;
    services.analytics.trackEvent(
      ...ANALYTICS_CUSTOM_EVENTS.SEARCH_SORTING,
      `Changed: ${axis ?? 'relevance'} (${order})`
    );
  };

  render() {
    const { isEmpty, hasFacets } = this;
    const { isBusy, placeholdersCount, response, pageRange, page } = stores.search;
    const { t } = stores.i18n;

    return (
      <Host class="results">
        <div class="results__filters">
          <mex-filters />
        </div>
        <div class="results__container" data-test="results">
          <div class="results__header">
            {!stores.filters.isEmpty && hasFacets && (
              <mex-filters-summary
                class="results__filters-summary"
                label={t('results.summary.filters')}
                testAttr="results:filter"
              />
            )}

            <div class="results__actions">
              <mex-results-download class="results__download" disabled={isBusy} />
              <mex-dropdown
                class="results__sorting"
                orientation="right"
                options={this.#sortingOptions}
                handleChange={this.changeSorting}
                value={stores.search.sorting}
                label={t('sorting.label')}
                testAttr="results:sorting"
                disabled={isEmpty}
              />
            </div>

            {stores.search.wasQueryCleaned && (
              <p class="results__help">
                {t('search.queryHelp')}
                <a {...href(getSearchUrl(true, SEARCH_PAGINATION_START, stores.search.cleanedQuery))}>
                  {stores.search.cleanedQuery}
                </a>
              </p>
            )}
            <h4 class="results__summary" data-test="results:summary">
              {isBusy ? <mex-placeholder text={this.summaryText} /> : <span innerHTML={this.summaryText} />}
            </h4>
          </div>

          <div class="results__items" data-test="results:items">
            <p aria-live="polite" class="u-visually-hidden">
              {isBusy ? t('results.busy') : ''}
            </p>
            {isEmpty && (
              <mex-not-found-state
                class="results__empty-state"
                testAttr="results:empty"
                caption={t('search.empty.title')}
                text={t('search.empty.text')}
              />
            )}

            {isBusy
              ? new Array(placeholdersCount).fill(null).map((_, i) => <mex-result key={i} class="results__item" />)
              : response.items?.map((item) => <mex-result key={item.itemId} class="results__item" item={item} />)}

            {pageRange.length > 1 && !!response && !isBusy && (
              <mex-pagination
                class="results__pagination"
                range={pageRange}
                current={page}
                handleClick={(page: number) => {
                  this.handlePagination?.(page);
                  services.analytics.trackEvent(...ANALYTICS_CUSTOM_EVENTS.SEARCH_NAVIGATION, 'Paginated', page);
                }}
                disabled={isBusy}
                testAttr="results:pagination"
              />
            )}
          </div>
        </div>
      </Host>
    );
  }
}
