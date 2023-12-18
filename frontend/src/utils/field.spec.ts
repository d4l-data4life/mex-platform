jest.mock('stencil-router-v2');

import * as config from 'config';
import { EntityTypeName } from 'config/entity-types';
import { FieldValueDatePrecisionLevel } from 'config/fields';
import { SidebarFeature } from 'config/item';
import { Field, FieldEntityVirtualType, FieldImportance, FieldRenderer } from 'models/field';
import { Item } from 'services/item';
import stores from 'stores';
import {
  translateFieldName,
  translateFieldDescription,
  translateFieldValue,
  translateFieldValueDescription,
  getDisplayValues,
  getConcatDisplayValue,
  getTextMatches,
  getValue,
  getValues,
  formatValue,
  normalizeKey,
  formatDate,
  hasValueChanged,
  calculateCompleteness,
  hasValue,
  normalizeConceptId,
  denormalizeConceptId,
  getInvolvedFields,
  getRawValues,
  getConcatenator,
  filterDuplicates,
  getLangAttrIfForeign,
  sortAndFilterValues,
  getDocumentationProp,
} from './field';

const { FIELDS, ITEM_CONFIG, SEARCH_CONFIG } = config;

const ITEM_BASE = {
  itemId: 'item-foo',
  entityType: 'Resource',
};
const AUTHOR_1 = {
  fieldName: 'author',
  fieldValue: 'Jane',
};
const AUTHOR_2 = {
  fieldName: 'author',
  fieldValue: 'Doe',
};
const AUTHOR_3 = {
  fieldName: 'author',
  fieldValue: 'John',
};
const LABEL = {
  fieldName: 'label',
  fieldValue: 'Foo Bar Baz with John Doe.',
  language: 'en',
};
const TIME = {
  fieldName: 'time',
  fieldValue: '2012-08-17',
};
const CONCEPT = {
  fieldName: 'concept_field',
  fieldValue: 'https://concept.foo/value-concept-foo',
};

const ITEM = { ...ITEM_BASE, values: [AUTHOR_1, AUTHOR_2, LABEL, CONCEPT, TIME] } as Item;
const HIGHLIGHTS = [
  {
    itemId: 'baz',
    matches: [
      {
        fieldName: 'author',
        snippets: ['Thomas \ue000Doe\ue001'],
      },
    ],
  },
  {
    itemId: 'item-foo',
    matches: [
      {
        fieldName: 'author',
        snippets: ['\ue000Doe\ue001'],
      },
      {
        fieldName: 'label',
        snippets: ['Baz with John \ue000Doe\ue001'],
        language: 'en',
      },
      {
        fieldName: 'label',
        snippets: ['Baz mit John \ue000Doe\ue001'],
        language: 'de',
      },
    ],
  },
];

FIELDS.entityName = new Field({
  name: 'entityName',
  renderer: FieldRenderer.none,
  importance: FieldImportance.mandatory,
  isEnumerable: true,
  isVirtual: false,
});

FIELDS.author = new Field({
  name: 'author',
  renderer: FieldRenderer.plain,
  importance: FieldImportance.recommended,
  isEnumerable: false,
  isVirtual: false,
});

FIELDS.label = new Field({
  name: 'label',
  renderer: FieldRenderer.title,
  importance: FieldImportance.mandatory,
  isEnumerable: false,
  isVirtual: false,
});

FIELDS.email = new Field({
  name: 'email',
  renderer: FieldRenderer.plain,
  importance: FieldImportance.optional,
  isEnumerable: false,
  isVirtual: false,
});

FIELDS.time = new Field({
  name: 'time',
  renderer: FieldRenderer.time,
  importance: FieldImportance.recommended,
  isEnumerable: false,
  isVirtual: false,
});

FIELDS.virtual = new Field({
  name: 'virtual',
  renderer: FieldRenderer.none,
  importance: FieldImportance.none,
  isEnumerable: false,
  isVirtual: true,
  resolvesTo: ['my_awesome_field', 'time'],
});

FIELDS.myAwesomeField = new Field({
  name: 'my_awesome_field',
  renderer: FieldRenderer.none,
  importance: FieldImportance.mandatory,
  isEnumerable: false,
  isVirtual: false,
});

FIELDS.conceptField = new Field({
  name: 'concept_field',
  renderer: FieldRenderer.plain,
  importance: FieldImportance.none,
  isEnumerable: true,
  isVirtual: false,
  vocabulary: {
    de: [
      {
        conceptId: 'foo-id',
        label: 'Foo label',
        description: '',
      },
    ],
    en: [
      {
        conceptId: 'https://concept-prefix/bar-id',
        label: 'Bar label',
        description: '',
      },
    ],
  },
});

FIELDS.documentedField = new Field({
  name: 'documented_field',
  renderer: FieldRenderer.plain,
  importance: FieldImportance.none,
  isEnumerable: true,
  isVirtual: false,
  documentation: {
    de: {
      descriptionText: [
        { entityType: FieldEntityVirtualType.all, text: 'generic description text DE' },
        { entityType: 'Resource', text: 'resource description text DE' },
      ],
      exampleValue: [{ entityType: 'Resource', text: 'resource example value DE' }],
      furtherInformation: [],
    },
    en: {
      exampleValue: [{ entityType: FieldEntityVirtualType.all, text: 'generic example value EN' }],
      furtherInformation: null,
      displayFormats: [
        { entityType: FieldEntityVirtualType.all, text: 'generic display formats EN' },
        { entityType: 'Source', text: 'source display formats EN' },
      ],
    },
  },
});

const createSidebarFeatureConfig = (entityType: string, feature: SidebarFeature, field: Field) => ({
  entityType: entityType as EntityTypeName,
  feature,
  field,
});

const createEntityTypeConfig = (name: string) => ({
  name,
  config: {
    isFocal: true,
    isAggregatable: false,
  },
});

describe('field util', () => {
  let translationsSpy, languageSpy, highlightsSpy, concatenationSpy;

  beforeAll(() => {
    translationsSpy = jest.spyOn(stores.i18n, 't');
    languageSpy = jest.spyOn(stores.i18n, 'language', 'get');
    highlightsSpy = jest.spyOn(stores.search, 'highlights', 'get');
    concatenationSpy = jest.spyOn(config, 'ITEM_FIELDS_CONCATENATOR', 'get');

    config.ENTITY_TYPES.Resource = createEntityTypeConfig('Resource');
    config.ENTITY_TYPES.Source = createEntityTypeConfig('Source');
    config.ENTITY_TYPES.Person = createEntityTypeConfig('Person');
  });

  beforeEach(() => {
    jest.clearAllMocks();
    concatenationSpy.mockReturnValue({
      [FieldRenderer.time]: ' - ',
      default: ', ',
    });
  });

  describe('normalizeKey()', () => {
    it('removes colons from the field name', () => expect(normalizeKey('team:name')).toBe('teamName'));
    it('converts snake_case to camelCase field names and forces first character to be lowercase', () =>
      expect(normalizeKey('My_awesome_field')).toBe('myAwesomeField'));
    it('removes the target field name prefix from linked fields', () => {
      expect(normalizeKey('fooField__label')).toBe('fooFieldLabel');
      expect(normalizeKey('fooField__barField')).toBe('fooFieldBarField');
    });
  });

  describe('normalizeConceptId()', () => {
    it('removes the concept prefix from the value if exists', () => {
      jest
        .spyOn(config, 'FIELD_CONCEPT_PREFIXES')
        .mockReturnValue(['https://concept.foo/', 'https://concept-two.foo/']);
      expect(normalizeConceptId('value-foo')).toBe('value-foo');
      expect(normalizeConceptId('https://concept.foo/value-bar')).toBe('value-bar');
      expect(normalizeConceptId('https://concept-two.foo/value-baz')).toBe('value-baz');
      expect(normalizeConceptId(null)).toBe(null);
    });
  });

  describe('denormalizeConceptId()', () => {
    it('adds the concept prefix to the value', () => {
      config.FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX = 'https://concept.bar/';

      expect(denormalizeConceptId('foo')).toBe('https://concept.bar/foo');
      expect(denormalizeConceptId('bar')).toBe('https://concept.bar/bar');
      expect(denormalizeConceptId(null)).toBe(null);
    });

    it('does not add a concept prefix if a field is given and it is not enumerable', () => {
      config.FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX = 'https://concept.bar/';
      expect(denormalizeConceptId('foo', FIELDS.label)).toBe('foo');
    });

    it('adds a config prefix if a field is given and enumerable (prefix according to match, default: DEFAULT_FIELD_CONCEPT_PREFIX)', () => {
      config.FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX = 'https://concept.bar/';
      expect(denormalizeConceptId('foo', FIELDS.conceptField)).toBe('https://concept.bar/foo');
      expect(denormalizeConceptId('foo-id', FIELDS.conceptField)).toBe('https://concept.bar/foo-id');

      jest
        .spyOn(config, 'FIELD_CONCEPT_PREFIXES')
        .mockReturnValue([config.FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX, 'https://concept-prefix/']);
      expect(denormalizeConceptId('bar-id', FIELDS.conceptField)).toBe('https://concept-prefix/bar-id');

      jest.spyOn(config, 'FIELD_CONCEPT_PREFIXES').mockReturnValue([config.FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX]);
      expect(denormalizeConceptId('bar-id', FIELDS.conceptField)).toBe('https://concept.bar/bar-id');
    });
  });

  describe('translateFieldName()', () => {
    it('detects pluralization by looking for values and returns translated field name', () => {
      translationsSpy.mockImplementation((key) => `translated ${key}`);

      let result = translateFieldName(FIELDS.author, { ...ITEM_BASE, values: [AUTHOR_1] });
      expect(result).toEqual(expect.stringContaining('translated fields.labels.author'));
      expect(result).toEqual(expect.stringContaining('singular'));

      expect(translationsSpy).toHaveBeenCalledWith(expect.stringContaining('fields.labels.author'));
      translationsSpy.mockClear();

      result = translateFieldName(FIELDS.author, { ...ITEM_BASE, values: [AUTHOR_1, LABEL] });
      expect(result).toEqual(expect.stringContaining('singular'));
      translationsSpy.mockClear();

      result = translateFieldName(FIELDS.author, { ...ITEM_BASE, values: [AUTHOR_1, AUTHOR_2] });
      expect(result).toEqual(expect.stringContaining('plural'));
      translationsSpy.mockClear();
    });

    it('returns translation for singular field name when item is not given', () => {
      translateFieldName(FIELDS.label);
      expect(translationsSpy).toHaveBeenCalledWith(expect.stringContaining('singular'));
      translationsSpy.mockClear();
    });

    it('takes into account the item entity type and first tries if specific translation exists for entity type, then for "all"', () => {
      translationsSpy.mockImplementation((key) => `translated ${key}`);

      translateFieldName(FIELDS.author);
      expect(translationsSpy).toHaveBeenCalledWith(expect.stringContaining('.all.'));
      expect(translationsSpy).not.toHaveBeenCalledWith(expect.stringContaining('.Resource.'));

      translationsSpy.mockClear();
      translateFieldName(FIELDS.author, ITEM);
      expect(translationsSpy).toHaveBeenCalledWith(expect.stringContaining('.Resource.'));
      expect(translationsSpy).not.toHaveBeenCalledWith(expect.stringContaining('.all.'));

      translationsSpy.mockClear().mockImplementation((key) => key);
      translateFieldName(FIELDS.author, ITEM);
      expect(translationsSpy).toHaveBeenCalledWith(expect.stringContaining('.Resource.'));
      expect(translationsSpy).toHaveBeenCalledWith(expect.stringContaining('.all.'));
    });

    it('returns the normalized untranslated field name if no translation exists', () => {
      translationsSpy.mockReturnValueOnce('My awesome test field');
      expect(translateFieldName(FIELDS.myAwesomeField)).toBe('My awesome test field');

      translationsSpy.mockImplementationOnce((key: string) => key);
      expect(translateFieldName(FIELDS.myAwesomeField)).toBe('myAwesomeField');
    });

    it('allows to enforce pluralization', () => {
      translationsSpy.mockImplementation((key) => `translated ${key}`);
      expect(translateFieldName(FIELDS.myAwesomeField, null, true)).toBe(
        'translated fields.labels.myAwesomeField.all.plural'
      );
      translationsSpy.mockClear();
    });

    it('allows to specify a custom fallback in case there is no translation', () => {
      translationsSpy.mockImplementation((key) => key);
      expect(translateFieldName(FIELDS.myAwesomeField, null, false, 'custom fallback')).toBe('custom fallback');
      translationsSpy.mockClear();
    });
  });

  describe('translateFieldDescription()', () => {
    it('behaves like translateFieldName() but returns a field description and does not offer pluralization options nor fallback', () => {
      // simulate that there is no description translation for email field:
      translationsSpy.mockImplementation((key) => (key.includes('email') ? key : `translated ${key}`));

      expect(translateFieldDescription(FIELDS.label)).toBe('translated fields.descriptions.label.all.text');
      expect(translateFieldDescription(FIELDS.label, ITEM)).toBe('translated fields.descriptions.label.Resource.text');
      expect(translateFieldDescription(FIELDS.email)).toBe(undefined);

      translationsSpy.mockClear();
    });
  });

  describe('translateFieldValue()', () => {
    it('returns the untranslated value if it is no concept ID nor entity name', () => {
      expect(translateFieldValue('foo')).toBe('foo');
    });

    it('normalizes the concept ID or entity name and then returns its value translation', () => {
      translationsSpy.mockImplementation((key: string) => `translated ${key}`);
      jest.spyOn(config, 'FIELD_CONCEPT_PREFIXES').mockReturnValue(['https://concept.prefix/']);

      expect(translateFieldValue('https://concept.prefix/foo')).toBe('translated vocabulary.foo.label');
      expect(translateFieldValue('https://concept.prefix/bar')).toBe('translated vocabulary.bar.label');
    });

    it('returns the raw normalized concept ID or entity name if no translation exists', () => {
      translationsSpy.mockImplementation((key: string) => key);
      jest.spyOn(config, 'FIELD_CONCEPT_PREFIXES').mockReturnValue(['https://concept.prefix/']);

      expect(translateFieldValue('https://concept.prefix/foo')).toBe('foo');
      expect(translateFieldValue('Bar', FIELDS.entityName)).toBe('Bar');
    });
  });

  describe('translateFieldValueDescription()', () => {
    it('behaves like translateFieldValue() but returns value descriptions', () => {
      jest.spyOn(config, 'FIELD_CONCEPT_PREFIXES').mockReturnValue(['https://concept.prefix/']);

      translationsSpy.mockImplementation((key: string) => `translated ${key}`);
      expect(translateFieldValueDescription('foo')).toBe('foo');
      expect(translateFieldValueDescription('https://concept.prefix/foo')).toBe(
        'translated vocabulary.foo.description'
      );
      expect(translateFieldValueDescription('Baz', FIELDS.entityName)).toBe('translated vocabulary.Baz.description');

      translationsSpy.mockImplementation((key: string) => key);
      expect(translateFieldValueDescription('https://concept.prefix/foo')).toBe('foo');
      expect(translateFieldValueDescription('Bar', FIELDS.entityName)).toBe('Bar');
    });
  });

  describe('formatDate()', () => {
    const date = new Date('Thu Feb 17 2022 14:22:02 GMT+0100');

    beforeEach(() => {
      stores.i18n.language = 'de';
    });

    it('returns localized date string from date', () => {
      expect(formatDate(date)).toBe('17. Feb. 2022, 13:22:02');

      stores.i18n.language = 'en';
      expect(formatDate(date)).toContain('Feb 17, 2022, 01:22:02');
    });

    it('provides an argument to specify the precision level to be rendered', () => {
      expect(formatDate(date, FieldValueDatePrecisionLevel.DAY)).toBe('17. Feb. 2022');
      expect(formatDate(date, FieldValueDatePrecisionLevel.MONTH)).toBe('Feb. 2022');
      expect(formatDate(date, FieldValueDatePrecisionLevel.YEAR)).toBe('2022');
      expect(formatDate(date, FieldValueDatePrecisionLevel.NONE)).toBe(undefined);

      stores.i18n.language = 'en';
      expect(formatDate(date, FieldValueDatePrecisionLevel.DAY)).toBe('Feb 17, 2022');
      expect(formatDate(date, FieldValueDatePrecisionLevel.MONTH)).toBe('Feb 2022');
    });

    it('provides an argument to change the rendered timezone (default: "UTC")', () => {
      expect(formatDate(date, FieldValueDatePrecisionLevel.TIME, 'Europe/Berlin')).toBe('17. Feb. 2022, 14:22:02');
      expect(formatDate(date, FieldValueDatePrecisionLevel.DAY, 'Pacific/Norfolk')).toBe('18. Feb. 2022');
    });

    it('returns iso string if date.toLocaleString() is not available', () => {
      jest.spyOn(date, 'toLocaleString').mockImplementation(() => {
        throw new Error('not implemented');
      });
      expect(formatDate(date)).toBe('2022-02-17T13:22:02.000Z');
    });
  });

  describe('getTextMatches()', () => {
    it('extracts search highlight matches for given field name', () => {
      languageSpy.mockReturnValueOnce(null);

      const authorMatches = getTextMatches(FIELDS.author, HIGHLIGHTS);
      expect(authorMatches).toEqual(['Thomas \ue000Doe\ue001', '\ue000Doe\ue001']);
    });

    it('filters out matches that do not belong to currently set user language (if FEATURE_FLAGS.FIELD_TRANSLATIONS is on)', () => {
      config.FEATURE_FLAGS.FIELD_TRANSLATIONS = true;

      languageSpy.mockReturnValueOnce(null);
      expect(getTextMatches(FIELDS.label, HIGHLIGHTS)).toEqual([]);

      languageSpy.mockReturnValueOnce('de');
      expect(getTextMatches(FIELDS.label, HIGHLIGHTS)).toEqual(['Baz mit John \ue000Doe\ue001']);

      languageSpy.mockReturnValueOnce('en');
      expect(getTextMatches(FIELDS.label, HIGHLIGHTS)).toEqual(['Baz with John \ue000Doe\ue001']);

      config.FEATURE_FLAGS.FIELD_TRANSLATIONS = false;
      expect(getTextMatches(FIELDS.label, HIGHLIGHTS)).toEqual([
        'Baz with John \ue000Doe\ue001',
        'Baz mit John \ue000Doe\ue001',
      ]);
    });
  });

  describe('filterDuplicates()', () => {
    it('filters out duplicates from an array, to be passed as filter arg', () => {
      expect(['foo', 'bar', 'foo', 'baz', 'bar'].filter(filterDuplicates)).toEqual(['foo', 'bar', 'baz']);
    });
  });

  describe('sortAndFilterValues()', () => {
    it('sorts values by their place property (if given, else place = 0)', () => {
      const valueA = { place: 2, display: 'Value A' };
      const valueB = { place: 0, display: 'Value B' };
      const valueC = { place: 3, display: 'Value C' };
      const valueD = { display: 'Value D' };
      expect(sortAndFilterValues([valueA, valueB, valueC, valueD])).toEqual([valueB, valueD, valueA, valueC]);
      expect(sortAndFilterValues([valueC, valueD, valueA])).toEqual([valueD, valueA, valueC]);
      expect(sortAndFilterValues([valueC, AUTHOR_1])).toEqual([AUTHOR_1, valueC]);
    });

    it('filters out values in foreign languages if feature flag is active', () => {
      // const LABEL is declared above and has language 'en'
      const valueDE = { language: 'de', display: 'Value DE' };
      const valueNeutral = { display: 'Value neutral language' };
      const valueEN = { language: 'en', display: 'Value EN' };

      stores.i18n.language = 'de';
      expect(sortAndFilterValues([valueDE, valueNeutral, valueEN, LABEL])).toEqual([
        valueDE,
        valueNeutral,
        valueEN,
        LABEL,
      ]);
      expect(sortAndFilterValues([valueDE, valueNeutral, valueEN, LABEL], true)).toEqual([valueDE, valueNeutral]);

      stores.i18n.language = 'en';
      expect(sortAndFilterValues([valueDE, valueNeutral, valueEN, LABEL], true)).toEqual([
        valueNeutral,
        valueEN,
        LABEL,
      ]);

      stores.i18n.language = 'no';
      expect(sortAndFilterValues([valueDE, valueNeutral, valueEN, LABEL], true)).toEqual([valueNeutral]);
    });

    it('filters out duplicates', () => {
      const valueA = { display: 'Value A' };
      const valueB = { display: 'Value B' };
      expect(sortAndFilterValues([valueA, valueB, valueA])).toEqual([valueA, valueB]);
    });
  });

  describe('getValues()', () => {
    it('extracts item values for given field name', () => {
      expect(getValues(FIELDS.author, ITEM)).toEqual([AUTHOR_1, AUTHOR_2]);
      expect(getValues(FIELDS.label, ITEM)).toEqual([LABEL]);
    });

    it('filters out values that do not belong to currently set user language (if FEATURE_FLAGS.FIELD_TRANSLATIONS is on)', () => {
      config.FEATURE_FLAGS.FIELD_TRANSLATIONS = true;

      const AUTHOR_1_DE = { ...AUTHOR_2, language: 'de' };
      const AUTHOR_2_EN = { ...AUTHOR_2, language: 'en' };
      const itemWithLanguageValues = { ...ITEM_BASE, values: [AUTHOR_1, AUTHOR_1_DE, AUTHOR_2_EN] };

      languageSpy.mockReturnValueOnce(null);
      expect(getValues(FIELDS.author, itemWithLanguageValues)).toEqual([AUTHOR_1]);

      languageSpy.mockReturnValueOnce('de');
      expect(getValues(FIELDS.author, itemWithLanguageValues)).toEqual([AUTHOR_1, AUTHOR_1_DE]);

      languageSpy.mockReturnValueOnce('en');
      expect(getValues(FIELDS.author, itemWithLanguageValues)).toEqual([AUTHOR_1, AUTHOR_2_EN]);

      config.FEATURE_FLAGS.FIELD_TRANSLATIONS = false;
      expect(getValues(FIELDS.author, itemWithLanguageValues)).toEqual([AUTHOR_1, AUTHOR_1_DE, AUTHOR_2_EN]);
    });

    it('falls back to showing all values when no values are found for user language (if FEATURE_FLAGS.FIELD_TRANSLATIONS is on)', () => {
      config.FEATURE_FLAGS.FIELD_TRANSLATIONS = true;

      const AUTHOR_1_EN = { ...AUTHOR_1, language: 'en' };
      const AUTHOR_2_EN = { ...AUTHOR_2, language: 'en' };
      const itemWithLanguageValues = { ...ITEM_BASE, values: [AUTHOR_1_EN, AUTHOR_2_EN] };

      languageSpy.mockReturnValueOnce('de');
      expect(getValues(FIELDS.author, itemWithLanguageValues)).toEqual([AUTHOR_1_EN, AUTHOR_2_EN]);

      config.FEATURE_FLAGS.FIELD_TRANSLATIONS = false;
      expect(getValues(FIELDS.author, itemWithLanguageValues)).toEqual([AUTHOR_1_EN, AUTHOR_2_EN]); // same as above
    });

    it('filters out duplicates', () => {
      const itemWithDuplicateValues = { ...ITEM_BASE, values: [AUTHOR_1, AUTHOR_2, AUTHOR_1, AUTHOR_1] };
      languageSpy.mockReturnValueOnce(null);
      expect(getValues(FIELDS.author, itemWithDuplicateValues)).toEqual([AUTHOR_1, AUTHOR_2]);
    });
  });

  describe('getValue()', () => {
    it('extracts the first value for given field name', () => {
      expect(getValue(FIELDS.author, ITEM)).toBe('Jane');
    });
  });

  describe('hasValue()', () => {
    it('returns if item has a value for given field name', () => {
      expect(hasValue(FIELDS.author, ITEM)).toBe(true);
      expect(hasValue(FIELDS.myAwesomeField, ITEM)).toBe(false);
    });
  });

  describe('formatValue()', () => {
    it('concatenates multiple values and returns a comma-separated string', () => {
      expect(formatValue(['foo', 'bar', 'baz'])).toBe('foo, bar, baz');
      expect(formatValue([])).toBe('');
    });

    it('supports concatenating with another string based on renderer if field is given', () => {
      concatenationSpy.mockReturnValue({
        [FieldRenderer.time]: ' and ',
        default: ' | ',
      });

      expect(formatValue(['foo', 'bar', 'baz'], FIELDS.author)).toBe('foo | bar | baz');
      expect(formatValue(['foo', 'bar', 'baz'], FIELDS.time)).toBe('foo and bar and baz');
      expect(formatValue(['foo', 'bar'], FIELDS.label)).toBe('foo | bar');
    });

    it('escapes html characters', () => {
      expect(formatValue(['foo', '<script>alert("Hello")</script>'])).toBe(
        'foo, &lt;script&gt;alert(&quot;Hello&quot;)&lt;/script&gt;'
      );
    });

    it('supports rendering dates with time renderer', () => {
      stores.i18n.language = 'de';

      // only YYYY(-DD)?(-MM)?(THH:MM.*)? supported, no time rendered, everything else falls back to raw value
      expect(formatValue(['2022-03-10T09:22:13.835Z', '2022-04'], FIELDS.time)).toBe('10. März 2022 - Apr. 2022');
      expect(formatValue(['May 1980', '2022-04-06'], FIELDS.time)).toBe('May 1980 - 06. Apr. 2022');
      expect(formatValue(['Unknown', '2021'], FIELDS.time)).toBe('Unknown - 2021');
      expect(formatValue(['1990-05-01T10:09:49.062Z', 'now'], FIELDS.time)).toBe('01. Mai 1990 - now');
      expect(formatValue(['1980-05', 'Yesterday'], FIELDS.time)).toBe('Mai 1980 - Yesterday');
    });

    it('translates values from concept pref labels', () => {
      jest.spyOn(config, 'FIELD_CONCEPT_PREFIXES').mockReturnValue(['concept/prefix/']);
      translationsSpy.mockImplementation((key: string) => `translated ${key}`);
      expect(formatValue(['concept/prefix/open'])).toBe(`translated vocabulary.open.label`);
    });
  });

  describe('hasValueChanged()', () => {
    it('returns false when only one item is given', () => {
      expect(hasValueChanged(FIELDS.author, ITEM)).toBe(false);
      expect(hasValueChanged(FIELDS.author, null, ITEM)).toBe(false);
    });

    it('returns false when values of both items for given field are equal', () => {
      expect(hasValueChanged(FIELDS.author, ITEM, { ...ITEM, values: [AUTHOR_2, AUTHOR_1] } as Item)).toBe(false);
    });

    it('returns true when values of both items for given field differ', () => {
      expect(hasValueChanged(FIELDS.author, ITEM, { ...ITEM, values: [AUTHOR_1, AUTHOR_2, AUTHOR_3] } as Item)).toBe(
        true
      );

      expect(
        hasValueChanged(FIELDS.author, ITEM, {
          ...ITEM,
          values: [
            AUTHOR_1,
            {
              fieldName: 'author',
              fieldValue: 'Micael',
            },
          ],
        } as Item)
      ).toBe(true);
    });
  });

  describe('getLangAttrIfForeign()', () => {
    describe('summarize argument set to false', () => {
      it('returns for every value its language if derives from user language (else undefined)', () => {
        languageSpy.mockReturnValue('de');
        expect(getLangAttrIfForeign([FIELDS.author, FIELDS.label], ITEM, false)).toEqual([undefined, undefined, 'en']);
        expect(getLangAttrIfForeign([FIELDS.label, FIELDS.author], ITEM, false)).toEqual(['en', undefined, undefined]);

        languageSpy.mockReturnValue('en');
        expect(getLangAttrIfForeign([FIELDS.author, FIELDS.label], ITEM, false)).toEqual([
          undefined,
          undefined,
          undefined,
        ]);
      });
    });

    describe('summarize argument set to true', () => {
      it('returns the common deriving language if all derive from user language (else undefined)', () => {
        languageSpy.mockReturnValue('de');
        expect(getLangAttrIfForeign([FIELDS.author, FIELDS.label], ITEM, true)).toEqual(undefined);
        expect(getLangAttrIfForeign([FIELDS.label], ITEM, true)).toEqual('en');

        languageSpy.mockReturnValue('en');
        expect(getLangAttrIfForeign([FIELDS.label], ITEM, true)).toEqual(undefined);
      });
    });
  });

  describe('getRawValues()', () => {
    it('returns the raw (non-formatted) values for given fields', () => {
      config.FIELDS_CONFIG.DEFAULT_FIELD_CONCEPT_PREFIX = 'https://concept.foo/';

      expect(getRawValues([FIELDS.label], ITEM)).toStrictEqual(['Foo Bar Baz with John Doe.']);
      expect(getRawValues([FIELDS.author, FIELDS.label], ITEM)).toStrictEqual([
        'Jane',
        'Doe',
        'Foo Bar Baz with John Doe.',
      ]);
      expect(getRawValues([FIELDS.conceptField], ITEM)).toStrictEqual(['https://concept.foo/value-concept-foo']);
      expect(getRawValues([FIELDS.time], ITEM)).toStrictEqual(['2012-08-17']);
    });

    it('removes duplicates', () => {
      expect(getRawValues([FIELDS.author, FIELDS.author], ITEM)).toStrictEqual(['Jane', 'Doe']);
    });
  });

  describe('getDisplayValues()', () => {
    beforeAll(() => {
      languageSpy.mockReturnValue('en');
    });

    it('has dependencies to the stored search highlights', () => {
      getDisplayValues([FIELDS.label], ITEM);
      expect(highlightsSpy).toHaveBeenCalled();
    });

    it('returns the item values for given fields with the text matches', () => {
      highlightsSpy.mockImplementationOnce(() => HIGHLIGHTS);
      expect(getDisplayValues([FIELDS.label], ITEM)).toStrictEqual(['Foo Bar Baz with John <em>Doe</em>.']);

      highlightsSpy.mockImplementationOnce(() => HIGHLIGHTS);
      expect(getDisplayValues([FIELDS.author], ITEM)).toStrictEqual(['Jane', '<em>Doe</em>']);

      highlightsSpy.mockImplementationOnce(() => []);
      expect(getDisplayValues([FIELDS.label, FIELDS.author], ITEM)).toStrictEqual([
        'Foo Bar Baz with John Doe.',
        'Jane',
        'Doe',
      ]);
    });

    it("formats the values according to the field's formatter", () => {
      translationsSpy.mockImplementation((key) => `translated ${key}`);
      jest.spyOn(config, 'FIELD_CONCEPT_PREFIXES').mockReturnValue(['https://concept.foo/']);

      expect(getDisplayValues([FIELDS.conceptField], ITEM)).toStrictEqual([
        'translated vocabulary.value-concept-foo.label',
      ]);
      expect(getDisplayValues([FIELDS.time], ITEM)).toStrictEqual(['Aug 17, 2012']);
    });

    it('provides an argument to skip highlighting the matches', () => {
      highlightsSpy.mockImplementationOnce(() => HIGHLIGHTS);
      expect(getDisplayValues([FIELDS.label], ITEM, false)).toStrictEqual(['Foo Bar Baz with John Doe.']);
    });
  });

  describe('getConcatDisplayValue()', () => {
    it('concatenates display values of multiple fields', () => {
      expect(getConcatDisplayValue([FIELDS.author, FIELDS.label], ITEM, '-', false)).toBe(
        'Jane, Doe, Foo Bar Baz with John Doe.'
      );
    });

    it('uses concatenation separator based on renderer of first field', () => {
      concatenationSpy.mockReturnValue({
        [FieldRenderer.title]: ' | ',
        default: ', ',
      });

      expect(getConcatDisplayValue([FIELDS.author, FIELDS.label], ITEM, '-', false)).toBe(
        'Jane, Doe, Foo Bar Baz with John Doe.'
      );
      expect(getConcatDisplayValue([FIELDS.label, FIELDS.author], ITEM, '-', false)).toBe(
        'Foo Bar Baz with John Doe. | Jane | Doe'
      );
    });

    it('supports providing a fallback in case the value is empty', () => {
      expect(getConcatDisplayValue([FIELDS.myAwesomeField], ITEM, 'my fallback')).toBe('my fallback');
      expect(getConcatDisplayValue([FIELDS.myAwesomeField], ITEM)).toBe(
        stores.i18n.t(config.FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY)
      );
    });

    it('provides an argument to skip highlighting the matches', () => {
      highlightsSpy.mockImplementationOnce(() => HIGHLIGHTS);
      expect(getConcatDisplayValue([FIELDS.label], ITEM, '-', false)).toBe('Foo Bar Baz with John Doe.');
    });
  });

  describe('getConcatenator()', () => {
    it('returns the concatenator of the renderer of the first field (fallback: ITEM_FIELDS_CONCATENATOR.default)', () => {
      concatenationSpy.mockReturnValue({
        [FieldRenderer.title]: ' | ',
        default: ', ',
      });

      expect(getConcatenator([FIELDS.label])).toBe(' | ');
      expect(getConcatenator([FIELDS.email])).toBe(', ');
      expect(getConcatenator([FIELDS.time])).toBe(', ');
      expect(getConcatenator([FIELDS.label, FIELDS.email])).toBe(' | ');
      expect(getConcatenator([FIELDS.email, FIELDS.label])).toBe(', ');
    });
  });

  describe('getInvolvedFields()', () => {
    it('returns rendered + sidebar resolved fields by entity type and removes duplicates', () => {
      ITEM_CONFIG.DISPLAYED_FIELDS = {
        Source: [],
        Resource: [FIELDS.author, FIELDS.myAwesomeField],
        Person: [FIELDS.label, FIELDS.email, FIELDS.virtual],
      };
      SEARCH_CONFIG.SIDEBAR_FEATURES = {
        Resource: [createSidebarFeatureConfig('resource', SidebarFeature.date, FIELDS.time)],
        default: [],
      };
      ITEM_CONFIG.SIDEBAR_FEATURES = {
        Source: [
          createSidebarFeatureConfig('source', SidebarFeature.displayField, FIELDS.label),
          createSidebarFeatureConfig('source', SidebarFeature.displayField, FIELDS.time),
        ],
        Resource: [createSidebarFeatureConfig('resource', SidebarFeature.displayField, FIELDS.email)],
        default: [createSidebarFeatureConfig('default', SidebarFeature.displayField, FIELDS.myAwesomeField)],
      };

      expect(getInvolvedFields('Source')).toEqual([FIELDS.label, FIELDS.time]);
      expect(getInvolvedFields('Resource')).toEqual([FIELDS.author, FIELDS.myAwesomeField, FIELDS.time, FIELDS.email]);
      expect(getInvolvedFields('Person')).toEqual([FIELDS.label, FIELDS.email, FIELDS.myAwesomeField, FIELDS.time]);
    });

    it('has an option to only return non-cached fields (removes fields solely used in sidebar contact form)', () => {
      ITEM_CONFIG.DISPLAYED_FIELDS = {
        Resource: [FIELDS.author, FIELDS.myAwesomeField],
      };
      SEARCH_CONFIG.SIDEBAR_FEATURES = { default: [] };
      ITEM_CONFIG.SIDEBAR_FEATURES = {
        Resource: [
          createSidebarFeatureConfig('resource', SidebarFeature.contactForm, FIELDS.email),
          createSidebarFeatureConfig('resource', SidebarFeature.contactForm, FIELDS.myAwesomeField),
        ],
      };

      expect(getInvolvedFields('Resource', false)).toEqual([FIELDS.author, FIELDS.myAwesomeField]);
    });
  });

  describe('getDocumentationProp()', () => {
    it('returns documentation texts by prop name and entity type context', () => {
      languageSpy.mockReturnValue('de');

      expect(getDocumentationProp(FIELDS.time, 'exampleValue')).toBe(null);

      expect(getDocumentationProp(FIELDS.documentedField, 'descriptionText')).toBe('generic description text DE');
      expect(getDocumentationProp(FIELDS.documentedField, 'descriptionText', FieldEntityVirtualType.all)).toBe(
        'generic description text DE'
      );
      expect(getDocumentationProp(FIELDS.documentedField, 'exampleValue')).toBe(null);
      expect(getDocumentationProp(FIELDS.documentedField, 'exampleValue', 'Resource')).toBe(
        'resource example value DE'
      );
      expect(getDocumentationProp(FIELDS.documentedField, 'furtherInformation')).toBe(null);
      expect(getDocumentationProp(FIELDS.documentedField, 'displayFormats')).toBe(null);

      languageSpy.mockReturnValue('en');
      expect(getDocumentationProp(FIELDS.documentedField, 'descriptionText')).toBe(null);
      expect(getDocumentationProp(FIELDS.documentedField, 'descriptionText', FieldEntityVirtualType.all)).toBe(null);
      expect(getDocumentationProp(FIELDS.documentedField, 'exampleValue')).toBe('generic example value EN');
      expect(getDocumentationProp(FIELDS.documentedField, 'exampleValue', 'Resource')).toBe('generic example value EN');
      expect(getDocumentationProp(FIELDS.documentedField, 'furtherInformation')).toBe(null);
      expect(getDocumentationProp(FIELDS.documentedField, 'displayFormats')).toBe('generic display formats EN');
    });
  });

  describe('calculateCompleteness()', () => {
    it('calculates a score by getting the involved fields, comparing against the populated fields, and applying a score weight per field category', () => {
      ITEM_CONFIG.DISPLAYED_FIELDS = {
        Resource: [FIELDS.label, FIELDS.virtual],
      };
      SEARCH_CONFIG.SIDEBAR_FEATURES = {
        Resource: [createSidebarFeatureConfig('resource', SidebarFeature.displayField, FIELDS.myAwesomeField)],
        default: [],
      };
      ITEM_CONFIG.SIDEBAR_FEATURES = {
        Resource: [createSidebarFeatureConfig('resource', SidebarFeature.contactForm, FIELDS.email)],
        default: [],
      };

      let mockedWeights = { mandatory: 50, recommended: 31, optional: 19 };
      config.FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS = mockedWeights;

      expect(calculateCompleteness(ITEM)).toStrictEqual({
        missingFields: { mandatory: [FIELDS.myAwesomeField], recommended: [], optional: [FIELDS.email] },
        populatedFields: { mandatory: [FIELDS.label], recommended: [FIELDS.time], optional: [] },
        scores: { mandatory: 25, recommended: 31, optional: 0, total: 56 },
        weights: mockedWeights,
      });

      ITEM_CONFIG.DISPLAYED_FIELDS = {
        Resource: [FIELDS.label, FIELDS.author, FIELDS.email, FIELDS.time],
      };

      mockedWeights = { mandatory: 50, recommended: 40, optional: 10 };
      config.FIELDS_CONFIG.METADATA_COMPLETENESS_WEIGHTS = mockedWeights;

      expect(calculateCompleteness(ITEM)).toEqual({
        missingFields: { mandatory: [FIELDS.myAwesomeField], recommended: [], optional: [FIELDS.email] },
        populatedFields: { mandatory: [FIELDS.label], recommended: [FIELDS.author, FIELDS.time], optional: [] },
        scores: { mandatory: 25, recommended: 40, optional: 0, total: 65 },
        weights: mockedWeights,
      });
    });
  });
});
