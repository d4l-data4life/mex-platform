import * as E from "fp-ts/lib/Either";
import { Command } from "commander";
import notifier from "node-notifier";

import { constantPause } from "../../utils/attempt";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";
import { loadData } from "../../utils/data-loader";
import { Item_T } from "../../api/defs";
import { addItems, removeAllItems, updateSolrIndex } from "../../api/mex-api";

export default function (): Command {
    return new Command("load")
        .description("load full data set")
        .requiredOption("-i, --items <items...>", "item file; can be specified multiple times")
        .option("--bulk", "use bulk loading", false)
        .option("--soft-duplicates", "use duplication detection that allow return to previous states", false)
        .action(handler);
}

type Options = {
    items: string[];
    bulk: boolean;
    softDuplicates: boolean;
};

const LOAD_BATCH_SIZE = 500;

async function handler(options: Options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);

    out.log(`Loading items from the following files:\n${options.items.join("\n")}`);
    const items = loadData(options.items, Item_T);
    if (E.isLeft(items)) {
        out.error(items.left);
        process.exit(1);
    }

    out.log("MEx Core service origin:", cfg.mexOrigin());

    out.log(`Preparing to set up system with:
        - ${items.right.length} items`);
    // TODO - put back
    //- ${fieldDefs.right.length} custom fields

    if (cfg.dryRun()) {
        out.log("dry run: exiting here");
        process.exit(0);
    }

    const config = { mexOrigin: cfg.mexOrigin(), jwt, loggingTraceSecret: cfg.loggingTraceSecret() };
    try {
        // eslint-disable-next-line no-console
        out.log("--- Starting system set-up ---");

        // Remove old data and configs elements from MEx DB
        out.log("Removing existing items");
        await removeAllItems(config);

        out.log("Adding items");
        const d1 = new Date();
        await addItems(config, items.right, cfg, options.bulk ? LOAD_BATCH_SIZE : -1, options.softDuplicates);
        const d2 = new Date();
        out.log(`----\nUpload took: ${((d2.getTime() - d1.getTime()) / 1000.0).toFixed(0)} s\n----`);

        out.log("Waiting to let item aggregation finish");
        await constantPause(5 * 1000)();

        // Update Solr index
        out.log("Updating search index");
        await updateSolrIndex(config);

        out.log("--- System set-up complete ---");
    } catch (e) {
        out.log(`Failed with error ${e}`);
    }

    notifier.notify({
        title: "MEx data load",
        message: `${items.right.length} items loaded`,
    });
}
