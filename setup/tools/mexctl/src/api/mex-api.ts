import cliProgress from "cli-progress";

import { pipe } from "fp-ts/lib/function";
import * as A from "fp-ts/lib/Array";

import { SearchConfig, EntityType, FieldDef, Item_A, EntityTypesList } from "./defs";

import { getAccessToken } from "../utils/jwt";
import { Config } from "../utils/config";
import { waitForJobEnd } from "../utils/job-waiter";
import { request } from "../utils/http";
import { wait } from "../utils/attempt";

const MS_FOR_30MIN = 30 * 60 * 1000;

export type ItemResponse = {
    itemId: string;
    createdAt: string;
    entityType: string;
    owner: string;
    businessId: string;
};

export type ListItemsResponse = {
    items: ItemResponse[];
};

export type ApiConfig = {
    mexOrigin: string;
    jwt: string;
    loggingTraceSecret: string;
};

export async function removeAllEntities(config: ApiConfig) {
    const bar = new cliProgress.SingleBar({});
    try {
        const response = await request({
            method: "GET",
            uri: config.mexOrigin + "/cms/entity_types",
            headers: {
                Authorization: `Bearer ${config.jwt}`,
                "X-MEx-Trace": config.loggingTraceSecret,
            },
        });

        if (response.statusCode !== 200) {
            throw new Error("error requesting entity types, status code " + response.statusCode);
        }

        const entityNames = (response.body.obj as EntityTypesList).entityTypes.map((e) => e.name);
        bar.start(entityNames.length, 0);

        for (const en of entityNames) {
            if (en === "id" || en === "entityName" || en === "createdAt" || en === "businessId") {
                continue;
            }
            console.log("Deleting entity: " + en);
            await request({
                method: "DELETE",
                uri: `${config.mexOrigin}/cms/api/v1/entity_types/${en}`,
                verifyCerts: false,
                headers: {
                    Authorization: `Bearer ${config.jwt}`,
                    "X-MEx-Trace": config.loggingTraceSecret,
                },
            });

            bar.increment();
        }
    } catch (e) {
        console.log(`Entity deletion failed with error: ${e}`);
        throw e;
    }

    bar.stop();
}

export async function removeAllSearchConfigs(config: ApiConfig) {
    try {
        await request({
            method: "DELETE",
            uri: `${config.mexOrigin}/cms/api/v1/search_configs`,
            verifyCerts: false,
            headers: {
                Authorization: `Bearer ${config.jwt}`,
                "X-MEx-Trace": config.loggingTraceSecret,
            },
        });
    } catch (e) {
        console.log(`search config deletion failed: ${e}`);
    }
}

export async function removeAllFields(config: ApiConfig) {
    try {
        await request({
            method: "DELETE",
            uri: config.mexOrigin + "/cms/api/v1/field_defs",
            verifyCerts: false,
            headers: {
                Authorization: `Bearer ${config.jwt}`,
                "X-MEx-Trace": config.loggingTraceSecret,
            },
        });
    } catch (e) {
        console.log(`field deletion failed: ${e}`);
    }
}

async function getItemIds(config: ApiConfig): Promise<string[]> {
    const response = await request({
        method: "GET",
        uri: config.mexOrigin + "/api/v0/metadata/items",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${config.jwt}`,
            "X-MEx-Trace": config.loggingTraceSecret,
        },
    });

    if (response.statusCode !== 200) {
        throw Error("Fetching item IDs failed - status: " + response.statusCode);
    }

    const returnedItems = response.body.obj as ListItemsResponse;
    return returnedItems.items.map((i) => i.itemId);
}

async function getItemBusinessIds(config: ApiConfig): Promise<string[]> {
    const response = await request({
        method: "GET",
        uri: config.mexOrigin + "/api/v0/metadata/items",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${config.jwt}`,
            "X-MEx-Trace": config.loggingTraceSecret,
        },
    });

    if (response.statusCode !== 200) {
        throw Error("Fetching item IDs failed - status: " + response.statusCode);
    }

    const returnedItems = response.body.obj as ListItemsResponse;
    return returnedItems.items.map((i) => i.businessId);
}

export async function removeAllItems(config: ApiConfig) {
    const response = await request({
        method: "DELETE",
        uri: config.mexOrigin + "/api/v0/metadata/all_items",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${config.jwt}`,
            "X-MEx-Trace": config.loggingTraceSecret,
        },
    });

    if (response.statusCode !== 204) {
        throw Error(`Failed to delete items - HTTP status code ${response.statusCode} returned`);
    }

    const itemIds = await getItemIds(config);
    if (itemIds.length > 0) {
        console.log(`warning: remaining items: ${itemIds.length}`);
    }
}

function ms(dur: [number, number]): number {
    return Math.round((dur[0] * 1e9 + dur[1]) / 1e6);
}
/**
 * Adds entity definitions to PostgreSQL, including one entity type for each link field for the relation type.
 *
 * @param config
 * @param entityDefs Array of field definitions
 */
export async function addEntities(config: ApiConfig, entityDefs: EntityType[]) {
    const bar = new cliProgress.SingleBar({});
    bar.start(entityDefs.length, 0);

    try {
        for (const entityDef of entityDefs) {
            await request({
                method: "POST",
                uri: config.mexOrigin + "/cms/api/v1/entity_types",
                verifyCerts: false,
                headers: {
                    Authorization: `Bearer ${config.jwt}`,
                    "X-MEx-Trace": config.loggingTraceSecret,
                },
                body: entityDef,
            });

            bar.increment();
        }
    } catch (e) {
        console.log(`Entity creation failed with error: ${e}`);
    }

    bar.stop();
}

/**
 * Adds field definitions to PostgreSQL.
 *
 * @param config
 * @param fieldDefs Array of field definitions
 */
export async function addFields(config: ApiConfig, fieldDefs: FieldDef[]) {
    const bar = new cliProgress.SingleBar({});
    bar.start(fieldDefs.length, 0);

    try {
        for (const fieldDef of fieldDefs) {
            await request({
                method: "POST",
                uri: config.mexOrigin + "/cms/api/v1/field_defs",
                verifyCerts: false,
                headers: {
                    Authorization: `Bearer ${config.jwt}`,
                    "X-MEx-Trace": config.loggingTraceSecret,
                },
                body: fieldDef,
            });

            bar.increment();
        }
    } catch (e) {
        console.log(`Field creation failed with error: ${e}`);
    }
    bar.stop();
}

/**
 * Adds search configs to PostgreSQL.
 *
 * @param config
 * @param searchConfigs Array of search configs
 */
export async function addSearchConfigs(config: ApiConfig, searchConfigs: SearchConfig[]) {
    const bar = new cliProgress.SingleBar({});
    bar.start(searchConfigs.length, 0);

    try {
        for (const sConfig of searchConfigs) {
            const response = await request({
                method: "POST",
                uri: `${config.mexOrigin}/cms/api/v1/search_configs`,
                verifyCerts: false,
                headers: {
                    Authorization: `Bearer ${config.jwt}`,
                    "X-MEx-Trace": config.loggingTraceSecret,
                },
                body: sConfig,
            });

            if (response.statusCode !== 201) {
                throw new Error(`Could not create search config - status code ${response.statusCode} returned`);
            }

            bar.increment();
        }
    } catch (e) {
        console.log(`Search config creation failed with error: ${e}`);
    }
    bar.stop();
}

/**
 * Adds items to PostgreSQL
 *
 * @param config API configuration
 * @param items Array of items to add
 * @param cfg Test run options
 * @param bulkSize No. of items uploaded per bulk request (default: -1 = bulk upload not used)
 * @param useSoftDuplicates If true, more precise duplicate detection will be used (default: false)
 */
export async function addItems(config: ApiConfig, items: Item_A[], cfg: Config, bulkSize: number = -1, useSoftDuplicates: boolean = false) {
    const bar = new cliProgress.SingleBar({});
    bar.start(items.length, 0);

    if (bulkSize === -1) {
        try {
            for (const item of items) {
                await request({
                    method: "POST",
                    uri: config.mexOrigin + "/api/v0/metadata/items",
                    verifyCerts: false,
                    headers: {
                        Authorization: `Bearer ${config.jwt}`,
                        "X-MEx-Trace": config.loggingTraceSecret,
                    },
                    body: { item },
                });

                bar.increment();
            }
        } catch (e) {
            console.log(`Item creation failed with error: ${e}`);
        }
    } else {
        try {
            const batches = pipe(items, A.chunksOf(bulkSize));
            let lastTokenUpdateTime = Date.now();
            let uploadTimeInSec: number;
            let batchStartTime: number;
            const dupAlg = useSoftDuplicates ? "LATEST_ONLY" : "SIMPLE";
            for (const batch of batches) {
                batchStartTime = Date.now();
                const response = await request({
                    method: "POST",
                    uri: config.mexOrigin + "/api/v0/metadata/items_bulk",
                    headers: {
                        Authorization: `Bearer ${config.jwt}`,
                        "X-MEx-Trace": config.loggingTraceSecret,
                    },
                    body: {
                        items: batch,
                        duplicateAlgorithm: dupAlg,
                    },
                });

                if (response.statusCode >= 400) {
                    console.error(`Error uploading items: ${response.body.str}`);
                    continue;
                }

                const responseJson = response.body.obj;
                const jobError = await waitForJobEnd(config.mexOrigin, config.jwt, responseJson.jobId);
                if (jobError !== "") {
                    throw new Error(jobError);
                }

                const timeNow = Date.now();
                uploadTimeInSec = (timeNow - batchStartTime) / 1000.0;

                const jobItemsResponse = await request({
                    method: "GET",
                    uri: `${config.mexOrigin}/api/v0/jobs/${responseJson.jobId}/items`,
                    headers: {
                        Authorization: `Bearer ${config.jwt}`,
                        "X-MEx-Trace": config.loggingTraceSecret,
                    },
                });
                const jobItemsResponseJson = jobItemsResponse.body.obj;
                const noCreated = jobItemsResponseJson.itemIds.filter((id) => id !== "{}").length; // Ignore dummy value which is contained in the list for internal reasons

                console.log(
                    `${batch.length} items uploaded, resulting in ${noCreated} new stored items - upload time: ${uploadTimeInSec.toFixed(0)} s`,
                );
                bar.increment(batch.length);

                // Brute-force solution for preventing timeout during long uploads:
                // renew token if older than a specific interval
                if (timeNow - lastTokenUpdateTime > MS_FOR_30MIN) {
                    console.log("Refreshing token");
                    config.jwt = await getAccessToken(cfg.effectiveInstance());
                    lastTokenUpdateTime = timeNow;
                }
            }
        } catch (e) {
            console.log(`Item creation failed with error: ${e}`);
        }
    }
    bar.stop();
}

/**
 * Delete old Solr index and re-index items
 * @param config
 */
export async function updateSolrIndex(config: ApiConfig) {
    // delete current solr schema
    await request({
        method: "DELETE",
        uri: config.mexOrigin + "/api/v0/metadata/index",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${config.jwt}`,
            "X-MEx-Trace": config.loggingTraceSecret,
        },
    });

    // Hard wait do ensure index is really empty
    await wait(3000);

    // reindex solr
    const updateIndexJob = await request({
        method: "PUT",
        uri: config.mexOrigin + "/api/v0/metadata/index",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${config.jwt}`,
            "X-MEx-Trace": config.loggingTraceSecret,
        },
    });

    const jobError = await waitForJobEnd(config.mexOrigin, config.jwt, updateIndexJob.body.obj.jobId);
    if (jobError !== "") {
        throw new Error(jobError);
    }
}

/**
 * Rebuild Solr schema
 * @param config
 * */
export async function rebuildSolrSchema(config: ApiConfig) {
    // Delete current Solr index (to allow changing schema)
    await request({
        method: "DELETE",
        uri: config.mexOrigin + "/api/v0/metadata/index",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${config.jwt}`,
            "X-MEx-Trace": config.loggingTraceSecret,
        },
    });

    // Hard wait do ensure index is really empty
    await wait(3000);

    // Rebuild Solr schema
    const rebuildSchemaJob = await request({
        method: "POST",
        uri: config.mexOrigin + "/api/v0/metadata/index",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${config.jwt}`,
            "X-MEx-Trace": config.loggingTraceSecret,
        },
    });

    const jobError = await waitForJobEnd(config.mexOrigin, config.jwt, rebuildSchemaJob.body.obj.jobId);
    if (jobError !== "") {
        throw new Error(jobError);
    }
}

export async function deleteSolrIndex(config: ApiConfig) {
    // Delete current Solr index
    const deleteResponse = await request({
        method: "DELETE",
        uri: config.mexOrigin + "/api/v0/metadata/index",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${config.jwt}`,
            "X-MEx-Trace": config.loggingTraceSecret,
        },
    });

    if (deleteResponse.statusCode !== 204) {
        throw new Error("index deletion failed: " + deleteResponse.body.str);
    }
}
