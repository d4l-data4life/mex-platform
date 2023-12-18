import { Command, Option } from "commander";

import { request } from "../../utils/http";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

export default function (): Command {
    return new Command("test")
        .description("test mesh resolution")
        .requiredOption("-n, --name <name>", "blob name")
        .requiredOption("-t, --type <type>", "blob type")
        .option("--gc", "run GC")
        .addOption(new Option("-l <mode>", "loading mode (0=in memory, 1=temp file)").choices(["0", "1"]).default("0"))
        .addOption(new Option("--bag <size>", "bag size").preset(5).argParser(parseInt))
        .addOption(new Option("--iter <number>", "iterations").preset(5).argParser(parseInt))
        .action(handler);
}

async function handler(options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);

    out.verbose(options);

    if (cfg.dryRun()) {
        out.log("dry run: abort here");
        process.exit(0);
    }

    const testBlobResponse = await request({
        method: "POST",
        uri: cfg.effectiveInstance().mexOrigin + "/api/v0/blobs/mesh",
        headers: {
            Authorization: `Bearer ${jwt}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
        body: {
            blobName: options.name,
            blobType: options.type,
            bagSize: options.bag,
            iterations: options.iter,
            showTerms: false,
            run_gc: options.gc,
            loading_mode: parseInt(options.l),
        },
    });

    if (testBlobResponse.statusCode !== 200) {
        out.error(testBlobResponse.body.str);
        process.exit(1);
    }

    out.json(testBlobResponse.body.obj);
}
