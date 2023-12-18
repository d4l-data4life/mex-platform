import { createStore } from '@stencil/store';
import { ITEM_GC_COUNT } from 'config';
import { Item } from 'services/item';

type ItemListener = (item: Item) => void;

interface StateType {
  [identifier: string]: Item;
}

const store = createStore<StateType>({});
export const history: string[] = [];

class ItemsStore {
  #listeners: { identifier: string; listener: ItemListener }[] = [];

  constructor() {
    store.on('set', (identifier: string, item: Item) =>
      this.#listeners.filter((entry) => entry.identifier === identifier).forEach(({ listener }) => listener(item))
    );
  }

  addListener(identifier: string, listener: ItemListener) {
    this.setListener(identifier, listener, false);
  }

  setListener(identifier: string, listener: ItemListener, removeOld: boolean = true) {
    removeOld && this.#listeners.some((entry) => entry.identifier === identifier) && this.removeListener(identifier);
    this.#listeners.push({ identifier, listener });
    listener(this.get(identifier));
  }

  removeListener(identifier: string) {
    this.#listeners = this.#listeners.filter((entry) => entry.identifier !== identifier);
  }

  get(identifier: string): Item {
    return store.get(identifier);
  }

  add(identifier: string, item: Item) {
    if (!this.get(identifier) && item) {
      store.set(identifier, item);
    }

    history.push(identifier);
    this.runGarbageCollection();
  }

  runGarbageCollection() {
    if (history.length <= ITEM_GC_COUNT) {
      return;
    }

    const identifier = history.shift();
    !history.includes(identifier) && store.set(identifier, null);
  }

  reset(): void {
    store.reset();
    history.forEach(() => history.shift());
  }
}

export default new ItemsStore();
