import notifications from './notifications';

describe('notifications store', () => {
  it('gets all items and adds an item, preventing duplicates', () => {
    expect(notifications.items).toEqual([]);
    notifications.add('foo');
    notifications.add('bar');
    notifications.add('foo');
    expect(notifications.items).toEqual(['foo', 'bar']);
  });

  it('adds a timeout to remove added notifications after 8s', () => {
    const timeoutSpy = jest.spyOn(window, 'setTimeout');
    expect(timeoutSpy).not.toHaveBeenCalled();

    notifications.add('foo');
    expect(timeoutSpy).toHaveBeenCalledWith(expect.anything(), 8000);
  });

  it('resets to the pristine state of the store', () => {
    expect(notifications.items.length).toBe(2);
    notifications.reset();
    expect(notifications.items.length).toBe(0);
  });
});
