jest.mock('stencil-router-v2');

import items, { history } from './items';
import * as config from 'config';

let itemCounter = 0;
const generateItem = () => {
  itemCounter++;
  return {
    itemId: `item-${itemCounter}`,
    businessId: `business-item-${itemCounter}`,
    entityType: 'Source',
    owner: null,
    createdAt: new Date().toISOString(),
    values: [],
  };
};

describe('items store', () => {
  beforeEach(() => {
    items.reset();
  });

  it('adds and gets an item to the store by identifier', () => {
    const item = generateItem();

    expect(items.get('foo-identifier')).toBe(undefined);
    items.add('foo-identifier', item);
    expect(items.get('foo-identifier')).toBe(item);
  });

  it('does not add item if empty or if an item was already added for the same identifier', () => {
    const item1 = generateItem();
    const item2 = generateItem();

    items.add('foo-identifier', null);
    expect(items.get('foo-identifier')).toBe(undefined);

    items.add('bar-identifier', item1);
    expect(items.get('bar-identifier')).toBe(item1);

    items.add('bar-identifier', item2);
    expect(items.get('bar-identifier')).toBe(item1);
  });

  it('runs the garbage collection to free up older cached items', () => {
    jest.spyOn(config, 'ITEM_GC_COUNT', 'get').mockReturnValue(2);

    const item1 = generateItem();
    const item2 = generateItem();
    const item3 = generateItem();
    const item4 = generateItem();

    items.reset();
    expect(history).toEqual([]);

    items.add('foo', item1);
    items.add('bar', item2);

    expect(items.get('foo')).toBe(item1);
    expect(history).toEqual(['foo', 'bar']);

    items.add('baz', item3);

    expect(items.get('foo')).toBe(null);
    expect(items.get('bar')).toBe(item2);
    expect(items.get('baz')).toBe(item3);
    expect(history).toEqual(['bar', 'baz']);

    items.add('foo', item4);
    expect(items.get('foo')).toBe(item4);
    expect(items.get('bar')).toBe(null);
    expect(items.get('baz')).toBe(item3);
    expect(history).toEqual(['baz', 'foo']);
  });

  it('resets the store', () => {
    const item = generateItem();

    items.reset();
    expect(history).toEqual([]);

    items.add('foo', item);
    expect(history).toEqual(['foo']);

    items.reset();
    expect(history).toEqual([]);
  });

  describe('setListener()', () => {
    it('sets a listener to be called with the initial value and when the store value for the given identifier has changed', () => {
      const listener = jest.fn();
      const item1 = generateItem();
      const item2 = generateItem();

      items.setListener('foo', listener);
      expect(listener).toHaveBeenCalledWith(undefined);

      items.add('foo', item1);
      items.add('bar', item2);
      expect(listener).toHaveBeenCalledWith(item1);
      expect(listener).not.toHaveBeenCalledWith(item2);
    });

    it('removes a previously set listener for the same identifier', () => {
      const listener1 = jest.fn();
      const listener2 = jest.fn();
      const item = generateItem();

      items.setListener('foo', listener1);
      items.setListener('foo', listener2);
      items.add('foo', item);

      expect(listener1).toHaveBeenCalledTimes(1);
      expect(listener2).toHaveBeenCalledTimes(2);
    });
  });

  describe('addListener()', () => {
    it('adds a listener (same behavior as setListener(), but does not remove previously set listeners for the same identifier)', () => {
      const listener1 = jest.fn();
      const listener2 = jest.fn();
      const listener3 = jest.fn();
      const item1 = generateItem();
      const item2 = generateItem();

      items.setListener('foo', listener1);
      items.addListener('foo', listener2);
      items.add('foo', item1);
      items.addListener('foo', listener3);
      items.add('foo', item2);

      expect(listener1).toHaveBeenCalledTimes(2);
      expect(listener2).toHaveBeenCalledTimes(2);
      expect(listener3).toHaveBeenCalledTimes(1);
    });
  });

  describe('removeListener()', () => {
    it('removes a listener for given identifier', () => {
      const listener = jest.fn();
      const item = generateItem();

      items.setListener('foo', listener);
      items.removeListener('foo');
      items.add('foo', item);
      expect(listener).not.toHaveBeenCalledWith(item);
    });
  });
});
