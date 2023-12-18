import { Command } from "commander";

import { PathReporter } from "io-ts/PathReporter";
import * as t from "io-ts";
import * as fs from "fs";
import * as yaml from "js-yaml";
import * as E from "fp-ts/lib/Either";

export const Instance_T = t.type({
    description: t.string,
    mexOrigin: t.string,
    clientId: t.string,
    clientSecret: t.string,
    tokenUrl: t.string,
    tokenContentType: t.string,
    loggingTraceSecret: t.string,
    configApiKey: t.string,
});

export type Instance_A = t.TypeOf<typeof Instance_T>;

export const Config_T = t.type({
    instances: t.record(t.string, Instance_T),
    defaultInstance: t.string,
});

export type Config_A = t.TypeOf<typeof Config_T>;

export class Config {
    private _config: Config_A;
    private _instance: Instance_A;

    private _dryRun: boolean;

    constructor(command: Command) {
        const optionValues = command.optsWithGlobals();

        let configContent = "";
        if (optionValues.config.startsWith("@")) {
            // load file
            const fileName = optionValues.config.substring(1);
            const stat = fs.statSync(fileName);
            if (!stat.isFile()) {
                throw new Error(`${fileName} is not a file`);
            }
            configContent = fs.readFileSync(fileName).toString();
        } else {
            // interpret string as config itself
            configContent = optionValues.config;
        }

        const config = Config_T.decode(yaml.load(configContent));
        if (E.isLeft(config)) {
            throw new Error(PathReporter.report(config).join("\n"));
        }

        this._config = config.right;

        const instanceName = optionValues.instance ?? this._config.defaultInstance;
        this._instance = this._config.instances[instanceName];
        if (typeof this._instance === "undefined") {
            throw new Error("unknown instance: " + instanceName);
        }

        this._dryRun = optionValues.dry === true;
    }

    public effectiveInstance(): Instance_A {
        return this._instance;
    }

    public dryRun(): boolean {
        return this._dryRun;
    }

    public mexOrigin(): string {
        return this._instance.mexOrigin;
    }

    public loggingTraceSecret(): string {
        return this._instance.loggingTraceSecret;
    }

    public allInstances(): Record<string, Instance_A> {
        return JSON.parse(JSON.stringify(this._config.instances));
    }
}
