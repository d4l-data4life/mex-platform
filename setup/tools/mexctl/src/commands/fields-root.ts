import { Command } from "commander";

import getFieldsListCommand from "./fields/fields-list";

export default function (): Command {
    return new Command("fields").description("field commands").addCommand(getFieldsListCommand());
}
