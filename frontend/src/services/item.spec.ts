jest.mock('stencil-router-v2');

import { API_URL } from 'config';
import * as fetchClient from 'utils/fetch-client';
import itemService, { Item, Version, VersionsResponse } from './item';

const requestSpy = jest.spyOn(fetchClient, 'request');
requestSpy.mockImplementation(async () => [null, new Headers()] as [any, Headers]);

const EXAMPLE_ITEM: Item = {
  itemId: 'foo-id',
  businessId: 'foo-business-id',
  entityType: 'Resource',
  owner: 'system',
  createdAt: '2022-04-13T12:17:16.676Z',
  values: [
    {
      itemValueId: 'cf15cbcc-a327-29eb-a32f',
      fieldName: 'author',
      fieldValue: 'John Doe',
      language: 'en',
      place: 0,
      version: 2122,
    },
    {
      itemValueId: 'e2d7aa35-629d-2eea-3c99',
      fieldName: 'title',
      fieldValue: 'Archery and its contributions to health',
      language: 'en',
      place: 1,
      version: 3419,
    },
  ],
};

const EXAMPLE_VERSIONS_RESPONSE: VersionsResponse = {
  versions: [
    {
      createdAt: '2022-02-17T10:20:07.983Z',
      itemId: 'item-bar',
      versionDesc: 'v2',
    },
    {
      createdAt: '2021-11-09T15:33:54.166Z',
      itemId: 'item-foo',
      versionDesc: 'v1',
    },
  ],
};

const EXAMPLE_VERSIONS: Version[] = [
  {
    createdAt: new Date('2022-02-17T10:20:07.983Z'),
    itemId: 'item-bar',
    versionDesc: 'v2',
  },
  {
    createdAt: new Date('2021-11-09T15:33:54.166Z'),
    itemId: 'item-foo',
    versionDesc: 'v1',
  },
];

describe('item service', () => {
  describe('fetch()', () => {
    it('fetches metadata item data by item ID', async () => {
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_ITEM, new Headers()]);

      expect(await itemService.fetch('foo-id')).toBe(EXAMPLE_ITEM);
      expect(requestSpy).toHaveBeenCalledWith({
        authorized: true,
        method: 'GET',
        url: `${API_URL}/metadata/items/foo-id`,
      });
    });

    it('throws an error when the request fails', async () => {
      const consoleSpy = jest.spyOn(console, 'error');
      consoleSpy.mockImplementationOnce(() => {});
      requestSpy.mockImplementationOnce(() => {
        throw new Error('Network error');
      });

      await expect(itemService.fetch('foo-id')).rejects.toThrow();
    });
  });

  describe('fetchVersions()', () => {
    beforeEach(() => {
      jest.clearAllMocks();
    });

    it('fetches versions by item ID and converts createdAt to Date', async () => {
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_VERSIONS_RESPONSE, new Headers()]);

      expect(await itemService.fetchVersions('item-foo')).toEqual(EXAMPLE_VERSIONS);
      expect(requestSpy).toHaveBeenCalledWith({
        authorized: true,
        method: 'POST',
        url: `${API_URL}/metadata/items/item-foo/versions`,
      });
    });

    it('throws an error when the request fails', async () => {
      const consoleSpy = jest.spyOn(console, 'error');
      consoleSpy.mockImplementationOnce(() => {});
      requestSpy.mockImplementationOnce(() => {
        throw new Error('Network error');
      });

      await expect(itemService.fetchVersions('item-foo')).rejects.toThrow();
    });
  });

  describe('fetchLatestItemIdByIdentifier()', () => {
    beforeEach(() => {
      jest.clearAllMocks();
    });

    it('fetches the versions of an item by its business identifier and returns the item ID of the latest version', async () => {
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_VERSIONS_RESPONSE, new Headers()]);

      expect(await itemService.fetchLatestItemIdByIdentifier('identifier-foo')).toEqual('item-bar');
      expect(requestSpy).toHaveBeenCalledWith({
        authorized: true,
        method: 'GET',
        url: `${API_URL}/metadata/versions/identifier-foo`,
      });
    });

    it('throws an error when the request fails', async () => {
      const consoleSpy = jest.spyOn(console, 'error');
      consoleSpy.mockImplementationOnce(() => {});
      requestSpy.mockImplementationOnce(() => {
        throw new Error('Network error');
      });

      await expect(itemService.fetchLatestItemIdByIdentifier('identifier-foo')).rejects.toThrow();
    });
  });

  describe('resolveLink()', () => {
    beforeEach(() => {
      jest.clearAllMocks();
    });

    it('calls fetchLatestItemIdByIdentifier() with the identifier and then fetch() with the item ID', async () => {
      const itemServiceFetchLatestItemIdByIdentifierSpy = jest
        .spyOn(itemService, 'fetchLatestItemIdByIdentifier')
        .mockImplementation(async () => 'bar-item-id');
      const itemServiceFetchSpy = jest.spyOn(itemService, 'fetch').mockImplementation(async () => EXAMPLE_ITEM);

      expect(await itemService.resolveLink('foo-identifier')).toBe(EXAMPLE_ITEM);
      expect(itemServiceFetchLatestItemIdByIdentifierSpy).toHaveBeenCalledWith('foo-identifier');
      expect(itemServiceFetchSpy).toHaveBeenCalledWith('bar-item-id');
    });

    it('does not throw an error when one of the requests fails but returns null instead', async () => {
      jest.spyOn(itemService, 'fetchLatestItemIdByIdentifier').mockImplementation(async () => {
        throw new Error('some error');
      });

      expect(await itemService.resolveLink('foo-identifier')).toBe(null);
    });
  });
});
