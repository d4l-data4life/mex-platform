import { API_URL, CONFIG_URL } from 'config';
import { get, post } from 'utils/fetch-client';
import { ContentUnitText, ContentUnitList, ContentUnitInfobox } from './content';
import { convertFormDataToSubmit } from 'utils/form';

export interface FormData {
  [key: string]: string[];
}

export enum FormBlockType {
  text = 'text',
  list = 'list',
  infobox = 'infobox',
  formInputText = 'formInputText',
  formInputTextarea = 'formInputTextarea',
  formInputRadio = 'formInputRadio',
  formInputGroup = 'formInputGroup',
}

interface FormBlockInputBase {
  name: string;
  label: string;
  placeholder?: string;
  width: '1' | '2/3' | '1/2' | '1/3';
  required: boolean;
}

export interface FormBlockInputText extends FormBlockInputBase {
  inputType: 'text' | 'email';
}

export interface FormBlockInputTextarea extends FormBlockInputBase {}

export interface FormBlockInputRadioOption {
  value: string;
  text: string;
}

export interface FormBlockInputRadio extends FormBlockInputBase {
  options: FormBlockInputRadioOption[];
  default?: string;
}

export interface FormBlockInputGroupCondition {
  name: string;
  value: string;
}

export interface FormBlockInputGroup {
  headline?: string;
  repeatable: boolean;
  repetitionSettings?: {
    repeatLabels: boolean;
    addLabel: string;
    deleteLabel: string;
    initialCount?: number;
    maxCount?: number;
  };
  indent: boolean;
  blocks: FormBlock[];
  conditions: FormBlockInputGroupCondition[];
}

export type FormBlock = {
  type: FormBlockType;
  parent?: FormBlockInputGroup;
} & (
  | ContentUnitText
  | ContentUnitList
  | ContentUnitInfobox
  | FormBlockInputText
  | FormBlockInputTextarea
  | FormBlockInputRadio
  | FormBlockInputGroup
);

export interface FormStep {
  headline?: string;
  subHeadline?: string;
  blocks: FormBlock[];
}

export interface FormConfig {
  [language: string]: {
    title: string;
    steps: FormStep[];
    mailTemplate: string;
  };
}

export class FormService {
  async fetch(formId: string): Promise<FormConfig> {
    const [response] = await get<FormConfig>({
      url: `${CONFIG_URL}/forms/${formId}`,
    });

    return response;
  }

  async submit({
    form,
    language,
    formData,
    contextItemId,
    recipientItemId,
  }: {
    form: FormConfig;
    language: string;
    formData: FormData;
    contextItemId: string;
    recipientItemId: string;
  }): Promise<void> {
    await post({
      url: `${API_URL}/notify`,
      authorized: true,
      body: JSON.stringify({
        templateInfo: {
          templateName: form[language].mailTemplate,
          contextItemId,
          recipientItemId,
        },
        formData: convertFormDataToSubmit(formData, form[language].steps),
      }),
    });
  }
}

export default new FormService();
