import { Resource } from 'i18next';
import {
  ANALYTICS_CONFIG,
  BROWSE_CONFIG,
  CONFIG_URL,
  FIELDS,
  FIELDS_CONFIG,
  FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY,
  ITEM_CONFIG,
  LANGUAGE_CODES,
  NAVIGATION_CONFIG,
  SEARCH_CONFIG,
  ENTITY_TYPES,
  HOME_CONFIG,
} from 'config';
import { BrowseItemConfigType } from 'config/browse';
import { SidebarFeature, SidebarFeatureConfigActionItem, SidebarFeatureConfigIconMapItem } from 'config/item';
import { get } from 'utils/fetch-client';
import { denormalizeConceptId } from 'utils/field';
import { Field, FieldConfig, configurePendingFields, FieldEntityVirtualType } from 'models/field';
import { addUiTranslations, addFieldTranslations, addAdditionalTranslations } from 'stores/i18n';
import { unflattenEntityTypeFieldsMapping, unflattenEntityTypeSidebarFeaturesMapping } from 'utils/config';
import stores from 'stores';
import { AnalyticsServiceProvider } from 'config/analytics';
import { EntityType, EntityTypeName } from 'config/entity-types';
import { DashboardMetricMethod } from 'config/home';

export interface EntityTypeFieldsFlatMapItem {
  entityType: EntityTypeName | FieldEntityVirtualType;
  field: string;
  targetField?: string;
}

export type EntityTypeFieldsFlatMap = EntityTypeFieldsFlatMapItem[];

export interface EntityTypeSidebarFeaturesFlatMapItem {
  entityType: EntityTypeName | FieldEntityVirtualType;
  feature: SidebarFeature;
  field?: string;
  displayField?: string;
  key?: string;
  actions?: SidebarFeatureConfigActionItem[];
  iconMap?: SidebarFeatureConfigIconMapItem[];
}

export type EntityTypeSidebarFeaturesFlatMap = EntityTypeSidebarFeaturesFlatMapItem[];

interface ConfigResponse {
  config: {
    search: {
      limit: number;
      facetLimit: number;
      paginationRangeCount: number;
      queryMaxEditDistance: 0 | 1;
      queryUseNGramField: boolean;
      displayedFields: EntityTypeFieldsFlatMap;
      sidebarFeatures: EntityTypeSidebarFeaturesFlatMap;
      ordinalAxes: {
        name: string;
        uiField: string;
      }[];
      facets: {
        axis: string;
        type: 'exact' | 'yearRange' | 'hierarchy';
        entityType: EntityTypeName;
        linkField: string;
        displayField: string;
        minLevel: number;
        maxLevel: number;
      }[];
      foci: (string | null)[];
      sorting: {
        axis: string | null;
        order: 'asc' | 'desc';
      }[];
    };
    browse: {
      tabs: {
        key: string;
        type: keyof BrowseItemConfigType;
        axis: string;
        entityType?: EntityTypeName;
        linkField?: string;
        displayField?: string;
        minLevel?: number;
        maxLevel?: number;
      }[];
    };
    item: {
      displayedFields: EntityTypeFieldsFlatMap;
      sidebarFeatures: EntityTypeSidebarFeaturesFlatMap;
      dedicatedViewSupportedEntityTypes: EntityTypeName[];
      relatedSearch: {
        entityType: EntityTypeName;
        targetEntityType: EntityTypeName;
        linkedField: string;
        limit: number;
      }[];
    };
    languageSwitcher?: {
      languages: string[];
    };
    analytics?: {
      provider: 'matomo';
      matomoSiteId?: number;
      ignoreDnt: boolean;
    };
    home: {
      supportEmailAddress?: string;
      chart?: {
        axis: string;
        greenColorBucket: string;
        redColorBucket: string;
      };
      dashboardMetrics: {
        entityType: EntityTypeName;
        axis: string;
        method: DashboardMetricMethod;
      }[];
      latestUpdateAxis?: string;
    };
    fields: {
      defaultFieldConceptPrefix?: string;
      externalFieldConceptPrefixes: string[];
      emptyValueFallback?: {
        [language: string]: string;
      };
      completenessRankingWeights: {
        mandatory: number;
        recommended: number;
        optional: number;
      };
    };
  };
}

export class ConfigService {
  async fetch(): Promise<ConfigResponse> {
    const [response] = await get<ConfigResponse>({
      url: `${CONFIG_URL}/`,
    });

    return response;
  }

  async fetchFields(): Promise<FieldConfig[]> {
    const [fields] = await get<FieldConfig[]>({
      url: `${CONFIG_URL}/fields`,
    });

    return fields;
  }

  async fetchEntityTypes(): Promise<{ entityTypes: EntityType[] }> {
    const [response] = await get<{ entityTypes: EntityType[] }>({
      url: `${CONFIG_URL}/entity_types`,
    });

    return response;
  }

  async loadTranslations(): Promise<void> {
    const [resources] = await get<Resource>({
      url: `${CONFIG_URL}/translations`,
    });

    addUiTranslations(resources);
  }

  async load(): Promise<void> {
    const fields = await this.fetchFields();
    fields.forEach((config) => (FIELDS[config.name] = new Field(config)));
    configurePendingFields();

    const { entityTypes } = await this.fetchEntityTypes();
    entityTypes?.forEach((entityType) => (ENTITY_TYPES[entityType.name] = entityType));

    const { config } = await this.fetch();
    const {
      search,
      browse,
      item,
      languageSwitcher,
      analytics,
      home,
      fields: {
        emptyValueFallback,
        defaultFieldConceptPrefix,
        externalFieldConceptPrefixes,
        completenessRankingWeights,
      },
    } = config;

    FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX = defaultFieldConceptPrefix ?? '';
    FIELDS_CONFIG.EXTERNAL_FIELD_CONCEPT_PREFIXES = externalFieldConceptPrefixes ?? [];
    FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS =
      completenessRankingWeights ?? FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS;

    addFieldTranslations();

    addAdditionalTranslations(FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY, {
      de: emptyValueFallback?.de ?? '--',
      en: emptyValueFallback?.en ?? '--',
    });

    SEARCH_CONFIG.FACETS_LIMIT = search.facetLimit;
    SEARCH_CONFIG.LIMIT = search.limit;
    stores.search.limit = search.limit;
    SEARCH_CONFIG.QUERY_MAX_EDIT_DISTANCE = search.queryMaxEditDistance;
    SEARCH_CONFIG.QUERY_USE_NGRAM_FIELD = search.queryUseNGramField;
    SEARCH_CONFIG.PAGINATION_RANGE_COUNT = search.paginationRangeCount;
    SEARCH_CONFIG.DISPLAYED_FIELDS = unflattenEntityTypeFieldsMapping(search.displayedFields);
    SEARCH_CONFIG.ORDINAL_AXES = search.ordinalAxes
      .map(({ name, uiField }) => ({ name, uiField: FIELDS[uiField] }))
      .filter(({ uiField }) => uiField);

    SEARCH_CONFIG.FACETS = search.facets

      .map(({ axis, type, entityType, linkField, displayField, minLevel, maxLevel }) => ({
        axis: SEARCH_CONFIG.ORDINAL_AXES.find(({ name }) => name === axis),
        type,
        ...(entityType && linkField && displayField
          ? {
              entityType,
              linkField: FIELDS[linkField],
              displayField: FIELDS[displayField],
              minLevel,
              maxLevel,
            }
          : {}),
      }))
      .filter(({ type, entityType }) => type !== 'hierarchy' || (entityType && !!ENTITY_TYPES[entityType]))
      .filter(({ axis }) => axis);
    SEARCH_CONFIG.FOCI = search.foci;
    SEARCH_CONFIG.SORTING_OPTIONS = search.sorting.map(({ axis, order }) => ({
      axis: SEARCH_CONFIG.ORDINAL_AXES.find(({ name }) => name === axis) ?? null,
      order,
    }));
    SEARCH_CONFIG.SIDEBAR_FEATURES = unflattenEntityTypeSidebarFeaturesMapping(search.sidebarFeatures);

    BROWSE_CONFIG.TABS = browse.tabs
      .map((tab) => {
        const axis = SEARCH_CONFIG.ORDINAL_AXES.find(({ name }) => name === tab.axis);
        return {
          type: BrowseItemConfigType[tab.type],
          key: axis?.uiField,
          axis,
          entityType: tab.entityType,
          linkField: tab.linkField ? FIELDS[tab.linkField] : undefined,
          displayField: tab.displayField ? FIELDS[tab.displayField] : undefined,
          minLevel: tab.minLevel,
          maxLevel: tab.maxLevel,
        };
      })
      .filter(({ entityType }) => !entityType || !!ENTITY_TYPES[entityType])
      .filter(({ type, axis }) => type && axis?.uiField.isInitialized);

    ITEM_CONFIG.DISPLAYED_FIELDS = unflattenEntityTypeFieldsMapping(item.displayedFields);
    ITEM_CONFIG.DISPLAYED_FIELDS.default = Object.values(FIELDS);
    ITEM_CONFIG.SIDEBAR_FEATURES = unflattenEntityTypeSidebarFeaturesMapping(
      search.sidebarFeatures.concat(item.sidebarFeatures)
    );

    ITEM_CONFIG.DEDICATED_VIEW_SUPPORTED_ENTITY_TYPES = item.dedicatedViewSupportedEntityTypes
      .filter((entityType) => !!ENTITY_TYPES[entityType])
      .filter(Boolean);
    ITEM_CONFIG.RELATED_RESULTS_CONFIG = item.relatedSearch
      .filter(
        ({ entityType, targetEntityType, linkedField, limit }) =>
          !!ENTITY_TYPES[entityType] &&
          !!ENTITY_TYPES[targetEntityType] &&
          FIELDS[linkedField].isInitialized &&
          Number.isInteger(limit)
      )
      .reduce(
        (config, { entityType, targetEntityType, linkedField, limit }) =>
          Object.assign(config, {
            [entityType]: {
              targetEntityType,
              linkedField: FIELDS[linkedField],
              limit,
            },
          }),
        {}
      );

    const homeConfigChartAxis =
      home.chart?.axis && SEARCH_CONFIG.ORDINAL_AXES.find(({ name }) => name === home.chart.axis);
    HOME_CONFIG.CHART = homeConfigChartAxis
      ? {
          axis: homeConfigChartAxis,
          greenColorBucket: denormalizeConceptId(home.chart.greenColorBucket, homeConfigChartAxis.uiField),
          redColorBucket: denormalizeConceptId(home.chart.redColorBucket, homeConfigChartAxis.uiField),
        }
      : null;
    HOME_CONFIG.DASHBOARD_METRIC_CONFIGS =
      home.dashboardMetrics
        ?.map(({ entityType, axis, method }) => ({
          entityType: ENTITY_TYPES[entityType],
          axis: SEARCH_CONFIG.ORDINAL_AXES.find(({ name }) => name === axis),
          method,
        }))
        .filter(({ entityType, axis, method }) => entityType && axis && method) ?? [];
    HOME_CONFIG.LATEST_UPDATE_AXIS =
      home.latestUpdateAxis && SEARCH_CONFIG.ORDINAL_AXES.find(({ name }) => name === home.latestUpdateAxis);

    NAVIGATION_CONFIG.SUPPORT_EMAIL = home.supportEmailAddress;
    NAVIGATION_CONFIG.LANGUAGE_SWITCHER_LANGUAGES = (languageSwitcher?.languages ?? []).filter((language) =>
      LANGUAGE_CODES.includes(language)
    );

    ANALYTICS_CONFIG.PROVIDER = analytics?.provider === 'matomo' ? AnalyticsServiceProvider.MATOMO : null;
    ANALYTICS_CONFIG.SITE_ID = analytics?.matomoSiteId;
    ANALYTICS_CONFIG.IGNORE_DNT = !!analytics?.ignoreDnt;
  }
}

export default new ConfigService();
