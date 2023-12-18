jest.mock('stencil-router-v2');

import { SEARCH_CONFIG } from 'config';
import * as config from 'config';
import search, { SearchResults } from './search';

const paginationStartSpy = jest.spyOn(config, 'SEARCH_PAGINATION_START', 'get');
paginationStartSpy.mockImplementation(() => 1);

const getItemValue = (fieldName: string, prefix: string) => ({
  fieldName,
  fieldValue: `${prefix.slice(0, 1).toUpperCase()}${prefix.slice(1)} ${fieldName
    .slice(0, 1)
    .toUpperCase()}${fieldName.slice(1)}`,
});

const ITEM_0 = {
  itemId: 'item-0',
  entityType: 'Source',
  values: [getItemValue('title', 'foo'), getItemValue('author', 'Jane Doe'), getItemValue('author', 'John Doe')],
};

const ITEM_1 = {
  itemId: 'item-1',
  entityType: 'Resource',
  values: [getItemValue('title', 'bar'), getItemValue('author', 'John Doe')],
};

const ITEM_2 = {
  itemId: 'item-2',
  entityType: 'Resource',
  values: [getItemValue('title', 'baz')],
};

const EXAMPLE_SEARCH_FACETS = [
  {
    type: 'exact' as 'exact',
    axis: 'author',
    bucketNo: 2,
    buckets: [
      { value: 'John Doe', count: 2 },
      { value: 'Jane Doe', count: 1 },
    ],
  },
];

const EXAMPLE_SEARCH_HIGHLIGHTS = [
  {
    itemId: 'foo',
    matches: [
      {
        fieldName: 'author',
        snippets: ['Thomas <em>Doe</em>'],
      },
    ],
  },
];

const EXAMPLE_SEARCH_RESPONSE: SearchResults = {
  numFound: 3,
  numFoundExact: true,
  start: 0,
  maxScore: 0,
  items: [ITEM_0, ITEM_1, ITEM_2],
  facets: EXAMPLE_SEARCH_FACETS,
  highlights: EXAMPLE_SEARCH_HIGHLIGHTS,
  diagnostics: { cleanedQuery: 'foo bar baz', queryWasCleaned: true },
};

describe('search store', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('returns and stores a search response', () => {
    expect(search.response).toBe(undefined);
    search.response = EXAMPLE_SEARCH_RESPONSE;
    expect(search.response).toBe(EXAMPLE_SEARCH_RESPONSE);
    search.response = undefined;
    expect(search.response).toBe(undefined);
  });

  it('has an isBusy getter (bool, initially true)', () => {
    expect(typeof search.isBusy).toBe('boolean');
    expect(search.isBusy).toBe(true);
  });

  it('has an isBusy setter (only effective when a response is stored)', () => {
    search.isBusy = false;
    expect(search.isBusy).toBe(true);

    search.response = EXAMPLE_SEARCH_RESPONSE;
    expect(search.isBusy).toBe(false);

    search.isBusy = true;
    expect(search.isBusy).toBe(true);
  });

  it('returns and stores the search query', () => {
    expect(search.query).toBe('');
    search.query = 'foo';
    expect(search.query).toBe('foo');
  });

  it('returns the search offset, reflecting the search response (default: 0)', () => {
    search.response = undefined;
    expect(search.offset).toBe(0);

    search.response = { ...EXAMPLE_SEARCH_RESPONSE, start: 7 };
    expect(search.offset).toBe(7);
  });

  it('returns and stores the search limit (default: SEARCH_CONFIG.LIMIT const)', () => {
    expect(search.limit).toBe(SEARCH_CONFIG.LIMIT);

    search.limit = SEARCH_CONFIG.LIMIT + 3;
    expect(search.limit).toBe(SEARCH_CONFIG.LIMIT + 3);
  });

  it('returns the total search results count (from the search response)', () => {
    search.response = undefined;
    expect(search.count).toBe(undefined);

    search.response = EXAMPLE_SEARCH_RESPONSE;
    expect(search.count).toBe(3);
  });

  it('provides the number of placeholders to be rendered while fetching the search results', () => {
    search.response = undefined;
    expect(search.placeholdersCount).toBe(search.limit);

    search.response = EXAMPLE_SEARCH_RESPONSE;
    expect(search.placeholdersCount).toBe(3);

    search.response = { ...EXAMPLE_SEARCH_RESPONSE, start: 2 };
    expect(search.placeholdersCount).toBe(1);
  });

  it('returns the current page based on offset, limit and SEARCH_PAGINATION_START const', () => {
    search.limit = 9;
    search.response = EXAMPLE_SEARCH_RESPONSE;
    expect(search.page).toBe(1);
    expect(paginationStartSpy).toHaveBeenCalled();

    search.response = { ...EXAMPLE_SEARCH_RESPONSE, start: 2 };
    expect(search.page).toBe(1);

    search.limit = 2;
    expect(search.page).toBe(2);

    search.limit = 1;
    expect(search.page).toBe(3);
  });

  it('provides the page range to be rendered based on current page, SEARCH_CONFIG.PAGINATION_RANGE_COUNT and SEARCH_PAGINATION_START consts', () => {
    search.response = EXAMPLE_SEARCH_RESPONSE;
    search.limit = 3;
    SEARCH_CONFIG.PAGINATION_RANGE_COUNT = 5;
    expect(search.pageRange).toEqual([1]);
    expect(paginationStartSpy).toHaveBeenCalled();

    search.limit = 2;
    expect(search.pageRange).toEqual([1, 2]);

    search.limit = 1;
    expect(search.pageRange).toEqual([1, 2, 3]);

    SEARCH_CONFIG.PAGINATION_RANGE_COUNT = 2;
    expect(search.pageRange).toEqual([1, 2]);

    search.response = { ...EXAMPLE_SEARCH_RESPONSE, start: 1 };
    expect(search.pageRange).toEqual([2, 3]);

    SEARCH_CONFIG.PAGINATION_RANGE_COUNT = 1;
    expect(search.pageRange).toEqual([2]);
  });

  it('returns and stores the facets (default: [])', () => {
    search.facets = undefined;
    expect(search.facets).toEqual([]);

    search.facets = EXAMPLE_SEARCH_FACETS;
    expect(search.facets).toBe(EXAMPLE_SEARCH_FACETS);

    search.facets = [];
    expect(search.facets).toEqual([]);
  });

  it('adds buckets to a facet by field name', () => {
    search.facets = EXAMPLE_SEARCH_FACETS;

    const oldBuckets = search.facets.find(({ axis }) => axis === 'author').buckets;
    const newBuckets = [
      { value: 'Foo new', count: 7 },
      { value: 'Bar new', count: 2 },
    ];
    search.addFacetBuckets('author', newBuckets);

    expect(search.facets.find(({ axis }) => axis === 'author').buckets).toEqual([...oldBuckets, ...newBuckets]);
  });

  it('does not add buckets to a facet when it is unknown', () => {
    search.facets = EXAMPLE_SEARCH_FACETS;

    search.addFacetBuckets('unknownField', [{ value: 'Foo', count: 3 }]);
    expect(search.facets).toEqual(EXAMPLE_SEARCH_FACETS);
  });

  it('returns the facet offset for subsequent calls to fetch more buckets', () => {
    search.facets = EXAMPLE_SEARCH_FACETS;

    expect(search.getFacetOffset('author')).toBe(2);
    expect(search.getFacetOffset('unknown-field')).toBe(0);
  });

  it('returns if there are more buckets to fetch for a facet', () => {
    search.facets = EXAMPLE_SEARCH_FACETS;

    expect(search.hasMoreBuckets('author')).toBe(false);
    expect(search.hasMoreBuckets('unknown-field')).toBe(false);

    search.facets = [{ ...EXAMPLE_SEARCH_FACETS[0], bucketNo: 3 }];
    expect(search.hasMoreBuckets('author')).toBe(true);
  });

  it('returns and stores the search focus (default: null)', () => {
    search.focus = undefined;
    expect(search.focus).toBe(null);

    search.focus = 'contact';
    expect(search.focus).toBe('contact');
  });

  it('returns and stores the sorting (default: SEARCH_CONFIG.SORTING_OPTIONS[0])', () => {
    search.sorting = undefined;
    expect(search.sorting).toBe(SEARCH_CONFIG.SORTING_OPTIONS[0]);

    const exampleSorting = { axis: { name: 'foo' }, order: 'asc' as 'asc' };
    search.sorting = exampleSorting;
    expect(search.sorting).toBe(exampleSorting);

    search.sorting = SEARCH_CONFIG.SORTING_OPTIONS[0];
    expect(search.sorting).toBe(SEARCH_CONFIG.SORTING_OPTIONS[0]);
  });

  it('returns the search highlights from the response', () => {
    search.response = EXAMPLE_SEARCH_RESPONSE;
    expect(search.highlights).toBe(EXAMPLE_SEARCH_HIGHLIGHTS);

    search.response = { ...EXAMPLE_SEARCH_RESPONSE, highlights: [] };
    expect(search.highlights).toEqual([]);
  });

  it('returns if the search query was cleaned', () => {
    search.response = EXAMPLE_SEARCH_RESPONSE;
    expect(search.wasQueryCleaned).toBe(true);

    search.response = { ...EXAMPLE_SEARCH_RESPONSE, diagnostics: { queryWasCleaned: false, cleanedQuery: '' } };
    expect(search.wasQueryCleaned).toBe(false);

    search.response = null;
    expect(search.wasQueryCleaned).toBe(false);
  });

  it('returns cleaned query', () => {
    search.response = EXAMPLE_SEARCH_RESPONSE;
    expect(search.cleanedQuery).toBe('foo bar baz');

    search.response = {
      ...EXAMPLE_SEARCH_RESPONSE,
      diagnostics: { queryWasCleaned: true, cleanedQuery: 'blah blah blah' },
    };
    expect(search.cleanedQuery).toBe('blah blah blah');
  });
});
