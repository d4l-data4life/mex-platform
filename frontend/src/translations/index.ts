import TRANSLATIONS_DE from './de.json';
import TRANSLATIONS_EN from './en.json';

/**
 * This returns the core translations. Those are a minimal set of keys necessary while the application
 * is being bootstrapped and to show the error modal if fetching the translations/config fails.
 */

export default {
  de: { ui: TRANSLATIONS_DE },
  en: { ui: TRANSLATIONS_EN },
};
