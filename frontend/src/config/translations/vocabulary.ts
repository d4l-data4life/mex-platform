import { FIELDS } from 'config/fields';
import { normalizeConceptId } from 'utils/field';

const getTranslations = (languageCode: string) => {
  return Object.values(FIELDS)
    .filter((field) => field.isEnumerable)
    .map((field) => field.vocabulary?.[languageCode] ?? [])
    .flat()
    .reduce(
      (translations, { conceptId, label, description }) =>
        Object.assign(translations, { [normalizeConceptId(conceptId)]: { label, description } }),
      {}
    );
};

export default () => {
  return {
    de: getTranslations('de'),
    en: getTranslations('en'),
  };
};
