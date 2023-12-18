import { Command } from "commander";

import { request } from "../../utils/http";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

export default function (): Command {
    return new Command("list").description("list all blobs").action(handler);
}

async function handler(options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);

    if (cfg.dryRun()) {
        out.log("dry run: abort here");
        process.exit(0);
    }

    const listBlobsResponse = await request({
        method: "GET",
        uri: cfg.effectiveInstance().mexOrigin + "/api/v0/blobs",
        headers: {
            Authorization: `Bearer ${jwt}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
    });

    if (listBlobsResponse.statusCode !== 200) {
        out.error(listBlobsResponse.body.str);
        process.exit(1);
    }

    out.log(listBlobsResponse.body.obj.blobInfos);
}
