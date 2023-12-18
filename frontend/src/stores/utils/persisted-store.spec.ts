import createPersistedStore from './persisted-store';
import { Persistor } from './persistor';

interface TestStore {
  foo: string;
  bar: number;
  baz: boolean;
  foo_2?: string;
}

const namespace = 'test';
const initialData: TestStore = {
  foo: 'foo',
  bar: 42,
  baz: true,
};

describe('persisted store', () => {
  let persistor;

  beforeEach(() => {
    persistor = new Persistor(null);
  });

  it('uses initial data if persistor does not contain data', () => {
    const store = createPersistedStore<TestStore>(persistor, namespace, initialData);

    expect(store.state.foo).toBe('foo');
    expect(store.state.bar).toBe(42);
    expect(store.state.baz).toBe(true);
  });

  it('overrides initial data with data from persistor', () => {
    let store = createPersistedStore<TestStore>(persistor, namespace, initialData);

    store.set('foo', 'bar');
    store.set('bar', 0);
    store.set('baz', false);

    store = createPersistedStore<TestStore>(persistor, namespace, initialData);

    expect(store.state.foo).toBe('bar');
    expect(store.state.bar).toBe(0);
    expect(store.state.baz).toBe(false);
  });

  it('supports partial overrides from persistor', () => {
    let store = createPersistedStore<TestStore>(persistor, namespace, initialData);

    expect(persistor.getKeys()).toEqual([]);
    store.set('foo', 'bar');
    expect(persistor.getKeys()).toEqual([`${namespace}__foo`]);

    store = createPersistedStore<TestStore>(persistor, namespace, initialData);

    expect(store.state.foo).toBe('bar');
    expect(store.state.bar).toBe(42);
    expect(store.state.baz).toBe(true);
  });

  it('supports persisting / resetting with nullish values', () => {
    const store = createPersistedStore<TestStore>(persistor, namespace, initialData);

    store.set('foo', 'bar');
    expect(persistor.getKeys()).toEqual([`${namespace}__foo`]);
    store.set('foo', undefined);
    expect(store.state.foo).toBe(undefined);
    expect(persistor.getKeys()).toEqual([]);

    store.set('foo', 'bar');
    expect(persistor.getKeys()).toEqual([`${namespace}__foo`]);
    store.set('foo', null);
    expect(store.state.foo).toBe(null);
    expect(persistor.getKeys()).toEqual([]);
  });

  it('resets to the initial data', () => {
    const store = createPersistedStore<TestStore>(persistor, namespace, initialData);

    store.set('foo', 'bar');
    expect(store.state).toEqual({ ...initialData, foo: 'bar' });

    store.reset();
    expect(store.state).toEqual(initialData);
  });

  it('includes and resets even persisted non-initial data when it belongs to the same namespace', () => {
    // previously persisted data
    persistor.set(`${namespace}__foo_3`, '"bar_3"');

    const store = createPersistedStore<TestStore>(persistor, namespace, initialData);
    store.set('foo_2', 'bar_2');
    // foo_3 is also included in the state, because it was persisted with the same namespace
    expect(store.state).toEqual({ ...initialData, foo_2: 'bar_2', foo_3: 'bar_3' });
    // both new keys are persisted
    expect(persistor.getKeys()).toEqual([`${namespace}__foo_3`, `${namespace}__foo_2`]);

    store.reset();
    // new keys are removed both from the state and the storage
    expect(store.state).toEqual(initialData);
    expect(persistor.getKeys()).toEqual([]);
  });

  it('does not include data from a foreign namespace', () => {
    persistor.set('storea__foo', '1');
    persistor.set('storea__fau', '"bar"');
    persistor.set('storeb__foo', '9');

    const storeA = createPersistedStore<any>(persistor, 'storea', {});
    storeA.set('bar', 'foo');
    storeA.set('baz', 'baz');

    const storeB = createPersistedStore<any>(persistor, 'storeb', {});
    storeB.set('bar', 'baz');

    expect(storeA.get('foo')).toBe(1);
    expect(storeA.get('bar')).toBe('foo');
    expect(storeA.get('baz')).toBe('baz');
    expect(storeA.get('fau')).toBe('bar');

    expect(storeB.get('foo')).toBe(9);
    expect(storeB.get('bar')).toBe('baz');
    expect(storeB.get('baz')).toBe(undefined);
    expect(storeB.get('fau')).toBe(undefined);
  });
});
