import { Command } from "commander";

import getEntitiesListCommand from "./entities/entities-list";
import getEntitiesAppendCommand from "./entities/entities-append";

export default function (): Command {
    return new Command("entities").description("entity commands").addCommand(getEntitiesListCommand()).addCommand(getEntitiesAppendCommand());
}
