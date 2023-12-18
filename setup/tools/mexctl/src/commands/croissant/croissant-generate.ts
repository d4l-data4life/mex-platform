import { Command } from "commander";
import { glob } from "glob";
import * as yaml from "js-yaml";
import * as fs from "fs";
import * as path from "path";

import { Output } from "../../utils/out";
import { EntityType, EntityTypesList, Field, FieldDef, FieldDefsList, Item } from "../../api/defs";

export default function (): Command {
    return new Command("generate")
        .description("generate MEx config from Croissant metadata.json")
        .requiredOption("--in <dir>", "folder containing subfolders with metadata.json files in them")
        .requiredOption("--config-out <dir>", "folder to write MEx config files to")
        .requiredOption("--items-out <dir>", "folder to write the items data to")
        .action(handler);
}

const DEFAULT_FIELD_DEFS: FieldDef[] = [
    {
        name: "createdAt",
        kind: "timestamp",
        indexDef: {
            multiValued: false,
            ext: [],
        },
    },
    {
        name: "entityName",
        kind: "string",
        indexDef: {
            multiValued: false,
            ext: [],
        },
    },
    {
        name: "id",
        kind: "string",
        indexDef: {
            multiValued: false,
            ext: [],
        },
    },
    {
        name: "identifier",
        kind: "string",
        indexDef: {
            multiValued: false,
            ext: [],
        },
    },
    {
        name: "businessId",
        kind: "string",
        indexDef: {
            multiValued: false,
            ext: [],
        },
    },
];

const DEFAULT_FIELDS: Field[] = [
    {
        name: "createdAt",
        renderer: "time",
        kind: "timestamp",
        importance: "none",
        isVirtual: false,
        isMultiValued: true,
        isEnumerable: false,
        documentation: null,
    },
    {
        name: "entityName",
        renderer: "none",
        kind: "string",
        importance: "mandatory",
        isVirtual: false,
        isMultiValued: true,
        isEnumerable: true,
        documentation: null,
    },
    {
        name: "identifier",
        renderer: "none",
        kind: "string",
        importance: "none",
        isVirtual: false,
        isMultiValued: false,
        isEnumerable: false,
        documentation: null,
    },
];

async function handler(options, command: Command): Promise<void> {
    const out = new Output(command);
    out.verbose(options);

    const files = await glob("**/metadata.json", {
        cwd: options.in,
    });

    const entityTypes: EntityType[] = [
        {
            name: "Dataset",
            config: {
                isFocal: true,
                businessIdFieldName: "identifier",
            },
        },
        {
            name: "FileObject",
            config: {
                isFocal: true,
                businessIdFieldName: "identifier",
            },
        },
        {
            name: "FileSet",
            config: {
                isFocal: true,
                businessIdFieldName: "identifier",
            },
        },
        {
            name: "RecordSet",
            config: {
                isFocal: true,
                businessIdFieldName: "identifier",
            },
        },
        {
            name: "Field",
            config: {
                isFocal: true,
                businessIdFieldName: "identifier",
            },
        },
    ];

    const fieldDefs: FieldDef[] = [
        ...DEFAULT_FIELD_DEFS,
        {
            name: "description",
            kind: "text",
            indexDef: {
                multiValued: true,
                ext: [],
            },
        },
        {
            name: "name",
            kind: "text",
            indexDef: {
                multiValued: false,
                ext: [],
            },
        },
        {
            name: "license",
            kind: "text",
            indexDef: {
                multiValued: true,
                ext: [],
            },
        },
        {
            name: "url",
            kind: "text",
            indexDef: {
                multiValued: true,
                ext: [],
            },
        },
    ];

    const fields: Field[] = [
        ...DEFAULT_FIELDS,
        {
            name: "description",
            renderer: "description",
            kind: "text",
            importance: "mandatory",
            isVirtual: false,
            isMultiValued: true,
            isEnumerable: false,
            documentation: null,
        },
        {
            name: "name",
            renderer: "title",
            kind: "text",
            importance: "mandatory",
            isVirtual: false,
            isMultiValued: true,
            isEnumerable: false,
            documentation: null,
        },
    ];

    const datasetList: Item[] = [];
    const recordSetList: Item[] = [];

    for (const file of files) {
        const metadata = JSON.parse(fs.readFileSync(path.resolve(options.in, file)).toString());
        out.json(metadata.name);

        datasetList.push({
            entityType: "Dataset",
            businessId: metadata.name,
            values: [
                { fieldName: "name", fieldValue: metadata.name },
                { fieldName: "description", fieldValue: metadata.description },
                ...(metadata.license instanceof Array ? metadata.license : [metadata.license]).map((l) => ({
                    fieldName: "license",
                    fieldValue: l ?? "n/a",
                })),
                { fieldName: "url", fieldValue: metadata.url ?? "n/a" },
            ],
        });

        for (const recordSet of metadata.recordSet ?? []) {
            recordSetList.push({
                entityType: "RecordSet",
                businessId: `${metadata.name}-${recordSet.name}`,
                values: [
                    { fieldName: "name", fieldValue: recordSet.name },
                    { fieldName: "description", fieldValue: recordSet.description ?? "n/a" },
                ],
            });
        }
    }

    writeEntityTypes(options.configOut, { entityTypes });
    writeFieldDefs(options.configOut, { fieldDefs });
    writeFields(options.configOut, fields);

    writeItems(options.itemsOut, "dataset", datasetList);
    writeItems(options.itemsOut, "recordset", recordSetList);
}

function writeEntityTypes(rootPath: string, entityTypesList: EntityTypesList) {
    mkdirpSync(path.resolve(rootPath, "entity_types"));

    for (const entityType of entityTypesList.entityTypes) {
        mkdirpSync(path.resolve(rootPath, "entity_types", entityType.name.toLowerCase()));
        fs.writeFileSync(path.resolve(rootPath, "entity_types", entityType.name.toLowerCase(), "index.json"), JSON.stringify(entityType));
    }

    fs.writeFileSync(path.resolve(rootPath, "entity_types", "index.json"), JSON.stringify(entityTypesList));
}

function writeFieldDefs(rootPath: string, fieldDefsList: FieldDefsList) {
    mkdirpSync(path.resolve(rootPath, "field_defs"));

    for (const fieldDef of fieldDefsList.fieldDefs) {
        mkdirpSync(path.resolve(rootPath, "field_defs", fieldDef.name.toLowerCase()));
        fs.writeFileSync(path.resolve(rootPath, "field_defs", fieldDef.name.toLowerCase(), "index.json"), JSON.stringify(fieldDef));
    }

    fs.writeFileSync(path.resolve(rootPath, "field_defs", "index.json"), JSON.stringify(fieldDefsList));
}

function writeFields(rootPath: string, fields: Field[]) {
    mkdirpSync(path.resolve(rootPath, "fields"));

    for (const field of fields) {
        mkdirpSync(path.resolve(rootPath, "fields", field.name.toLowerCase()));
        fs.writeFileSync(path.resolve(rootPath, "fields", field.name.toLowerCase(), "index.json"), JSON.stringify(field));
    }

    fs.writeFileSync(path.resolve(rootPath, "fields", "index.json"), JSON.stringify(fields));
}

function writeItems(itemsPath: string, itemsType: string, items: Item[]) {
    mkdirpSync(itemsPath);
    fs.writeFileSync(path.resolve(itemsPath, `items.${itemsType}.yaml`), yaml.dump(items));
}

function mkdirpSync(dir: string) {
    if (fs.existsSync(dir)) {
        return;
    }
    fs.mkdirSync(dir);
}
