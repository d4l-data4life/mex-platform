import * as os from "os";
import * as path from "path";

import { Command, Option } from "commander";

import getAuthRootCommand from "./commands/auth-root";
import getBlobsRootCommand from "./commands/blobs-root";
import getIndexRootCommand from "./commands/index-root";
import getDataRootCommand from "./commands/data-root";
import getEntitiesRootCommand from "./commands/entities-root";
import getFieldsRootCommand from "./commands/fields-root";
import getConfigRootCommand from "./commands/config-root";
import getSearchConfigRootCommand from "./commands/searchconfig-root";
import getTreeRootCommand from "./commands/tree-root";
import getItemsRootCommand from "./commands/items-root";
import getCroissantRootCommand from "./commands/croissant-root";

const program = new Command();

program
    .name("mexctl")
    .description("MEx command line tool")
    .version(process.env["npm_package_version"] ?? "unknown")
    .option("--dry", "dry run")
    .option("-v, --verbose", "verbose output")
    .addOption(new Option(" --format <output>", "output format").choices(["text", "json", "yaml"]).default("text"))
    .option("--config <file>", "config file content or @filename", "@" + path.resolve(os.homedir(), ".mexctl/config.yaml"))
    .option("--instance <instance>", "MEx instance key");

program.addCommand(getAuthRootCommand());
program.addCommand(getBlobsRootCommand());
program.addCommand(getConfigRootCommand());
program.addCommand(getIndexRootCommand());
program.addCommand(getDataRootCommand());
program.addCommand(getEntitiesRootCommand());
program.addCommand(getFieldsRootCommand());
program.addCommand(getSearchConfigRootCommand());
program.addCommand(getTreeRootCommand());
program.addCommand(getItemsRootCommand());
program.addCommand(getCroissantRootCommand());

(async function () {
    await program.parseAsync(process.argv);
})();
