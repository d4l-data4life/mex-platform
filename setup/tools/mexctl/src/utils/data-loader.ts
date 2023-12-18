import * as fs from "fs";
import * as yaml from "js-yaml";

import * as t from "io-ts";
import { PathReporter } from "io-ts/PathReporter";
import * as E from "fp-ts/Either";

export function loadData<A, O, I>(files: string[], typeClass: t.Type<A, O, I>): E.Either<Error, A[]> {
    const ret: A[] = [];
    for (let file of files) {
        let obj: unknown;
        try {
            const buf = fs.readFileSync(file);
            obj = yaml.load(buf.toString());
        } catch (err) {
            return E.left(err);
        }

        const val = t.array(typeClass).decode(obj);
        if (E.isLeft(val)) {
            return E.left(new Error(file + ": " + PathReporter.report(val).join("\n")));
        }

        ret.push(...val.right);
    }

    return E.right(ret);
}
