import * as E from "fp-ts/lib/Either";
import { Command } from "commander";
import notifier from "node-notifier";

import { constantPause } from "../../utils/attempt";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";
import { loadData } from "../../utils/data-loader";
import { Item_T } from "../../api/defs";
import { addItems, updateSolrIndex } from "../../api/mex-api";

export default function (): Command {
    return new Command("append")
        .description("add new data and re-index (existing data and configuration is left unchanged)")
        .requiredOption("-i, --items <items...>", "item file; can be specified multiple times")
        .option("--bulk", "use bulk loading", false)
        .option("--soft-duplicates", "use duplication detection that allow return to previous states", true)
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
    out.log(`Preparing to add ${items.right.length} items to system`);

    if (cfg.dryRun()) {
        out.log("dry run: exiting here");
        process.exit(0);
    }

    const config = { mexOrigin: cfg.mexOrigin(), jwt, loggingTraceSecret: cfg.loggingTraceSecret() };
    try {
        // eslint-disable-next-line no-console
        out.log("--- Starting addition of new data ---");

        out.log("Adding new items");
        const d1 = new Date();
        await addItems(config, items.right, cfg, options.bulk ? LOAD_BATCH_SIZE : -1, options.softDuplicates);
        const d2 = new Date();
        out.log(`----\nUpload took: ${((d2.getTime() - d1.getTime()) / 1000.0).toFixed(0)} s\n----`);

        out.log("Waiting to let item aggregation finish");
        await constantPause(5 * 1000)();

        // Update Solr index
        out.log("Updating search index");
        await updateSolrIndex(config);

        out.log("--- Addition of new data complete ---");
    } catch (e) {
        out.log(`Failed with error ${e}`);
    }

    notifier.notify({
        title: "MEx data appended",
        message: `${items.right.length} new items loaded`,
    });
}
