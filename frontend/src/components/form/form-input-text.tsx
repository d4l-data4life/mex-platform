import { Component, Host, Prop, State, h } from '@stencil/core';

@Component({
  tag: 'mex-form-input-text',
  styleUrl: 'form-input.css',
})
export class FormInputTextComponent {
  #inputEl?: HTMLInputElement | HTMLTextAreaElement;

  @Prop() name: string;
  @Prop() value?: string = '';
  @Prop() index?: number = 0;
  @Prop() type?: string = 'text';
  @Prop() label?: string;
  @Prop() showLabel?: boolean = true;
  @Prop() placeholder?: string;
  @Prop() ariaLabelAttr?: string;
  @Prop() testAttr?: string;
  @Prop() required?: boolean = false;
  @Prop() width: '1' | '2/3' | '1/2' | '1/3' = '1';

  @Prop() handleFocus?: () => void;
  @Prop() handleBlur?: () => void;
  @Prop() handleInput?: (value: string) => void;

  @State() isPristine: boolean = true;

  get id() {
    return `${this.name}-${this.index}`;
  }

  async autoResize() {
    if (!this.#inputEl || this.type !== 'textarea') {
      return;
    }

    this.#inputEl.style.setProperty('--scroll-height', '0px');
    await new Promise((resolve) => requestAnimationFrame(resolve));
    this.#inputEl.style.setProperty('--scroll-height', this.#inputEl.scrollHeight + 'px');
  }

  componentDidLoad() {
    this.autoResize();
  }

  render() {
    const baseClass = 'form-input';
    const Tag = this.type === 'textarea' ? 'textarea' : 'input';

    return (
      <Host
        class={`${baseClass} ${baseClass}--${this.type} ${this.required ? `${baseClass}--required` : ''}`}
        style={{ width: `calc(${this.width} * 100% - ${this.width !== '1' ? 'var(--spacing-medium) / 2' : '0px'})` }}
      >
        {this.showLabel && !!this.label && (
          <label class={`${baseClass}__label`} htmlFor={this.id}>
            {this.label}
          </label>
        )}
        <Tag
          ref={(el) => (this.#inputEl = el)}
          id={this.id}
          name={this.name}
          type={this.type}
          class={`${baseClass}__input ${this.isPristine ? `${baseClass}__input--pristine` : ''}`}
          placeholder={this.placeholder}
          title={this.showLabel ? null : this.label}
          aria-label={this.ariaLabelAttr}
          value={this.value}
          onFocus={this.handleFocus}
          onBlur={() => {
            this.isPristine = false;
            this.handleBlur?.();
          }}
          onInput={(event) => {
            this.autoResize();
            this.handleInput?.((event.target as HTMLInputElement | HTMLTextAreaElement).value);
          }}
          data-test={this.testAttr}
          required={this.required}
          onInvalid={() => (this.isPristine = false)}
        />
      </Host>
    );
  }
}
