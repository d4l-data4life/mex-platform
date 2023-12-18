import { attempt } from "./attempt";
import { Instance_A } from "./config";
import { request } from "./http";

export async function getAccessToken(instance: Instance_A): Promise<string> {
    async function tokenWaiter(): Promise<string> {
        if (instance.tokenContentType === "application/json") {
            const requestBody = {
                client_id: instance.clientId,
                scope: [`api://${instance.clientId}/.default`].join(" "),
                grant_type: "client_credentials",
                client_secret: instance.clientSecret,
            };

            const response = await request({
                method: "POST",
                uri: instance.tokenUrl,
                verifyCerts: false,
                headers: {
                    "content-type": "application/json",
                    accept: "application/json",
                },
                body: requestBody,
            });

            if (response.statusCode !== 200) {
                console.error(response.body);
                throw new Error("token endpoint not successful: " + response.body);
            }
            return response.body.obj.access_token;
        } else if (instance.tokenContentType === "application/x-www-form-urlencoded") {
            const requestBody = [
                { p: "client_id", v: instance.clientId },
                { p: "scope", v: [`api://${instance.clientId}/.default`].join(" ") },
                { p: "grant_type", v: "client_credentials" },
                { p: "client_secret", v: instance.clientSecret },
            ]
                .map((q) => q.p + "=" + encodeURIComponent(q.v))
                .join("&");

            const response = await request({
                method: "POST",
                uri: instance.tokenUrl,
                verifyCerts: false,
                headers: {
                    "content-type": "application/x-www-form-urlencoded",
                    accept: "application/json",
                    "content-length": "" + requestBody.length,
                },
                body: requestBody,
            });

            if (response.statusCode !== 200) {
                throw new Error("token endpoint not successful: " + response.body);
            }
            return response.body.obj.access_token;
        } else {
            throw new Error("unsupported content type: " + instance.tokenContentType);
        }
    }

    const silencer = {
        error: (...data: any[]) => {
            console.error(...data);
        },
        log: (...data: any[]) => {
            // console.log(...data);
        },
    };

    return attempt(tokenWaiter, {
        description: "acquiring JWT",
        console: silencer,
    });
}
