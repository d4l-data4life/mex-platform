import { Command } from "commander";

export interface Output {
    log(...data: any[]): void;
    dir(...data: any[]): void;
    json(data: any): void;
    error(...data: any[]): void;
    verbose(...data: any[]): void;
}

export class Output {
    private verboseEnabled: boolean = false;

    constructor(command: Command) {
        this.verboseEnabled = command.optsWithGlobals().verbose == true;
    }

    public log(...data: any[]): void {
        console.log(...data);
    }

    public dir(data: any): void {
        console.dir(data, { depth: 5 });
    }

    public error(...data: any[]): void {
        console.error(...data);
    }

    public verbose(...data: any[]): void {
        if (this.verboseEnabled) {
            console.error(...data);
        }
    }

    public json(data: any): void {
        console.log(JSON.stringify(data));
    }
}
