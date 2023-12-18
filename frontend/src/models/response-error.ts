export class ResponseError extends Error {
  readonly status?: number;

  constructor(responseBody: string, status?: number) {
    super();
    this.message = responseBody;
    this.status = status;
  }
}
