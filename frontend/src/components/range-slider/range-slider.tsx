import { Component, h, Host, Prop, State, Watch } from '@stencil/core';
import stores from 'stores';
import { buildAscNumSequence } from 'utils/search';

@Component({
  tag: 'mex-range-slider',
  styleUrl: 'range-slider.css',
})
export class RangeSliderComponent {
  #sliderBarEl: HTMLDivElement;
  #stepWidth: number;
  #xStart: number;

  @Prop() min: number;
  @Prop() max: number;
  @Prop() value: number[];
  @Prop() handleChange: (value: number[]) => void;
  @Prop() handleDrag?: (value: number[]) => void;
  @Prop() disabled = false;
  @Prop() highlightActiveRange = false;
  @Prop() mode: 'point' | 'frame' = 'point';
  @Prop() testAttr?: string;

  @State() inputValue: number[];
  @State() focusedIndex: number;
  @State() isDragging = false;

  @Watch('value')
  onValueChange() {
    this.updateValue();
  }

  constructor() {
    this.startDrag = this.startDrag.bind(this);
    this.drag = this.drag.bind(this);
    this.endDrag = this.endDrag.bind(this);
    this.point = this.point.bind(this);
    this.onKeyDown = this.onKeyDown.bind(this);
  }

  get stepsCount() {
    return this.max - this.min + 1;
  }

  get knobsCount() {
    return this.inputValue?.length ?? 1;
  }

  get activeRangeStyle() {
    const { knobsCount } = this;
    const [start, end] = [this.getValuePerc(0), this.getValuePerc(knobsCount - 1)];
    return { left: `${start}%`, right: `${(100 - parseFloat(end)).toFixed(4)}%` };
  }

  getValuePerc(index: number) {
    const { min, stepsCount, mode } = this;

    if (mode === 'point') {
      return (((this.inputValue[index] - min) / (stepsCount - 1)) * 100).toFixed(4);
    }

    return (((this.inputValue[index] + index - min) / stepsCount) * 100).toFixed(4);
  }

  getLabel(index: number) {
    const { t } = stores.i18n;
    const { knobsCount, mode } = this;

    if (mode === 'frame') {
      return t(`range.${index ? 'end' : 'start'}`);
    }

    if (knobsCount === 1) {
      return t('range.one');
    }

    return t('range.many', { index: index + 1, count: knobsCount });
  }

  isKnobDragging(index: number) {
    return this.focusedIndex === index && this.isDragging;
  }

  getX(event: MouseEvent | TouchEvent): number {
    const coords = 'changedTouches' in event ? event.changedTouches[0] : event;
    const { clientX } = coords ?? {};
    return clientX;
  }

  startDrag(event: MouseEvent | TouchEvent) {
    const { mode } = this;
    const { left, right } = this.#sliderBarEl.getBoundingClientRect();
    const sliderBarKnobEl = event.target as HTMLElement;

    if (sliderBarKnobEl.ariaDisabled === 'true') {
      return;
    }

    const index = parseInt(sliderBarKnobEl.dataset?.index, 10);

    if (mode === 'point') {
      this.#stepWidth = (right - left) / (this.stepsCount - 1);
      this.#xStart = left;
    } else {
      this.#stepWidth = (right - left) / this.stepsCount;
      this.#xStart = left + index * this.#stepWidth;
    }

    this.focusedIndex = index;
    this.isDragging = true;
    this.attachEvents();
  }

  drag(event: MouseEvent | TouchEvent) {
    const { min, stepsCount } = this;
    const x = this.getX(event);
    const value = min + Math.min(Math.max(Math.round((x - this.#xStart) / this.#stepWidth), 0), stepsCount - 1);
    this.setValue(value);
  }

  endDrag() {
    this.removeEvents();
    this.submitValue();
    this.isDragging = false;
  }

  preventPageSwipe(event: MouseEvent | TouchEvent) {
    event.preventDefault();
    event.stopPropagation();
  }

  point(event: MouseEvent | TouchEvent) {
    this.preventPageSwipe(event);

    if (this.isDragging || !this.#sliderBarEl) {
      return;
    }

    const x = this.getX(event);
    const knobs = [...(this.#sliderBarEl.querySelectorAll('button[data-index]') as unknown as HTMLElement[])];
    const closestKnob = knobs.reduce(
      (result, knob) => {
        const { left, width } = knob.getBoundingClientRect();
        const distance = Math.abs(x - (left + width / 2));
        return distance < result.distance ? { distance, knob } : result;
      },
      { distance: Infinity, knob: null }
    ).knob;

    if (!closestKnob) {
      return;
    }

    Object.defineProperty(event, 'target', {
      value: closestKnob,
      writable: false,
    });

    this.startDrag(event);
    this.drag(event);
  }

  attachEvents() {
    document.addEventListener('mouseup', this.endDrag, { passive: true });
    document.addEventListener('touchend', this.endDrag, { passive: true });
    document.addEventListener('mousemove', this.drag, { passive: true });
    document.addEventListener('touchmove', this.drag, { passive: true });
    this.#sliderBarEl?.addEventListener('mousedown', this.preventPageSwipe, { passive: false });
    this.#sliderBarEl?.addEventListener('touchstart', this.preventPageSwipe, { passive: false });
  }

  removeEvents() {
    document.removeEventListener('mouseup', this.endDrag);
    document.removeEventListener('touchend', this.endDrag);
    document.removeEventListener('mousemove', this.drag);
    document.removeEventListener('touchmove', this.drag);
    this.#sliderBarEl?.removeEventListener('mousedown', this.preventPageSwipe);
    this.#sliderBarEl?.removeEventListener('touchstart', this.preventPageSwipe);
  }

  setValue(value, index = this.focusedIndex) {
    const { inputValue, min, max } = this;
    if (index === null || inputValue?.[index] === value) {
      return;
    }

    this.inputValue = buildAscNumSequence(
      this.value.map((v, i) => (index === i ? value : v)),
      index,
      min,
      max
    );
    this.handleDrag?.(this.inputValue);
  }

  updateValue() {
    this.setValue(null, -1);
  }

  submitValue() {
    if (this.value.some((v, i) => this.inputValue[i] !== v)) {
      this.handleChange?.(this.inputValue);
    }
  }

  onKeyDown(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      this.submitValue();
    }

    if (event.key === 'ArrowUp' || event.key === 'ArrowLeft') {
      event.preventDefault();
      this.setValue(this.inputValue[this.focusedIndex] - 1);
    }

    if (event.key === 'ArrowDown' || event.key === 'ArrowRight') {
      event.preventDefault();
      this.setValue(this.inputValue[this.focusedIndex] + 1);
    }
  }

  componentWillLoad() {
    this.updateValue();

    if (this.mode === 'frame' && this.knobsCount !== 2) {
      throw new Error('range slider error: frame mode is only supported with two knobs!');
    }
  }

  disconnectedCallback() {
    this.removeEvents();
  }

  render() {
    const { disabled, testAttr } = this;

    return (
      <Host
        class={{ 'range-slider': true, 'range-slider--disabled': disabled }}
        onMouseDown={this.point}
        onTouchStart={this.point}
        data-test={testAttr}
      >
        <div class="range-slider__bar" ref={(el) => (this.#sliderBarEl = el)}>
          {this.highlightActiveRange && <div class="range-slider__activeRange" style={this.activeRangeStyle} />}

          {new Array(this.knobsCount).fill(null).map((_, index) => (
            <div
              class={{
                'range-slider__knob': true,
                'range-slider__knob--dragging': this.isKnobDragging(index),
              }}
              key={index}
              style={{ left: `${this.getValuePerc(index)}%` }}
              onMouseDown={this.startDrag}
              onTouchStart={this.startDrag}
              onFocus={() => (this.focusedIndex = index)}
              onBlur={() => {
                this.submitValue();
                this.focusedIndex = null;
              }}
              onKeyDown={this.onKeyDown}
              data-index={index}
              title={this.getLabel(index)}
              tabIndex={disabled ? -1 : 0}
              draggable={false}
              role="slider"
              aria-valuemin={this.min}
              aria-valuemax={this.max}
              aria-valuenow={this.inputValue[index]}
              aria-disabled={disabled ? 'true' : 'false'}
              data-test={testAttr && `${testAttr}:knob`}
              data-test-key={testAttr && index}
            />
          ))}
        </div>

        {this.isDragging && <div class="range-slider__backdrop" />}
      </Host>
    );
  }
}
