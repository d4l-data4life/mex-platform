import { Command } from "commander";

import showTreeCommand from "./tree/tree-show";

export default function (): Command {
    return new Command("tree").description("tree commands").addCommand(showTreeCommand());
}
