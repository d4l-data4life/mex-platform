import * as t from "io-ts";

export const EntityType_T = t.type({
    name: t.string,
    config: t.partial({
        businessIdFieldName: t.string,
        isAggregatable: t.boolean,
        aggregationEntityType: t.string,
        aggregationAlgorithm: t.string,
        partitionField: t.string,
        duplicateStrategy: t.string,
        isFocal: t.boolean,
    }),
});

export type EntityType_A = t.TypeOf<typeof EntityType_T>;
export type EntityType = t.OutputOf<typeof EntityType_T>;

export const EntityTypesList_T = t.type({
    entityTypes: t.array(EntityType_T),
});

export type EntityTypesList_A = t.TypeOf<typeof EntityTypesList_T>;
export type EntityTypesList = t.OutputOf<typeof EntityTypesList_T>;

export const SearchConfig_T = t.type({
    type: t.string,
    name: t.string,
    fields: t.array(t.string),
});

export type SearchConfig_A = t.TypeOf<typeof SearchConfig_T>;
export type SearchConfig = t.OutputOf<typeof SearchConfig_T>;

export const FieldDefIndexDefExtLink_T = t.intersection([
    t.type({
        "@type": t.literal("type.googleapis.com/mex.v0.IndexDefExtLink"),
        relationType: t.string,
    }),
    t.partial({
        linkedTargetFields: t.array(t.string),
    }),
]);

export type FieldDefIndexDefExtLink_A = t.TypeOf<typeof FieldDefIndexDefExtLink_T>;
export type FieldDefIndexDefExtLink = t.OutputOf<typeof FieldDefIndexDefExtLink_T>;

export const FieldDefIndexDefExtHierarchy_T = t.type({
    "@type": t.literal("type.googleapis.com/mex.v0.IndexDefExtHierarchy"),
    codeSystemNameOrNodeEntityType: t.string,
    linkFieldName: t.string,
    displayFieldName: t.string,
});

export type FieldDefIndexDefExtHierarchy_A = t.TypeOf<typeof FieldDefIndexDefExtHierarchy_T>;
export type FieldDefIndexDefExtHierarchy = t.OutputOf<typeof FieldDefIndexDefExtHierarchy_T>;

export const FieldDef_T = t.intersection([
    t.type({
        name: t.string,
        kind: t.string,
        indexDef: t.partial({
            multiValued: t.boolean,
            ext: t.array(t.union([FieldDefIndexDefExtLink_T, FieldDefIndexDefExtHierarchy_T])),
        }),
    }),
    t.partial({
        displayId: t.string,
    }),
]);

export type FieldDef_A = t.TypeOf<typeof FieldDef_T>;
export type FieldDef = t.OutputOf<typeof FieldDef_T>;

export const FieldDefsList_T = t.type({
    fieldDefs: t.array(FieldDef_T),
});

export type FieldDefsList_A = t.TypeOf<typeof FieldDefsList_T>;
export type FieldDefsList = t.OutputOf<typeof FieldDefsList_T>;

export const Item_T = t.intersection([
    t.type({
        entityType: t.string,
        values: t.array(
            t.intersection([
                t.type({
                    fieldName: t.string,
                    fieldValue: t.string,
                }),
                t.partial({
                    language: t.string,
                }),
            ]),
        ),
    }),
    t.partial({
        businessId: t.string,
    }),
]);

export type Item_A = t.TypeOf<typeof Item_T>;
export type Item = t.OutputOf<typeof Item_T>;

export const Field_T = t.type({
    name: t.string,
    kind: t.string,
    renderer: t.string,
    importance: t.string,
    isVirtual: t.boolean,
    isMultiValued: t.boolean,
    isEnumerable: t.boolean,
    documentation: t.null,
});

export type Field_A = t.TypeOf<typeof Field_T>;
export type Field = t.OutputOf<typeof Field_T>;
