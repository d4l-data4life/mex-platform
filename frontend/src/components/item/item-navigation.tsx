import { Component, Event, EventEmitter, h, Prop, State, Watch } from '@stencil/core';
import { ANALYTICS_CUSTOM_EVENTS, FIELDS, SEARCH_PAGINATION_START, SEARCH_QUERY_EVERYTHING } from 'config';
import { FieldValueDatePrecisionLevel } from 'config/fields';
import services from 'services';
import { Item, Version } from 'services/item';
import stores from 'stores';
import { catchRetryableAction } from 'utils/error';
import { formatDate, formatValue } from 'utils/field';
import { getSearchUrl } from 'utils/search';

@Component({
  tag: 'mex-item-navigation',
  styleUrl: 'item-navigation.css',
})
export class ItemNavigationComponent {
  @Prop() itemId: string;
  @Prop() item?: Item;
  @Prop() versions?: Version[];
  @Prop() expandVersions = false;
  @Prop() breadcrumb: string;

  @State() isBusy = false;
  @State() temporaryNum?: number;

  @Event() changeItemId: EventEmitter;
  @Event() toggleVersions: EventEmitter;

  @Watch('itemId')
  watchIdHandler() {
    requestAnimationFrame(
      () => this.currentIndex === -1 && catchRetryableAction(async () => await this.restorePaginationReference())
    );
  }

  get currentIndex() {
    const items = stores.search.response?.items ?? [];
    return items.indexOf(items.find((item) => this.isOneOfItemIds(item.itemId)));
  }

  get currentNum() {
    return (stores.search.offset ?? 0) + this.currentIndex;
  }

  get paginationRange() {
    return new Array(stores.search.count).fill(null).map((_, index) => index);
  }

  get versionOptions() {
    const versions = this.versions ?? [];
    return versions.map(({ itemId, createdAt, versionDesc }, index) => ({
      value: itemId,
      label: `${versionDesc ?? `v${versions.length - index}`} – ${formatDate(createdAt)}`,
    }));
  }

  get latestVersion() {
    return this.versions?.[0];
  }

  get itemNavigationBreadcrumbs() {
    const { item } = this;

    return [
      {
        label: stores.i18n.t('item.navigation.search'),
        url: getSearchUrl(false, stores.search.page),
        icon: 'back',
        testAttr: 'item:navigation:search:link',
      },
      ...(item
        ? [
            {
              label: stores.i18n.t('item.navigation.metadata', {
                entityType: formatValue([item.entityType], FIELDS.entityName),
              }),
              icon: `entity:${item.entityType}`,
            },
          ]
        : []),
    ];
  }

  isOneOfItemIds(itemId: string) {
    const itemIds = this.versions?.map((version) => version.itemId) ?? [this.itemId];
    return itemIds.includes(itemId);
  }

  async paginate(num: number) {
    const { limit } = stores.search;
    const isNext = num > this.currentNum;
    const offset = Math.floor(num / limit) * limit;
    this.isBusy = true;

    if (offset !== stores.search.offset) {
      await services.search.updateResults({
        query: stores.search.query || SEARCH_QUERY_EVERYTHING,
        offset,
        reset: false,
      });
    }

    const newItemId = stores.search.response?.items?.[num % limit]?.itemId;
    this.changeItemId.emit(newItemId);

    this.isBusy = false;

    services.analytics.trackEvent(...ANALYTICS_CUSTOM_EVENTS.ITEM_SEARCH_RESULT_PAGINATION, isNext ? 'Next' : 'Prev');
  }

  async restorePaginationReference() {
    const { offset, limit } = stores.search;
    const query = stores.search.query || SEARCH_QUERY_EVERYTHING;
    const filters = stores.filters.all;

    if (
      offset &&
      this.isOneOfItemIds(
        (await services.search.fetchResults({ query, offset: offset - 1, filters, limit: 1 }))?.items.shift()?.itemId
      )
    ) {
      await services.search.updateResults({ query, offset: offset - limit, reset: false });
      return;
    }

    if (
      this.isOneOfItemIds(
        (await services.search.fetchResults({ query, offset: offset + limit, filters, limit: 1 }))?.items.shift()
          ?.itemId
      )
    ) {
      await services.search.updateResults({ query, offset: offset + limit, reset: false });
    }
  }

  render() {
    const { t } = stores.i18n;
    const { count } = stores.search;
    const { currentIndex, currentNum, latestVersion, itemNavigationBreadcrumbs } = this;

    return (
      <nav class="item-navigation" data-test="item:navigation">
        <div class="item-navigation__container">
          <mex-breadcrumbs class="item-navigation__breadcrumbs" items={itemNavigationBreadcrumbs} />
          {!!latestVersion && (
            <span class="item-navigation__hint" data-test="item:navigation:versions:latest:date">
              {t('item.navigation.versions.latest', {
                date: formatDate(latestVersion.createdAt, FieldValueDatePrecisionLevel.DAY),
              })}
            </span>
          )}
          <mex-toggle
            class="item-navigation__toggle"
            label={t('item.navigation.versions.toggle')}
            active={this.expandVersions}
            toggleHandler={(active) => {
              this.toggleVersions.emit(active);

              services.analytics.trackEvent(
                ...ANALYTICS_CUSTOM_EVENTS.ITEM_VERSION,
                active ? 'Versions expanded' : 'Versions collapsed'
              );
            }}
            testAttr="item:navigation:versions:toggle"
          />

          <mex-accordion class="item-navigation__row" expanded={this.expandVersions}>
            <div class="item-navigation__card">
              <legend class="item-navigation__legend" data-test="item:navigation:versions:legend">
                <span class="u-highlight-2">{t('item.navigation.versions.highlight')}</span>
                {' = '}
                {t('item.navigation.versions.legend')}
              </legend>

              <strong>{t('item.navigation.versions.dropdown')}:</strong>

              {this.versions ? (
                <mex-dropdown
                  class="item-navigation__versions"
                  toggleClass="item-navigation__select"
                  orientation="left"
                  options={this.versionOptions}
                  handleChange={(itemId) => {
                    this.changeItemId.emit(itemId);

                    const index = this.versions?.indexOf(this.versions.find((version) => version.itemId === itemId));
                    services.analytics.trackEvent(...ANALYTICS_CUSTOM_EVENTS.ITEM_VERSION, 'Changed', index + 1);
                  }}
                  value={this.itemId}
                  disabled={!this.expandVersions}
                  withLabelsTranslation={false}
                  testAttr="item:navigation:versions:dropdown"
                />
              ) : (
                '…'
              )}
            </div>
          </mex-accordion>
        </div>

        {!!count && currentIndex >= 0 && (
          <mex-pagination
            class="item-navigation__pagination"
            range={this.paginationRange}
            current={currentNum}
            handleClick={(num) => catchRetryableAction(async () => await this.paginate(num))}
            disabled={this.isBusy}
            ariaLabelAttr={t('item.navigation.pagination')}
            testAttr="item:navigation:search:pagination"
            renderSlot
          >
            <div>
              {t('item.navigation.resultIndex', {
                count,
                index: currentNum + SEARCH_PAGINATION_START,
              })}
            </div>
          </mex-pagination>
        )}
      </nav>
    );
  }
}
