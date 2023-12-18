import { createStore } from '@stencil/store';
import { EntityType } from './entity-types';
import { OrdinalAxis } from './search';

export interface DashboardChartConfig {
  axis: OrdinalAxis;
  greenColorBucket: string;
  redColorBucket: string;
}

export enum DashboardMetricMethod {
  NUMBER_OF_BUCKETS = 'bucketNo',
  BUCKET = 'bucket',
}

export interface DashboardMetricConfig {
  entityType: EntityType;
  axis: OrdinalAxis;
  method: DashboardMetricMethod;
}

interface HomeConfig {
  CHART?: DashboardChartConfig;
  DASHBOARD_METRIC_CONFIGS: DashboardMetricConfig[];
  LATEST_UPDATE_AXIS?: OrdinalAxis;
}

const store = createStore<HomeConfig>({
  DASHBOARD_METRIC_CONFIGS: [],
});

export default store.state;
