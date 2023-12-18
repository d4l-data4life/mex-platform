import { Component, h, Host, State } from '@stencil/core';
import { FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY, HOME_CONFIG, SEARCH_QUERY_EVERYTHING } from 'config';
import { DashboardMetricMethod } from 'config/home';
import { OrdinalAxis } from 'config/search';
import services from 'services';
import stores from 'stores';
import { SearchResults } from 'stores/search';
import { catchRetryableAction } from 'utils/error';
import { translateFieldName, translateFieldValue, getConcatDisplayValue, normalizeConceptId } from 'utils/field';

@Component({
  tag: 'mex-dashboard',
  styleUrl: 'dashboard.css',
})
export class DashboardComponent {
  @State() results?: SearchResults;

  get isBusy() {
    return !this.results;
  }

  get lastUpdate() {
    return getConcatDisplayValue(
      [HOME_CONFIG.LATEST_UPDATE_AXIS?.uiField],
      this.results?.items[0],
      stores.i18n.t(FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY),
      false
    );
  }

  get chartBuckets() {
    return HOME_CONFIG.CHART && this.results?.facets?.find(({ axis }) => axis === HOME_CONFIG.CHART.axis.name)?.buckets;
  }

  get chartGreenCount() {
    return this.chartBuckets?.find(({ value }) => value === HOME_CONFIG.CHART.greenColorBucket)?.count;
  }

  get chartRedCount() {
    return this.chartBuckets?.find(({ value }) => value === HOME_CONFIG.CHART.redColorBucket)?.count;
  }

  get chartGreenColorPerc() {
    const { chartGreenCount = 0, chartRedCount = 0 } = this;
    return chartGreenCount && Math.round((chartGreenCount / (chartGreenCount + chartRedCount)) * 100);
  }

  getEntityTypeCount(axis: OrdinalAxis, method: DashboardMetricMethod, bucketName?: string) {
    const facet = this.results?.facets?.find((facet) => facet.axis === axis.name);
    if (!facet?.buckets) {
      return 0;
    }

    if (method === DashboardMetricMethod.NUMBER_OF_BUCKETS) {
      return facet.bucketNo ?? 0;
    }

    return facet.buckets.find(({ value }) => value === bucketName)?.count ?? 0;
  }

  componentDidLoad() {
    const latestUpdateAxis = HOME_CONFIG.LATEST_UPDATE_AXIS;
    const axes = HOME_CONFIG.DASHBOARD_METRIC_CONFIGS.map(({ axis }) => axis);

    catchRetryableAction(async () => {
      this.results = await services.search.fetchResults({
        query: SEARCH_QUERY_EVERYTHING,
        limit: latestUpdateAxis ? 1 : 0,
        fields: latestUpdateAxis ? [latestUpdateAxis.uiField] : [],
        sorting: {
          axis: latestUpdateAxis ?? null,
          order: 'desc',
        },
        highlightFields: [],
        facets: [...axes, ...(HOME_CONFIG.CHART ? [HOME_CONFIG.CHART.axis] : [])]
          .filter(Boolean)
          .map((axis) => ({ axis, type: 'exact' })),
        facetsLimit: 5,
      });
    });
  }

  render() {
    const { isBusy, lastUpdate, chartGreenColorPerc, results } = this;
    const { t } = stores.i18n;

    const { CHART } = HOME_CONFIG;

    return (
      <Host class="dashboard">
        <h2 class="dashboard__title u-underline-2">{t('dashboard.title')}</h2>
        {!!HOME_CONFIG.LATEST_UPDATE_AXIS && (
          <p class="dashboard__hint">
            {t('dashboard.updatedAt')} {isBusy ? <mex-placeholder text="xx.xx.xxxx" /> : lastUpdate}
          </p>
        )}
        {new Array(Math.ceil(HOME_CONFIG.DASHBOARD_METRIC_CONFIGS.length / 2)).fill(null).map((_, row) => (
          <div class="dashboard__tiles" key={`dashboard-metrics-${row}`}>
            {HOME_CONFIG.DASHBOARD_METRIC_CONFIGS.slice(row * 2, row * 2 + 2).map(
              ({ entityType, axis, method }, index) => (
                <span class="dashboard__tile" key={`dashboard-metric-${row}-${index}`}>
                  <strong>
                    <mex-icon-entity
                      entityName={entityType.name}
                      attrs={{ class: 'u-underline-3', classes: 'icon--large' }}
                    />
                    {isBusy ? <mex-placeholder text=".." /> : this.getEntityTypeCount(axis, method, entityType.name)}
                  </strong>
                  {t(`dashboard.metrics.${normalizeConceptId(entityType.name)}`)}
                </span>
              )
            )}
          </div>
        ))}

        {!!CHART && !!results?.numFound && (
          <div class="dashboard__charts">
            <div class="dashboard__chart">
              <h3 class="dashboard__label">{translateFieldName(CHART.axis.uiField)}</h3>
              {isBusy ? (
                <mex-placeholder />
              ) : (
                <div class="dashboard__chart-tiles">
                  <div class="dashboard__chart-tile">
                    <div
                      class="dashboard__donut-chart"
                      style={{
                        '--dashboard-chart-first-perc': `${chartGreenColorPerc}%`,
                        '--dashboard-chart-first-color': 'var(--color-green-light)',
                        '--dashboard-chart-second-color': 'var(--color-red-light)',
                      }}
                    />
                  </div>
                  <div class="dashboard__chart-tile">
                    <legend class="dashboard__chart-legend">
                      <ul>
                        <li style={{ '--dashboard-chart-legend-color': 'var(--color-green-light)' }}>
                          {chartGreenColorPerc}% {translateFieldValue(CHART.greenColorBucket, CHART.axis.uiField)}
                        </li>
                        <li style={{ '--dashboard-chart-legend-color': 'var(--color-red-light)' }}>
                          {100 - chartGreenColorPerc}% {translateFieldValue(CHART.redColorBucket, CHART.axis.uiField)}
                        </li>
                      </ul>
                    </legend>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}
      </Host>
    );
  }
}
