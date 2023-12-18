export class AttemptCircuitBreaker extends Error {
    public constructor(public cause: Error) {
        super(cause.message);
    }
}

export class AttemptError extends Error {}

export type AttemptMaker<T> = () => Promise<T>;

export type PauseStrategy = () => Promise<number>;

export interface Console {
    error(...data: any[]): void;
    log(...data: any[]): void;
}

export interface AttemptOptions {
    description: string;
    maxAttempts?: number;
    pauseStrategy?: PauseStrategy;
    console?: Console;
}

export const DEFAULT_PAUSE_MILLIS = 2000;
export const DEFAULT_NUMBER_OF_ATTEMPTS = 10;

export function constantPause(millis: number = DEFAULT_PAUSE_MILLIS): PauseStrategy {
    return function (): Promise<number> {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        return new Promise(function (resolve, _reject) {
            setTimeout(() => resolve(millis), millis);
        });
    };
}

export function exponentialPause(startMillis = 10): PauseStrategy {
    let duration = startMillis / 2;
    return function (): Promise<number> {
        duration *= 2;
        return constantPause(duration)();
    };
}

export const wait = (millis: number) => constantPause(millis)();

export async function attempt<T>(maker: AttemptMaker<T>, options?: AttemptOptions): Promise<T> {
    if (typeof options === "undefined") {
        options = { description: "generic attempt", console };
    }

    options.maxAttempts = options.maxAttempts || DEFAULT_NUMBER_OF_ATTEMPTS;
    options.console = options.console || console;

    const pauseStrategy = typeof options.pauseStrategy !== "undefined" ? options.pauseStrategy : constantPause();

    let i = 0;
    while (i < options.maxAttempts) {
        try {
            const result = await maker();
            options.console.log(`${options.description}: attempt ${i + 1}/${options.maxAttempts}: success`);
            return result;
        } catch (error) {
            options.console.log(`${options.description}: attempt ${i + 1}/${options.maxAttempts} pending: ${error.message}`);
            if (error instanceof AttemptCircuitBreaker) {
                options.console.error(`${options.description}: stopping already after ${options.maxAttempts} attempts: ${error.message}`);
                throw error;
            }
            await pauseStrategy();
        }

        i += 1;
    }

    throw new AttemptError(`all ${options.maxAttempts} ${options.description} attempts failed`);
}
