import { FunctionalComponent, h } from '@stencil/core';
import { getClass, getClasses } from 'utils/content';
import {
  ContentUnitCode,
  ContentUnitCompletenessDocumentationTable,
  ContentUnitEntityTypeHeadline,
  ContentUnitFieldDescription,
  ContentUnitHeading,
  ContentUnitImage,
  ContentUnitInfobox,
  ContentUnitList,
  ContentUnitMarkdown,
  ContentUnitTable,
  ContentUnitText,
} from 'services/content';
import { CONFIG_URL, ENTITY_TYPES, FIELDS } from 'config';
import { Item } from 'services/item';
import stores from 'stores';
import { formatValue, getDocumentationProp, translateFieldDescription, translateFieldName } from 'utils/field';
import { Field, FieldKind, FieldRenderer } from 'models/field';

export const ContentBlockHeading: FunctionalComponent<{
  content: ContentUnitHeading;
  ref?: (el: HTMLElement) => void;
}> = ({ content, ref }) => {
  const { id, level = 'h1', text = '' } = content;
  const TagName = level;
  return <TagName id={id} ref={ref} innerHTML={text} />;
};

export const ContentBlockText: FunctionalComponent<{
  content: ContentUnitText;
}> = ({ content }) =>
  content.text?.indexOf('<p>') === 0 ? <div innerHTML={content.text} /> : <p innerHTML={content.text} />;

export const ContentBlockList: FunctionalComponent<{
  content: ContentUnitList;
}> = ({ content }) => <div innerHTML={content.text} />;

export const ContentBlockMarkdown: FunctionalComponent<{
  content: ContentUnitMarkdown;
}> = ({ content }) => <div innerHTML={content.text} />;

export const ContentBlockCode: FunctionalComponent<{
  content: ContentUnitCode;
}> = ({ content }) => <code innerHTML={content.code} />;

export const ContentBlockInfobox: FunctionalComponent<{ content: ContentUnitInfobox }> = ({ content }) => (
  <div class={getClasses(['infobox-container', `infobox-container--${content.style}`])} innerHTML={content.text}>
    {content.style === 'info' && (
      <mex-icon-info
        class={getClass('infobox-container--icon-info')}
        classes="icon--large u-underline-3"
        hasBackground
      />
    )}
  </div>
);

export const ContentBlockTable: FunctionalComponent<{ content: ContentUnitTable; props?: object }> = ({
  content,
  props = {},
}) => {
  const { hasheadercolumn, hasheaderrow, table } = content;
  const headerData = hasheaderrow && table?.shift();

  const getColSpan = function (index, rowLength) {
    const isLast = index === rowLength - 1;
    const diff = (Math.max(...table.map((row) => row?.length)) ?? 0) - rowLength;
    return isLast ? diff + 1 : 1;
  };

  return (
    <table class={getClass('table')} {...props}>
      <colgroup>
        <col class={getClass('table-first-col')} />
      </colgroup>
      {hasheaderrow && (
        <tr>
          {headerData?.map((headerHTML, index) => {
            const classes = ['table-cell', 'table-cell--header'];
            return (
              <th colSpan={getColSpan(index, headerData.length)} class={getClasses(classes)} innerHTML={headerHTML} />
            );
          })}
        </tr>
      )}
      {table?.map((rowData) => (
        <tr>
          {rowData.map((cellHTML, cellIndex) => {
            const classes = [
              'table-cell',
              ...(hasheadercolumn && cellIndex === 0 && rowData.length > 1 ? ['table-cell--header'] : []),
            ];
            return (
              <td colSpan={getColSpan(cellIndex, rowData.length)} class={getClasses(classes)} innerHTML={cellHTML} />
            );
          })}
        </tr>
      ))}
    </table>
  );
};

export const ContentBlockFigure: FunctionalComponent<{ content: ContentUnitImage }> = ({ content }) => {
  const {
    alt = '',
    caption = '',
    link = '',
    src = '',
    width = 'auto',
    height = 'auto',
    min_width = 'auto',
    min_height = 'auto',
  } = content as ContentUnitImage;

  const url = src.replace(/.*\/cms/, CONFIG_URL);

  let img = (
    <img loading="lazy" src={url} alt={alt} style={{ width, height, minWidth: min_width, minHeight: min_height }} />
  );

  if (!!link) {
    img = (
      <a href={link} target="_blank" rel="noopener noreferrer">
        {img}
      </a>
    );
  }

  return (
    <figure class={getClass('illustration')}>
      {img}
      {!!caption && <figcaption innerHTML={caption} />}
    </figure>
  );
};

export const ContentBlockEntityTypeHeadline: FunctionalComponent<{
  content: ContentUnitEntityTypeHeadline;
  ref?: (el: HTMLElement) => void;
}> = ({ content, ref }) => {
  const { id, link, entitytype: entityType } = content as ContentUnitEntityTypeHeadline;
  return (
    <mex-item-entity-headline
      id={id}
      ref={ref}
      class={getClass('entity-type-headline')}
      item={{ entityType } as Item}
      link={link ? `/${link}` : null}
    />
  );
};

const getFieldType = (field: Field) => {
  const { isEnumerable, renderer, kind } = field?.config ?? {};
  if (isEnumerable) {
    return 'controlledList';
  }

  if (renderer === FieldRenderer.link) {
    return 'url';
  }

  switch (kind) {
    case FieldKind.link:
      return 'id';
    case FieldKind.timestamp:
      return 'date';
    case FieldKind.string:
    case FieldKind.text:
      return 'text';
    case FieldKind.coding:
      return 'coding';
  }
};

export const ContentBlockFieldDescription: FunctionalComponent<{
  content: ContentUnitFieldDescription;
  ref?: (el: HTMLElement) => void;
}> = ({ content, ref }) => {
  const field = FIELDS[content.field];
  const entityType = ENTITY_TYPES[content.entitytype];
  if (!field?.isInitialized) {
    return null;
  }

  const t = (key: string) => stores.i18n.t(`content.documentation.${key}`);
  const item = { entityType: entityType?.name, values: [] } as Item;
  const fieldLabel = translateFieldName(field, item);
  const [descriptionText, exampleValue, furtherInformation, displayFormats] = [
    'descriptionText',
    'exampleValue',
    'furtherInformation',
    'displayFormats',
  ].map((prop) => getDocumentationProp(field, prop, entityType?.name));
  const description = descriptionText ?? translateFieldDescription(field, item);
  const vocabulary = field.vocabulary?.[stores.i18n.language]?.map(({ label }) => label);
  const { importance, isMultiValued } = field.config;
  const fieldType = getFieldType(field);
  const tags = [`${importance}Importance`, `${isMultiValued ? 'multi' : 'single'}Value`, `${fieldType}Type`]
    .map((key) => t(key))
    .filter((translation) => translation !== '-');

  const table = [
    [fieldLabel],
    [t('name'), field.name],
    ...(description ? [[t('description'), description]] : []),
    ...(exampleValue ? [[t('exampleValue'), exampleValue]] : []),
    [t('properties'), tags.map((tag) => `<mex-tag text="${tag}" closable=false></mex-tag>`).join('')],
    ...(vocabulary ? [[t('vocabulary'), `<ul>${vocabulary.map((item) => `<li>${item}</li>`).join('')}</ul>`]] : []),
    ...(furtherInformation ? [[t('furtherInformation'), furtherInformation]] : []),
    ...(displayFormats ? [[t('displayFormats'), displayFormats]] : []),
  ];

  return (
    <ContentBlockTable
      content={{ hasheadercolumn: true, hasheaderrow: true, table }}
      props={{
        id: content.id,
        ref,
        class: getClasses(['table', 'field-description']),
      }}
    />
  );
};

export const ContentBlockCompletenessDocumentationTable: FunctionalComponent<{
  content: ContentUnitCompletenessDocumentationTable;
}> = ({ content }) => {
  const t = (key: string) => stores.i18n.t(`content.completeness.${key}`);
  const { label, description, entityTypes } = content;

  const IMPORTANCE_LEVELS = Object.keys(entityTypes?.[0].fieldsByImportance) ?? [];
  const IMPORTANCE_TAG_CLASSES = {
    mandatory: 'red',
    recommended: 'green',
    optional: 'yellow',
  };

  const table = [
    [label],
    [description],
    [
      t('entityType'),
      ...IMPORTANCE_LEVELS.map(
        (importance) =>
          `<span class="mex-text-tag mex-text-tag--color-${IMPORTANCE_TAG_CLASSES[importance]}">${t(importance)}</span>`
      ),
    ],
    ...entityTypes.map(({ name, comment, fieldsByImportance }) => {
      const renderFields = (importance) => fieldsByImportance[importance].map((field) => `<li>${field}</li>`).join('');

      return [
        formatValue([name], FIELDS.entityName),
        ...(comment
          ? [comment]
          : IMPORTANCE_LEVELS.map((importance) =>
              fieldsByImportance[importance]?.length > 0 ? `<ul>${renderFields(importance)}</ul>` : 'n/a'
            )),
      ];
    }),
  ];
  return (
    <ContentBlockTable
      content={{ hasheadercolumn: true, hasheaderrow: true, table }}
      props={{
        class: getClasses(['table', 'completeness-documentation-table']),
      }}
    />
  );
};
