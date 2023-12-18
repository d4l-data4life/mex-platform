export enum FieldEntityVirtualType {
  all = 'all',
  unknown = 'unknown',
}

export enum FieldRenderer {
  title = 'title',
  description = 'description',
  plain = 'plain',
  reference = 'reference',
  link = 'link',
  time = 'time',
  entity = 'entity',
  bullets = 'bullets',
  none = 'none',
}

export enum FieldKind {
  string = 'string',
  text = 'text',
  link = 'link',
  timestamp = 'timestamp',
  hierarchy = 'hierarchy',
  coding = 'coding',
}

export enum FieldImportance {
  mandatory = 'mandatory',
  recommended = 'recommended',
  optional = 'optional',
  none = 'none',
}

export interface FieldConfigLabel {
  entityType: FieldEntityVirtualType | string;
  singular: string;
  plural?: string;
}

export interface FieldConfigDescription {
  entityType: FieldEntityVirtualType | string;
  text: string;
}

export interface FieldConfigDocumentation {
  descriptionText?: FieldConfigDescription[];
  exampleValue?: FieldConfigDescription[];
  furtherInformation?: FieldConfigDescription[];
  displayFormats?: FieldConfigDescription[];
}

interface FieldConfigVocabulary {
  conceptId: string;
  label: string;
  description: string;
}

export interface FieldConfig {
  name: string;
  renderer: FieldRenderer;
  kind?: FieldKind;
  importance: FieldImportance;
  isVirtual: boolean;
  resolvesTo?: string[];
  isEnumerable: boolean;
  isMultiValued?: boolean;
  linkedField?: string;
  useTranslationsFrom?: string;
  label?: {
    [language: string]: FieldConfigLabel[];
  };
  description?: {
    [language: string]: FieldConfigDescription[];
  };
  vocabulary?: {
    [language: string]: FieldConfigVocabulary[];
  };
  documentation?: {
    [language: string]: FieldConfigDocumentation;
  };
}

export interface PendingFieldConfig {
  name: string;
  isPending: true;
  [key: string]: any;
}

let pendingFields: Field[] = [];
let registeredFields: Field[] = [];

export const configurePendingFields = () => {
  pendingFields.forEach((field) => {
    const config = registeredFields.find(({ name }) => name === field.name)?.config;
    config && field.configure(config);
  });
  pendingFields = pendingFields.filter((field) => !field.isInitialized);
};

export class Field {
  config: FieldConfig | PendingFieldConfig;
  #isSearchField: boolean;
  #customLinkedField?: Field;
  #nameAppendix?: string;

  identity: Symbol;

  constructor(
    config: FieldConfig | PendingFieldConfig,
    isSearchField: boolean = false,
    customLinkedField?: Field,
    nameAppendix?: string
  ) {
    this.configure(config, isSearchField, customLinkedField, nameAppendix);
  }

  static createPending(name: string) {
    return new Field({
      name,
      isPending: true,
    });
  }

  configure(
    config: FieldConfig | PendingFieldConfig,
    isSearchField: boolean = this.#isSearchField,
    customLinkedField: Field = this.#customLinkedField,
    nameAppendix: string = this.#nameAppendix
  ) {
    this.config = config;
    this.#isSearchField = isSearchField;

    if (customLinkedField) {
      this.#customLinkedField = customLinkedField;
    }

    if (nameAppendix) {
      this.#nameAppendix = nameAppendix;
    }

    if (this.isInitialized && !customLinkedField && !nameAppendix) {
      registeredFields.push(this);
    }

    if (!this.isInitialized) {
      pendingFields.push(this);
    }

    this.identity = Symbol(this.linkedName);
  }

  get isInitialized() {
    return !('isPending' in this.config);
  }

  get name() {
    return this.config.name;
  }

  get linkedName() {
    const { name, linkedField } = this;
    const appendix = this.#nameAppendix ?? '';

    if (!this.#isSearchField || !linkedField) {
      return `${name}${appendix}`;
    }

    return `${name}__${linkedField.name}${appendix}`;
  }

  get renderer() {
    return this.config.renderer;
  }

  get importance() {
    return this.config.importance;
  }

  get isVirtual() {
    return this.config.isVirtual;
  }

  get isEnumerable() {
    return this.useTranslationsFrom?.isEnumerable ?? this.config.isEnumerable;
  }

  get resolvesTo(): Field[] {
    const resolvedFields = this.config.resolvesTo
      ?.map((name) => registeredFields.find((field) => field.name === name))
      .filter(Boolean)
      .map((field) =>
        this.#nameAppendix ? field.cloneAndSetNameAppendix(this.#nameAppendix, this.#isSearchField) : field
      );

    return resolvedFields?.length ? resolvedFields : [this];
  }

  get linkedField(): Field {
    const linkedFieldName = this.#customLinkedField?.name ?? this.config.linkedField;
    return linkedFieldName && registeredFields.find((field) => field.name === linkedFieldName);
  }

  get useTranslationsFrom() {
    const fieldName = this.config.useTranslationsFrom;
    return fieldName && registeredFields.find((field) => field.name === fieldName);
  }

  get label() {
    return this.useTranslationsFrom?.label ?? this.config.label;
  }

  get description() {
    return this.useTranslationsFrom?.description ?? this.config.description;
  }

  get vocabulary() {
    return this.useTranslationsFrom?.vocabulary ?? this.config.vocabulary;
  }

  get conceptIds() {
    const { isEnumerable, vocabulary, name } = this;

    if (name === 'entityName') {
      return []; // TODO edge case, to be cleaned up later when entityName has real concept IDs
    }

    return isEnumerable && vocabulary
      ? Object.keys(vocabulary)
          .flatMap((language) => vocabulary[language])
          .map(({ conceptId }) => conceptId)
          .filter(Boolean)
          .filter((conceptId, index, arr) => arr.indexOf(conceptId) === index)
      : [];
  }

  get documentation() {
    return this.config.documentation;
  }

  cloneAndLinkToField(linkedField: Field, isSearchField = true) {
    return new Field(this.config, isSearchField, linkedField);
  }

  cloneAndSetNameAppendix(appendix: string, isSearchField = true) {
    return new Field(this.config, isSearchField, null, appendix);
  }

  cloneAndSetRawValueMode() {
    return this.cloneAndSetNameAppendix('_raw_value'); // exclusive search feature
  }
}
