import {
  ROUTES,
  SEARCH_OPERATOR_PAD_MAP,
  SEARCH_FACETS_COMBINE_OPERATOR,
  SEARCH_PAGINATION_START,
  SEARCH_PARAM_SORTING_AXIS,
  SEARCH_PARAM_SORTING_ORDER,
  SEARCH_PARAM_PAGE,
  SEARCH_PARAM_FOCUS,
  SEARCH_PARAM_FILTER_PREFIX,
  SEARCH_FACETS_SINGLE_NODE_VALUE_PREFIX,
  SEARCH_CONFIG,
  FIELDS,
} from 'config';
import { PadCondition, SearchAxisConstraint } from 'config/search';
import stores from 'stores';
import { HierarchyNode, SearchResultsItem, SearchResultsItemValue } from 'stores/search';
import { filterDuplicates } from './field';

export const getFilterQueryParamName = (filterName: string): string => `${SEARCH_PARAM_FILTER_PREFIX}${filterName}[]`;
export const getQueryStringFromParams = (params: URLSearchParams) => params.toString().replace(/%5B%5D/g, '[]');

export const getSearchUrl = (resetFilters = false, page = SEARCH_PAGINATION_START, query = stores.search.query) => {
  const params = resetFilters ? new URLSearchParams() : stores.filters.toQueryParams();
  const { focus } = stores.search;

  const sorting = stores.search.sorting;
  if (!resetFilters && sorting !== SEARCH_CONFIG.SORTING_OPTIONS[0]) {
    params.append(SEARCH_PARAM_SORTING_AXIS, sorting.axis?.name);
    params.append(SEARCH_PARAM_SORTING_ORDER, sorting.order);
  }

  if (focus) {
    params.append(SEARCH_PARAM_FOCUS, focus);
  }

  page !== SEARCH_PAGINATION_START && params.append(SEARCH_PARAM_PAGE, String(page));

  params.sort();
  const queryStr = getQueryStringFromParams(params);

  return (
    (query ? ROUTES.SEARCH_QUERY.replace(':query', encodeURIComponent(query)) : ROUTES.SEARCH) +
    (queryStr ? `?${queryStr}` : '')
  );
};

export const addSearchOperator = (inputEl: HTMLInputElement, operator: string, pad = ' '): void => {
  const { value, selectionStart: cursorStart, selectionEnd: cursorEnd } = inputEl;
  const [operatorStart = '', operatorEnd = ''] = operator;
  const charLeft = value[cursorStart - 1];
  const charRight = operatorEnd ? value[cursorEnd] : value[cursorStart];
  const padConfig = SEARCH_OPERATOR_PAD_MAP[operator] ?? SEARCH_OPERATOR_PAD_MAP.default;
  const segments = operatorEnd
    ? [value.slice(0, cursorStart), value.slice(cursorStart, cursorEnd), value.slice(cursorEnd)]
    : [value.slice(0, cursorStart), '', value.slice(cursorStart)];

  if (operatorStart === charLeft && operatorEnd === charRight) {
    return inputEl.focus();
  }

  if (padConfig[0] === PadCondition.MUST && (charLeft ?? pad) !== pad) {
    segments[0] += pad;
  }

  if (padConfig[0] === PadCondition.MUST_NOT && charLeft === pad) {
    segments[0] = segments[0].slice(0, -1);
  }

  if (padConfig[1] === PadCondition.MUST && (charRight ?? (operatorEnd ? pad : '')) !== pad) {
    segments[2] = pad + segments[2];
  }

  if (padConfig[1] === PadCondition.MUST_NOT && charRight === pad) {
    segments[2] = segments[2].slice(1);
  }

  segments[0] += operatorStart;
  segments[1] += operatorEnd;

  const newValue = segments.join('');
  const cursorOffset = cursorStart === cursorEnd && operatorEnd ? -1 : 0;
  const newCursorPos = cursorEnd + (newValue.length - value.length) + cursorOffset;

  inputEl.value = newValue;
  inputEl.setSelectionRange(newCursorPos, newCursorPos);
  inputEl.dispatchEvent(new Event('input'));
  inputEl.focus();
};

/**
 * buildAscNumSequence arranges points of a range so that they do no overlap
 * (builds ascending sequence of numbers without reordering the input array).
 * Use cases: range filter and range slider component
 */
export const buildAscNumSequence = (
  value: number[],
  preservedIndex: number = -1,
  min: number = -Infinity,
  max: number = Infinity
): number[] => {
  return value
    .map((v, i, arr: number[]) => (preservedIndex === i ? v : Math.min(...arr.slice(i))))
    .map((v, i, arr: number[]) => (preservedIndex === i ? v : Math.max(...arr.slice(0, i + 1))))
    .map((v: number) => Math.min(Math.max(v, min), max));
};

export const createExactAxisConstraint = (name: string, allValues: string[]): SearchAxisConstraint => {
  const values = allValues.filter((value) => value.indexOf(SEARCH_FACETS_SINGLE_NODE_VALUE_PREFIX) !== 0);
  const singleNodeValues = allValues
    .filter((value) => value.indexOf(SEARCH_FACETS_SINGLE_NODE_VALUE_PREFIX) === 0)
    .map((value) => value.replace(SEARCH_FACETS_SINGLE_NODE_VALUE_PREFIX, ''));

  return {
    type: 'exact',
    axis: name,
    combineOperator: SEARCH_FACETS_COMBINE_OPERATOR,
    values,
    ...(!!singleNodeValues.length ? { singleNodeValues } : {}),
  };
};

export const formatResultsItemsAsCsv = (items: SearchResultsItem[]) => {
  const rows = items.map(({ entityType, values }) =>
    values.concat([{ fieldName: FIELDS.entityName.name, fieldValue: entityType } as SearchResultsItemValue]).reduce(
      (items, { fieldName, fieldValue }) =>
        Object.assign(items, {
          [fieldName]: (items[fieldName] ?? []).concat([fieldValue]).filter(filterDuplicates),
        }),
      {}
    )
  );

  const fieldNames = rows
    .flatMap((row) => Object.keys(row))
    .filter(filterDuplicates)
    .sort();

  const valueCounts = fieldNames.reduce(
    (counts, fieldName) =>
      Object.assign(counts, {
        [fieldName]: rows
          .map((row) => row[fieldName]?.length ?? 0)
          .sort((a, b) => b - a)
          .shift(),
      }),
    {}
  );

  const fillByFieldName = (fieldName: string, fillFn: (index: number) => string): string[] =>
    new Array(valueCounts[fieldName]).fill(null).map((_, index) => fillFn(index));

  const header = fieldNames.flatMap((fieldName) => fillByFieldName(fieldName, (index) => `${fieldName}[${index}]`));
  const body = rows.map((row) =>
    fieldNames.flatMap((fieldName) => fillByFieldName(fieldName, (index) => row[fieldName]?.[index] ?? ''))
  );

  return [header]
    .concat(body)
    .map((row) => row.map((value) => JSON.stringify(value)).join(','))
    .join('\r\n');
};

export const getHierarchyNodesFirstLevel = (nodes: HierarchyNode[] = [], minLevel: number = 0) =>
  Math.max(
    minLevel,
    nodes.reduce((lowest, node) => (node.depth < lowest && node.depth >= minLevel ? node.depth : lowest), Infinity)
  );
