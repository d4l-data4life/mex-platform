import { FIELDS } from 'config/fields';
import { FieldConfigDescription, FieldConfigLabel } from 'models/field';
import { normalizeKey } from 'utils/field';

const getTranslations = (attr: 'label' | 'description', languageCode: string) => {
  return Object.values(FIELDS)
    .map((field) => ({
      name: field.name,
      translationConfig: field[attr]?.[languageCode] ?? [],
    }))
    .reduce(
      (translations, { name, translationConfig }) =>
        Object.assign(
          translations,
          translationConfig.reduce((fieldTranslations, translation: FieldConfigLabel | FieldConfigDescription) => {
            return Object.assign(fieldTranslations, {
              [`${normalizeKey(name)}.${translation.entityType}`]: translation,
            });
          }, {})
        ),
      {}
    );
};

export default () => {
  return {
    de: { labels: getTranslations('label', 'de'), descriptions: getTranslations('description', 'de') },
    en: { labels: getTranslations('label', 'en'), descriptions: getTranslations('description', 'en') },
  };
};
