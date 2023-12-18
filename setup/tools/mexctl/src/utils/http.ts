import * as url from "url";
import * as http from "http";
import * as https from "https";

export interface Body {
    buf: Buffer;
    obj: any;
    str: string;
}

export interface Response {
    statusCode: number;
    headers: http.IncomingHttpHeaders;
    body: Body;
    duration: number;
}

export interface ReqOptions {
    socket?: string;
    method?: string;
    uri: string;
    headers?: http.OutgoingHttpHeaders;
    body?: unknown;

    verifyCerts?: boolean;
    expectJson?: boolean;
}

export async function request(opts: ReqOptions): Promise<Response> {
    // eslint-disable-next-line complexity
    return new Promise<Response>(function (resolve, reject) {
        if (typeof opts.method === "undefined") {
            opts.method = "GET";
        }

        let modu: any;
        let socketPath = "";
        let u: url.UrlWithStringQuery;
        try {
            u = url.parse(opts.uri);
        } catch (error) {
            reject(error);
            return;
        }

        if (typeof opts.socket === "string") {
            socketPath = opts.socket;
            modu = http;
        }

        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        if (u.protocol === "http:") {
            modu = http;
        } else if (u.protocol === "https:") {
            modu = https;
        } else {
            if (socketPath === "") {
                reject(new Error(`unsupported protocol: ${u.protocol}`));
                return;
            }
        }

        const startTime = process.hrtime();

        const req = modu.request(
            {
                method: opts.method,
                ...(socketPath !== ""
                    ? { socketPath: socketPath }
                    : { protocol: u.protocol, hostname: u.hostname, port: u.port }),
                path: u.path,
                headers: opts.headers,
                rejectUnauthorized: !!opts.verifyCerts,
            },
            (resp) => {
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                const chunks: any[] = [];
                resp.on("data", (chunk) => {
                    chunks.push(chunk);
                });

                resp.on("end", () => {
                    const buf = Buffer.concat(chunks);
                    const str = buf.toString();

                    let obj: any;
                    if (resp.headers["content-type"] === "application/json" || opts.expectJson) {
                        try {
                            if (str === "") {
                                obj = null;
                            } else {
                                obj = JSON.parse(str);
                            }
                        } catch (err) {
                            reject(new Error("expected JSON response, but got parse error: " + err));
                            return;
                        }
                    }

                    resolve({
                        statusCode: resp.statusCode,
                        headers: resp.headers,
                        body: {
                            buf,
                            obj,
                            str,
                        },
                        duration: millis(process.hrtime(startTime)),
                    });
                });
            },
        );

        req.on("error", reject);

        if (typeof opts.body !== "undefined") {
            if (opts.body instanceof Buffer) {
                req.write(opts.body);
            } else if (typeof opts.body === "string") {
                req.write(opts.body);
            } else if (opts.body instanceof Object) {
                req.write(JSON.stringify(opts.body));
            } else {
                req.write(opts.body);
            }
        }

        req.end();
    });
}

export function cookieMap(cookies: string[]): Record<string, string> {
    const map: Record<string, string> = {};
    for (const c of cookies) {
        const eqIdx = c.indexOf("=");
        map[c.substring(0, eqIdx)] = c;
    }
    return map;
}

export function millis(dur: [number, number]): number {
    return Math.round((dur[0] * 1e9 + dur[1]) / 1e6);
}

export function ms(start: [number, number]): number {
    return millis(process.hrtime(start));
}
