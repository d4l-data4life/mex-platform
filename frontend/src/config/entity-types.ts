import { createStore } from '@stencil/store';

export enum EntityTypeIcon {
  datum = 'datum',
  platform = 'platform',
  source = 'source',
  resource = 'resource',
}

export interface EntityTypeConfig {
  isFocal: boolean;
  isAggregatable: boolean;
  businessIdFieldName?: string;
  aggregationAlgorithm?: 'simple' | 'source_partition';
  aggregationEntityType?: EntityTypeName;
  duplicateStrategy?: 'keepall' | 'removeall';
  partitionFieldName?: string;
  icon?: EntityTypeIcon;
}

export type EntityTypeName = string;
export interface EntityType {
  name: EntityTypeName;
  config: EntityTypeConfig;
}

export interface EntityTypesConfig {
  [key: EntityTypeName]: EntityType;
}

const store = createStore<EntityTypesConfig>({});
export default store.state;
