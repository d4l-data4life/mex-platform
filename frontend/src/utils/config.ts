import { ENTITY_TYPES, FIELDS } from 'config';
import { EntityTypeFieldsMapping } from 'config/fields';
import { EntityTypeSidebarFeaturesMapping, SidebarFeature, SidebarFeatureConfig } from 'config/item';
import { Field, FieldEntityVirtualType } from 'models/field';
import {
  EntityTypeFieldsFlatMap,
  EntityTypeFieldsFlatMapItem,
  EntityTypeSidebarFeaturesFlatMap,
  EntityTypeSidebarFeaturesFlatMapItem,
} from 'services/config';

export type AggregatedSidebarFeatures = {
  feature: SidebarFeature | 'itemInfo';
  configs: AggregatedSidebarFeatureConfig[];
}[];

export interface AggregatedSidebarFeatureConfig extends SidebarFeatureConfig {
  index: number;
}

const unflattenMapping = <M = unknown, MI = unknown, FMI = unknown>(fm: FMI[], mapFn: (configItem: FMI) => MI): M => {
  return [...Object.keys(ENTITY_TYPES), 'default'].reduce((mapping, key) => {
    const configItems = fm.filter(
      (configItem: any) =>
        configItem.entityType === FieldEntityVirtualType.all ||
        (ENTITY_TYPES[configItem.entityType]?.name ?? 'default') === key
    );
    return Object.assign(mapping, {
      [key]: configItems.map(mapFn),
    });
  }, {}) as M;
};

export const unflattenEntityTypeFieldsMapping = (flatMapping: EntityTypeFieldsFlatMap): EntityTypeFieldsMapping => {
  return unflattenMapping<EntityTypeFieldsMapping, Field, EntityTypeFieldsFlatMapItem>(flatMapping, (configItem) => {
    return configItem.targetField
      ? FIELDS[configItem.field].cloneAndLinkToField(FIELDS[configItem.targetField])
      : FIELDS[configItem.field];
  });
};

export const unflattenEntityTypeSidebarFeaturesMapping = (
  flatMapping: EntityTypeSidebarFeaturesFlatMap
): EntityTypeSidebarFeaturesMapping => {
  return unflattenMapping<EntityTypeSidebarFeaturesMapping, SidebarFeatureConfig, EntityTypeSidebarFeaturesFlatMapItem>(
    flatMapping,
    ({ feature, field, displayField, key, actions, iconMap }) => {
      return {
        feature,
        ...(field ? { field: FIELDS[field] } : {}),
        ...(displayField ? { displayField: FIELDS[displayField] } : {}),
        ...(key ? { key } : {}),
        ...(actions
          ? {
              actions: actions.map((action) => ({
                ...action,
                ...(action.field ? { field: FIELDS[action.field as string] } : {}),
              })),
            }
          : {}),
        ...(iconMap ? { iconMap } : {}),
      };
    }
  );
};

const itemInfoSidebarFeatures = [SidebarFeature.completeness, SidebarFeature.date, SidebarFeature.accessRestriction];

export const aggregateSidebarFeatures = (configs: SidebarFeatureConfig[]) => {
  return configs.reduce((aggregatedConfigs, config, index) => {
    const latest = aggregatedConfigs[aggregatedConfigs.length - 1];
    const aggregatedConfig = { ...config, index };
    const feature = itemInfoSidebarFeatures.includes(config.feature) ? 'itemInfo' : config.feature;

    if (!latest || latest.feature !== feature) {
      return aggregatedConfigs.concat([{ feature, configs: [aggregatedConfig] }]);
    }

    latest.configs.push(aggregatedConfig);
    return aggregatedConfigs;
  }, [] as AggregatedSidebarFeatures);
};
