import { Component, getAssetPath, h, Host, Prop, State } from '@stencil/core';
import { ANALYTICS_CUSTOM_EVENTS, SEARCH_CONFIG } from 'config';
import services from 'services';
import stores from 'stores';
import { escapeHtml } from 'utils/field';
import { addSearchOperator } from 'utils/search';

@Component({
  tag: 'mex-search',
  styleUrl: 'search.css',
  assetsDirs: ['assets'],
})
export class SearchComponent {
  #inputEl?: HTMLInputElement;
  #searchFocusOptions = SEARCH_CONFIG.FOCI.map((value) => ({
    value,
    label: `search.focus.${value ?? 'all'}`,
  }));

  @Prop() expanded = false;
  @Prop() autofocus = false;
  @Prop() handleSearch?: (value: string, searchFocus: string) => void;
  @Prop() handleReset?: () => void;
  @Prop() value = '';
  @Prop() searchFocus = null;

  @State() inputValue = '';
  @State() updatedValue = '';
  @State() selectedSearchFocus = null;
  @State() isFocused = false;
  @State() areTipsExpanded = false;

  get pristine() {
    return !this.inputValue;
  }

  onInput(value: string) {
    this.inputValue = value;
  }

  updateInput() {
    this.onInput(this.value);
    this.updatedValue = this.value;
  }

  onKeyDown(key: KeyboardEvent['key']) {
    key === 'Escape' && !this.inputValue && this.#inputEl?.blur();
    key === 'Escape' && this.inputValue && this.onReset();
  }

  onReset() {
    this.onInput('');
    this.#inputEl?.focus();
    this.handleReset?.();
  }

  onSubmit(event: Event) {
    event.preventDefault();
    this.handleSearch?.(this.inputValue.trim(), this.selectedSearchFocus);

    services.analytics.trackEvent(
      ...(this.selectedSearchFocus ? ANALYTICS_CUSTOM_EVENTS.SEARCH_FOCUS : ANALYTICS_CUSTOM_EVENTS.SEARCH_PLAIN),
      'Initiated'
    );
  }

  onTipsClick(event: Event) {
    const target = event.target as HTMLElement;
    if (this.#inputEl && target.tagName.toLowerCase() === 'code') {
      addSearchOperator(this.#inputEl, target.innerText.trim().slice(0, 2));
    }
  }

  renderTipsHtml() {
    const html = stores.i18n.t('search.tips.text');
    return html
      .split(/\n+--\n+/)
      .map((tip) =>
        tip
          .split(/\n+/)
          .map((line = '', index) => {
            line = escapeHtml(line)
              .trim()
              .split('`')
              .map((item, index) =>
                !index ? item : index % 2 ? `<code data-test="search:tips:item:code">${item}` : `</code>${item}`
              )
              .join('');

            return index ? (index > 1 ? `<em>${line}</em>` : line) : `<strong>${line}</strong>`;
          })
          .join('\n')
      )
      .map((tip) => <div class="search__tip" innerHTML={tip} data-test="search:tips:item" />);
  }

  componentWillLoad() {
    this.selectedSearchFocus = this.searchFocus;
    this.updateInput();
  }

  componentDidLoad() {
    this.autofocus && setTimeout(() => this.#inputEl?.focus(), 100);
  }

  componentWillUpdate() {
    this.updatedValue !== this.value && this.updateInput();
  }

  render() {
    const { expanded, areTipsExpanded } = this;
    const { t } = stores.i18n;

    return (
      <Host
        class={{ search: true, 'search--expanded': expanded }}
        style={{
          '--background-pattern-topLeft': `url(${getAssetPath('./assets/search-background-pattern-1.svg')})`,
          '--background-pattern-bottomRight': `url(${getAssetPath('./assets/search-background-pattern-2.svg')})`,
        }}
        data-test="search"
      >
        <div class="search__wrapper">
          {expanded && (
            <h1 class="search__title" data-test="search:title">
              <span class="u-animated-stroke-underline-1">{t('search.title')}</span>
            </h1>
          )}
          {expanded && <p class="search__info" data-test="search:info" innerHTML={t('search.info')} />}

          <form
            class={{ search__form: true, 'search__form--focused': this.isFocused }}
            onReset={() => this.onReset()}
            onSubmit={(event) => this.onSubmit(event)}
            data-test="search:form"
          >
            <mex-dropdown
              class="search__focus"
              toggleClass="search__select"
              options={this.#searchFocusOptions}
              value={this.selectedSearchFocus}
              handleChange={(value) => {
                this.selectedSearchFocus = value;
                services.analytics.trackEvent(...ANALYTICS_CUSTOM_EVENTS.SEARCH_FOCUS, `Selected: ${value ?? 'all'}`);
              }}
              title={stores.i18n.t('search.fields')}
              testAttr="search:focus"
            />
            <input
              ref={(el) => (this.#inputEl = el)}
              type="text"
              class="search__input"
              placeholder={t('search.placeholder')}
              title={t('search.placeholder')}
              aria-label={t('search.placeholder')}
              value={this.inputValue}
              onFocus={() => (this.isFocused = true)}
              onBlur={() => (this.isFocused = false)}
              onInput={(event) => this.onInput((event.target as HTMLInputElement).value)}
              onKeyDown={(event) => this.onKeyDown(event.key)}
              data-test="search:input"
            />
            <button
              class="search__reset"
              type="reset"
              tabIndex={this.pristine ? -1 : 0}
              aria-hidden={String(this.pristine)}
              title={stores.i18n.t('search.reset')}
              data-test="search:reset"
            >
              <mex-icon-close-inline classes="icon--large" />
            </button>
            <button
              class="search__button"
              type="submit"
              title={stores.i18n.t('search.execute')}
              data-test="search:execute"
            >
              <mex-icon-search classes="icon--medium" />
            </button>
          </form>

          <div class="search__actions">
            <button
              class="search__action"
              onClick={() => (this.areTipsExpanded = !areTipsExpanded)}
              aria-expanded={String(areTipsExpanded)}
              data-test="search:tips:toggle"
            >
              {t('search.tips.title')} &#160;
              <mex-icon-chevron
                classes={`icon--inline icon--mirrorable ${areTipsExpanded ? 'icon--mirrored-vertical' : ''}`}
              />
            </button>
          </div>

          <mex-accordion expanded={areTipsExpanded} testAttr="search:tips">
            <div class="search__tips">
              <div class="search__tips-inner" onClick={(event) => this.onTipsClick(event)}>
                {stores.i18n.t('search.instructions')}
                {this.renderTipsHtml()}
              </div>
            </div>
          </mex-accordion>
        </div>
      </Host>
    );
  }
}
