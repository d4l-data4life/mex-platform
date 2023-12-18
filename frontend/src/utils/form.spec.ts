jest.mock('stencil-router-v2');

import stores from 'stores';
import { FormConfig, FormBlock, FormBlockType, FormBlockInputGroup, FormData } from 'services/form';
import { Item } from 'services/item';
import {
  addGroupFields,
  fillBlockVariableSlots,
  fillVariableSlots,
  getBlockKind,
  getGroupSetsCount,
  groupBlocks,
} from './form';
import { removeGroupFields } from './form';

const FORM: FormConfig = {
  en: {
    title: 'Test form',
    steps: [
      {
        headline: 'Test headline',
        subHeadline: 'Test sub-headline',
        blocks: [
          {
            type: FormBlockType.text,
            text: 'This is a test for context {{ context.headline }}',
          },
          {
            type: FormBlockType.list,
            text: '<ul><li>foo</li><li>bar</li><li>baz</li></ul>',
          },
          {
            type: FormBlockType.infobox,
            text: '<p>This really is a test. You are sending an email to {{ recipient.email }}</p>',
            style: 'info',
          },
          {
            type: FormBlockType.formInputText,
            name: 'firstName',
            label: 'First name',
            placeholder: 'Quincy',
            width: '1/2',
            required: true,
            inputType: 'text',
          },
          {
            type: FormBlockType.formInputText,
            name: 'lastName',
            label: 'Last name',
            placeholder: 'Doe',
            width: '1/2',
            required: true,
            inputType: 'text',
          },
          {
            type: FormBlockType.formInputTextarea,
            name: 'message',
            label: 'Your message',
            placeholder: 'Hello, I would like to ...',
            width: '1',
            required: false,
          },
          {
            type: FormBlockType.formInputRadio,
            name: 'showContactData',
            label: 'Do you want to add contact data, {{ form.firstName }}?',
            width: '1',
            required: false,
            options: [
              {
                value: 'no',
                text: 'No',
              },
              {
                value: 'yes',
                text: 'Yes, I, {{ form.firstName }} want to add contact data',
              },
            ],
            default: 'no',
          },
          {
            type: FormBlockType.formInputGroup,
            headline: 'Contact {{ index }}',
            repeatable: true,
            repetitionSettings: {
              repeatLabels: true,
              addLabel: 'Add contact data',
              deleteLabel: '',
              initialCount: 1,
              maxCount: 3,
            },
            indent: true,
            blocks: [
              {
                type: FormBlockType.formInputText,
                name: 'contactName',
                label: 'Name',
                placeholder: '{{ form.firstName }} {{ form.lastName }}',
                width: '1',
                required: false,
                inputType: 'text',
              },
              {
                type: FormBlockType.formInputText,
                name: 'contactEmail',
                label: 'Email address',
                width: '1',
                required: false,
                inputType: 'email',
              },
            ],
            conditions: [
              {
                name: 'showContactData',
                value: 'yes',
              },
            ],
          },
        ],
      },
    ],
    mailTemplate: 'foo-template',
  },
};

const formData: FormData = {
  firstName: ['Steph'],
  message: ['Hello!'],
  showContactData: ['yes'],
  contactName: ['Steph Doe', 'Oliver Doe'],
  contactEmail: ['steph@doe.org'],
};

const form = FORM.en;
const blocks = form.steps[0].blocks;

const createItemWithValues = (fieldName: string, fieldValues: string[]) =>
  ({
    values: fieldValues.map((fieldValue) => ({
      fieldName,
      fieldValue,
    })),
  } as Item);

describe('form util', () => {
  beforeAll(() => {
    jest.spyOn(stores.i18n, 'language', 'get').mockReturnValue('en');
  });

  describe('getBlockKind()', () => {
    it('returns the kind of the block (differentiation between content, input and group)', () => {
      expect(getBlockKind(blocks[0])).toBe('content');
      expect(getBlockKind(blocks[1])).toBe('content');
      expect(getBlockKind(blocks[2])).toBe('content');
      expect(getBlockKind(blocks[3])).toBe('input');
      expect(getBlockKind(blocks[5])).toBe('input');
      expect(getBlockKind(blocks[6])).toBe('input');
      expect(getBlockKind(blocks[7])).toBe('group');
      expect(getBlockKind({ ...blocks[0], type: 'foo' as any })).toBe('unknown');
    });
  });

  describe('groupBlocks', () => {
    it('groups blocks by their distinct kind', () => {
      const groups = groupBlocks(blocks);
      expect(groups[0]).toEqual([blocks[0], blocks[1], blocks[2]]);
      expect(groups[1]).toEqual([blocks[3], blocks[4], blocks[5], blocks[6]]);
      expect(groups[2]).toEqual([blocks[7]]);
    });
  });

  describe('getGroupSetsCount()', () => {
    it('counts the number of times a set of group blocks is present', () => {
      const group = blocks[7] as FormBlockInputGroup;
      expect(getGroupSetsCount(group, formData)).toBe(2);
      expect(getGroupSetsCount({ ...group, repeatable: false }, formData)).toBe(1); // group is not repeatable => 1 set of blocks
      expect(
        getGroupSetsCount(group, {
          ...formData,
          contactEmail: ['steph@doe.org', 'oliver@doe.org', 'max@doe.org'],
        })
      ).toBe(3);
      expect(
        getGroupSetsCount(group, {
          ...formData,
          contactName: [],
          contactEmail: [],
        })
      ).toBe(1); // data contradicts min count => min count is taken
      expect(
        getGroupSetsCount(group, {
          ...formData,
          contactEmail: ['steph@doe.org', 'oliver@doe.org', 'max@doe.org', 'trevor@doe.org'],
        })
      ).toBe(3); // data contradicts max count => max count is taken
      expect(getGroupSetsCount(blocks[0] as FormBlockInputGroup, formData)).toBe(0); // block is no group => 0 sets of blocks
    });
  });

  describe('fillVariableSlots()', () => {
    it('fills variable slots in a text with dynamic recipient, context, index or form data', () => {
      expect(fillVariableSlots('You have entered email address {{ form.contactEmail }}.', formData)).toBe(
        'You have entered email address steph@doe.org.'
      );
      expect(
        fillVariableSlots('This form goes to {{ recipient.name }}', formData, createItemWithValues('name', ['John']))
      ).toBe('This form goes to John');
      expect(
        fillVariableSlots(
          'Request {{ context.label }} data from {{ recipient.name }}.',
          formData,
          createItemWithValues('name', ['Mika']),
          createItemWithValues('label', ['population statistics'])
        )
      ).toBe('Request population statistics data from Mika.');
      expect(fillVariableSlots('Data set #{{ index }}', formData, null, null, 8)).toBe('Data set #9');
      expect(fillVariableSlots('This form goes to {{ recipient.name }}', formData)).toBe('This form goes to ');
    });
  });

  describe('fillBlockVariableSlots()', () => {
    it("fills variable slots in a block's text, label, placeholder or option texts", () => {
      const fillSlotsWithTestData = (block: FormBlock) =>
        fillBlockVariableSlots<any>(
          { ...block },
          { ...formData, lastName: ['Doe'] },
          createItemWithValues('email', ['caroline@doe.org']),
          createItemWithValues('headline', ['population statistics']),
          1
        );

      expect(fillSlotsWithTestData(blocks[0]).text).toBe('This is a test for context population statistics');
      expect(fillSlotsWithTestData(blocks[2]).text).toBe(
        '<p>This really is a test. You are sending an email to caroline@doe.org</p>'
      );
      expect(fillSlotsWithTestData(blocks[6]).label).toBe('Do you want to add contact data, Steph?');
      expect(fillSlotsWithTestData(blocks[6]).options[1].text).toBe('Yes, I, Steph want to add contact data');
      expect(fillSlotsWithTestData((blocks[7] as FormBlockInputGroup).blocks[0]).placeholder).toBe('Steph Doe');
    });
  });

  describe('convertFormDataToSubmit()', () => {
    test.todo('skipped to test because it is not finalized, yet');
  });

  describe('addGroupFields()', () => {
    it('adds a set of empty values of an input group to the form data', () => {
      const group = blocks[7] as FormBlockInputGroup;
      addGroupFields(group, formData);
      expect(formData.contactName).toEqual(['Steph Doe', 'Oliver Doe', '']);
      expect(formData.contactEmail).toEqual(['steph@doe.org', '', '']);
    });
  });

  describe('removeGroupFields()', () => {
    it('removes a set of values of an input group from the form data (by index)', () => {
      const group = blocks[7] as FormBlockInputGroup;
      removeGroupFields(group, formData, 0);
      expect(formData.contactName).toEqual(['Oliver Doe', '']);
      expect(formData.contactEmail).toEqual(['', '']);
    });
  });
});
