import { Component, h, Prop } from '@stencil/core';

let overallToggleCount = 0;

@Component({
  tag: 'mex-toggle',
  styleUrl: 'toggle.css',
})
export class ToggleComponent {
  #id: string;

  @Prop() label: string;
  @Prop() active: boolean;
  @Prop() toggleHandler: (isActive: boolean) => void;
  @Prop() disabled = false;
  @Prop() testAttr?: string;

  componentWillLoad() {
    overallToggleCount++;
    this.#id = `toggle-${overallToggleCount}`;
  }

  render() {
    const { active, label, disabled, testAttr } = this;

    return (
      <label id={this.#id} class={{ toggle: true, 'toggle--active': active, 'toggle--disabled': disabled }}>
        <span class="toggle__label">{label}</span>
        <button
          class="toggle__handle"
          onClick={() => this.toggleHandler(!active)}
          disabled={disabled}
          role="switch"
          aria-labelledby={this.#id}
          aria-checked={String(active)}
          data-test={testAttr}
          data-test-active={testAttr && String(active)}
        >
          <span class="toggle__knob">
            <mex-icon-close classes="icon--small toggle__icon toggle__icon--inactive" />
            <mex-icon-check classes="icon toggle__icon toggle__icon--active" />
          </span>
        </button>
      </label>
    );
  }
}
