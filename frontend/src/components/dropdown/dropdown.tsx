import { Component, h, Host, Method, Prop, State } from '@stencil/core';
import stores from 'stores';
import { IS_MOBILE, IS_TOUCH } from 'utils/device';

export interface DropdownOption {
  value: any;
  label: string;
}

@Component({
  tag: 'mex-dropdown',
  styleUrl: 'dropdown.css',
})
export class DropdownComponent {
  #dropdownRef: HTMLSelectElement;
  #accordionRef: HTMLMexAccordionElement;

  @Prop() options?: DropdownOption[];
  @Prop() handleChange?: (value: any) => void;
  @Prop() handleExpand?: () => void;
  @Prop() handleCollapse?: () => void;
  @Prop() label = '';
  @Prop() value?: any;
  @Prop() required = false;
  @Prop() orientation: 'left' | 'right' = 'left';
  @Prop() toggleClass = 'dropdown__toggle';
  @Prop() disabled = false;
  @Prop() testAttr?: string;
  @Prop() withLabelsTranslation = true;

  @State() isExpanded = false;
  @State() isFocused = false;
  @State() isFocusVisible = false;

  get selectedOption() {
    return this.options?.find(({ value }) => this.value === value) ?? this.options?.[0];
  }

  toggle(force = false) {
    if (force || !this.options || !IS_MOBILE || !IS_TOUCH) {
      this.isExpanded ? this.collapse() : this.expand();
    }
  }

  expand() {
    this.isExpanded = true;
    this.handleExpand?.();
  }

  @Method() async collapse() {
    this.isExpanded = false;
    this.handleCollapse?.();
  }

  change(event: Event) {
    const index = parseInt((event.target as HTMLInputElement).value, 10);
    this.select(this.options?.[index]);
  }

  select(option: DropdownOption, collapse = false) {
    collapse && this.collapse();
    this.handleChange?.(option?.value);
    this.#dropdownRef?.focus();
  }

  keydown(event) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      this.toggle(true);
    }

    if (event.key === 'Escape') {
      this.collapse();
    }

    // Mac OS accessibility fixes below (works out of the box in sane OSes)
    const { options, selectedOption } = this;
    if (!options) {
      return;
    }

    if (event.key === 'ArrowUp' || event.key === 'ArrowLeft') {
      event.preventDefault();
      this.select(options[options.indexOf(selectedOption) - 1]);
    }

    if (event.key === 'ArrowDown' || event.key === 'ArrowRight') {
      event.preventDefault();
      this.select(options[options.indexOf(selectedOption) + 1]);
    }
  }

  getTranslatedLabel(option) {
    return this.withLabelsTranslation ? stores.i18n.t(option?.label) : option?.label;
  }

  render() {
    const { selectedOption, options, orientation, label, testAttr } = this;

    const selectAttrs = {
      class: 'dropdown__input',
      onClick: () => {
        this.isFocusVisible = false;
        this.toggle();
      },
      onChange: (event) => this.change(event),
      onFocus: () => {
        this.isFocusVisible = true;
        this.isFocused = true;
      },
      onBlur: () => {
        this.isFocused = false;
        requestAnimationFrame(() => {
          if (
            this.isExpanded &&
            document.activeElement !== document.body &&
            !this.#accordionRef?.contains(document.activeElement)
          ) {
            this.collapse();
          }
        });
      },
      onKeyDown: (event) => this.keydown(event),
      disabled: this.disabled,
      'aria-expanded': String(this.isExpanded),
    };

    return (
      <Host
        class={{
          dropdown: true,
          'dropdown--focused': this.isFocused && this.isFocusVisible,
          'dropdown--expanded': this.isExpanded,
          'dropdown--disabled': this.disabled,
        }}
        data-test={testAttr}
        data-test-active={testAttr && String(this.isExpanded)}
      >
        {this.isExpanded && <div class="dropdown__backdrop" onClick={() => this.collapse()} />}
        <label class={this.toggleClass}>
          {!!label && <span class="dropdown__label">{label}</span>}
          <slot name="label">{this.getTranslatedLabel(selectedOption)}</slot>
          <mex-icon-chevron
            classes={`icon--inline icon--mirrorable ${this.isExpanded ? 'icon--mirrored-vertical' : ''}`}
          />
          {options ? (
            <select
              ref={(el) => (this.#dropdownRef = el)}
              {...selectAttrs}
              data-test={testAttr && `${testAttr}:select`}
              disabled={this.disabled}
            >
              {options.map((option, index) => (
                <option
                  key={option.value}
                  value={index}
                  selected={selectedOption === option}
                  data-test={testAttr && `${testAttr}:select:option`}
                  data-test-context={testAttr && typeof option.value === 'string' && option.value}
                  data-test-key={testAttr && index}
                  data-test-active={testAttr && String(selectedOption === option)}
                >
                  {this.getTranslatedLabel(option)}
                </option>
              ))}
            </select>
          ) : (
            <button {...selectAttrs} data-test={testAttr && `${testAttr}:toggle`} disabled={this.disabled}></button>
          )}
        </label>
        <div class={`dropdown__list dropdown__list--${orientation}`}>
          <mex-accordion
            expanded={this.isExpanded}
            testAttr={testAttr && `${testAttr}:list`}
            ref={(el) => (this.#accordionRef = el)}
          >
            <slot name="options">
              <ul class="dropdown__items" role="listbox">
                {options?.map((option, index) => (
                  <li
                    class={{ dropdown__item: true, 'dropdown__item--selected': selectedOption === option }}
                    key={option.value}
                    role="option"
                    onClick={() => this.select(option, true)}
                    data-test={testAttr && `${testAttr}:list:item`}
                    data-test-context={testAttr && typeof option.value === 'string' && option.value}
                    data-test-key={testAttr && index}
                    data-test-active={testAttr && String(selectedOption === option)}
                  >
                    {this.getTranslatedLabel(option)}
                  </li>
                ))}
              </ul>
            </slot>
          </mex-accordion>
        </div>
      </Host>
    );
  }
}
