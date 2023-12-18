import { Component, Fragment, h, Prop } from '@stencil/core';
import stores from 'stores';

export interface ModalContentsButton {
  label: string;
  clickHandler: () => void;
  modifier?: 'primary' | 'secondary' | 'tertiary';
  disabled?: boolean;
  testAttr?: string;
}

@Component({
  tag: 'mex-modal-contents',
  styleUrl: 'modal-contents.css',
})
export class ModalContentsComponent {
  @Prop() context?: any;
  @Prop() illustration?: 'error' | 'download' | 'confirm' | 'contact';
  @Prop() caption?: string;
  @Prop() text?: string;
  @Prop() buttons?: ModalContentsButton[];
  @Prop() progress?: number;

  render() {
    const { context, illustration, caption, text, buttons, progress } = this;
    const { t } = stores.i18n;

    return (
      <Fragment>
        {illustration === 'error' && <mex-illustration-error classes="modal-contents__icon" />}
        {illustration === 'download' && <mex-illustration-download classes="modal-contents__icon" />}
        {illustration === 'confirm' && <mex-illustration-confirm classes="modal-contents__icon" />}
        {illustration === 'contact' && <mex-illustration-contact classes="modal-contents__icon" />}

        {caption && <h3 class="modal-contents__title">{t(caption)}</h3>}
        {text && <p class="modal-contents__text">{t(text)}</p>}

        <slot />
        {Number.isInteger(progress) && <mex-progress-bar spaced progress={progress} />}

        {!!buttons?.length && (
          <div class="modal-contents__buttons">
            {buttons.map((button, index) => (
              <button
                key={index}
                ref={(el) => context && context.modalContentButtonRefs && (context.modalContentButtonRefs[index] = el)}
                class={`button ${button.modifier ? `button--${button.modifier}` : ''}`}
                onClick={button.clickHandler}
                data-test={button.testAttr}
                disabled={!!button.disabled}
              >
                {t(button.label)}
              </button>
            ))}
          </div>
        )}
      </Fragment>
    );
  }
}
