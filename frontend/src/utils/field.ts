import stores from 'stores';
import {
  HierarchyNodeDisplayValue,
  SearchResultsHighlight,
  SearchResultsItem,
  SearchResultsItemValue,
} from 'stores/search';
import { Item, ItemValue } from 'services/item';
import {
  FIELDS,
  ITEM_FIELDS_CONCATENATOR,
  FIELD_CONCEPT_PREFIXES,
  ITEM_CONFIG,
  SEARCH_CONFIG,
  FIELDS_CONFIG,
  FEATURE_FLAGS,
  FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY,
  ENTITY_TYPES,
} from 'config';
import { FieldValueDatePrecisionLevel } from 'config/fields';
import { Field, FieldEntityVirtualType, FieldImportance, FieldRenderer } from 'models/field';
import { EntityTypeName } from 'config/entity-types';
import { SidebarFeature } from 'config/item';

export const normalizeKey = (key: string) => {
  return key
    .split(/[:_]/)
    .map((part, index) =>
      index ? part.slice(0, 1).toUpperCase() + part.slice(1) : part.slice(0, 1).toLowerCase() + part.slice(1)
    )
    .join('');
};

export const normalizeConceptId = (id: string) =>
  FIELD_CONCEPT_PREFIXES().reduce((normalizedId, prefix) => normalizedId?.replace(prefix, '') ?? null, id);

export const denormalizeConceptId = (normalizedId: string, field?: Field) => {
  if (!normalizedId || (field && !field.isEnumerable)) {
    return normalizedId;
  }

  if (!field) {
    return FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX + normalizedId;
  }

  const { conceptIds } = field;
  const prefix =
    FIELD_CONCEPT_PREFIXES().find((prefix) => conceptIds.includes(prefix + normalizedId)) ??
    FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX;
  return prefix + normalizedId;
};

const translateFieldAttr = ({
  attr,
  field,
  item,
  forcePlural,
  fallback,
}: {
  attr: 'labels' | 'descriptions';
  field: Field;
  item?: Item | SearchResultsItem;
  forcePlural?: boolean;
  fallback?: string;
}) => {
  const supportsPlural = attr === 'labels';
  const values = item
    ? getValues(field, item)
        .map(({ fieldValue }) => fieldValue)
        .filter(filterDuplicates)
    : [];
  const usePlural = (supportsPlural && forcePlural) || values.length > 1;
  const normalizedKey = normalizeKey(field.name);
  const entityTypeKey =
    Object.keys(ENTITY_TYPES).find((key) => ENTITY_TYPES[key].name === item?.entityType) ?? FieldEntityVirtualType.all;
  const targetKey = supportsPlural ? (usePlural ? 'plural' : 'singular') : 'text';
  const translationKeys = [
    `fields.${attr}.${normalizedKey}.${entityTypeKey}.${targetKey}`,
    usePlural && `fields.${attr}.${normalizedKey}.${entityTypeKey}.singular`,
    `fields.${attr}.${normalizedKey}.all.${targetKey}`,
    usePlural && `fields.${attr}.${normalizedKey}.all.singular`,
  ].filter(Boolean);
  const translation = translationKeys.reduce((translation, key) => {
    if (translation) {
      return translation;
    }

    const translatedName = stores.i18n.t(key);
    return translatedName !== key ? translatedName : null;
  }, null);

  return translation ?? fallback;
};

export const translateFieldName = (
  field: Field,
  item?: Item | SearchResultsItem,
  forcePlural?: boolean,
  fallback?: string
) => translateFieldAttr({ attr: 'labels', field, item, forcePlural, fallback: fallback ?? normalizeKey(field.name) });

export const translateFieldDescription = (field: Field, item?: Item | SearchResultsItem) =>
  translateFieldAttr({ attr: 'descriptions', field, item });

const translateFieldVocabularyAttr = (attr: string, value: string, field?: Field) => {
  const isEntityName = field === FIELDS.entityName;
  if (!value || (!isEntityName && !FIELD_CONCEPT_PREFIXES().some((prefix) => value.includes(prefix)))) {
    return value;
  }

  const key = `vocabulary.${normalizeConceptId(value)}.${attr}`;
  const translation = stores.i18n.t(key);
  return key === translation ? normalizeConceptId(value) : translation;
};

export const translateFieldValue = (value: string, field?: Field) =>
  translateFieldVocabularyAttr('label', value, field);

export const translateFieldValueDescription = (value: string, field?: Field) =>
  translateFieldVocabularyAttr('description', value, field);

export const formatDate = (
  date: Date,
  precisionLevel: FieldValueDatePrecisionLevel = FieldValueDatePrecisionLevel.TIME,
  timeZone = 'UTC'
) => {
  if (precisionLevel === FieldValueDatePrecisionLevel.NONE) {
    return undefined;
  }

  if (precisionLevel === FieldValueDatePrecisionLevel.YEAR) {
    return String(date?.getFullYear());
  }

  try {
    return date.toLocaleString(stores.i18n.language, {
      timeZone,
      year: 'numeric',
      month: 'short',
      day: precisionLevel < FieldValueDatePrecisionLevel.DAY ? undefined : '2-digit',
      ...(precisionLevel === FieldValueDatePrecisionLevel.TIME
        ? {
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
          }
        : {}),
    });
  } catch (_) {
    return date?.toISOString();
  }
};

export const getTextMatches = (field: Field, highlights?: SearchResultsHighlight[]): string[] | undefined => {
  const currentLanguage = stores.i18n.language;
  return (
    highlights
      ?.flatMap(({ matches }) => matches ?? [])
      .filter((match) => !FEATURE_FLAGS.FIELD_TRANSLATIONS || !match.language || match.language === currentLanguage)
      .filter((match) => match.fieldName === field.linkedName)
      .flatMap((match) => match.snippets) ?? []
  );
};

export const filterDuplicates = (item: any, index: number, arr: any[]) => arr.indexOf(item) === index;

export const sortAndFilterValues = <T>(
  values?: (ItemValue | SearchResultsItemValue | HierarchyNodeDisplayValue)[],
  filterOutForeignTranslations = false
) => {
  const currentLanguage = stores.i18n.language;
  return values
    ?.sort(({ place: a = 0 }: any, { place: b = 0 }: any) => a - b)
    .filter((item) => !filterOutForeignTranslations || !item.language || item.language === currentLanguage)
    .filter(filterDuplicates) as unknown as T | undefined;
};

export const getValues = (field: Field, item: Item | SearchResultsItem) => {
  const allValues = item?.values.filter((item) => item.fieldName === field.linkedName);
  const filteredValues = <(ItemValue | SearchResultsItemValue)[]>(
    sortAndFilterValues(allValues, FEATURE_FLAGS.FIELD_TRANSLATIONS)
  );
  return filteredValues?.length ? filteredValues : allValues?.filter(filterDuplicates) ?? [];
};

export const getValue = (field: Field, item: Item | SearchResultsItem) => {
  return getValues(field, item)?.[0]?.fieldValue;
};

export const hasValue = (field: Field, item: Item | SearchResultsItem) => !!getValue(field, item);

export const formatValue = (values: string[], field?: Field) => {
  const renderer = field?.renderer ?? FieldRenderer.plain;
  if (renderer === FieldRenderer.time) {
    values = values.map((value) => {
      const precisionLevel = (value?.match(/^(\d{4})(-\d{2})?(-\d{2})?(T[\d:]{5,}.+)?$/)?.filter(Boolean) ?? []).length;
      const date = value && new Date(value);
      if (!date || precisionLevel === FieldValueDatePrecisionLevel.NONE || date.toString() === 'Invalid Date') {
        return value;
      }

      return formatDate(date, Math.min(precisionLevel, FieldValueDatePrecisionLevel.DAY));
    });
  }

  const translatedValues = values.map((value) => translateFieldValue(value, field));
  return escapeHtml(translatedValues.join(ITEM_FIELDS_CONCATENATOR[renderer] ?? ITEM_FIELDS_CONCATENATOR.default));
};

export const escapeHtml = (unsafe: string) => {
  return unsafe
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
};

export const getRawValues = (fields: Field[], item: Item | SearchResultsItem) =>
  fields
    .flatMap((field) => getValues(field, item))
    .map(({ fieldValue }) => fieldValue)
    .filter(filterDuplicates);

export const getDisplayValues = (fields: Field[], item: Item | SearchResultsItem, highlightValues = true): string[] => {
  const primaryField = fields[0];
  const values = getRawValues(fields, item);
  const matches = highlightValues
    ? fields
        .flatMap((field) => getTextMatches(field, stores.search.highlights))
        .filter(Boolean)
        .filter(filterDuplicates)
    : [];

  if (!matches.length) {
    return values.map((value) => formatValue([value], primaryField));
  }

  const HIGHLIGHT_START_CHARACTER = '\ue000';
  const HIGHLIGHT_END_CHARACTER = '\ue001';

  const highlightedValues = values.map((value) =>
    matches.reduce((result, match) => {
      const slice = match.replaceAll(HIGHLIGHT_START_CHARACTER, '').replaceAll(HIGHLIGHT_END_CHARACTER, '');
      return result.replaceAll(slice, match);
    }, value)
  );

  return highlightedValues
    .map((highlightedValue) =>
      formatValue([highlightedValue], primaryField)
        .replaceAll(HIGHLIGHT_START_CHARACTER, '<em>')
        .replaceAll(HIGHLIGHT_END_CHARACTER, '</em>')
    )
    .filter(Boolean);
};

export const getConcatenator = (fields: Field[]) => {
  const commonRenderer = fields[0]?.renderer ?? FieldRenderer.plain;
  return ITEM_FIELDS_CONCATENATOR[commonRenderer] ?? ITEM_FIELDS_CONCATENATOR.default;
};

export const getConcatDisplayValue = (
  fields: Field[],
  item: Item | SearchResultsItem,
  fallback = stores.i18n.t(FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY),
  highlightValues = true
): string => {
  return getDisplayValues(fields, item, highlightValues).join(getConcatenator(fields)) || fallback;
};

export const hasValueChanged = (field: Field, itemA: Item, itemB: Item = null) => {
  if (!itemA || !itemB) {
    return false;
  }

  const valuesA = getValues(field, itemA).map(({ fieldValue }) => fieldValue);
  const valuesB = getValues(field, itemB).map(({ fieldValue }) => fieldValue);
  return valuesA.length !== valuesB.length || valuesA.some((value) => !valuesB.includes(value));
};

export const getLangAttrIfForeign = (
  fields: Field[],
  item: Item | SearchResultsItem,
  summarize = true
): string | string[] => {
  const values = fields.flatMap((field) => getValues(field, item));
  const attrs = values.map(({ language }) => (!!language && language !== stores.i18n.language ? language : undefined));
  const filteredAttrs = attrs.filter(Boolean);
  return summarize ? (filteredAttrs.length === values.length ? filteredAttrs[0] : undefined) : attrs;
};

export const getInvolvedFields = (entityType: EntityTypeName, includeCached = true): Field[] => {
  const displayedFields = ITEM_CONFIG.DISPLAYED_FIELDS[entityType] ?? ITEM_CONFIG.DISPLAYED_FIELDS.default ?? [];
  const sidebarFields = (SEARCH_CONFIG.SIDEBAR_FEATURES[entityType] ?? SEARCH_CONFIG.SIDEBAR_FEATURES.default)
    .concat(ITEM_CONFIG.SIDEBAR_FEATURES[entityType] ?? ITEM_CONFIG.SIDEBAR_FEATURES.default)
    .filter(({ feature }) => includeCached || feature !== SidebarFeature.contactForm)
    .map(({ field }) => field)
    .filter(Boolean);

  return (
    displayedFields
      .concat(sidebarFields)
      .map((field) => field?.resolvesTo)
      .flat()
      .filter(Boolean)
      .filter(filterDuplicates) ?? []
  );
};

const intersect = (a: any[], b: any[], not = false) => a.filter((item) => (not ? !b.includes(item) : b.includes(item)));

export const getDocumentationProp = (field: Field, prop: string, entityTypeName?: EntityTypeName) => {
  const entries = field.documentation?.[stores.i18n.language]?.[prop] ?? [];

  return (
    entries.find(({ entityType }) => entityType === entityTypeName)?.text ??
    entries.find(({ entityType }) => entityType === FieldEntityVirtualType.all)?.text ??
    null
  );
};

export interface CompletenessData {
  weights: { mandatory: number; recommended: number; optional: number };
  scores: { mandatory: number; recommended: number; optional: number; total: number };
  missingFields: { mandatory: Field[]; recommended: Field[]; optional: Field[] };
  populatedFields: { mandatory: Field[]; recommended: Field[]; optional: Field[] };
}

export const calculateCompleteness = (item: Item | SearchResultsItem): CompletenessData => {
  const entityType = item.entityType;
  const involved = getInvolvedFields(entityType as EntityTypeName);
  const populated = involved.filter((field) => hasValue(field, item));

  const mandatoryFields = Object.values(FIELDS).filter(({ importance }) => importance === FieldImportance.mandatory);
  const recommendedFields = Object.values(FIELDS).filter(
    ({ importance }) => importance === FieldImportance.recommended
  );
  const optionalFields = Object.values(FIELDS).filter(({ importance }) => importance === FieldImportance.optional);

  const involvedMandatory = intersect(involved, mandatoryFields);
  const involvedRecommended = intersect(involved, recommendedFields);
  const involvedOptional = intersect(involved, optionalFields);

  const populatedMandatory = intersect(populated, involvedMandatory);
  const populatedRecommended = intersect(populated, involvedRecommended);
  const populatedOptional = intersect(populated, involvedOptional);

  const mandatoryScore = involvedMandatory.length
    ? Math.round(
        (populatedMandatory.length / involvedMandatory.length) * FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS.mandatory
      )
    : FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS.mandatory;
  const recommendedScore = involvedRecommended.length
    ? Math.round(
        (populatedRecommended.length / involvedRecommended.length) *
          FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS.recommended
      )
    : FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS.recommended;
  const optionalScore = involvedOptional.length
    ? Math.round(
        (populatedOptional.length / involvedOptional.length) * FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS.optional
      )
    : FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS.optional;

  return {
    weights: FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS,
    scores: {
      mandatory: mandatoryScore,
      recommended: recommendedScore,
      optional: optionalScore,
      total: mandatoryScore + recommendedScore + optionalScore,
    },
    missingFields: {
      mandatory: intersect(involvedMandatory, populatedMandatory, true),
      recommended: intersect(involvedRecommended, populatedRecommended, true),
      optional: intersect(involvedOptional, populatedOptional, true),
    },
    populatedFields: {
      mandatory: populatedMandatory,
      recommended: populatedRecommended,
      optional: populatedOptional,
    },
  };
};
