import { Item } from 'services/item';

export const serializeItemData = (item: Item) =>
  item?.values?.reduce(
    (data, { fieldName, fieldValue }) => Object.assign(data, { [fieldName]: [...(data[fieldName] ?? []), fieldValue] }),
    {}
  ) ?? {};
