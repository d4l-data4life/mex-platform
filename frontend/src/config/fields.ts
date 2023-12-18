import { createStore } from '@stencil/store';
import { Field } from 'models/field';
import { EntityTypeName } from './entity-types';

export enum FieldValueDatePrecisionLevel {
  NONE = 0,
  YEAR = 2,
  MONTH = 3,
  DAY = 4,
  TIME = 5,
}

export type EntityTypeFieldsMapping = {
  [key in EntityTypeName | 'default']?: Field[];
};

const fieldsStore = createStore<{ [name: string]: Field }>({});
export const FIELDS = new Proxy(fieldsStore.state, {
  get(obj, prop: string) {
    return obj[prop] ?? Field.createPending(prop);
  },
});

interface FieldsConfig {
  DEFAULT_FIELD_CONCEPT_PREFIX: string;
  EXTERNAL_FIELD_CONCEPT_PREFIXES: string[];
  METADATA_COMPLETENESS_WEIGHTS: {
    mandatory: number;
    recommended: number;
    optional: number;
  };
}
const store = createStore<FieldsConfig>({
  DEFAULT_FIELD_CONCEPT_PREFIX: '',
  EXTERNAL_FIELD_CONCEPT_PREFIXES: [],
  METADATA_COMPLETENESS_WEIGHTS: {
    mandatory: 60,
    recommended: 30,
    optional: 10,
  },
});

export const FIELD_CONCEPT_PREFIXES: () => string[] = () =>
  [store.get('DEFAULT_FIELD_CONCEPT_PREFIX')].concat(store.get('EXTERNAL_FIELD_CONCEPT_PREFIXES')).filter(Boolean);

export const FIELD_CONCEPT_IDS = () =>
  Object.values(fieldsStore.state)
    .flatMap((field) => field.conceptIds)
    .filter((conceptId, index, arr) => arr.indexOf(conceptId) === index);

export default store.state;
