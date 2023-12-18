import { createStore } from '@stencil/store';

interface StateType {
  notifications: string[];
}

const store = createStore<StateType>({
  notifications: [],
});

class NotificationsStore {
  get items() {
    return store.get('notifications');
  }

  add(key: string) {
    store.set(
      'notifications',
      this.items.concat([key]).filter((item, index, arr) => arr.indexOf(item) === index)
    );

    window.setTimeout(() => {
      store.set(
        'notifications',
        this.items.filter((item) => item !== key)
      );
    }, 8000);
  }

  reset(): void {
    store.reset();
  }
}

export default new NotificationsStore();
