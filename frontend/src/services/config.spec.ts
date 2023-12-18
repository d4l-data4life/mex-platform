jest.mock('stencil-router-v2');

import i18n from 'i18next';
import { Env } from '@stencil/core';
const CONFIG_URL = 'https://example.org/config';
Env.CONFIG_STATIC_URL = CONFIG_URL;

import stores from 'stores';
import * as fetchClient from 'utils/fetch-client';
import {
  SEARCH_CONFIG,
  BROWSE_CONFIG,
  ITEM_CONFIG,
  FIELDS,
  FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY,
  FIELDS_CONFIG,
  NAVIGATION_CONFIG,
  ANALYTICS_CONFIG,
  ENTITY_TYPES,
  HOME_CONFIG,
} from 'config';
import configService from './config';
import { Field } from 'models/field';
import { SidebarFeature } from 'config/item';
import { AnalyticsServiceProvider } from 'config/analytics';

const requestSpy = jest.spyOn(fetchClient, 'request');
requestSpy.mockImplementation(async () => [{} as any, new Headers()] as [any, Headers]);

const EXAMPLE_CORE_CONFIG_RESPONSE = {
  config: {
    search: {
      limit: 5,
      facetLimit: 8,
      paginationRangeCount: 3,
      queryMaxEditDistance: 2,
      queryUseNGramField: true,
      ordinalAxes: [
        {
          name: 'accessRestriction',
          uiField: 'accessRestriction',
        },
        {
          name: 'created',
          uiField: 'created',
        },
        {
          name: 'departmentOrUnit',
          uiField: 'departmentOrUnit',
        },
        {
          name: 'entityName',
          uiField: 'entityName',
        },
      ],
      displayedFields: [
        {
          entityType: 'all',
          field: 'label',
        },
        {
          entityType: 'Resource',
          field: 'keyword',
        },
        {
          entityType: 'unknown',
          field: 'identifier',
        },
      ],
      sidebarFeatures: [
        {
          entityType: 'Resource',
          feature: 'completeness',
        },
        {
          entityType: 'Resource',
          feature: 'date',
          field: 'created',
        },
        {
          entityType: 'Source',
          feature: 'completeness',
        },
      ],
      foci: [null, 'title'],
      sorting: [
        {
          axis: 'created',
          order: 'desc',
        },
      ],
      facets: [
        {
          axis: 'created',
          type: 'yearRange',
        },
        {
          axis: 'accessRestriction',
          type: 'exact',
        },
        {
          axis: 'departmentOrUnit',
          type: 'hierarchy',
          entityType: 'OrganizationalUnit',
          linkField: 'parentDepartment',
          displayField: 'label',
        },
      ],
    },
    browse: {
      tabs: [
        {
          axis: 'accessRestriction',
          type: 'facet',
        },
        {
          axis: 'departmentOrUnit',
          type: 'hierarchy',
          entityType: 'OrganizationalUnit',
          linkField: 'parentDepartment',
          displayField: 'label',
        },
      ],
    },
    item: {
      displayedFields: [
        {
          entityType: 'all',
          field: 'label',
        },
        {
          entityType: 'Resource',
          field: 'accessRestriction',
        },
        {
          entityType: 'Source',
          field: 'description',
        },
        {
          entityType: 'unknown',
          field: 'created',
        },
      ],
      dedicatedViewSupportedEntityTypes: ['Source', 'Resource'],
      sidebarFeatures: [
        {
          entityType: 'Source',
          feature: 'contactForm',
          field: 'contact',
          displayField: 'labelWithEmail',
        },
        {
          entityType: 'Resource',
          feature: 'displayField',
          field: 'identifier',
        },
      ],
      relatedSearch: [
        {
          entityType: 'Source',
          targetEntityType: 'Resource',
          linkedField: 'accessRestriction',
          limit: 3,
        },
      ],
    },
    analytics: {
      provider: 'matomo',
      matomoSiteId: 999,
      ignoreDnt: true,
    },
    languageSwitcher: {
      languages: ['de', 'en'],
    },
    home: {
      chart: {
        axis: 'accessRestriction',
        greenColorBucket: 'open',
        redColorBucket: 'restricted',
      },
      supportEmailAddress: 'test@test.test',
      dashboardMetrics: [
        {
          entityType: 'Resource',
          axis: 'entityName',
          method: 'bucket',
        },
      ],
      latestUpdateAxis: 'created',
    },
    fields: {
      completenessRankingWeights: {
        mandatory: 60,
        recommended: 25,
        optional: 15,
      },
      defaultFieldConceptPrefix: 'https://unit.test/concept/',
      externalFieldConceptPrefixes: ['https://www.some-other.org/concepts/'],
      emptyValueFallback: {
        de: '-',
        en: '-',
      },
    },
  },
};

const EXAMPLE_FIELDS_CONFIG_RESPONSE = [
  {
    name: 'accessRestriction',
    renderer: 'plain',
    importance: 'mandatory',
    isVirtual: false,
    isEnumerable: true,
    label: {
      de: [
        {
          entityType: 'all',
          singular: 'Zugriffsbeschränkung',
        },
      ],
      en: [
        {
          entityType: 'all',
          singular: 'Access restriction',
        },
      ],
    },
  },
  {
    name: 'created',
    renderer: 'time',
    importance: 'mandatory',
    isVirtual: false,
    isEnumerable: false,
    label: {
      de: [
        {
          entityType: 'all',
          singular: 'Erstellt',
        },
      ],
      en: [
        {
          entityType: 'all',
          singular: 'Created',
        },
      ],
    },
    description: {
      de: [
        {
          entityType: 'Resource',
          text: 'DE created field description for resources',
        },
      ],
      en: [
        {
          entityType: 'all',
          text: 'EN created field description for all entities',
        },
      ],
    },
  },
  {
    name: 'departmentOrUnit',
    renderer: 'reference',
    importance: 'mandatory',
    isVirtual: false,
    isEnumerable: false,
    linkedField: 'label',
  },
];

const EXAMPLE_TRANSLATIONS_RESPONSE = {
  de: {
    foo: {
      bar: 'DE foo.bar',
    },
    baz: 'DE baz',
  },
  en: {
    foo: {
      bar: 'EN foo.bar',
    },
    baz: 'EN baz',
  },
};

const EXAMPLE_ENTITY_TYPES_RESPONSE = {
  entityTypes: [
    { name: 'Resource', config: { isFocal: true, isAggregatable: false } },
    { name: 'Source', config: { isFocal: true, isAggregatable: false } },
    { name: 'OrganizationalUnit', config: { isFocal: true, isAggregatable: false } },
  ],
};

describe('config service', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('fetch()', () => {
    it('fetches the core config (search, browse, item)', async () => {
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_CORE_CONFIG_RESPONSE, new Headers()]);
      expect(await configService.fetch()).toBe(EXAMPLE_CORE_CONFIG_RESPONSE);
      expect(requestSpy).toHaveBeenCalledWith(expect.objectContaining({ method: 'GET', url: `${CONFIG_URL}/` }));
    });
  });

  describe('fetchFields()', () => {
    it('fetches the fields config', async () => {
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_FIELDS_CONFIG_RESPONSE, new Headers()]);
      expect(await configService.fetchFields()).toBe(EXAMPLE_FIELDS_CONFIG_RESPONSE);
      expect(requestSpy).toHaveBeenCalledWith(expect.objectContaining({ method: 'GET', url: `${CONFIG_URL}/fields` }));
    });
  });

  describe('fetchEntityTypes()', () => {
    it('fetches the entity types config', async () => {
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_ENTITY_TYPES_RESPONSE, new Headers()]);
      expect(await configService.fetchEntityTypes()).toBe(EXAMPLE_ENTITY_TYPES_RESPONSE);
      expect(requestSpy).toHaveBeenCalledWith(
        expect.objectContaining({ method: 'GET', url: `${CONFIG_URL}/entity_types` })
      );
    });
  });

  describe('loadTranslations()', () => {
    it('fetches and loads the UI translations', async () => {
      stores.i18n.language = 'de';
      expect(stores.i18n.t('foo.bar')).toBe('foo.bar');
      expect(stores.i18n.t('baz')).toBe('baz');

      requestSpy.mockImplementationOnce(async () => [EXAMPLE_TRANSLATIONS_RESPONSE, new Headers()]);
      await configService.loadTranslations();

      expect(requestSpy).toHaveBeenCalledWith(
        expect.objectContaining({ method: 'GET', url: `${CONFIG_URL}/translations` })
      );
      expect(stores.i18n.t('foo.bar')).toBe('DE foo.bar');
      expect(stores.i18n.t('baz')).toBe('DE baz');

      stores.i18n.language = 'en';
      expect(stores.i18n.t('foo.bar')).toBe('EN foo.bar');
      expect(stores.i18n.t('baz')).toBe('EN baz');
    });
  });

  describe('load()', () => {
    beforeEach(() => {
      jest.resetAllMocks();
      jest.clearAllMocks();

      requestSpy.mockImplementationOnce(async () => [EXAMPLE_FIELDS_CONFIG_RESPONSE, new Headers()]);
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_ENTITY_TYPES_RESPONSE, new Headers()]);
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_CORE_CONFIG_RESPONSE, new Headers()]);
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_TRANSLATIONS_RESPONSE, new Headers()]);
    });

    it('fetches core config and populates config stores accordingly', async () => {
      expect(FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX).toBe('');
      expect(FIELDS_CONFIG.EXTERNAL_FIELD_CONCEPT_PREFIXES).toEqual([]);
      expect(FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS).toEqual({ mandatory: 60, recommended: 30, optional: 10 });

      expect(Object.keys(ENTITY_TYPES)).toEqual([]);

      expect(SEARCH_CONFIG.DISPLAYED_FIELDS).toEqual({ default: [] });
      expect(SEARCH_CONFIG.FACETS).toEqual([]);
      expect(SEARCH_CONFIG.FACETS_LIMIT).toBe(20);
      expect(SEARCH_CONFIG.FOCI).toEqual([null]);
      expect(SEARCH_CONFIG.LIMIT).toBe(10);
      expect(SEARCH_CONFIG.ORDINAL_AXES).toEqual([]);
      expect(SEARCH_CONFIG.PAGINATION_RANGE_COUNT).toBe(4);
      expect(SEARCH_CONFIG.QUERY_MAX_EDIT_DISTANCE).toBe(0);
      expect(SEARCH_CONFIG.QUERY_USE_NGRAM_FIELD).toBe(false);
      expect(SEARCH_CONFIG.SIDEBAR_FEATURES).toEqual({ default: [] });
      expect(SEARCH_CONFIG.SORTING_OPTIONS).toEqual([{ axis: null, order: 'desc' }]);

      expect(BROWSE_CONFIG.TABS).toEqual([]);

      expect(ITEM_CONFIG.DISPLAYED_FIELDS).toEqual({ default: [] });
      expect(ITEM_CONFIG.RELATED_RESULTS_CONFIG).toEqual({});
      expect(ITEM_CONFIG.SIDEBAR_FEATURES).toEqual({ default: [] });
      expect(ITEM_CONFIG.DEDICATED_VIEW_SUPPORTED_ENTITY_TYPES).toEqual([]);

      expect(HOME_CONFIG.CHART).toBe(undefined);
      expect(HOME_CONFIG.DASHBOARD_METRIC_CONFIGS).toEqual([]);
      expect(HOME_CONFIG.LATEST_UPDATE_AXIS).toBe(undefined);

      expect(NAVIGATION_CONFIG.SUPPORT_EMAIL).toBe(undefined);
      expect(NAVIGATION_CONFIG.LANGUAGE_SWITCHER_LANGUAGES).toBe(undefined);

      expect(ANALYTICS_CONFIG.PROVIDER).toBe(undefined);
      expect(ANALYTICS_CONFIG.SITE_ID).toBe(undefined);
      expect(ANALYTICS_CONFIG.IGNORE_DNT).toBe(false);

      await configService.load();

      expect(FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX).toBe('https://unit.test/concept/');
      expect(FIELDS_CONFIG.EXTERNAL_FIELD_CONCEPT_PREFIXES).toEqual(['https://www.some-other.org/concepts/']);
      expect(FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS).toEqual({ mandatory: 60, recommended: 25, optional: 15 });

      expect(Object.keys(ENTITY_TYPES)).toEqual(['Resource', 'Source', 'OrganizationalUnit']);

      expect(SEARCH_CONFIG.DISPLAYED_FIELDS.Source?.[0]?.name).toBe('label');
      expect(SEARCH_CONFIG.FACETS[0]?.axis?.name).toBe('created');
      expect(SEARCH_CONFIG.FACETS[2]?.linkField?.name).toBe('parentDepartment');
      expect(SEARCH_CONFIG.FACETS_LIMIT).toBe(8);
      expect(SEARCH_CONFIG.FOCI).toEqual([null, 'title']);
      expect(SEARCH_CONFIG.LIMIT).toBe(5);
      expect(SEARCH_CONFIG.ORDINAL_AXES?.[0]?.name).toBe('accessRestriction');
      expect(SEARCH_CONFIG.PAGINATION_RANGE_COUNT).toBe(3);
      expect(SEARCH_CONFIG.QUERY_MAX_EDIT_DISTANCE).toBe(2);
      expect(SEARCH_CONFIG.QUERY_USE_NGRAM_FIELD).toBe(true);
      expect(SEARCH_CONFIG.SIDEBAR_FEATURES.Source).toEqual([
        {
          feature: 'completeness',
        },
      ]);
      expect(SEARCH_CONFIG.SORTING_OPTIONS[0]?.axis?.uiField?.name).toBe('created');

      expect(BROWSE_CONFIG.TABS?.[0]?.key?.name).toBe('accessRestriction');
      expect(BROWSE_CONFIG.TABS?.[1]?.entityType).toBe('OrganizationalUnit');

      expect(ITEM_CONFIG.DISPLAYED_FIELDS?.Resource?.map(({ name }) => name)).toEqual(['label', 'accessRestriction']);
      expect(ITEM_CONFIG.RELATED_RESULTS_CONFIG?.Source).toEqual({
        limit: 3,
        linkedField: FIELDS.accessRestriction,
        targetEntityType: 'Resource',
      });
      expect(ITEM_CONFIG.SIDEBAR_FEATURES?.Source?.[1]?.feature).toBe(SidebarFeature.contactForm);
      expect(ITEM_CONFIG.DEDICATED_VIEW_SUPPORTED_ENTITY_TYPES).toEqual(['Source', 'Resource']);

      expect(HOME_CONFIG.CHART?.axis?.name).toBe('accessRestriction');
      expect(HOME_CONFIG.CHART?.greenColorBucket).toBe('https://unit.test/concept/open');
      expect(HOME_CONFIG.DASHBOARD_METRIC_CONFIGS[0]?.axis.name).toBe('entityName');
      expect(HOME_CONFIG.DASHBOARD_METRIC_CONFIGS[0]?.entityType.name).toBe('Resource');
      expect(HOME_CONFIG.DASHBOARD_METRIC_CONFIGS[0]?.method).toBe('bucket');
      expect(HOME_CONFIG.LATEST_UPDATE_AXIS?.name).toBe('created');

      expect(NAVIGATION_CONFIG.SUPPORT_EMAIL).toBe('test@test.test');
      expect(NAVIGATION_CONFIG.LANGUAGE_SWITCHER_LANGUAGES).toEqual(['de', 'en']);

      expect(ANALYTICS_CONFIG.PROVIDER).toBe(AnalyticsServiceProvider.MATOMO);
      expect(ANALYTICS_CONFIG.SITE_ID).toBe(999);
      expect(ANALYTICS_CONFIG.IGNORE_DNT).toBe(true);
    });

    it('fetches the fields config and poplates the store accordingly', async () => {
      FIELDS.created = Field.createPending('created');
      expect(FIELDS.created.isInitialized).toBe(false);

      await configService.load();

      expect(FIELDS.created.isInitialized).toBe(true);
    });

    it('adds fields translations to the UI translations', async () => {
      stores.i18n.language = 'de';
      i18n.addResourceBundle('de', 'ui', {
        fields: undefined,
      });

      const KEY_ACCESS_RESTRICTION_LABEL = 'fields.labels.accessRestriction.all.singular';
      const KEY_CREATED_DESCRIPTION_RESOURCE = 'fields.descriptions.created.Resource.text';
      const KEY_CREATED_DESCRIPTION_ALL = 'fields.descriptions.created.all.text';

      expect(stores.i18n.t('fields')).toBe('fields');
      expect(stores.i18n.t(KEY_ACCESS_RESTRICTION_LABEL)).toBe(KEY_ACCESS_RESTRICTION_LABEL);

      await configService.load();

      expect(stores.i18n.t(KEY_ACCESS_RESTRICTION_LABEL)).toBe('Zugriffsbeschränkung');
      expect(stores.i18n.t(KEY_CREATED_DESCRIPTION_RESOURCE)).toBe('DE created field description for resources');
      expect(stores.i18n.t(KEY_CREATED_DESCRIPTION_ALL)).toBe(KEY_CREATED_DESCRIPTION_ALL);

      expect(stores.i18n.t(FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY)).toBe('-');

      stores.i18n.language = 'en';
      expect(stores.i18n.t(KEY_ACCESS_RESTRICTION_LABEL)).toBe('Access restriction');
      expect(stores.i18n.t(KEY_CREATED_DESCRIPTION_RESOURCE)).toBe(KEY_CREATED_DESCRIPTION_RESOURCE);
      expect(stores.i18n.t(KEY_CREATED_DESCRIPTION_ALL)).toBe('EN created field description for all entities');
    });
  });
});
