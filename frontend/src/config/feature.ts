import { createStore } from '@stencil/store';
import { NAVIGATION_CONFIG } from 'config';

interface FeatureFlags {
  FIELD_TRANSLATIONS: boolean;
  BROWSE_HIERARCHY_TRANSLATIONS: boolean;
  METADATA_COMPLETENESS_DETAILS: boolean;
  DEBUG_UNUSED_FIELDS: boolean;
}

export const store = createStore<FeatureFlags>({
  FIELD_TRANSLATIONS: false,
  BROWSE_HIERARCHY_TRANSLATIONS: true,
  METADATA_COMPLETENESS_DETAILS: true,
  DEBUG_UNUSED_FIELDS: false,
});

export const FEATURE_FLAG_LANGUAGE_SWITCHER = () =>
  !NAVIGATION_CONFIG.LANGUAGE_SWITCHER_LANGUAGES || NAVIGATION_CONFIG.LANGUAGE_SWITCHER_LANGUAGES.length > 1;

export default store.state;
