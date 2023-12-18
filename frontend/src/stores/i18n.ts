import { createStore } from '@stencil/store';
import i18n, { Resource, TFunction } from 'i18next';
import { LANGUAGE_CODES, LANGUAGE_QUERY_PARAM } from 'config';
import resources from 'translations';
import getFieldTranslations from 'config/translations/fields';
import getFieldVocabularyTranslations from 'config/translations/vocabulary';

const preferredLanguage = new URLSearchParams(document.location.search).get(LANGUAGE_QUERY_PARAM);

i18n.init({
  lng: LANGUAGE_CODES.includes(preferredLanguage) ? preferredLanguage : LANGUAGE_CODES[1],
  resources,
  ns: ['ui'],
  partialBundledLanguages: true,
});

interface StateType {
  language: string;
  t: TFunction;
}

const store = createStore<StateType>({
  language: i18n.language,
  t: i18n.t,
});

const updateHtmlLang = (language: string) => document.documentElement.setAttribute('lang', language);
i18n.on('languageChanged', updateHtmlLang);
updateHtmlLang(i18n.language);

type LanguageChangeListener = (language?: string) => void;
let listeners: LanguageChangeListener[] = [];

const setLanguage = (language: string) =>
  i18n.changeLanguage(language, () => {
    store.set('language', language);
    store.set('t', i18n.t.bind(i18n));
    listeners.forEach((listener) => listener(language));
  });

export const addUiTranslations = (resources: Resource) => {
  Object.keys(resources).forEach((lng) => i18n.addResourceBundle(lng, 'ui', resources[lng]));
  setLanguage(i18n.language); // update UI
};

export const addFieldTranslations = () => {
  const fieldTranslations = getFieldTranslations();
  const fieldVocabularyTranslations = getFieldVocabularyTranslations();
  i18n.addResourceBundle('de', 'ui', { fields: fieldTranslations.de });
  i18n.addResourceBundle('en', 'ui', { fields: fieldTranslations.en });
  i18n.addResourceBundle('de', 'ui', { vocabulary: fieldVocabularyTranslations.de });
  i18n.addResourceBundle('en', 'ui', { vocabulary: fieldVocabularyTranslations.en });
};

export const addAdditionalTranslations = (key: string, resources: any) => {
  i18n.addResourceBundle('de', 'ui', { [key]: resources?.de });
  i18n.addResourceBundle('en', 'ui', { [key]: resources?.en });
};

export default {
  addListener: (listener: LanguageChangeListener) => {
    listeners.push(listener);
  },
  removeListener: (listener: LanguageChangeListener) => {
    listeners = listeners.filter((item) => item !== listener);
  },
  get language() {
    return store.get('language');
  },
  set language(language: string) {
    setLanguage(language);
  },
  get t() {
    return function (...args) {
      const translation = store.get('t').apply(null, args);
      const key = args[0];

      if (key?.split('.')[0] !== 'fields' && key === translation && typeof global === 'undefined') {
        console.error(`i18n: ${key} has no translation`);
      }
      return translation;
    };
  },
  convertLinks(text: string) {
    return text.replace(/\[([^\]]+)\]\(([^\)]+)\)/gm, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>');
  },
};
