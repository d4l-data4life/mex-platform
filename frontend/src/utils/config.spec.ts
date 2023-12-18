jest.mock('stencil-router-v2');

import * as config from 'config';
import { EntityTypeName } from 'config/entity-types';
import { SidebarFeature } from 'config/item';
import { FieldEntityVirtualType } from 'models/field';
import {
  aggregateSidebarFeatures,
  unflattenEntityTypeFieldsMapping,
  unflattenEntityTypeSidebarFeaturesMapping,
} from './config';

const { FIELDS } = config;

const EXAMPLE_ENTITY_TYPE_FIELDS_FLAT_MAPPING = [
  { entityType: FieldEntityVirtualType.all, field: 'foo' },
  { entityType: 'Source' as EntityTypeName, field: 'bar' },
  { entityType: 'Resource' as EntityTypeName, field: 'baz' },
  { entityType: FieldEntityVirtualType.unknown, field: 'test' },
];

const EXAMPLE_ENTITY_TYPE_SIDEBAR_FEATURES_FLAT_MAPPING = [
  { entityType: FieldEntityVirtualType.all, feature: SidebarFeature.completeness },
  { entityType: 'Resource' as EntityTypeName, feature: SidebarFeature.date, field: 'created' },
  { entityType: 'Resource' as EntityTypeName, feature: SidebarFeature.displayField, field: 'identifier' },
  { entityType: 'Source' as EntityTypeName, feature: SidebarFeature.date, field: 'timeRange' },
  {
    entityType: 'Resource' as EntityTypeName,
    feature: SidebarFeature.accessRestriction,
    field: 'accessRestriction',
  },
  { entityType: FieldEntityVirtualType.unknown, feature: SidebarFeature.displayField, field: 'foo' },
];

const EXAMPLE_ENTITY_TYPE_SIDEBAR_FEATURES_MAPPING = {
  Datum: [{ feature: SidebarFeature.completeness }],
  Person: [{ feature: SidebarFeature.completeness }],
  Platform: [{ feature: SidebarFeature.completeness }],
  Resource: [
    { feature: SidebarFeature.completeness },
    { feature: SidebarFeature.date, field: FIELDS.created },
    { feature: SidebarFeature.displayField, field: FIELDS.identifier },
    { feature: SidebarFeature.accessRestriction, field: FIELDS.accessRestriction },
  ],
  Source: [{ feature: SidebarFeature.completeness }, { feature: SidebarFeature.date, field: FIELDS.timeRange }],
  OrganizationalUnit: [{ feature: SidebarFeature.completeness }],
  default: [{ feature: SidebarFeature.completeness }, { feature: SidebarFeature.displayField, field: FIELDS.foo }],
};

const createEntityTypeConfig = (name: string) => ({
  name,
  config: {
    isFocal: true,
    isAggregatable: false,
  },
});

describe('config util', () => {
  beforeAll(() => {
    config.ENTITY_TYPES.Datum = createEntityTypeConfig('Datum');
    config.ENTITY_TYPES.Person = createEntityTypeConfig('Person');
    config.ENTITY_TYPES.Platform = createEntityTypeConfig('Platform');
    config.ENTITY_TYPES.Resource = createEntityTypeConfig('Resource');
    config.ENTITY_TYPES.Source = createEntityTypeConfig('Source');
    config.ENTITY_TYPES.OrganizationalUnit = createEntityTypeConfig('OrganizationalUnit');
  });

  describe('unflattenEntityTypeFieldsMapping()', () => {
    it('converts a flat entity type field mapping to an object with entity types as key and array of fields as value', () => {
      const expected = {
        Datum: [FIELDS.foo],
        Person: [FIELDS.foo],
        Platform: [FIELDS.foo],
        Resource: [FIELDS.foo, FIELDS.baz],
        Source: [FIELDS.foo, FIELDS.bar],
        OrganizationalUnit: [FIELDS.foo],
        default: [FIELDS.foo, FIELDS.test],
      };

      const result = unflattenEntityTypeFieldsMapping(EXAMPLE_ENTITY_TYPE_FIELDS_FLAT_MAPPING);
      expect(JSON.stringify(result)).toBe(JSON.stringify(expected));
    });
  });

  describe('unflattenEntityTypeSidebarFeaturesMapping()', () => {
    it('converts a flat entity type sidebar feature mapping to an object with entity types as key and array of sidebar configs as value', () => {
      expect(
        JSON.stringify(unflattenEntityTypeSidebarFeaturesMapping(EXAMPLE_ENTITY_TYPE_SIDEBAR_FEATURES_FLAT_MAPPING))
      ).toBe(JSON.stringify(EXAMPLE_ENTITY_TYPE_SIDEBAR_FEATURES_MAPPING));
    });
  });

  describe('aggregateSidebarFeatures()', () => {
    it('aggregates multiple sequential sidebar features of the same type and adds an index', () => {
      // completeness, date and accessRestriction are grouped together

      const resultForResource = aggregateSidebarFeatures(EXAMPLE_ENTITY_TYPE_SIDEBAR_FEATURES_MAPPING.Resource);

      expect(resultForResource[0]?.feature).toBe('itemInfo');
      expect(JSON.stringify(resultForResource[0].configs)).toBe(
        JSON.stringify([
          { feature: SidebarFeature.completeness, index: 0 },
          {
            feature: SidebarFeature.date,
            field: FIELDS.created,
            index: 1,
          },
        ])
      );

      expect(resultForResource[1]?.feature).toBe(SidebarFeature.displayField);
      expect(JSON.stringify(resultForResource[1].configs)).toBe(
        JSON.stringify([{ feature: SidebarFeature.displayField, field: FIELDS.identifier, index: 2 }])
      );

      expect(resultForResource[2]?.feature).toBe('itemInfo');
      expect(JSON.stringify(resultForResource[2].configs)).toBe(
        JSON.stringify([{ feature: SidebarFeature.accessRestriction, field: FIELDS.accessRestriction, index: 3 }])
      );
    });
  });
});
