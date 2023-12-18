import { attempt, constantPause, Console } from "./attempt";
import { request } from "./http";

/**
 *
 * @param origin Base URL
 * @param jwt JWT bearer token
 * @param jobId ID to assign to async job
 * @param con Console for output
 */
export async function waitForJobEnd(origin: string, jwt: string, jobId: string, con: Console = console): Promise<string> {
    async function jobStatus(): Promise<string> {
        const response = await request({
            method: "GET",
            headers: {
                Authorization: `Bearer ${jwt}`,
            },
            uri: `${origin}/api/v0/jobs/${jobId}`,
        });

        if (response.statusCode !== 200) {
            throw new Error("error querying job status: " + jobId);
        }

        const body = response.body.obj;
        if (body.status !== "DONE") {
            throw new Error("job still running: " + jobId);
        }

        return body.error;
    }

    return attempt(jobStatus, {
        description: "job waiter: " + jobId,
        maxAttempts: 250,
        pauseStrategy: constantPause(2000),
        console: con,
    });
}
