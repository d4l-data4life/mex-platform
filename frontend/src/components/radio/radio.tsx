import { Component, h, Prop, State } from '@stencil/core';

@Component({
  tag: 'mex-radio',
  styleUrl: 'radio.css',
})
export class RadioComponent {
  @Prop() label: string;
  @Prop() name: string;
  @Prop() value: string | number;

  @Prop() classes = '';
  @Prop() checked: boolean = false;
  @Prop() disabled?: boolean = false;
  @Prop() required?: boolean = false;
  @Prop() testAttr?: string;

  @Prop() handleChange?: (value: string | number) => void;
  @Prop() handleInvalid?: (event: Event) => void;

  @State() focused = false;

  render() {
    const { checked, disabled, required, focused, label, name, value, testAttr } = this;

    return (
      <label
        class={{
          radio: true,
          [this.classes]: true,
          'radio--checked': checked,
          'radio--focused': focused,
          'radio--disabled': disabled,
        }}
      >
        <input
          class="radio__input"
          type="radio"
          name={name}
          value={value}
          defaultChecked={checked}
          checked={checked}
          disabled={disabled}
          required={required}
          onChange={() => this.handleChange?.(value)}
          onInvalid={(event) => this.handleInvalid?.(event)}
          onFocus={() => (this.focused = true)}
          onBlur={() => (this.focused = false)}
          data-test={testAttr}
        />
        <span class="radio__text">{label}</span>
      </label>
    );
  }
}
