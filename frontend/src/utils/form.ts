import { FIELDS } from 'config';
import { FormStep, FormBlock, FormBlockInputGroup, FormBlockType, FormData } from 'services/form';
import { Item } from 'services/item';
import { getConcatDisplayValue } from './field';

export const getBlockKind = ({ type }: FormBlock) => {
  switch (type) {
    case FormBlockType.text:
    case FormBlockType.infobox:
    case FormBlockType.list:
      return 'content';
    case FormBlockType.formInputText:
    case FormBlockType.formInputTextarea:
    case FormBlockType.formInputRadio:
      return 'input';
    case FormBlockType.formInputGroup:
      return 'group';
    default:
      return 'unknown';
  }
};

export const groupBlocks = (blocks: FormBlock[]) => {
  return blocks?.reduce((groups, block) => {
    if (!groups.length || getBlockKind(groups[groups.length - 1][0]) !== getBlockKind(block)) {
      groups.push([block]);
    } else {
      groups[groups.length - 1].push(block);
    }

    return groups;
  }, []);
};

export const getGroupSetsCount = (group: FormBlockInputGroup, formData: FormData): number => {
  if (!group.blocks?.length) {
    return 0;
  }

  if (!group.repeatable) {
    return 1;
  }

  const count =
    group.blocks
      .filter(({ name }: any) => name)
      .map(({ name }: any) => formData[name]?.length ?? 0)
      .sort()
      .pop() ?? 0;

  const {
    repetitionSettings: { initialCount, maxCount },
  } = group;
  return Math.min(maxCount ?? Infinity, Math.max(initialCount ?? 0, count));
};

export const fillVariableSlots = (
  text: string,
  formData: FormData,
  recipient?: Item,
  context?: Item,
  index?: number
) => {
  return text.replace(/[{]{2}\s?([a-zA-Z.]+)\s?[}]{2}/gm, (_, match) => {
    const [prefix, key] = match.split('.');
    const field = prefix !== 'form' ? FIELDS[key] : null;

    switch (prefix) {
      case 'form':
        return formData?.[key]?.filter(Boolean).join(', ') ?? '';
      case 'recipient':
        return recipient ? getConcatDisplayValue([field], recipient, '') : '';
      case 'context':
        return context ? getConcatDisplayValue([field], context, '') : '';
      case 'index':
        return String((index ?? 0) + 1);
      default:
        return '';
    }
  });
};

export const fillBlockVariableSlots = <T>(
  block: any,
  formData: FormData,
  recipient?: Item,
  context?: Item,
  index?: number
): T => {
  block.text && (block.text = fillVariableSlots(block.text, formData, recipient, context, index));
  block.label && (block.label = fillVariableSlots(block.label, formData, recipient, context, index));
  block.placeholder && (block.placeholder = fillVariableSlots(block.placeholder, formData, recipient, context, index));
  block.options &&
    (block.options = block.options.map((option) => ({
      ...option,
      text: option.text && fillVariableSlots(option.text, formData, recipient, context, index),
    })));
  return block;
};

export const isGroupShown = (group: FormBlockInputGroup, formData: FormData) => {
  return group.conditions.reduce((show, { name, value }) => show && !!formData[name]?.includes(value), true);
};

export const convertFormDataToSubmit = (
  formData: FormData,
  steps: FormStep[]
): { [key: string]: string | string[] } => {
  const inputBlocks = steps
    .flatMap(({ blocks }) =>
      blocks.concat(
        blocks
          .filter((block) => block.type === FormBlockType.formInputGroup)
          .map((block) => block as FormBlockInputGroup)
          .filter((group) => isGroupShown(group, formData))
          .flatMap((group) => group.blocks.map((block) => ({ ...block, parent: group })))
      )
    )
    .filter(Boolean)
    .filter(({ type }) =>
      [FormBlockType.formInputText, FormBlockType.formInputTextarea, FormBlockType.formInputRadio].includes(type)
    );

  return inputBlocks.reduce((data, { name, parent }: FormBlock & { name: string }) => {
    const valueCount = parent ? getGroupSetsCount(parent, formData) : 1;
    const values = new Array(valueCount).fill(null).map((_, index) => formData[name]?.[index] ?? '');
    return Object.assign(data, values?.some(Boolean) ? { [name]: parent?.repeatable ? values : values[0] } : {});
  }, {});
};

export const addGroupFields = (group: FormBlockInputGroup, formData: FormData) => {
  const { blocks } = group;
  const count = getGroupSetsCount(group, formData);

  new Array(count + 1).fill(null).forEach((_, index) =>
    blocks
      .filter(({ name }: any) => name)
      .forEach(({ name }: any) => {
        if (!formData[name]) {
          formData[name] = [];
        }

        formData[name][index] ??= '';
      })
  );

  return formData;
};

export const removeGroupFields = (group: FormBlockInputGroup, formData: FormData, index: number) => {
  const { blocks } = group;
  blocks
    .filter(({ name }: any) => name)
    .forEach(({ name }: any) => {
      formData[name] = [...(formData[name]?.filter((_, dataIndex) => dataIndex !== index) ?? [])];
    });

  return formData;
};
