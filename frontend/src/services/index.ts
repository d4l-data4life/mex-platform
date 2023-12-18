import analyticsService from './analytics';
import authService from './auth';
import configService from './config';
import contentService from './content';
import formService from './form';
import itemService from './item';
import navigationService from './navigation';
import searchService from './search';

const services = {
  analytics: analyticsService,
  auth: authService,
  config: configService,
  content: contentService,
  form: formService,
  item: itemService,
  navigation: navigationService,
  search: searchService,
};

export type Services = typeof services;

export default services;
