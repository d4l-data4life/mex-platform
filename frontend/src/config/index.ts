import { Env } from '@stencil/core';
import NavigationConfig from './navigation';
import FieldsConfig from './fields';
import HomeConfig from './home';
import BrowseConfig from './browse';
import ItemConfig from './item';
import SearchConfig from './search';
import AnalyticsConfig from './analytics';
import FeatureFlags from './feature';
import EntityTypes from './entity-types';

const environmentMapping = {
  // add your hostname here to map to an environment, like this:
  // 'staging.mex.xyz': 'staging',
  default: 'development',
};

export const LANGUAGES = {
  de: 'Deutsch',
  en: 'English',
};
export const LANGUAGE_CODES = Object.keys(LANGUAGES);
export const { API_URL, APP_VERSION, CONFIG_STATIC_URL, CONFIG_CMS_URL } = Env;
export const ENVIRONMENT = environmentMapping[document.location.hostname] ?? environmentMapping.default;

export const LANGUAGE_QUERY_PARAM = 'lng';

import { IS_PREVIEW_MODE } from './preview';
export { IS_PREVIEW_MODE } from './preview';

export const CONFIG_URL = IS_PREVIEW_MODE ? CONFIG_CMS_URL || CONFIG_STATIC_URL : CONFIG_STATIC_URL;
export const CONTENT_URL = `${CONFIG_URL}/pages`;

export {
  AUTH_PERSIST_TOKENS,
  AUTH_PROVIDER,
  AUTH_CLIENT_ID,
  AUTH_SCOPE,
  AUTH_AUTHORIZE_URI,
  AUTH_TOKEN_URI,
  AUTH_LOGOUT_URI,
  AUTH_CHALLENGE_VERIFIER_LENGTH,
  AUTH_CHALLENGE_VERIFIER_MASK,
} from './auth';

export { ROUTES, URL_METADATA_COMPLETENESS_INFO } from './navigation';
export const NAVIGATION_CONFIG = NavigationConfig;

export const HOME_CONFIG = HomeConfig;

export const ITEM_CONFIG = ItemConfig;
export { ITEM_FIELDS_CONCATENATOR, ITEM_GC_COUNT } from './item';

export const FEATURE_FLAGS = FeatureFlags;
export { FEATURE_FLAG_LANGUAGE_SWITCHER } from './feature';

export const SEARCH_CONFIG = SearchConfig;
export {
  SEARCH_FACETS_COMBINE_OPERATOR,
  SEARCH_FACETS_MAX_LIMIT,
  SEARCH_FACETS_SINGLE_NODE_LABEL_KEY,
  SEARCH_FACETS_SINGLE_NODE_VALUE_PREFIX,
  SEARCH_INVISIBLE_FACETS,
  SEARCH_OPERATOR_PAD_MAP,
  SEARCH_PAGINATION_START,
  SEARCH_QUERY_EVERYTHING,
  SEARCH_PARAM_PAGE,
  SEARCH_PARAM_SORTING_AXIS,
  SEARCH_PARAM_SORTING_ORDER,
  SEARCH_PARAM_FOCUS,
  SEARCH_PARAM_FILTER_PREFIX,
  SEARCH_DOWNLOAD_CONFIG,
  SEARCH_TIMEOUT,
} from './search';

export const BROWSE_CONFIG = BrowseConfig;
export { BROWSE_FACETS_LIMIT } from './browse';

export const ANALYTICS_CONFIG = AnalyticsConfig;
export {
  ANALYTICS_IS_DNT,
  ANALYTICS_IS_ENABLED,
  ANALYTICS_CUSTOM_EVENTS,
  ANALYTICS_HEARTBEAT_TRACKING_INTERVAL,
} from './analytics';

export const ENTITY_TYPES = EntityTypes;

export const FIELDS_CONFIG = FieldsConfig;
export { FIELDS, FIELD_CONCEPT_IDS, FIELD_CONCEPT_PREFIXES } from './fields';
export const FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY = 'fields.emptyValueFallback';
