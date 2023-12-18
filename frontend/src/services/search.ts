import {
  API_URL,
  FIELDS,
  SEARCH_CONFIG,
  SEARCH_FACETS_COMBINE_OPERATOR,
  SEARCH_FACETS_MAX_LIMIT,
  SEARCH_INVISIBLE_FACETS,
  SEARCH_QUERY_EVERYTHING,
  SEARCH_TIMEOUT,
} from 'config';
import { EntityTypeName } from 'config/entity-types';
import { SearchFacet, SearchAxisConstraint, SearchSorting } from 'config/search';
import { Field, FieldRenderer } from 'models/field';
import stores from 'stores';
import { FilterItems } from 'stores/filters';
import { Hierarchy, HierarchyNode, SearchResults } from 'stores/search';
import { post } from 'utils/fetch-client';
import { filterDuplicates } from 'utils/field';
import { createExactAxisConstraint } from 'utils/search';

interface HierarchyResponse {
  nodes: HierarchyNode[];
}

export class SearchService {
  private getAxisConstraint(name: string, values: string[]) {
    const facet = SEARCH_CONFIG.FACETS.concat(SEARCH_INVISIBLE_FACETS()).find((facet) => facet.axis.name === name);

    switch (facet?.type) {
      case 'exact':
      case 'hierarchy':
        return createExactAxisConstraint(name, values);

      case 'yearRange':
        try {
          const [minYear, maxYear] = values[0].split('-');
          const minDate = new Date(`${minYear}-01-01T00:00:00Z`);
          const maxDate = new Date(`${maxYear}-12-31T23:59:59Z`);

          if (isNaN(minDate.getTime()) || isNaN(maxDate.getTime()) || minDate > maxDate) {
            throw new Error('Invalid range');
          }

          return {
            type: 'stringRange',
            axis: name,
            combineOperator: SEARCH_FACETS_COMBINE_OPERATOR,
            stringRanges: [
              {
                min: minDate.toISOString(),
                max: maxDate.toISOString(),
              },
            ],
          };
        } catch (_) {
          return null;
        }
      default:
        return null;
    }
  }

  private getFacetConfig({ axis: { name: axis }, type }: SearchFacet, offset: number = 0, limit?: number) {
    switch (type) {
      case 'exact':
        return { axis, type, offset, limit };
      case 'hierarchy':
        return { axis, type: 'exact', offset, limit: SEARCH_FACETS_MAX_LIMIT };
      default:
        return { axis, type };
    }
  }

  private getSearchFields(): Field[] {
    const rawFields = Object.values(FIELDS)
      .filter((field) => field.renderer === FieldRenderer.time)
      .map((field) => field.cloneAndSetRawValueMode());

    return Object.values(FIELDS)
      .concat(rawFields)
      .concat(Object.values(SEARCH_CONFIG.DISPLAYED_FIELDS).flat())
      .flatMap((field) => field.resolvesTo)
      .filter((field, index, arr) => arr.indexOf(field) === index);
  }

  async updateResults({ query = SEARCH_QUERY_EVERYTHING, offset = 0, reset = false }): Promise<void> {
    if (reset) {
      stores.search.response = null;
      stores.search.facets = null;
    }

    stores.search.isBusy = true;

    const response = await this.fetchResults({
      query,
      offset,
      filters: stores.filters.all,
      searchFocus: stores.search.focus,
    });
    stores.search.response = response;
    stores.search.facets = response?.facets;
    stores.search.isBusy = false;
  }

  async fetchResults({
    query,
    offset = 0,
    limit = null,
    axisConstraints = null,
    filters = null,
    facets = SEARCH_CONFIG.FACETS,
    facetsOffset = 0,
    facetsLimit = SEARCH_CONFIG.FACETS_LIMIT,
    fields = this.getSearchFields(),
    highlightFields = null,
    searchFocus = null,
    maxEditDistance = SEARCH_CONFIG.QUERY_MAX_EDIT_DISTANCE,
    useNgramField = SEARCH_CONFIG.QUERY_USE_NGRAM_FIELD,
    sorting = stores.search.sorting,
    timeout = SEARCH_TIMEOUT,
  }: {
    query: string;
    offset?: number;
    limit?: number;
    axisConstraints?: SearchAxisConstraint[];
    filters?: FilterItems;
    facets?: SearchFacet[];
    facetsOffset?: number;
    facetsLimit?: number;
    fields?: Field[];
    highlightFields?: Field[];
    searchFocus?: string;
    maxEditDistance?: number;
    useNgramField?: boolean;
    sorting?: SearchSorting;
    timeout?: number;
  }): Promise<SearchResults> {
    const url = `${API_URL}/query/search`;
    const [response] = await post<SearchResults>({
      url,
      timeout,
      body: JSON.stringify({
        query,
        offset,
        limit: limit ?? stores.search.limit,
        sorting: sorting?.axis ? { ...sorting, axis: sorting.axis.name } : null,
        facets: facets.map((facet) => this.getFacetConfig(facet, facetsOffset, facetsLimit)),
        fields: fields.map(({ linkedName }) => linkedName).filter(filterDuplicates),
        searchFocus,
        maxEditDistance,
        useNgramField,
        ...(highlightFields
          ? { highlightFields: highlightFields.map(({ linkedName }) => linkedName) }
          : { autoHighlight: true }),
        ...(axisConstraints || filters?.length
          ? {
              axisConstraints:
                axisConstraints ??
                filters.map(([name, values]) => this.getAxisConstraint(name, values)).filter(Boolean),
            }
          : {}),
      }),
      authorized: true,
    });

    return response;
  }

  async fetchHierarchy(
    key: Field,
    entityType: EntityTypeName,
    linkField: Field,
    displayField: Field,
    displayFieldLanguage?: string
  ): Promise<Hierarchy> {
    const { hierarchies } = stores.search;
    if (hierarchies[key.name]) {
      return hierarchies[key.name];
    }

    const url = `${API_URL}/metadata/tree`;
    const [response] = await post<HierarchyResponse>({
      url,
      body: JSON.stringify({
        nodeEntityType: entityType,
        linkFieldName: linkField.name,
        displayFieldName: displayField.name,
        displayFieldLanguage,
      }),
      authorized: true,
    });

    const hierarchy = { key, nodes: response?.nodes ?? [] };
    hierarchies[key.name] = hierarchy;
    return hierarchy;
  }
}

export default new SearchService();
