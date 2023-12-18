import { FieldEntityVirtualType, Field, FieldImportance, FieldRenderer } from './field';

const EXAMPLE_FIELD_FOO = new Field({
  name: 'foo',
  renderer: FieldRenderer.title,
  importance: FieldImportance.mandatory,
  isVirtual: false,
  isEnumerable: false,
  label: {
    de: [{ entityType: FieldEntityVirtualType.all, singular: 'Foo label [de]' }],
    en: [{ entityType: FieldEntityVirtualType.all, singular: 'Foo label [en]' }],
  },
  description: {
    de: [{ entityType: FieldEntityVirtualType.all, text: 'Foo description [de]' }],
    en: [{ entityType: FieldEntityVirtualType.all, text: 'Foo description [en]' }],
  },
});

const EXAMPLE_FIELD_BAR = new Field(
  {
    name: 'bar',
    renderer: FieldRenderer.plain,
    importance: FieldImportance.recommended,
    isVirtual: false,
    isEnumerable: true,
    vocabulary: {
      de: [
        {
          conceptId: 'concept-1',
          label: 'Label for concept 1 [DE]',
          description: 'Description for concept 1 [DE]',
        },
      ],
      en: [
        {
          conceptId: 'concept-1',
          label: 'Label for concept 1 [EN]',
          description: 'Description for concept 1 [EN]',
        },
        {
          conceptId: 'concept-2',
          label: 'Label for concept 2 [EN]',
          description: 'Description for concept 2 [EN]',
        },
      ],
    },
  },
  true,
  null,
  '_test'
);

const EXAMPLE_FIELD_BAZ = new Field(
  {
    name: 'baz',
    renderer: FieldRenderer.description,
    importance: FieldImportance.optional,
    isVirtual: false,
    isEnumerable: false,
    linkedField: 'foo',
  },
  true
);

const EXAMPLE_FIELD_VIRTUAL = new Field({
  name: 'virtual',
  renderer: FieldRenderer.plain,
  importance: FieldImportance.none,
  isVirtual: true,
  isEnumerable: false,
  resolvesTo: ['foo', 'baz'],
  useTranslationsFrom: 'foo',
});

const EXAMPLE_FIELD_PENDING = Field.createPending('pending');

describe('field model', () => {
  it("returns the field's name and linked name (if applicable with appendix)", () => {
    expect(EXAMPLE_FIELD_FOO.name).toBe('foo');
    expect(EXAMPLE_FIELD_BAR.name).toBe('bar');
    expect(EXAMPLE_FIELD_BAZ.name).toBe('baz');
    expect(EXAMPLE_FIELD_VIRTUAL.name).toBe('virtual');
    expect(EXAMPLE_FIELD_PENDING.name).toBe('pending');

    /**
     * Linked name is specific to search - refers to field y in resolved entity of field x
     * (fallback: field name).
     * Can contain an appendix to refer to specific versions of the same field (e.g. raw)
     */
    expect(EXAMPLE_FIELD_FOO.linkedName).toBe('foo');
    expect(EXAMPLE_FIELD_BAR.linkedName).toBe('bar_test');
    expect(EXAMPLE_FIELD_BAZ.linkedName).toBe('baz__foo');
    expect(EXAMPLE_FIELD_PENDING.linkedName).toBe('pending');
  });

  it('returns if the field is initialized (fully configured)', () => {
    expect(EXAMPLE_FIELD_FOO.isInitialized).toBe(true);
    expect(EXAMPLE_FIELD_PENDING.isInitialized).toBe(false);
  });

  it("returns the field's renderer", () => {
    expect(EXAMPLE_FIELD_FOO.renderer).toBe(FieldRenderer.title);
    expect(EXAMPLE_FIELD_BAR.renderer).toBe(FieldRenderer.plain);
    expect(EXAMPLE_FIELD_BAZ.renderer).toBe(FieldRenderer.description);
    expect(EXAMPLE_FIELD_VIRTUAL.renderer).toBe(FieldRenderer.plain);
  });

  it("returns the field's metadata completeness importance", () => {
    expect(EXAMPLE_FIELD_FOO.importance).toBe(FieldImportance.mandatory);
    expect(EXAMPLE_FIELD_BAR.importance).toBe(FieldImportance.recommended);
    expect(EXAMPLE_FIELD_BAZ.importance).toBe(FieldImportance.optional);
    expect(EXAMPLE_FIELD_VIRTUAL.importance).toBe(FieldImportance.none);
  });

  it('returns if the field is virtual (resolves to one or more other fields)', () => {
    expect(EXAMPLE_FIELD_FOO.isVirtual).toBe(false);
    expect(EXAMPLE_FIELD_VIRTUAL.isVirtual).toBe(true);
  });

  it('returns the resolved fields (if virtual field, else field itself)', () => {
    expect(EXAMPLE_FIELD_FOO.resolvesTo).toEqual([EXAMPLE_FIELD_FOO]);
    expect(EXAMPLE_FIELD_VIRTUAL.resolvesTo).toEqual([EXAMPLE_FIELD_FOO, EXAMPLE_FIELD_BAZ]);
  });

  it('returns the linked field (if specified)', () => {
    expect(EXAMPLE_FIELD_FOO.linkedField).toBe(undefined);
    expect(EXAMPLE_FIELD_BAZ.linkedField).toBe(EXAMPLE_FIELD_FOO);
  });

  it('returns the field translations are drawn from (if specified)', () => {
    expect(EXAMPLE_FIELD_FOO.useTranslationsFrom).toBe(undefined);
    expect(EXAMPLE_FIELD_VIRTUAL.useTranslationsFrom).toBe(EXAMPLE_FIELD_FOO);
  });

  it("returns the field's label mapping", () => {
    expect(EXAMPLE_FIELD_FOO.label?.de?.[0]?.singular).toBe('Foo label [de]');
    expect(EXAMPLE_FIELD_BAR.label).toBe(undefined);
    expect(EXAMPLE_FIELD_VIRTUAL.label?.en?.[0]?.singular).toBe('Foo label [en]');
  });

  it("returns the field's description mapping", () => {
    expect(EXAMPLE_FIELD_FOO.description?.en?.[0]?.text).toBe('Foo description [en]');
    expect(EXAMPLE_FIELD_BAR.description).toBe(undefined);
    expect(EXAMPLE_FIELD_VIRTUAL.description?.de?.[0]?.text).toBe('Foo description [de]');
  });

  it("returns the field's vocabulary mapping", () => {
    expect(EXAMPLE_FIELD_FOO.vocabulary).toBe(undefined);
    expect(EXAMPLE_FIELD_BAR.vocabulary?.de?.[0]?.conceptId).toBe('concept-1');
  });

  it("returns the concept IDs extracted from the field's vocabulary", () => {
    expect(EXAMPLE_FIELD_FOO.conceptIds).toEqual([]);
    expect(EXAMPLE_FIELD_BAR.conceptIds).toEqual(['concept-1', 'concept-2']);
  });

  describe('configure()', () => {
    it('allows to configure the field (also allows to override previous config)', () => {
      const field = new Field(
        {
          name: 'test',
          renderer: FieldRenderer.plain,
          importance: FieldImportance.none,
          isVirtual: false,
          isEnumerable: false,
        },
        false,
        null,
        '_raw'
      );

      expect(field.linkedName).toBe('test_raw');
      field.configure({ ...field.config, name: 'test2' }, true, EXAMPLE_FIELD_FOO);
      expect(field.linkedName).toBe('test2__foo_raw');
    });
  });

  describe('createPending()', () => {
    it('creates a pending field (only name, without config)', () => {
      // pending fields can be used before initialization is complete - will
      // be configured as soon as configure() is called with a config
      expect(Field.createPending('myPendingField').name).toBe('myPendingField');
    });
  });

  describe('cloneAndLinkToField()', () => {
    it('allows to link one field to another, converts to search field (returns cloned instance)', () => {
      const field = EXAMPLE_FIELD_FOO.cloneAndLinkToField(EXAMPLE_FIELD_BAZ);
      expect(EXAMPLE_FIELD_FOO.linkedName).toBe('foo');
      expect(field.linkedName).toBe('foo__baz');
    });
  });

  describe('cloneAndSetNameAppendix()', () => {
    it("allows to add an appendix to a field's linkedName, converts to search field (returns cloned instance)", () => {
      const field = EXAMPLE_FIELD_FOO.cloneAndSetNameAppendix('appendix');
      expect(EXAMPLE_FIELD_FOO.linkedName).toBe('foo');
      expect(field.linkedName).toBe('fooappendix');
    });
  });

  describe('cloneAndSetRawValueMode()', () => {
    it('sets the field appendix to "_raw_value" (search-specific feature, returns cloned instance)', () => {
      expect(EXAMPLE_FIELD_FOO.cloneAndSetRawValueMode().linkedName).toBe('foo_raw_value');
      expect(EXAMPLE_FIELD_FOO.linkedName).toBe('foo');
    });
  });
});
