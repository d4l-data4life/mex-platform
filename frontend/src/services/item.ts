import { API_URL } from 'config';
import { EntityTypeName } from 'config/entity-types';
import { get, post } from 'utils/fetch-client';

export interface ItemValue {
  fieldName: string;
  fieldValue: string;
  itemValueId: string;
  language: string;
  place: number;
  version: number;
}

export interface Item {
  itemId: string;
  entityType: EntityTypeName;
  owner: string;
  createdAt: string;
  businessId: string;
  values: ItemValue[];
}

export interface Version {
  itemId: string;
  createdAt: Date;
  versionDesc?: string;
}

export interface VersionsResponse {
  versions: {
    itemId: string;
    createdAt: string;
    versionDesc: string;
  }[];
}

export class ItemService {
  async fetch(itemId: string): Promise<Item> {
    const [response] = await get<Item>({
      url: `${API_URL}/metadata/items/${itemId}`,
      authorized: true,
    });

    return response;
  }

  protected sortVersions(response: VersionsResponse): Version[] {
    return (response?.versions ?? [])
      .map(({ createdAt, itemId, versionDesc }) => ({
        itemId,
        createdAt: new Date(createdAt),
        versionDesc,
      }))
      .sort((a, b) => (a.createdAt > b.createdAt ? -1 : 1));
  }

  async fetchVersions(itemId: string): Promise<Version[]> {
    const [response] = await post<VersionsResponse>({
      url: `${API_URL}/metadata/items/${itemId}/versions`,
      authorized: true,
    });

    return this.sortVersions(response);
  }

  async fetchLatestItemIdByIdentifier(identifier: string): Promise<string> {
    const [response] = await get<VersionsResponse>({
      url: `${API_URL}/metadata/versions/${encodeURIComponent(identifier)}`,
      authorized: true,
    });

    return this.sortVersions(response)?.[0]?.itemId;
  }

  async resolveLink(identifier: string): Promise<Item> {
    try {
      const itemId = await this.fetchLatestItemIdByIdentifier(identifier);
      return itemId ? await this.fetch(itemId) : null;
    } catch {
      return null;
    }
  }
}

export default new ItemService();
