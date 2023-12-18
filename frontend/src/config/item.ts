import { createStore } from '@stencil/store';
import { Field, FieldRenderer } from 'models/field';
import { EntityTypeName } from './entity-types';
import { FIELDS, EntityTypeFieldsMapping } from './fields';

export interface ItemRelatedResultsConfigItem {
  targetEntityType: EntityTypeName;
  linkedField: Field;
  limit: number;
}

export enum SidebarFeature {
  completeness = 'completeness',
  date = 'date',
  accessRestriction = 'accessRestriction',
  contactForm = 'contactForm',
  displayField = 'displayField',
}

export enum SidebarFeatureContactAction {
  email = 'email',
  form = 'form',
}

export interface SidebarFeatureConfigActionItem {
  type: SidebarFeatureContactAction;
  field?: string | Field;
  form?: string;
}

export interface SidebarFeatureConfigIconMapItem {
  conceptid: string;
  icon: 'locked' | 'unlocked';
}

export interface SidebarFeatureConfig {
  feature: SidebarFeature;
  field?: Field;
  displayField?: Field;
  key?: string;
  actions?: SidebarFeatureConfigActionItem[];
  iconMap?: SidebarFeatureConfigIconMapItem[];
}

export type EntityTypeSidebarFeaturesMapping = {
  [key in EntityTypeName | 'default']?: SidebarFeatureConfig[];
};

export const ITEM_GC_COUNT: number = 30;

export const ITEM_FIELDS_CONCATENATOR: {
  [renderer: string]: string;
} = {
  [FieldRenderer.title]: ' | ',
  [FieldRenderer.time]: ' - ',
  [FieldRenderer.bullets]: '',
  default: ', ',
};

export interface ItemConfig {
  DISPLAYED_FIELDS: EntityTypeFieldsMapping;
  RELATED_RESULTS_CONFIG: {
    [key in EntityTypeName]?: ItemRelatedResultsConfigItem;
  };
  SIDEBAR_FEATURES: EntityTypeSidebarFeaturesMapping;
  DEDICATED_VIEW_SUPPORTED_ENTITY_TYPES: EntityTypeName[];
}

// Defaults are for unit tests
const store = createStore<ItemConfig>({
  DISPLAYED_FIELDS: {
    default: Object.values(FIELDS),
  },
  RELATED_RESULTS_CONFIG: {},
  SIDEBAR_FEATURES: {
    default: [],
  },
  DEDICATED_VIEW_SUPPORTED_ENTITY_TYPES: [],
});

export default store.state;
