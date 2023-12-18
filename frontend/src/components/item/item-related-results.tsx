import { Component, h, Prop, State } from '@stencil/core';
import { href } from 'stencil-router-v2';
import { FIELDS, ITEM_CONFIG, ROUTES, SEARCH_QUERY_EVERYTHING } from 'config';
import services from 'services';
import { Item } from 'services/item';
import stores from 'stores';
import { SearchResults } from 'stores/search';
import { catchRetryableAction } from 'utils/error';
import { normalizeKey } from 'utils/field';
import { createExactAxisConstraint, getFilterQueryParamName, getQueryStringFromParams } from 'utils/search';

@Component({
  tag: 'mex-item-related-results',
  styleUrl: 'item-related-results.css',
})
export class ItemRelatedResultsComponent {
  @Prop() item: Item;

  @State() results?: SearchResults;

  get isShown() {
    return !!this.config && !!this.results?.items?.length;
  }

  get config() {
    return ITEM_CONFIG.RELATED_RESULTS_CONFIG[this.item.entityType];
  }

  get totalCount() {
    const limit = this.config?.limit ?? 0;
    return this.results?.numFound ?? limit;
  }

  get moreCount() {
    const limit = this.config?.limit ?? 0;
    return Math.max(this.totalCount - limit, 0);
  }

  get moreUrl() {
    const { targetEntityType, linkedField } = this.config;
    const businessIdentifier = this.item.businessId;
    const params = new URLSearchParams();

    params.append(getFilterQueryParamName(FIELDS.entityName.linkedName), targetEntityType);
    params.append(getFilterQueryParamName(linkedField.linkedName), businessIdentifier);

    return `${ROUTES.SEARCH}?${getQueryStringFromParams(params)}`;
  }

  async populateResults() {
    const { config } = this;
    if (!config) {
      return;
    }

    const { targetEntityType, linkedField, limit } = config;
    const businessIdentifier = this.item.businessId;
    if (!businessIdentifier) {
      return;
    }

    this.results = await services.search.fetchResults({
      query: SEARCH_QUERY_EVERYTHING,
      limit,
      axisConstraints: [
        createExactAxisConstraint(FIELDS.entityName.linkedName, [targetEntityType]),
        createExactAxisConstraint(linkedField.linkedName, [businessIdentifier]),
      ],
      facets: [],
      highlightFields: [],
    });
  }

  componentWillLoad() {
    catchRetryableAction(async () => await this.populateResults());
  }

  render() {
    const { totalCount, moreCount } = this;

    return (
      this.isShown && (
        <div class="item-related-results" data-test="item:relatedResults">
          <h4 class="item-related-results__title u-underline-2">
            {stores.i18n.t(`item.relatedResults.${normalizeKey(this.config.targetEntityType)}`, { count: totalCount })}
          </h4>
          {this.results.items?.map((item) => (
            <mex-result key={item.itemId} class="item-related-results__item" item={item} />
          ))}
          {!!moreCount && (
            <a class="item-related-results__more-link" {...href(this.moreUrl)} data-test="item:relatedResults:moreLink">
              {stores.i18n.t('item.relatedResults.more', { totalCount, moreCount })}
              <mex-icon-arrow classes="icon--medium" />
            </a>
          )}
        </div>
      )
    );
  }
}
