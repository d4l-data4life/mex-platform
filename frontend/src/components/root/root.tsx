import { Component, Event, EventEmitter, Fragment, FunctionalComponent, h, Host, Listen, State } from '@stencil/core';
import { match, Route } from 'stencil-router-v2';
import { ANALYTICS_IS_ENABLED, ROUTES } from 'config';
import stores from 'stores';
import { IS_DESKTOP, IS_MOBILE, IS_POINTER, IS_TOUCH } from 'utils/device';
import services from 'services';
import { catchRetryableAction } from 'utils/error';

export interface ModalData {
  Contents: FunctionalComponent;
  handleClose?: () => void;
  handleSetFocus?: (closeRef?: HTMLButtonElement) => void;
}

@Component({
  tag: 'mex-root',
  styleUrl: 'root.css',
})
export class RootComponent {
  @State() url: URL = stores.router.url;
  @State() hasTranslations = false;
  @State() hasConfig = false;
  @State() hasStickyFooter = false;
  @State() hasDisabledScrolling = false;
  @State() modalData?: ModalData;

  @Event() closeModal: EventEmitter;
  @Event() closeError: EventEmitter;

  @Listen('staticPageTitleChanged')
  staticPageTitleChangedListener(event: CustomEvent) {
    const title = stores.i18n.t('navigation.title');
    const staticPageTitle = (event.detail as string) || '';
    this.updatePageTitle(`${title} - ${staticPageTitle}`);
  }

  @Listen('stickyFooterEnabled')
  stickyFooterEnabledListener(event: CustomEvent) {
    this.hasStickyFooter = (event.detail as boolean) || false;
  }

  @Listen('scrollingDisabled')
  scrollingDisabledListener(event: CustomEvent) {
    this.hasDisabledScrolling = (event.detail as boolean) || false;
  }

  @Listen('showModal')
  showModal(event: CustomEvent & { detail: ModalData }) {
    const { Contents, handleClose, handleSetFocus } = event.detail;

    this.closeError.emit();
    this.modalData = {
      Contents,
      handleSetFocus: (closeRef?: HTMLButtonElement) => {
        handleSetFocus ? handleSetFocus(closeRef) : closeRef?.focus();
      },
      handleClose: () => {
        handleClose?.();
        this.modalData = null;
      },
    };
  }

  @Listen('closeModal')
  closeModalHandler() {
    this.modalData?.handleClose();
    this.closeError.emit();
  }

  constructor() {
    this.handleLanguageChange = this.handleLanguageChange.bind(this);
  }

  get isInitialized() {
    return this.hasConfig && this.hasTranslations;
  }

  get showTrackingBanner() {
    return ANALYTICS_IS_ENABLED() && stores.auth.isAuthenticated && !stores.analytics.hasChosen;
  }

  get pageTitle() {
    const { t } = stores.i18n;
    const { pathname } = this.url;
    const title = t('navigation.title');

    if (pathname === ROUTES.ROOT) {
      return `${title} – ${t('navigation.pages.home')}`;
    }

    if (match(ROUTES.SEARCH)(pathname)) {
      return `${title} – ${t('navigation.pages.search')}`;
    }

    if (match(ROUTES.ITEM)(pathname)) {
      return `${title} – ${t('navigation.pages.details')}`;
    }

    return title;
  }

  updatePageTitle(customTitle?: string) {
    if (!this.hasTranslations) {
      return;
    }

    document.title = customTitle || this.pageTitle;
  }

  async loadTranslations() {
    await catchRetryableAction(async () => {
      await services.config.loadTranslations();
      this.hasTranslations = true;
    }, false);

    this.updatePageTitle();
  }

  loadConfig() {
    catchRetryableAction(async () => {
      await services.config.load();
      this.hasConfig = true;
    }, false);

    this.loadNavigation();
  }

  loadNavigation() {
    catchRetryableAction(async () => {
      await services.navigation.load();
    }, false);
  }

  handleLanguageChange() {
    this.updatePageTitle();
  }

  componentWillLoad() {
    stores.router.onChange('url', (url: URL) => {
      this.url = url;
      window.scrollTo(0, 0);
      this.closeModal.emit();
      this.updatePageTitle();

      if (!match(ROUTES.CONTENT_PAGE)(url.pathname)) {
        this.hasStickyFooter = false;
        this.hasDisabledScrolling = false;
      }

      services.analytics.trackImpression();
    });

    stores.i18n.addListener(this.handleLanguageChange);

    this.loadTranslations();
    this.loadConfig();
  }

  componentDidLoad() {
    services.analytics.trackImpression();
  }

  disconnectedCallback() {
    stores.i18n.removeListener(this.handleLanguageChange);
  }

  render() {
    const Router = stores.router;
    const { showTrackingBanner, hasStickyFooter, hasDisabledScrolling, modalData, isInitialized } = this;
    const { isAuthenticated } = stores.auth;

    if (hasStickyFooter) {
      document.documentElement.scrollTop = 0;
    }

    return (
      <Host
        class={{
          root: true,
          'feature--trackingBanner': showTrackingBanner,
          'feature--sticky-footer': hasStickyFooter,
          'feature--disabled-scrolling': hasDisabledScrolling,
          'device--mobile': IS_MOBILE,
          'device--desktop': IS_DESKTOP,
          'device--pointer': IS_POINTER,
          'device--touch': IS_TOUCH,
        }}
      >
        <mex-header class="root__header" />

        {modalData && (
          <mex-modal class="root__modal" handleClose={modalData.handleClose} handleSetFocus={modalData.handleSetFocus}>
            <modalData.Contents />
          </mex-modal>
        )}
        <mex-error class="root__modal" />

        <main class="root__content">
          <Router.Switch>
            {isAuthenticated && isInitialized && (
              <Fragment>
                <Route path={ROUTES.ROOT} render={() => <mex-view-home />} />
                <Route path={ROUTES.AUTH} to={ROUTES.ROOT} />
                <Route path={ROUTES.SEARCH} render={() => <mex-view-search url={this.url} />} />
                <Route path={match(ROUTES.SEARCH_QUERY)} render={() => <mex-view-search url={this.url} />} />
                <Route path={match(ROUTES.ITEM)} render={({ id }) => <mex-view-item itemId={id} />} />
                <Route
                  path={match(ROUTES.CONTENT_PAGE)}
                  render={({ pageid }) => <mex-view-content-page pageId={pageid} />}
                />
              </Fragment>
            )}

            {!isAuthenticated && (
              <Route
                path={(path) => ![ROUTES.AUTH, ROUTES.LOGOUT].includes(path)}
                to={(path) => {
                  stores.auth.requestedRoute = path + this.url.search;
                  return ROUTES.AUTH;
                }}
              />
            )}

            <Route path={ROUTES.AUTH} render={() => <mex-view-auth url={this.url} />} />
            <Route path={ROUTES.LOGOUT} render={() => <mex-view-logout />} />
            {isInitialized && <Route path={() => true} render={() => <mex-view-not-found />} />}
          </Router.Switch>

          {isAuthenticated && !isInitialized && (
            <div class="view view--center">
              <mex-logo class="root__loader" loader />
            </div>
          )}
        </main>

        {isInitialized && <mex-footer class="root__footer" isThin={this.hasStickyFooter} />}

        <mex-notifications />
        {showTrackingBanner && isInitialized && <mex-tracking-banner />}
      </Host>
    );
  }
}
