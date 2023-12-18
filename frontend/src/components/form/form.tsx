import { Component, Element, Event, EventEmitter, Fragment, Host, Prop, State, h } from '@stencil/core';
import services from 'services';
import { ContentUnitText, ContentUnitList, ContentUnitInfobox } from 'services/content';
import {
  FormBlockType,
  FormBlock,
  FormConfig,
  FormStep,
  FormBlockInputText,
  FormBlockInputTextarea,
  FormBlockInputGroup,
  FormBlockInputRadio,
} from 'services/form';
import { Item } from 'services/item';
import stores from 'stores';
import { ContentBlockInfobox, ContentBlockList, ContentBlockText } from 'components/content-blocks/content-blocks';
import { catchRetryableAction } from 'utils/error';
import {
  addGroupFields,
  fillBlockVariableSlots,
  fillVariableSlots,
  getBlockKind,
  getGroupSetsCount,
  groupBlocks,
  isGroupShown,
  removeGroupFields,
} from 'utils/form';

@Component({
  tag: 'mex-form',
  styleUrl: 'form.css',
})
export class FormComponent {
  #renderedPromiseResolve?: (value: unknown) => void;

  @Element() hostEl: HTMLElement;

  @Prop() formId: string;
  @Prop() formKey?: string;
  @Prop() recipient?: Item;
  @Prop() context?: Item;
  @Prop() embedded?: boolean;

  @State() config: FormConfig;
  @State() currentFormStep: FormStep;
  @State() formData: { [key: string]: string[] } = {};
  @State() isSubmitting: boolean = false;

  @Event() closeModal: EventEmitter;
  @Event() scrollModalToTop: EventEmitter;

  get currentConfig() {
    return this.config?.[stores.i18n.language];
  }

  get currentFormStepIndex() {
    return this.currentConfig?.steps?.indexOf(this.currentFormStep) ?? 0;
  }

  fillVariableSlots(text: string, index?: number) {
    const { formData, recipient, context } = this;
    return fillVariableSlots(text, formData, recipient, context, index);
  }

  fillBlockVariableSlots<T>(block, index?: number): T {
    const { formData, recipient, context } = this;
    return fillBlockVariableSlots<T>(block, formData, recipient, context, index);
  }

  goToPreviousStep() {
    const previousStep = this.currentConfig.steps[this.currentFormStepIndex - 1];
    if (!previousStep) {
      return this.embedded && this.closeModal.emit();
    }

    this.currentFormStep = previousStep;

    this.focusFirstInput();
  }

  goToNextStep() {
    this.focusFirstInput();

    const nextStep = this.currentConfig.steps[this.currentFormStepIndex + 1];
    if (!nextStep) {
      return catchRetryableAction(
        async () => {
          this.isSubmitting = true;

          await services.form.submit({
            form: this.config,
            language: stores.i18n.language,
            formData: this.formData,
            contextItemId: this.context?.itemId,
            recipientItemId: this.recipient?.itemId,
          });

          const { formKey } = this;
          stores.notifications.add(`${formKey ? `${formKey}Form` : 'form'}.submit:success`);
          this.closeModal.emit();
        },
        true,
        () => (this.isSubmitting = false)
      );
    }

    this.currentFormStep = nextStep;
  }

  async focusFirstInput() {
    await new Promise((resolve) => (this.#renderedPromiseResolve = resolve));
    (this.hostEl.querySelector('input, textarea, button[type=submit]') as HTMLInputElement)?.focus?.();
    this.scrollModalToTop.emit();
  }

  handleSubmit(event: Event) {
    event.preventDefault();
    this.goToNextStep();
  }

  handleLanguageChange = () => {
    // we do not know for certain if form is the same for all languages - jump back to first step
    this.currentFormStep = this.currentConfig.steps[0];
    this.focusFirstInput();
  };

  componentWillLoad() {
    stores.i18n.addListener(this.handleLanguageChange);
  }

  componentDidLoad() {
    catchRetryableAction(async () => {
      this.config = await services.form.fetch(this.formId);
      this.currentFormStep = this.currentConfig.steps[0];
      this.focusFirstInput();
    });
  }

  componentDidRender() {
    this.#renderedPromiseResolve?.(true);
  }

  disconnectedCallback() {
    stores.i18n.removeListener(this.handleLanguageChange);
  }

  renderTextInputBlock(block: FormBlockInputText | FormBlockInputTextarea, key: string, index = 0, showLabel = true) {
    return (
      <mex-form-input-text
        key={key}
        index={index}
        name={block.name}
        value={this.formData[block.name]?.[index]}
        type={(block as FormBlockInputText).inputType ?? 'textarea'}
        label={block.label}
        showLabel={showLabel}
        placeholder={block.placeholder}
        width={block.width}
        required={block.required}
        testAttr={`form.input.${block.name}`}
        ariaLabelAttr={block.label}
        handleInput={(value) => {
          if (!this.formData[block.name]) {
            this.formData[block.name] = [];
          }

          this.formData[block.name][index] = value;
          this.formData = { ...this.formData }; // re-render
        }}
      />
    );
  }

  renderRadioInputBlock(block: FormBlockInputRadio, key: string, index = 0, showLabel = true) {
    return (
      <mex-form-input-radio
        key={key}
        name={block.name}
        options={block.options}
        value={this.formData[block.name]?.[index]}
        default={block.default}
        label={block.label}
        showLabel={showLabel}
        width={block.width}
        required={block.required}
        testAttr={`form.input.${block.name}`}
        handleChange={(value) => {
          if (!this.formData[block.name]) {
            this.formData[block.name] = [];
          }

          this.formData[block.name][index] = value;
          this.formData = { ...this.formData }; // re-render
        }}
      />
    );
  }

  renderGroupBlocks(group: FormBlockInputGroup, key: string, index: number = 0) {
    return groupBlocks(group.blocks).map((blocks) => {
      return this.renderBlocks(
        blocks,
        key,
        index,
        group.repeatable ? group.repetitionSettings.repeatLabels || !index : true
      );
    });
  }

  renderGroup(group: FormBlockInputGroup, key: string) {
    const { formData } = this;
    const { blocks, repeatable, headline, indent } = group;
    const { repeatLabels, maxCount, addLabel, deleteLabel } = group.repetitionSettings ?? {};
    const countOfSets = getGroupSetsCount(group, formData);
    const canAdd = group.repeatable && (maxCount ?? Infinity) > countOfSets;

    return (
      isGroupShown(group, formData) &&
      !!blocks && (
        <Fragment>
          {new Array(countOfSets).fill(null).map((_, index) => (
            <div
              class={{
                form__group: true,
                'form__group--multiline': repeatable && repeatLabels,
                'form__group--indent': indent,
              }}
            >
              {!repeatable && headline && (
                <h5 class="form__group-headline">{this.fillVariableSlots(headline, index)}</h5>
              )}
              {this.renderGroupBlocks(group, key, index)}
              {repeatable && (
                <div class="form__group-actions">
                  {headline && <h5 class="form__group-headline">{this.fillVariableSlots(headline, index)}</h5>}
                  <button
                    type="button"
                    class="button button--tertiary"
                    onClick={() => (this.formData = removeGroupFields(group, { ...formData }, index))}
                  >
                    <mex-icon-trash classes="icon--large icon--inline" />
                    <span>{deleteLabel}</span>
                  </button>
                </div>
              )}
            </div>
          ))}
          {canAdd && (
            <button
              type="button"
              class="form__add button button--secondary"
              onClick={() => (this.formData = addGroupFields(group, { ...formData }))}
            >
              {addLabel ?? stores.i18n.t('form.add')}
            </button>
          )}
        </Fragment>
      )
    );
  }

  renderBlock(block: FormBlock, key: string, index = 0, showLabel = true) {
    const childKey = `${key}-${(block as any).name ?? 'misc'}-${index}`;

    switch (block.type) {
      case FormBlockType.text:
        return <ContentBlockText content={this.fillBlockVariableSlots<ContentUnitText>({ ...block }, index)} />;
      case FormBlockType.list:
        return <ContentBlockList content={this.fillBlockVariableSlots<ContentUnitList>({ ...block }, index)} />;
      case FormBlockType.infobox:
        return <ContentBlockInfobox content={this.fillBlockVariableSlots<ContentUnitInfobox>({ ...block }, index)} />;
      case FormBlockType.formInputText:
        return this.renderTextInputBlock(
          this.fillBlockVariableSlots<FormBlockInputText>({ ...block }, index),
          childKey,
          index,
          showLabel
        );
      case FormBlockType.formInputTextarea:
        return this.renderTextInputBlock(
          this.fillBlockVariableSlots<FormBlockInputTextarea>({ ...block }, index),
          childKey,
          index,
          showLabel
        );
      case FormBlockType.formInputRadio:
        return this.renderRadioInputBlock(
          this.fillBlockVariableSlots<FormBlockInputRadio>({ ...block }, index),
          childKey,
          index,
          showLabel
        );
      case FormBlockType.formInputGroup:
        return this.renderGroup(block as unknown as FormBlockInputGroup, childKey);
    }
  }

  renderBlocks(blocks: FormBlock[], key: string, index = 0, showLabels = true) {
    const blockKind = getBlockKind(blocks[0]);

    return (
      <div key={key} class={`form__blocks form__blocks--${blockKind}`}>
        {blocks.map((block) => this.renderBlock(block, key, index, showLabels))}
      </div>
    );
  }

  render() {
    const { currentConfig: config, currentFormStep, currentFormStepIndex, embedded, isSubmitting } = this;
    const groupedBlocks = groupBlocks(this.currentFormStep?.blocks);
    const isLastStep = currentFormStepIndex === (config?.steps?.length ?? 0) - 1;
    const { language } = stores.i18n;

    return (
      <Host class="form">
        {!config && <mex-placeholder lines={10} />}

        {!!config && !!currentFormStep && !isSubmitting && (
          <Fragment>
            {!!config.title && (
              <h2 class="form__title">{this.fillVariableSlots(config.title, currentFormStepIndex)}</h2>
            )}
            {!!currentFormStep.headline && (
              <h3 class="form__headline">{this.fillVariableSlots(currentFormStep.headline, currentFormStepIndex)}</h3>
            )}
            {!!currentFormStep.subHeadline && (
              <h4
                class="form__sub-headline"
                innerHTML={this.fillVariableSlots(currentFormStep.subHeadline, currentFormStepIndex)}
              />
            )}

            <form onSubmit={(event) => this.handleSubmit(event)}>
              <div class="form__container">
                {groupedBlocks?.map((blocks, index) =>
                  this.renderBlocks(blocks, `${language}-step-${currentFormStepIndex}-blocks-${index}`)
                )}
              </div>
              <div class="form__actions">
                <button
                  class="button button--secondary"
                  type="button"
                  onClick={() => this.goToPreviousStep()}
                  disabled={!currentFormStepIndex && !embedded}
                >
                  <mex-icon-arrow classes="icon--inline icon--no-vertical-offset icon--medium icon--mirrored-horizontal" />
                  <span>{stores.i18n.t('form.back')}</span>
                </button>
                <button class="button" type="submit">
                  <span>{stores.i18n.t(isLastStep ? 'form.submit' : 'form.next')}</span>
                  {!isLastStep && <mex-icon-arrow classes="icon--inline icon--no-vertical-offset icon--medium" />}
                </button>
              </div>
            </form>
          </Fragment>
        )}

        {this.isSubmitting && <mex-logo class="form__busy" loader />}
      </Host>
    );
  }
}
