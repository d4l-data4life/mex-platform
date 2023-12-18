import { createStore } from '@stencil/store';
import { SEARCH_CONFIG, SEARCH_PAGINATION_START } from 'config';
import { EntityTypeName } from 'config/entity-types';
import { SearchSorting } from 'config/search';
import { Field } from 'models/field';

export interface SearchResultsItemValue {
  fieldName: string;
  fieldValue: string;
  language?: string;
}

export interface SearchResultsItem {
  itemId: string;
  entityType: EntityTypeName;
  createdAt?: string;
  values: SearchResultsItemValue[];
}

export interface SearchResultsFacetBucketHierarchyInfo {
  '@type': string;
  parentValue: string;
  display: string;
  depth: number;
}

export interface SearchResultsFacetBucket {
  count: number;
  value: string;
  hierarchyInfo?: SearchResultsFacetBucketHierarchyInfo;
}

export interface SearchResultsFacet {
  type: 'exact' | 'yearRange';
  axis: string;
  bucketNo: number;
  buckets: SearchResultsFacetBucket[];
}

export interface SearchResultsHighlightMatch {
  fieldName: string;
  snippets: string[];
  language?: string;
}

export interface SearchResultsHighlight {
  itemId: string;
  matches: SearchResultsHighlightMatch[];
}

export interface SearchResultsDiagnostics {
  cleanedQuery: string;
  queryWasCleaned: boolean;
}

export interface SearchResults {
  numFound: number;
  numFoundExact: boolean;
  start: number;
  maxScore: number;
  items: SearchResultsItem[];
  facets?: SearchResultsFacet[];
  highlights?: SearchResultsHighlight[];
  diagnostics?: SearchResultsDiagnostics;
}

interface StateType {
  isBusy: boolean;
  query: string;
  limit: number;
  response?: SearchResults;
  facets?: SearchResultsFacet[];
  sorting?: SearchSorting;
  focus?: string;
}

const store = createStore<StateType>({
  isBusy: true,
  query: '',
  limit: SEARCH_CONFIG.LIMIT,
});

export interface HierarchyNodeDisplayValue {
  place?: number;
  display: string;
  language?: string;
}

export interface HierarchyNode {
  depth: number;
  nodeId: string;
  display?: HierarchyNodeDisplayValue[];
  parentNodeId?: string;
}

export interface Hierarchy {
  key: Field;
  nodes: HierarchyNode[];
}

const hierarchyStore = createStore<{ [key: string]: Hierarchy }>({});

class SearchStore {
  get isBusy() {
    return store.get('isBusy') || !this.response;
  }

  set isBusy(isBusy: boolean) {
    store.set('isBusy', isBusy);
  }

  get query() {
    return store.get('query');
  }

  set query(query: string) {
    store.set('query', query);
  }

  get offset() {
    return this.response?.start ?? 0;
  }

  get limit() {
    return store.get('limit');
  }

  set limit(limit: number) {
    store.set('limit', limit);
  }

  get count() {
    return this.response?.numFound;
  }

  get placeholdersCount() {
    return Math.min(this.limit, (this.count ?? Infinity) - this.offset);
  }

  get page() {
    const min = SEARCH_PAGINATION_START;
    return Math.floor(this.offset / this.limit) + min;
  }

  get pageRange() {
    const { page } = this;

    const min = SEARCH_PAGINATION_START;
    const max = Math.max(Math.ceil((this.count ?? 0) / this.limit) + (min - 1), min);
    const offset = Math.floor((SEARCH_CONFIG.PAGINATION_RANGE_COUNT - 1) / 2);
    const start = Math.max(min, page - offset - Math.max(0, offset - (max - page)));
    const end = Math.min(max, start + (SEARCH_CONFIG.PAGINATION_RANGE_COUNT - 1));

    return new Array(end - start + 1).fill(null).map((_, index) => start + index);
  }

  get response() {
    return store.get('response');
  }

  set response(response: SearchResults) {
    store.set('response', response);
  }

  get facets() {
    return store.get('facets') ?? [];
  }

  set facets(facets: SearchResultsFacet[]) {
    store.set('facets', facets);
  }

  addFacetBuckets(axis: string, newBuckets: SearchResultsFacetBucket[]) {
    const { facets } = this;
    const match = facets.find((facet) => facet.axis === axis);

    store.set(
      'facets',
      facets.map((facet) => {
        if (facet !== match) {
          return facet;
        }

        const buckets = facet.buckets.concat(newBuckets ?? []);
        const values = buckets.map(({ value }) => value);
        return { ...facet, buckets: buckets.filter((bucket, index) => values.indexOf(bucket.value) === index) };
      })
    );
  }

  getFacetOffset(axis: string) {
    return this.facets.find((facet) => facet.axis === axis)?.buckets.length ?? 0;
  }

  hasMoreBuckets(axis: string) {
    return this.getFacetOffset(axis) < this.facets?.find((facet) => facet.axis === axis)?.bucketNo;
  }

  get focus() {
    return store.get('focus') ?? null;
  }

  set focus(focus: string) {
    store.set('focus', focus);
  }

  get sorting() {
    return store.get('sorting') ?? SEARCH_CONFIG.SORTING_OPTIONS[0];
  }

  set sorting(sorting: SearchSorting) {
    store.set('sorting', sorting);
  }

  get highlights() {
    return this.response?.highlights;
  }

  get wasQueryCleaned() {
    return !!this.response?.diagnostics?.queryWasCleaned;
  }

  get cleanedQuery() {
    return this.response?.diagnostics?.cleanedQuery;
  }

  get hierarchies() {
    return hierarchyStore.state;
  }
}

export default new SearchStore();
