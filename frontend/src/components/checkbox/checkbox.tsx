import { Component, h, Prop, State } from '@stencil/core';
import { CheckStatusEnum } from 'utils/filters';

@Component({
  tag: 'mex-checkbox',
  styleUrl: 'checkbox.css',
})
export class CheckboxComponent {
  @Prop() label: string;
  @Prop() secondaryText?: string;

  @Prop() classes = '';
  @Prop() checked: CheckStatusEnum = CheckStatusEnum.UNCHECKED;
  @Prop() disabled = false;
  @Prop() required = false;
  @Prop() handleChange?: (checked: CheckStatusEnum) => void;
  @Prop() testAttr?: string;

  @State() focused = false;

  onChange() {
    this.handleChange?.(this.checked);
  }

  render() {
    const { checked, disabled, focused, label, required, secondaryText, testAttr } = this;

    const isChecked = checked === CheckStatusEnum.CHECKED;

    return (
      <label
        class={{
          checkbox: true,
          [this.classes]: true,
          'checkbox--checked': isChecked,
          'checkbox--semi-checked': checked === CheckStatusEnum.SEMI,
          'checkbox--focused': focused,
          'checkbox--disabled': disabled,
        }}
      >
        <mex-icon-check class="checkbox__check" />
        <input
          class="checkbox__input"
          type="checkbox"
          defaultChecked={isChecked}
          checked={isChecked}
          disabled={disabled}
          required={required}
          onChange={() => this.onChange()}
          onFocus={() => (this.focused = true)}
          onBlur={() => (this.focused = false)}
          data-test={testAttr}
        />
        <span class="checkbox__text checkbox__text--primary">{label}</span>
        {!!secondaryText && <span class="checkbox__text checkbox__text--secondary">{secondaryText}</span>}
      </label>
    );
  }
}
