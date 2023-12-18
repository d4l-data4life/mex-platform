import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-accordion',
  styleUrl: 'accordion.css',
})
export class AccordionComponent {
  #outerEl: HTMLElement;
  #innerEl: HTMLElement;
  #animationRemovalTimeout: number;

  @Prop() expanded = false;
  @Prop() testAttr?: string;

  async adjustOffset(expand: boolean) {
    expand && this.#outerEl.classList.remove('accordion--animated');
    expand && (await new Promise<void>((resolve) => window.requestAnimationFrame(() => resolve())));
    this.#innerEl.style.marginTop = `${this.#innerEl.offsetHeight * -1}px`;
    await new Promise<void>((resolve) => window.requestAnimationFrame(() => resolve()));
    expand && this.#outerEl.classList.add('accordion--animated');
    expand && (await new Promise<void>((resolve) => window.requestAnimationFrame(() => resolve())));
  }

  removeAnimation() {
    const rawDuration = getComputedStyle(this.#innerEl).getPropertyValue('--duration-medium');
    const numDuration = parseFloat(rawDuration);
    const parsedDuration = rawDuration.includes('ms') ? numDuration : numDuration * 1000;
    this.#animationRemovalTimeout = window.setTimeout(
      () => this.#outerEl.classList.remove('accordion--animated'),
      parsedDuration
    );
  }

  async update() {
    window.clearTimeout(this.#animationRemovalTimeout);
    await this.adjustOffset(this.expanded);
    this.#outerEl.classList[this.expanded ? 'add' : 'remove']('accordion--expanded');
    !this.expanded && this.removeAnimation();
  }

  componentDidRender() {
    this.update();
  }

  render() {
    const { testAttr } = this;

    return (
      <Host
        class="accordion"
        ref={(el) => (this.#outerEl = el)}
        aria-hidden={String(!this.expanded)}
        data-test={testAttr}
        data-test-active={testAttr && String(this.expanded)}
      >
        <div class="accordion__inner" ref={(el) => (this.#innerEl = el)}>
          <slot />
        </div>
      </Host>
    );
  }
}
