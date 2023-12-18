export interface PersistableData {
  [key: string]: string;
}

export class Persistor {
  private cachedData: PersistableData = {};
  private storage: Storage;

  constructor(storage: Storage) {
    this.storage = storage;
  }

  get(key: string): string | undefined {
    if (key in this.cachedData) {
      return this.cachedData[key];
    }

    return this.getStoreValue(key);
  }

  set(key: string, value?: string): void {
    if (value === null || value === undefined || value === 'null') {
      delete this.cachedData[key];
      this.removeStoreValue(key);
    } else {
      this.cachedData[key] = value;
      this.setStoreValue(key, value);
    }
  }

  getKeys(): string[] {
    try {
      return new Array(this.storage.length)
        .fill(undefined)
        .map((_, index) => this.storage.key(index))
        .filter(Boolean);
    } catch (_) {
      return Object.keys(this.cachedData);
    }
  }

  private getStoreValue(key: string): string | undefined {
    try {
      const value = this.storage?.getItem(key);
      if (typeof value === 'string') {
        this.cachedData[key] = value;
      }

      return value;
    } catch (_) {
      return undefined;
    }
  }

  private setStoreValue(key: string, value: string) {
    try {
      this.storage?.setItem(key, value);
    } catch (_) {}
  }

  private removeStoreValue(key: string) {
    try {
      this.storage?.removeItem(key);
    } catch (_) {}
  }
}
