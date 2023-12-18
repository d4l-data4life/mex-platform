import { createRouter } from 'stencil-router-v2';

import analytics from './analytics';
import auth from './auth';
import filters from './filters';
import i18n from './i18n';
import items from './items';
import notifications from './notifications';
import search from './search';

const stores = {
  analytics,
  auth,
  filters,
  i18n,
  items,
  notifications,
  router: createRouter(),
  search,
};

export type Stores = typeof stores;

export default stores;
