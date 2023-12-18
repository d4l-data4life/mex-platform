import { createStore } from '@stencil/store';
import { FIELDS, FIELD_CONCEPT_IDS, SEARCH_CONFIG, SEARCH_INVISIBLE_FACETS } from 'config';
import { getFilterQueryParamName } from 'utils/search';
import { denormalizeConceptId, normalizeConceptId } from 'utils/field';

export type FilterItems = [name: string, values: string[]][];

type FilterListener = (values?: string[]) => void;

interface StateType {
  [name: string]: string[];
}

const store = createStore<StateType>({});

class FiltersStore {
  #listeners: { name: string; listener: FilterListener }[] = [];

  constructor() {
    store.on('set', (name: string, values: string[]) =>
      this.#listeners.filter((item) => item.name === name).forEach(({ listener }) => listener(values))
    );
  }

  get all(): FilterItems {
    return Object.entries({ ...store.state }).filter(([_, values]) => values.length);
  }

  get isEmpty() {
    return !this.all.length;
  }

  addListener(name: string, listener: FilterListener) {
    this.#listeners.push({ name, listener });
  }

  removeListener(listener: FilterListener) {
    this.#listeners = this.#listeners.filter((item) => item.listener !== listener);
  }

  get(name: string): string[] {
    return store.get(name) ?? [];
  }

  set(name: string, values: string[]) {
    store.set(name, values);
  }

  add(name: string, value: string): void {
    store.set(
      name,
      this.get(name)
        .concat([value])
        .filter(Boolean)
        .filter((item, index, arr) => arr.indexOf(item) === index)
    );
  }

  remove(name: string, value: string): void {
    store.set(
      name,
      this.get(name).filter((item) => item !== value)
    );
  }

  reset(name?: string): void {
    if (name) {
      store.set(name, []);
    } else {
      store.reset();
    }
  }

  toQueryParams(): URLSearchParams {
    const params = new URLSearchParams();
    this.all.forEach(([name, values]) =>
      values.forEach((value) => params.append(getFilterQueryParamName(name), normalizeConceptId(value)))
    );

    return params;
  }

  fromQueryParams(params: URLSearchParams): void {
    store.reset();

    const conceptIdValues = FIELD_CONCEPT_IDS().map(normalizeConceptId);

    SEARCH_CONFIG.FACETS.concat(SEARCH_INVISIBLE_FACETS()).forEach(({ axis: { name } }) =>
      store.set(
        name,
        params
          .getAll(getFilterQueryParamName(name))
          .filter(Boolean)
          .map((value) => (conceptIdValues.includes(value) ? denormalizeConceptId(value, FIELDS[name]) : value))
      )
    );
  }
}

export default new FiltersStore();
