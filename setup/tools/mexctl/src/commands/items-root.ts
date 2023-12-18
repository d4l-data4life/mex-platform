import { Command } from "commander";

import getItemsDeleteCommand from "./items/items-delete";

export default function (): Command {
    return new Command("items")
        .description("items commands")
        .addCommand(getItemsDeleteCommand());
}
