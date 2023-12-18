import { Command } from "commander";
import * as tar from "tar-fs";
import * as uuid from "uuid";
import { Readable, Stream } from "stream";

import { request } from "../../utils/http";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";
import { getAccessToken } from "../../utils/jwt";
import { waitForJobEnd } from "../../utils/job-waiter";

export default function (): Command {
    return new Command("send")
        .description("send a config from file system")
        .requiredOption("-d, --dir <folder name>", "directory to TAR and send")
        .action(handler);
}

async function handler(options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(options);
    out.verbose(jwt);

    const tr = await streamToBuffer(tar.pack(options.dir));
    out.log("bytes: " + tr.length);
    const newConfigHash = uuid.v4();

    const updateConfigResponse = await request({
        method: "POST",
        uri: cfg.mexOrigin() + "/api/v0/config/update",
        headers: {
            Authorization: `apikey ${cfg.effectiveInstance().configApiKey}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
        body: {
            cannedConfig: {
                tarData: tr.toString("base64"),
                configHash: newConfigHash,
            },
        },
    });

    if (updateConfigResponse.statusCode !== 200) {
        out.error(updateConfigResponse.body.str);
        process.exit(1);
    }

    const jobError = await waitForJobEnd(cfg.mexOrigin(), jwt, updateConfigResponse.body.obj.jobId);
    if (jobError !== "") {
        throw new Error(jobError);
    }
}

export async function streamToBuffer(s: Stream): Promise<Buffer> {
    return new Promise<Buffer>((resolve, reject) => {
        const buf: any[] = [];
        s.on("data", (data) => buf.push(data));
        s.on("end", () => resolve(Buffer.concat(buf)));
        s.on("error", (err) => reject(err));
    });
}
