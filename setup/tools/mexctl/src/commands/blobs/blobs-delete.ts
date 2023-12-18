import { Command } from "commander";

import { request } from "../../utils/http";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

export default function (): Command {
    return new Command("delete")
        .description("delete a blob")
        .requiredOption("-n, --name <name>", "blob name")
        .requiredOption("-t, --type <type>", "blob type")
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

    const createBlobResponse = await request({
        method: "DELETE",
        uri: `${cfg.effectiveInstance().mexOrigin}/api/v0/blobs/${options.name}?blobType=${options.type}`,
        headers: {
            Authorization: `Bearer ${jwt}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
    });

    if (createBlobResponse.statusCode !== 204) {
        out.error(createBlobResponse.body.str);
        process.exit(1);
    }

    out.log(createBlobResponse.body.str);
}
