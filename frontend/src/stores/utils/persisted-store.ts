import { createStore, ObservableMap } from '@stencil/store';
import { Persistor } from './persistor';

export default function createPersistedStore<T>(
  persistor: Persistor,
  namespace: string,
  initialData: T
): ObservableMap<T> {
  function namespacedKey(key: string): string {
    return `${namespace}__${key}`;
  }

  function getPersistedKeys() {
    return persistor
      .getKeys()
      .filter((key) => key.indexOf(prefix) === 0)
      .map((key) => key.replace(prefix, ''));
  }

  const prefix = namespacedKey('');
  const persistedData = getPersistedKeys()
    .concat(Object.keys(initialData))
    .filter((key, index, arr) => arr.indexOf(key) === index)
    .reduce((data, key) => {
      try {
        const rawValue = persistor.get(namespacedKey(key));
        const value = typeof rawValue === 'string' ? JSON.parse(rawValue) : rawValue;
        if (value !== null && value !== undefined) {
          return Object.assign(data, { [key]: value });
        }
      } catch (_) {}
      return data;
    }, {});

  const store = createStore<T>(initialData);
  Object.keys(persistedData).forEach((key) => store.set(key as keyof T & string, persistedData[key]));

  const actions = {
    set(key: any, value: any) {
      try {
        persistor.set(namespacedKey(key), JSON.stringify(value));
      } catch (_) {}
    },
    reset() {
      getPersistedKeys().forEach((key) => {
        actions.set(key, initialData[key]);
      });
    },
  };

  store.use(actions);

  return store;
}
