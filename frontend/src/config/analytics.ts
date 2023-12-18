import { createStore } from '@stencil/store';

export enum AnalyticsServiceProvider {
  MATOMO,
}

interface AnalyticsConfig {
  PROVIDER?: AnalyticsServiceProvider.MATOMO;
  SITE_ID?: number;
  IGNORE_DNT: boolean;
}

const store = createStore<AnalyticsConfig>({
  IGNORE_DNT: false,
});

export const ANALYTICS_IS_DNT: () => boolean = () => !store.get('IGNORE_DNT') && navigator.doNotTrack === '1';
export const ANALYTICS_IS_ENABLED: () => boolean = () =>
  !ANALYTICS_IS_DNT() && store.get('PROVIDER') === AnalyticsServiceProvider.MATOMO && !!store.get('SITE_ID');

export const ANALYTICS_HEARTBEAT_TRACKING_INTERVAL: number = 20000;
export const ANALYTICS_CUSTOM_EVENTS: {
  [key: string]: [category: string, action: string];
} = {
  HOME_SUPPORT_FORM: ['Home', 'Support form'], // Button clicked
  SEARCH_PLAIN: ['Search', 'Plain'], // Initiated
  SEARCH_FOCUS: ['Search', 'Focus'], // Initiated, Selected: <selectedFocus>
  SEARCH_BROWSING: ['Search', 'Browsing'], // Initiated: <fieldName>
  SEARCH_FILTER: ['Search', 'Filter'], // Switched on:|Switched off:|Set range:|Removed: <filterName>
  SEARCH_NAVIGATION: ['Search', 'Navigation'], // Reset, Paginated > (int)<pageNumber>, Filter expanded|collapsed|extended, Hierarchy filter subtree expanded|collapsed
  SEARCH_SORTING: ['Search', 'Sorting'], // Changed: <fieldName> (asc|desc)
  ITEM_CONTACT_FORM: ['Expanded view', 'Contact form'], // Button clicked
  ITEM_SEARCH_RESULT_PAGINATION: ['Expanded view', 'Search pagination'], // Next|Prev
  ITEM_VERSION: ['Expanded view', 'Version'], // Versions expanded|collapsed, Changed > (int)<index>
};

export enum MatomoParams {
  SITE_ID = 'idsite',
  USER_ID = 'uid',
  VISITOR_ID = '_id',
  RECORD = 'rec',
  API_VERSION = 'apiv',
  NONCE = 'rand',
  REFERRER_URL = 'urlref',
  SCREEN_RESOLUTION = 'res',
  LOCAL_TIME_HOUR = 'h',
  LOCAL_TIME_MINUTE = 'm',
  LOCAL_TIME_SECOND = 's',
  SEND_IMAGE = 'send_image',
  PING = 'ping',
  CUSTOM_DIMENSION = 'dimension',
  PAGE_TITLE = 'action_name',
  URL = 'url',
  EVENT_CATEGORY = 'e_c',
  EVENT_ACTION = 'e_a',
  EVENT_NAME = 'e_n',
  EVENT_VALUE = 'e_v',
}

export enum MatomoFlag {
  NO,
  YES,
}

export const MATOMO_URL = '/matomo.php';
export const MATOMO_CUSTOM_DIMENSION_PROJECT = 13;
export const MATOMO_CUSTOM_DIMENSION_PRODUCT = 11;

export default store.state;
