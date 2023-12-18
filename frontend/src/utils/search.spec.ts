jest.mock('stencil-router-v2');

import * as config from 'config';
import { PadCondition } from 'config/search';
import stores from 'stores';
import { SearchResultsItem } from 'stores/search';
import {
  addSearchOperator,
  buildAscNumSequence,
  createExactAxisConstraint,
  getFilterQueryParamName,
  getQueryStringFromParams,
  getSearchUrl,
  formatResultsItemsAsCsv,
  getHierarchyNodesFirstLevel,
} from './search';

const { ROUTES, SEARCH_CONFIG, SEARCH_PAGINATION_START } = config;

describe('search util', () => {
  describe('getFilterQueryParamName()', () => {
    it('returns SEARCH_PARAM_FILTER_PREFIX + filterName + "[]"', () => {
      jest.spyOn(config, 'SEARCH_PARAM_FILTER_PREFIX', 'get').mockReturnValueOnce('filter.');

      expect(getFilterQueryParamName('foo')).toBe('filter.foo[]');
    });
  });

  describe('getQueryStringFromParams()', () => {
    it('converts url params to string and replaces encoded array notation', () => {
      const params = new URLSearchParams();
      params.append('filters.author[]', 'John Doe');
      params.append('filters.author[]', 'Jane Doe');
      params.append('filters.year[]', '2020');

      expect(getQueryStringFromParams(params)).toBe(
        'filters.author[]=John+Doe&filters.author[]=Jane+Doe&filters.year[]=2020'
      );
    });
  });

  describe('getSearchUrl()', () => {
    let filtersSpy, sortingSpy, querySpy, focusSpy, searchParamFilterPrefixSpy;

    beforeAll(() => {
      filtersSpy = jest.spyOn(stores.filters, 'all', 'get');
      sortingSpy = jest.spyOn(stores.search, 'sorting', 'get');
      querySpy = jest.spyOn(stores.search, 'query', 'get');
      focusSpy = jest.spyOn(stores.search, 'focus', 'get');
      searchParamFilterPrefixSpy = jest.spyOn(config, 'SEARCH_PARAM_FILTER_PREFIX', 'get');
    });

    beforeEach(() => {
      jest.resetAllMocks();

      filtersSpy.mockReturnValue([]);
      sortingSpy.mockReturnValue(SEARCH_CONFIG.SORTING_OPTIONS[0]);
      querySpy.mockReturnValue('');
      focusSpy.mockReturnValue(null);
      searchParamFilterPrefixSpy.mockReturnValue('');
    });

    it('has dependencies to the filters store and search + query in the search store', () => {
      expect(filtersSpy).not.toHaveBeenCalled();
      expect(sortingSpy).not.toHaveBeenCalled();
      expect(querySpy).not.toHaveBeenCalled();

      getSearchUrl();

      expect(filtersSpy).toHaveBeenCalled();
      expect(sortingSpy).toHaveBeenCalled();
      expect(querySpy).toHaveBeenCalled();
    });

    it('reflects selected filters unless resetFilters argument is true', () => {
      const dummyFilters = [
        ['foo', ['a', 'b', 'c']],
        ['bar', ['d']],
      ];

      filtersSpy.mockReturnValueOnce(dummyFilters);
      expect(getSearchUrl()).toBe(`${ROUTES.SEARCH}?bar[]=d&foo[]=a&foo[]=b&foo[]=c`);

      filtersSpy.mockReturnValueOnce(dummyFilters);
      expect(getSearchUrl(true)).toBe(ROUTES.SEARCH);
    });

    it('reflects selected sorting unless resetFilters argument is true or first (default) sorting option is choosen', () => {
      const dummySorting = {
        axis: { name: 'foo' },
        order: 'asc',
      };

      sortingSpy.mockReturnValueOnce(dummySorting);
      expect(getSearchUrl()).toBe(`${ROUTES.SEARCH}?sorting.axis=foo&sorting.order=asc`);

      sortingSpy.mockReturnValueOnce(dummySorting);
      expect(getSearchUrl(true)).toBe(ROUTES.SEARCH);

      sortingSpy.mockReturnValueOnce(SEARCH_CONFIG.SORTING_OPTIONS[0]);
      expect(getSearchUrl()).toBe(ROUTES.SEARCH);
    });

    it('reflects the (uri-encoded) query', () => {
      querySpy.mockReturnValueOnce('');
      expect(getSearchUrl()).toBe(ROUTES.SEARCH);

      querySpy.mockReturnValueOnce('foo bar');
      expect(getSearchUrl()).toBe(ROUTES.SEARCH_QUERY.replace(':query', 'foo%20bar'));

      querySpy.mockReturnValueOnce('https://evil.foo:8080/bar');
      expect(getSearchUrl()).toBe(ROUTES.SEARCH_QUERY.replace(':query', 'https%3A%2F%2Fevil.foo%3A8080%2Fbar'));
    });

    it('reflects the selected search focus', () => {
      focusSpy.mockReturnValueOnce(null);
      querySpy.mockReturnValueOnce('foo');
      expect(getSearchUrl()).toBe(ROUTES.SEARCH_QUERY.replace(':query', 'foo'));

      focusSpy.mockReturnValueOnce('title');
      querySpy.mockReturnValueOnce('bar');
      expect(getSearchUrl()).toBe(ROUTES.SEARCH_QUERY.replace(':query', 'bar') + '?focus=title');

      focusSpy.mockReturnValueOnce('keyword');
      querySpy.mockReturnValueOnce('baz');
      expect(getSearchUrl()).toBe(ROUTES.SEARCH_QUERY.replace(':query', 'baz') + '?focus=keyword');
    });

    it('reflects page when given as argument (unless when page is 1)', () => {
      expect(getSearchUrl(false, 2)).toBe(`${ROUTES.SEARCH}?page=2`);
      expect(getSearchUrl(false, 44)).toBe(`${ROUTES.SEARCH}?page=44`);
      expect(getSearchUrl(false, 1)).toBe(ROUTES.SEARCH);
    });

    it('combines filters, sorting, query, focus and page in uri', () => {
      sortingSpy.mockReturnValueOnce({
        axis: { name: 'bar' },
        order: 'desc',
      });
      filtersSpy.mockReturnValueOnce([['foo', ['test']]]);
      querySpy.mockReturnValueOnce('baz');
      focusSpy.mockReturnValueOnce('contributor');

      expect(getSearchUrl(false, 5)).toBe(
        `${ROUTES.SEARCH_QUERY.replace(
          ':query',
          'baz'
        )}?focus=contributor&foo[]=test&page=5&sorting.axis=bar&sorting.order=desc`
      );
    });

    it('allows to specify a custom query', () => {
      expect(getSearchUrl(false, SEARCH_PAGINATION_START, 'foo bar')).toBe(`${ROUTES.SEARCH}/foo%20bar`);
    });
  });

  describe('addSearchOperator()', () => {
    let mockedInputEl, valueGetter, valueSetter, selectionStartGetter, selectionEndGetter, padMapSpy;

    beforeEach(() => {
      jest.clearAllMocks();

      valueGetter = jest.fn(() => '');
      valueSetter = jest.fn();
      selectionStartGetter = jest.fn(() => 0);
      selectionEndGetter = jest.fn(() => 0);

      mockedInputEl = {
        get value() {
          return valueGetter();
        },
        set value(value: string) {
          valueSetter(value);
        },
        get selectionStart() {
          return selectionStartGetter();
        },
        get selectionEnd() {
          return selectionEndGetter();
        },
        setSelectionRange: jest.fn(() => {}) as any,
        dispatchEvent: jest.fn(() => {}) as any,
        focus: jest.fn(() => {}) as any,
      } as HTMLInputElement;

      padMapSpy = jest.spyOn(config, 'SEARCH_OPERATOR_PAD_MAP', 'get');
      padMapSpy.mockImplementation(() => ({
        default: [PadCondition.MAY, PadCondition.MAY],
      }));
    });

    it('gets the input value and adds the operator', () => {
      expect(valueGetter).not.toHaveBeenCalled();
      expect(valueSetter).not.toHaveBeenCalled();

      addSearchOperator(mockedInputEl, '*');

      expect(valueGetter).toHaveBeenCalled();
      expect(valueSetter).toHaveBeenCalledWith('*');
    });

    it('adds the operator at the position of the input cursor', () => {
      valueGetter.mockReturnValue('This is a test');

      selectionStartGetter.mockReturnValueOnce(4);
      selectionEndGetter.mockReturnValueOnce(4);
      addSearchOperator(mockedInputEl, '*');
      expect(valueSetter).toHaveBeenCalledWith('This* is a test');

      selectionStartGetter.mockReturnValueOnce(14);
      selectionEndGetter.mockReturnValueOnce(14);
      addSearchOperator(mockedInputEl, '!!');
      expect(valueSetter).toHaveBeenCalledWith('This is a test!!');
    });

    it('adds the single-character operator at the beginning of the input selection range', () => {
      valueGetter.mockReturnValue('This is a test');

      selectionStartGetter.mockReturnValueOnce(10);
      selectionEndGetter.mockReturnValueOnce(14);
      addSearchOperator(mockedInputEl, '#');
      expect(valueSetter).toHaveBeenCalledWith('This is a #test');
    });

    it('wraps the multi-character operator around the input selection range', () => {
      valueGetter.mockReturnValue('This is a test');

      selectionStartGetter.mockReturnValueOnce(8);
      selectionEndGetter.mockReturnValueOnce(14);
      addSearchOperator(mockedInputEl, '""');
      expect(valueSetter).toHaveBeenCalledWith('This is "a test"');
    });

    it('adds or removes a padding according to a set of conditions', () => {
      const testForPadding = (start, end, operator, expectedText) => {
        valueSetter.mockClear();
        selectionStartGetter.mockReturnValueOnce(start);
        selectionEndGetter.mockReturnValueOnce(end);
        addSearchOperator(mockedInputEl, operator);
        expect(valueSetter).toHaveBeenCalledWith(expectedText);
      };

      padMapSpy.mockImplementation(() => ({
        '/': [PadCondition.MUST, PadCondition.MUST],
        ':': [PadCondition.MUST_NOT, PadCondition.MUST],
        '#': [PadCondition.MUST, PadCondition.MUST_NOT],
        '""': [PadCondition.MUST, PadCondition.MUST],
        '!': [PadCondition.MUST_NOT, PadCondition.MUST_NOT],
        default: [PadCondition.MAY, PadCondition.MAY],
      }));

      valueGetter.mockReturnValue('To benot to be');

      testForPadding(5, 5, '/', 'To be / not to be');
      testForPadding(8, 8, '/', 'To benot / to be');
      testForPadding(3, 3, '/', 'To / benot to be');
      testForPadding(5, 9, '/', 'To be / not to be');
      testForPadding(14, 14, '/', 'To benot to be / ');

      valueGetter.mockReturnValueOnce('To be  not to be');
      testForPadding(6, 6, '/', 'To be / not to be');

      valueGetter.mockReturnValue('Unit tests are great');

      testForPadding(13, 13, ':', 'Unit tests ar: e great');
      testForPadding(14, 14, ':', 'Unit tests are: great');
      testForPadding(15, 15, ':', 'Unit tests are: great');

      testForPadding(13, 13, '#', 'Unit tests ar #e great');
      testForPadding(14, 14, '#', 'Unit tests are #great');
      testForPadding(15, 15, '#', 'Unit tests are #great');

      testForPadding(13, 13, '""', 'Unit tests ar "" e great');
      testForPadding(14, 20, '""', 'Unit tests are " great"');
      testForPadding(15, 20, '""', 'Unit tests are "great"');
      testForPadding(5, 10, '""', 'Unit "tests" are great');
      testForPadding(0, 10, '""', '"Unit tests" are great');
      testForPadding(20, 20, '""', 'Unit tests are great ""');

      testForPadding(9, 9, '?', 'Unit test?s are great');
      testForPadding(20, 20, '?', 'Unit tests are great?');

      valueGetter.mockReturnValueOnce('onelasttest');
      testForPadding(3, 7, '""', 'one "last" test');

      valueGetter.mockReturnValueOnce('very  last test');
      testForPadding(5, 5, '!', 'very!last test');
    });

    it('prevents repeated use when same operator is already present at position', () => {
      valueGetter.mockReturnValue('foo ""');
      selectionStartGetter.mockReturnValue(5);
      selectionEndGetter.mockReturnValue(5);
      addSearchOperator(mockedInputEl, '""');
      expect(valueSetter).not.toHaveBeenCalled();
    });

    it('allows to specify a custom padding character', () => {
      padMapSpy.mockImplementation(() => ({
        default: [PadCondition.MUST, PadCondition.MUST],
      }));
      valueGetter.mockReturnValue('testcase');
      selectionStartGetter.mockReturnValue(4);
      selectionEndGetter.mockReturnValue(4);

      addSearchOperator(mockedInputEl, '/', '#');
      expect(valueSetter).toHaveBeenCalledWith('test#/#case');

      addSearchOperator(mockedInputEl, ':', '?');
      expect(valueSetter).toHaveBeenCalledWith('test?:?case');
    });

    it('sets a new input cursor position', () => {
      addSearchOperator(mockedInputEl, '*');
      expect(mockedInputEl.setSelectionRange).toHaveBeenCalledWith(1, 1);
      addSearchOperator(mockedInputEl, '!!');
      // middle of the multi-character operator, if no text was selected
      expect(mockedInputEl.setSelectionRange).toHaveBeenCalledWith(1, 1);

      selectionStartGetter.mockReturnValue(0);
      selectionEndGetter.mockReturnValue(8);
      valueGetter.mockReturnValue('testcase');
      addSearchOperator(mockedInputEl, '""');
      // end of the multi-character operator, if text was selected
      expect(mockedInputEl.setSelectionRange).toHaveBeenCalledWith(10, 10);
    });

    it('dispatches an input event', () => {
      expect(mockedInputEl.dispatchEvent).not.toHaveBeenCalled();
      addSearchOperator(mockedInputEl, '*');
      expect(mockedInputEl.dispatchEvent).toHaveBeenCalled();
    });

    it('focuses the input element', () => {
      expect(mockedInputEl.focus).not.toHaveBeenCalled();
      addSearchOperator(mockedInputEl, '*');
      expect(mockedInputEl.focus).toHaveBeenCalled();
    });
  });

  describe('buildAscNumSequence()', () => {
    it('builds an ascending sequence of numbers without reordering the input array', () => {
      expect(buildAscNumSequence([5])).toEqual([5]);
      expect(buildAscNumSequence([5, 4])).toEqual([4, 4]);
      expect(buildAscNumSequence([3, 4, 1, 9])).toEqual([1, 1, 1, 9]);
      expect(buildAscNumSequence([2, 3, 7, 6])).toEqual([2, 3, 6, 6]);
      expect(buildAscNumSequence([5, 4, 3, 2])).toEqual([2, 2, 2, 2]);
      expect(buildAscNumSequence([5, 6, 7, 8])).toEqual([5, 6, 7, 8]);
      expect(buildAscNumSequence([3, 4, 4, 6, 5])).toEqual([3, 4, 4, 5, 5]);
    });

    it('allows to specify one index that must be preserved', () => {
      expect(buildAscNumSequence([3, 4, 1, 9], 0)).toEqual([3, 3, 3, 9]);
      expect(buildAscNumSequence([3, 4, 1, 9], 1)).toEqual([1, 4, 4, 9]);
      expect(buildAscNumSequence([9, 8, 7, 6], 2)).toEqual([6, 6, 7, 7]);
    });

    it('allows to specify a total min and max value', () => {
      const NO_PRESERVED_INDEX = -1;

      // min value
      expect(buildAscNumSequence([-5, 6, 7], NO_PRESERVED_INDEX, 0)).toEqual([0, 6, 7]);
      expect(buildAscNumSequence([5, 99], NO_PRESERVED_INDEX, 12)).toEqual([12, 99]);
      expect(buildAscNumSequence([6, 7], NO_PRESERVED_INDEX, 8)).toEqual([8, 8]);

      // min + max value
      expect(buildAscNumSequence([2, 3, 4, 5, 6], NO_PRESERVED_INDEX, 3, 5)).toEqual([3, 3, 4, 5, 5]);
      expect(buildAscNumSequence([50, 100], NO_PRESERVED_INDEX, 50, 60)).toEqual([50, 60]);
    });

    it('combines all of the above', () => {
      expect(buildAscNumSequence([88, 6, -5, 7, 10, 9], 4, 0, 5)).toEqual([0, 0, 0, 5, 5, 5]);
      expect(buildAscNumSequence([88, 77], 0, 0, 80)).toEqual([80, 80]);
    });
  });

  describe('createExactAxisConstraint()', () => {
    it('creates an axis constraint', () => {
      expect(createExactAxisConstraint('foo', ['bar', 'baz'])).toEqual({
        axis: 'foo',
        combineOperator: config.SEARCH_FACETS_COMBINE_OPERATOR,
        type: 'exact',
        values: ['bar', 'baz'],
      });
    });

    it('creates an axis constraint for single node values if filter value contains prefix', () => {
      jest.spyOn(config, 'SEARCH_FACETS_SINGLE_NODE_VALUE_PREFIX', 'get').mockReturnValue('direct:');

      expect(createExactAxisConstraint('foo', ['direct:bar', 'baz'])).toEqual({
        axis: 'foo',
        combineOperator: config.SEARCH_FACETS_COMBINE_OPERATOR,
        type: 'exact',
        values: ['baz'],
        singleNodeValues: ['bar'],
      });
      expect(createExactAxisConstraint('foo', ['direct:bar', 'direct:baz'])).toEqual({
        axis: 'foo',
        combineOperator: config.SEARCH_FACETS_COMBINE_OPERATOR,
        type: 'exact',
        values: [],
        singleNodeValues: ['bar', 'baz'],
      });
    });
  });

  describe('formatResultsItemsAsCsv()', () => {
    it('formats and returns given search results items as multi-lined csv string', () => {
      const MOCKED_SEARCH_RESULTS_ITEM_1 = {
        itemId: 'item-foo',
        entityType: 'Source',
        values: [
          {
            fieldName: 'author',
            fieldValue: 'Jane Doe',
          },
          {
            fieldName: 'author',
            fieldValue: 'John Doe',
          },
          {
            fieldName: 'label',
            fieldValue: 'Top secret project',
          },
          {
            fieldName: 'accessRestriction',
            fieldValue: 'restricted',
          },
        ],
      } as SearchResultsItem;
      const MOCKED_SEARCH_RESULTS_ITEM_2 = {
        itemId: 'item-bar',
        entityType: 'Source',
        values: [
          {
            fieldName: 'label',
            fieldValue: 'Yet another project',
          },
          {
            fieldName: 'author',
            fieldValue: 'Maxim Doe',
          },
        ],
      } as SearchResultsItem;

      expect(formatResultsItemsAsCsv([MOCKED_SEARCH_RESULTS_ITEM_1, MOCKED_SEARCH_RESULTS_ITEM_2])).toBe(
        '"accessRestriction[0]","author[0]","author[1]","entityName[0]","label[0]"\r\n' +
          '"restricted","Jane Doe","John Doe","Source","Top secret project"\r\n' +
          '"","Maxim Doe","","Source","Yet another project"'
      );
    });
  });

  describe('getHierarchyNodesFirstLevel', () => {
    it('returns the lowest available level number of given nodes while respecting a min level (default: 0)', () => {
      const nodes = [
        { depth: -1, nodeId: 'node-a' },
        { depth: 1, nodeId: 'node-b' },
        { depth: 2, nodeId: 'node-c' },
        { depth: 4, nodeId: 'node-d' },
      ];
      expect(getHierarchyNodesFirstLevel(nodes)).toBe(1);
      expect(getHierarchyNodesFirstLevel(nodes, 0)).toBe(1);
      expect(getHierarchyNodesFirstLevel(nodes, 1)).toBe(1);
      expect(getHierarchyNodesFirstLevel(nodes, 2)).toBe(2);
      expect(getHierarchyNodesFirstLevel(nodes, 3)).toBe(4);
      expect(getHierarchyNodesFirstLevel(nodes, 4)).toBe(4);
      expect(getHierarchyNodesFirstLevel(nodes, 5)).toBe(Infinity);
    });
  });
});
