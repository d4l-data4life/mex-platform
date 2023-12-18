jest.mock('stencil-router-v2');

import searchService from './search';
import * as fetchClient from 'utils/fetch-client';
import * as config from 'config';
import { SearchFacet } from 'config/search';
import stores from 'stores';
import { Field, FieldImportance, FieldRenderer } from 'models/field';

const { API_URL, SEARCH_CONFIG } = config;
const requestSpy = jest.spyOn(fetchClient, 'request');
requestSpy.mockImplementation(async () => [EXAMPLE_SEARCH_RESPONSE, new Headers()] as [any, Headers]);

const EXAMPLE_SEARCH_RESPONSE = {}; // stub

const expectInBody = (str: string, not = false) => {
  expect(requestSpy).toHaveBeenCalledWith(
    expect.objectContaining({
      authorized: true,
      method: 'POST',
      url: `${API_URL}/query/search`,
      body: not ? expect.not.stringContaining(str) : expect.stringContaining(str),
    })
  );
};

const expectNotInBody = (str: string) => expectInBody(str, true);

const createField = (name: string) =>
  new Field({
    name,
    renderer: FieldRenderer.plain,
    importance: FieldImportance.none,
    isEnumerable: false,
    isVirtual: false,
  });

const createAxis = (name: string) => {
  return { name, uiField: createField(name) };
};

describe('search service', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('fetchResults()', () => {
    it('fetches search results', async () => {
      await searchService.fetchResults({ query: 'contributions' });
      expectInBody('"query":"contributions",');
    });

    it('allows specifying an offset (default: 0)', async () => {
      await searchService.fetchResults({ query: '*' });
      expectInBody('"offset":0,');

      await searchService.fetchResults({ query: '*', offset: 9633 });
      expectInBody('"offset":9633,');
    });

    it('allows specifying a limit (default: store limit)', async () => {
      await searchService.fetchResults({ query: '*' });
      expectInBody(`"limit":${stores.search.limit},`);

      await searchService.fetchResults({ query: '*', limit: 7392 });
      expectInBody('"limit":7392,');
    });

    it('allows specifying exact axis constraints (default: none)', async () => {
      SEARCH_CONFIG.FACETS = [{ axis: createAxis('author'), type: 'exact' }];
      jest.spyOn(config, 'SEARCH_FACETS_COMBINE_OPERATOR', 'get').mockReturnValue('and');

      await searchService.fetchResults({ query: '*' });
      expectNotInBody('"axisConstraints"');

      await searchService.fetchResults({
        query: '*',
        filters: [['author', ['Jane Doe']]],
      });
      expectInBody(
        '"axisConstraints":[{"type":"exact","axis":"author","combineOperator":"and","values":["Jane Doe"]}]'
      );
    });

    it('allows specifying yearRange axis constraints (default: none)', async () => {
      SEARCH_CONFIG.FACETS = [{ axis: createAxis('createdAt'), type: 'yearRange' }];
      jest.spyOn(config, 'SEARCH_FACETS_COMBINE_OPERATOR', 'get').mockReturnValue('or');

      await searchService.fetchResults({ query: '*' });
      expectNotInBody('"axisConstraints"');

      await searchService.fetchResults({
        query: '*',
        filters: [['createdAt', ['some invalid string']]],
      });
      expectNotInBody('"axisConstraints"');

      await searchService.fetchResults({
        query: '*',
        filters: [['createdAt', ['1975-2025']]],
      });
      expectInBody(
        '"axisConstraints":[{"type":"stringRange","axis":"createdAt","combineOperator":"or","stringRanges":[{"min":"1975-01-01T00:00:00.000Z","max":"2025-12-31T23:59:59.000Z"}]}]'
      );
    });

    it('allows specifying facets to be fetched (default: SEARCH_CONFIG.FACETS const)', async () => {
      SEARCH_CONFIG.FACETS = [{ axis: createAxis('barfoo'), type: 'exact' }];

      await searchService.fetchResults({ query: '*' });
      expectInBody(`"facets":[{"axis":"barfoo",`);

      await searchService.fetchResults({
        query: '*',
        facets: [{ axis: createAxis('foobar'), type: 'exact' }],
      });
      expectInBody('"facets":[{"axis":"foobar",');

      await searchService.fetchResults({
        query: '*',
        facets: [{ axis: createAxis('bazfoo__label'), type: 'exact' }],
      });
      expectInBody('"facets":[{"axis":"bazfoo__label",');
    });

    it('allows specifying fields to be fetched (default: all fields)', async () => {
      Object.keys(config.FIELDS).forEach((key) => {
        delete config.FIELDS[key];
      });
      config.FIELDS.fieldFromUnitTest = createField('fieldFromUnitTest');
      await searchService.fetchResults({ query: '*' });
      expectInBody(`"fields":["fieldFromUnitTest"]`);

      await searchService.fetchResults({ query: '*', fields: [createField('foobar')] });
      expectInBody('"fields":["foobar"]');

      await searchService.fetchResults({
        query: '*',
        fields: [createField('bazfaz').cloneAndLinkToField(createField('linkedField'))],
      });
      expectInBody('"fields":["bazfaz__linkedField"]');
    });

    it('allows specifying fields for which highlights should be fetched (default: null, will set autoHighlight parameter)', async () => {
      await searchService.fetchResults({ query: '*' });
      expectInBody('"autoHighlight":true');
      expectNotInBody('"highlightFields":');

      await searchService.fetchResults({ query: '*', highlightFields: [createField('foobar')] });
      expectInBody('"highlightFields":["foobar"]');
      expectNotInBody('"autoHighlight":');
    });

    it('allows specifying an offset for the facet buckets to be fetched (default: 0)', async () => {
      const EXAMPLE_FACET = { axis: createAxis('fooField'), type: 'exact' } as SearchFacet;
      SEARCH_CONFIG.FACETS = [EXAMPLE_FACET];

      await searchService.fetchResults({ query: '*' });
      expectInBody(`"facets":[{"axis":"${EXAMPLE_FACET.axis.name}","type":"${EXAMPLE_FACET.type}","offset":0,`);

      await searchService.fetchResults({ query: '*', facetsOffset: 3499 });
      expectInBody(`"facets":[{"axis":"${EXAMPLE_FACET.axis.name}","type":"${EXAMPLE_FACET.type}","offset":3499,`);
    });

    it('allows specifying a limit for the facet buckets to be fetched (default: SEARCH_CONFIG.FACETS_LIMIT const)', async () => {
      const EXAMPLE_FACET = { axis: createAxis('fooField'), type: 'exact' } as SearchFacet;
      SEARCH_CONFIG.FACETS = [EXAMPLE_FACET];

      await searchService.fetchResults({ query: '*' });
      expectInBody(`"limit":${SEARCH_CONFIG.FACETS_LIMIT}`);

      await searchService.fetchResults({ query: '*', facetsLimit: 7247 });
      expectInBody('"limit":7247');
    });

    it('allows specifying a search focus (default: none)', async () => {
      await searchService.fetchResults({ query: '*' });
      expectInBody('"searchFocus":null');

      await searchService.fetchResults({ query: '*', searchFocus: 'contributor' });
      expectInBody('"searchFocus":"contributor"');
    });

    it('allows specifying a sorting (default: stores.search.sorting)', async () => {
      await searchService.fetchResults({ query: '*' });
      expectInBody('"sorting":null,');

      await searchService.fetchResults({ query: '*', sorting: { axis: createAxis('foo'), order: 'asc' } });
      expectInBody('"sorting":{"axis":"foo","order":"asc"},');
    });

    it('allows specifying a query max edit distance (default: SEARCH_CONFIG.QUERY_MAX_EDIT_DISTANCE const)', async () => {
      await searchService.fetchResults({ query: '*' });
      expectInBody(`"maxEditDistance":${SEARCH_CONFIG.QUERY_MAX_EDIT_DISTANCE},`);

      await searchService.fetchResults({ query: '*', maxEditDistance: 2 });
      expectInBody('"maxEditDistance":2,');
    });

    it('allows specifying whether to use the N-gram search field (default: SEARCH_CONFIG.QUERY_USE_NGRAM_FIELD const)', async () => {
      await searchService.fetchResults({ query: '*' });
      expectInBody(`"useNgramField":${SEARCH_CONFIG.QUERY_USE_NGRAM_FIELD},`);

      // We pass the negated default value as an argument to ensure that we test a non-default value
      await searchService.fetchResults({ query: '*', useNgramField: !SEARCH_CONFIG.QUERY_USE_NGRAM_FIELD });
      expectInBody(`"useNgramField":${!SEARCH_CONFIG.QUERY_USE_NGRAM_FIELD},`);
    });

    it('allows specifying a timeout for the request', async () => {
      await searchService.fetchResults({ query: '*', timeout: 30 });
      expect(requestSpy).toHaveBeenCalledWith(expect.objectContaining({ timeout: 30 }));

      await searchService.fetchResults({ query: '*' });
      expect(requestSpy).toHaveBeenCalledWith(expect.not.objectContaining({ timeout: 30 }));
    });

    it('throws an error when the request fails', async () => {
      const consoleSpy = jest.spyOn(console, 'error');
      consoleSpy.mockImplementationOnce(() => {});
      requestSpy.mockImplementationOnce(() => {
        throw new Error('Network error');
      });

      await expect(searchService.fetchResults({ query: '*' })).rejects.toThrow();
    });
  });

  describe('updateResults()', () => {
    let searchResponseSetterSpy, searchFacetsSetterSpy, searchIsBusySetterSpy, filtersAllSpy;

    beforeAll(() => {
      searchResponseSetterSpy = jest.spyOn(stores.search, 'response', 'set');
      searchFacetsSetterSpy = jest.spyOn(stores.search, 'facets', 'set');
      searchIsBusySetterSpy = jest.spyOn(stores.search, 'isBusy', 'set');
      filtersAllSpy = jest.spyOn(stores.filters, 'all', 'get');
    });

    it('has dependencies to the search and filters store', async () => {
      await searchService.updateResults({});
      expect(searchResponseSetterSpy).toHaveBeenCalledTimes(1);
      expect(searchFacetsSetterSpy).toHaveBeenCalledTimes(1);
      expect(searchIsBusySetterSpy).toHaveBeenCalledTimes(2);
      expect(filtersAllSpy).toHaveBeenCalled();
    });

    it('fetches search results and updates the search store with the response', async () => {
      expect(await searchService.updateResults({})).toBe(undefined);
      expect(requestSpy).toHaveBeenCalled();
      expect(searchResponseSetterSpy).toHaveBeenCalled();
    });

    it('resets the response and facets in the search store when reset argument is true', async () => {
      await searchService.updateResults({ reset: true });
      expect(searchResponseSetterSpy).toHaveBeenCalledTimes(2);
      expect(searchResponseSetterSpy).toHaveBeenCalledWith(null);
      expect(searchResponseSetterSpy).toHaveBeenCalledWith(EXAMPLE_SEARCH_RESPONSE);
      expect(searchFacetsSetterSpy).toHaveBeenCalledTimes(2);
      expect(searchFacetsSetterSpy).toHaveBeenCalledWith(null);
      expect(searchFacetsSetterSpy).toHaveBeenCalledWith(undefined);
    });

    it('allows specifying the query and propagates it to the request', async () => {
      await searchService.updateResults({ query: 'foo bar baz' });
      expectInBody('"query":"foo bar baz",');
    });

    it('allows specifying the offset and propagates it to the request', async () => {
      await searchService.updateResults({ offset: 8533 });
      expectInBody('"offset":8533,');
    });

    it('propagates filters from the store to the request', async () => {
      SEARCH_CONFIG.FACETS = [{ axis: createAxis('author'), type: 'exact' }];

      filtersAllSpy.mockImplementation(() => []);
      await searchService.updateResults({});
      expectNotInBody('"axisConstraints"');

      jest.clearAllMocks();
      filtersAllSpy.mockImplementation(() => [['author', ['Eva Doe']]]);
      await searchService.updateResults({});
      expectInBody('"axisConstraints"');
      expectInBody('"Eva Doe"');
    });
  });

  describe('fetchHierarchy()', () => {
    it('fetches hierarchy tree nodes', async () => {
      await searchService.fetchHierarchy(
        createField('example-key'),
        'Person',
        createField('foo-link-field'),
        createField('bar-display-field'),
        'baz-language'
      );

      expect(requestSpy).toHaveBeenCalledWith({
        authorized: true,
        method: 'POST',
        url: `${API_URL}/metadata/tree`,
        body: expect.stringContaining(
          `"nodeEntityType":"Person","linkFieldName":"foo-link-field","displayFieldName":"bar-display-field","displayFieldLanguage":"baz-language"`
        ),
      });
    });
  });
});
