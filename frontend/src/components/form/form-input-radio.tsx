import { Component, Host, Prop, State, h } from '@stencil/core';
import { FormBlockInputRadioOption } from 'services/form';

@Component({
  tag: 'mex-form-input-radio',
  styleUrl: 'form-input.css',
})
export class FormInputTextComponent {
  @Prop() name: string;
  @Prop() options: FormBlockInputRadioOption[];
  @Prop() value?: string;
  @Prop() default?: string;
  @Prop() label?: string;
  @Prop() showLabel?: boolean = true;
  @Prop() testAttr?: string;
  @Prop() required?: boolean = false;
  @Prop() width: '1' | '2/3' | '1/2' | '1/3' = '1';

  @Prop() handleFocus?: () => void;
  @Prop() handleBlur?: () => void;
  @Prop() handleChange?: (value: string) => void;

  @State() isPristine: boolean = true;

  componentWillLoad() {
    this.default && !this.value && this.handleChange?.(this.default);
  }

  render() {
    const baseClass = 'form-input';

    return (
      <Host
        class={`${baseClass} ${baseClass}--radio ${this.required ? `${baseClass}--required` : ''}`}
        style={{ width: `calc(${this.width} * 100% - ${this.width !== '1' ? 'var(--spacing-medium) / 2' : '0px'})` }}
      >
        {this.showLabel && !!this.label && <label class={`${baseClass}__label`}>{this.label}</label>}

        <div class={`${baseClass}__options`}>
          {this.options.map((option, index) => (
            <mex-radio
              classes={`${baseClass}__option ${this.isPristine ? `${baseClass}__option--pristine` : ''}`}
              key={`${this.name}-${index}`}
              label={option.text ?? option.value}
              name={this.name}
              required={this.required}
              value={option.value}
              checked={this.value === option.value}
              handleChange={(value) => this.handleChange?.(String(value))}
              handleInvalid={() => (this.isPristine = false)}
              testAttr={this.testAttr}
            />
          ))}
        </div>
      </Host>
    );
  }
}
