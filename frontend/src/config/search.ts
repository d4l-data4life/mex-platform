import { createStore } from '@stencil/store';

import { Field } from 'models/field';

import { EntityTypeName } from './entity-types';
import { EntityTypeFieldsMapping } from './fields';
import ITEM_CONFIG, { EntityTypeSidebarFeaturesMapping } from './item';

export interface OrdinalAxis {
  name: string;
  uiField?: Field;
}

export interface SearchFacet {
  axis: OrdinalAxis;
  type: 'exact' | 'yearRange' | 'hierarchy';
  entityType?: EntityTypeName;
  linkField?: Field;
  displayField?: Field;
  minLevel?: number;
  maxLevel?: number;
}

export interface SearchAxisConstraint {
  type: 'exact' | 'stringRange';
  axis: string;
  combineOperator: string;
  values?: string[];
  singleNodeValues?: string[];
  stringRanges?: {
    min: string;
    max: string;
  }[];
}

export interface SearchSorting {
  axis: OrdinalAxis | null;
  order: 'asc' | 'desc';
}

export interface SearchConfig {
  DISPLAYED_FIELDS: EntityTypeFieldsMapping;
  FACETS: SearchFacet[];
  FACETS_LIMIT: number;
  FOCI: (string | null)[];
  LIMIT: number;
  ORDINAL_AXES: OrdinalAxis[];
  PAGINATION_RANGE_COUNT: number;
  QUERY_MAX_EDIT_DISTANCE: 0 | 1;
  QUERY_USE_NGRAM_FIELD: boolean;
  SIDEBAR_FEATURES: EntityTypeSidebarFeaturesMapping;
  SORTING_OPTIONS: SearchSorting[];
}

// Defaults are for unit tests
const store = createStore<SearchConfig>({
  DISPLAYED_FIELDS: {
    default: [],
  },
  FACETS: [],
  FACETS_LIMIT: 20,
  FOCI: [null],
  LIMIT: 10,
  ORDINAL_AXES: [],
  PAGINATION_RANGE_COUNT: 4,
  QUERY_MAX_EDIT_DISTANCE: 0,
  QUERY_USE_NGRAM_FIELD: false,
  SIDEBAR_FEATURES: {
    default: [],
  },
  SORTING_OPTIONS: [{ axis: null, order: 'desc' }],
});

export const SEARCH_INVISIBLE_FACETS: () => SearchFacet[] = () =>
  Object.values(ITEM_CONFIG.RELATED_RESULTS_CONFIG).map(({ linkedField }) => ({
    axis: {
      name: linkedField.linkedName,
      uiField: linkedField,
    },
    type: 'exact',
  }));

export const SEARCH_FACETS_COMBINE_OPERATOR: 'and' | 'or' = 'or';
export const SEARCH_FACETS_MAX_LIMIT: number = 1000; // hard limit set by backend
export const SEARCH_FACETS_SINGLE_NODE_LABEL_KEY: string = 'filters.hierarchyDirect';
export const SEARCH_FACETS_SINGLE_NODE_VALUE_PREFIX: string = 'direct:';
export const SEARCH_QUERY_EVERYTHING = '*';
export const SEARCH_PAGINATION_START: number = 1;
export const SEARCH_PARAM_PAGE: string = 'page';
export const SEARCH_PARAM_SORTING_AXIS: string = 'sorting.axis';
export const SEARCH_PARAM_SORTING_ORDER: string = 'sorting.order';
export const SEARCH_PARAM_FOCUS: string = 'focus';
export const SEARCH_PARAM_FILTER_PREFIX: string = 'f.';
export const SEARCH_TIMEOUT: number = 5000;

export enum PadCondition {
  MAY,
  MUST,
  MUST_NOT,
}
export const SEARCH_OPERATOR_PAD_MAP: {
  [operator: string]: PadCondition[];
} = {
  '*': [PadCondition.MAY, PadCondition.MAY],
  '""': [PadCondition.MUST, PadCondition.MUST],
  '+': [PadCondition.MUST, PadCondition.MUST],
  '|': [PadCondition.MUST, PadCondition.MUST],
  '-': [PadCondition.MUST, PadCondition.MUST_NOT],
  '()': [PadCondition.MUST, PadCondition.MUST],
  default: [PadCondition.MAY, PadCondition.MAY],
};

export interface SearchDownloadConfig {
  key: string;
  limit: number;
  entityType?: EntityTypeName;
}

export const SEARCH_DOWNLOAD_CONFIG: SearchDownloadConfig[] = [
  {
    key: 'results',
    limit: SEARCH_FACETS_MAX_LIMIT,
  },
];

export default store.state;
