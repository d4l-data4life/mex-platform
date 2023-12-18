import { createStore } from '@stencil/store';

export interface LinkItem {
  label: string;
  url: string;
  testAttr?: string;
}

export interface NavigationConfig {
  [language: symbol]: {
    HEADER: LinkItem[];
    FOOTER: LinkItem[];
    LEARNING_ENVIRONMENT: LinkItem[];
    SERVICES: LinkItem[];
  };
  SUPPORT_EMAIL?: string;
  LANGUAGE_SWITCHER_LANGUAGES?: string[];
}

export const store = createStore<NavigationConfig>({});

export const ROUTES = {
  ROOT: '/',
  AUTH: '/auth',
  SEARCH: '/search',
  SEARCH_QUERY: '/search/:query',
  ITEM: '/items/:id',
  LOGOUT: '/logout',
  CONTENT_PAGE: '/pages/:pageid',
};

export const URL_METADATA_COMPLETENESS_INFO = 'https://www.data4life.care';

export default store.state;
