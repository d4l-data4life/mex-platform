import { Command } from "commander";

import getBlobsCreateCommand from "./blobs/blobs-create";
import getBlobsListCommand from "./blobs/blobs-list";
import getBlobsDeleteCommand from "./blobs/blobs-delete";
import getBlobsTestCommand from "./blobs/blobs-test";

export default function (): Command {
    return new Command("blobs")
        .description("blobs commands")
        .addCommand(getBlobsCreateCommand())
        .addCommand(getBlobsListCommand())
        .addCommand(getBlobsDeleteCommand())
        .addCommand(getBlobsTestCommand());
}
