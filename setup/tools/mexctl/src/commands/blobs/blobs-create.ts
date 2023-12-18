import * as fs from "fs";
import * as z from "zlib";
import { Command } from "commander";

import { request } from "../../utils/http";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

export default function (): Command {
    return new Command("create")
        .description("create/overwrite a blob")
        .requiredOption("-n, --name <name>", "blob name")
        .requiredOption("-t, --type <type>", "blob type")
        .requiredOption("-f, --file <file>", "file to store as blob")
        .option("-z, --zip", "compress content via zlib deflate")
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

    let data = fs.readFileSync(options.file);
    out.log("file size:", data.length);
    if (options.zip) {
        out.log("deflating");
        data = z.deflateSync(data);
    }
    out.log("data size:", data.length);

    const chunkSize = 1024 * 1024;
    let offset = 0;
    while (offset < data.length) {
        out.log("sending chunk", offset);
        const createBlobResponse = await request({
            method: "POST",
            uri: cfg.effectiveInstance().mexOrigin + "/api/v0/blobs",
            headers: {
                Authorization: `Bearer ${jwt}`,
                "X-MEx-Trace": cfg.loggingTraceSecret(),
            },
            body: {
                blobName: options.name,
                blobType: options.type,
                append: true,
                data: data.subarray(offset, offset + chunkSize).toString("base64"),
            },
        });

        if (createBlobResponse.statusCode !== 201) {
            out.error(createBlobResponse.body.str);
            process.exit(1);
        }

        offset += chunkSize;
    }

    out.log("done");
}
