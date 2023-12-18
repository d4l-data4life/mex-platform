import { createStore } from '@stencil/store';
import { Field } from 'models/field';
import { EntityTypeName } from './entity-types';
import { OrdinalAxis } from './search';

export enum BrowseItemConfigType {
  facet = 'facet',
  hierarchy = 'hierarchy',
}

export interface BrowseItemConfig {
  key: Field;
  type: BrowseItemConfigType;
  axis: OrdinalAxis;
  entityType?: EntityTypeName;
  linkField?: Field;
  displayField?: Field;
  minLevel?: number;
  maxLevel?: number;
  enableSingleNodeVersion?: boolean;
}

export interface BrowseConfig {
  TABS: BrowseItemConfig[];
}

// Defaults are for unit tests
const store = createStore<BrowseConfig>({
  TABS: [],
});

export const BROWSE_FACETS_LIMIT = 1000;

export default store.state;
