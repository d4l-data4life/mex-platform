jest.mock('stencil-router-v2');

import * as config from 'config';
import { Field, FieldImportance, FieldRenderer } from 'models/field';
import filters from './filters';

const { SEARCH_CONFIG, FIELDS } = config;

FIELDS.author = new Field({
  name: 'author',
  renderer: FieldRenderer.plain,
  importance: FieldImportance.recommended,
  isEnumerable: true,
  isVirtual: false,
});

FIELDS.keyword = new Field({
  name: 'keyword',
  renderer: FieldRenderer.plain,
  importance: FieldImportance.optional,
  isEnumerable: false,
  isVirtual: false,
});

FIELDS.invisible = new Field({
  name: 'invisible',
  renderer: FieldRenderer.none,
  importance: FieldImportance.none,
  isEnumerable: false,
  isVirtual: false,
});

describe('filters store', () => {
  it('has an isEmpty getter (bool)', () => {
    expect(typeof filters.isEmpty).toBe('boolean');
  });

  it('has a getter for all items, returning an array of filter items', () => {
    expect(filters.all).toBeInstanceOf(Array);
  });

  it('adds a value to a filter item by field name', () => {
    expect(filters.isEmpty).toBe(true);
    expect(filters.all.length).toBe(0);

    filters.add('author', 'foo');
    expect(filters.isEmpty).toBe(false);
    expect(filters.all.length).toBe(1);

    filters.add('author', 'bar');
    expect(filters.all.length).toBe(1);

    filters.add('keyword', 'baz');
    expect(filters.all.length).toBe(2);
  });

  it('gets filter values by field name', () => {
    expect(filters.isEmpty).toBe(false);
    expect(filters.get('author')).toEqual(['foo', 'bar']);
  });

  it('does not allow to add duplicated values to a filter item', () => {
    expect(filters.get('keyword')).toEqual(['baz']);
    filters.add('keyword', 'baz');
    expect(filters.get('keyword')).toEqual(['baz']);
  });

  it('removes a value from a filter item by field name', () => {
    expect(filters.get('author')).toEqual(['foo', 'bar']);
    expect(filters.get('keyword')).toEqual(['baz']);

    filters.remove('author', 'foo');
    expect(filters.get('author')).toEqual(['bar']);

    filters.remove('keyword', 'baz');
    expect(filters.get('keyword')).toEqual([]);
  });

  it('resets to the pristine state of the store', () => {
    expect(filters.isEmpty).toBe(false);

    filters.reset();
    expect(filters.isEmpty).toBe(true);
    expect(filters.all.length).toBe(0);
  });

  it('resets an individual filter item by field name', () => {
    filters.add('keyword', 'foo');
    filters.add('keyword', 'bar');
    expect(filters.get('keyword')).toEqual(['foo', 'bar']);
    filters.reset('keyword');
    expect(filters.get('keyword')).toEqual([]);
  });

  it('exports the filters to query params', () => {
    jest.spyOn(config, 'SEARCH_PARAM_FILTER_PREFIX', 'get').mockReturnValue('');

    filters.reset();
    filters.add('author', 'foo');
    filters.add('author', 'bar');
    filters.add('keyword', 'baz');
    const queryParams = filters.toQueryParams();
    expect(queryParams.toString()).toBe('author%5B%5D=foo&author%5B%5D=bar&keyword%5B%5D=baz');
  });

  it('normalizes concept IDs on export', () => {
    jest.spyOn(config, 'SEARCH_PARAM_FILTER_PREFIX', 'get').mockReturnValue('');
    jest.spyOn(config, 'FIELD_CONCEPT_PREFIXES').mockReturnValue(['https://unit.test/concepts/']);

    filters.reset();
    filters.add('author', 'foo');
    filters.add('keyword', 'https://unit.test/concepts/baz');
    const queryParams = filters.toQueryParams();
    expect(queryParams.toString()).toBe('author%5B%5D=foo&keyword%5B%5D=baz');
  });

  it('imports query params to filters, allowing only known ones (including invisible/internal options)', () => {
    jest.spyOn(config, 'SEARCH_PARAM_FILTER_PREFIX', 'get').mockReturnValue('f.');
    SEARCH_CONFIG.FACETS = [
      { axis: { name: 'author' }, type: 'exact' },
      { axis: { name: 'keyword' }, type: 'exact' },
    ];
    jest.spyOn(config, 'SEARCH_INVISIBLE_FACETS').mockReturnValueOnce([{ axis: { name: 'invisible' }, type: 'exact' }]);

    filters.reset();

    filters.add('author', 'foo');
    filters.add('author', 'foo 2.0');
    filters.add('keyword', 'bar');
    filters.add('invisible', 'baz');

    expect(filters.get('author')).toEqual(['foo', 'foo 2.0']);
    expect(filters.get('keyword')).toEqual(['bar']);
    expect(filters.get('invisible')).toEqual(['baz']);

    const queryParams = new URLSearchParams(
      'f.author[]=foo&f.author[]=bar&f.keyword[]=baz&f.unknownfield[]=evil&f.invisible[]=invisible+option'
    );
    filters.fromQueryParams(queryParams);

    expect(filters.get('author')).toEqual(['foo', 'bar']);
    expect(filters.get('keyword')).toEqual(['baz']);
    expect(filters.get('unknownfield')).toEqual([]);
    expect(filters.get('invisible')).toEqual(['invisible option']);
  });

  it('denormalizes known concept IDs on import', () => {
    config.FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX = 'https://unit.test/concepts/';
    jest
      .spyOn(config, 'FIELD_CONCEPT_IDS')
      .mockReturnValue(['https://unit.test/concepts/bar', 'https://other.test/concepts/baz']);
    jest.spyOn(config, 'SEARCH_PARAM_FILTER_PREFIX', 'get').mockReturnValue('f.');
    SEARCH_CONFIG.FACETS = [{ axis: { name: 'author' }, type: 'exact' }];

    const queryParams = new URLSearchParams('f.author[]=foo&f.author[]=bar&f.author[]=baz');
    filters.fromQueryParams(queryParams);

    expect(filters.get('author')).toEqual(['foo', 'https://unit.test/concepts/bar', 'baz']);
  });
});
