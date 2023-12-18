export const IS_MOBILE = /[^\w](android|iphone|ipad|ipod)[^\w].*[^\w]mobile[^\w]/i.test(navigator.userAgent);
export const IS_TOUCH = 'ontouchstart' in window || navigator.maxTouchPoints > 0;
export const IS_DESKTOP = !IS_MOBILE;
export const IS_POINTER = !IS_TOUCH;

type MobileViewportChangesListener = (isMobileViewport: boolean) => void;

class MobileViewportChanges {
  #isMobileViewport: boolean;
  #listeners: MobileViewportChangesListener[] = [];

  constructor() {
    const tabletMinWidth = getComputedStyle(document.documentElement).getPropertyValue('--viewport-tablet');
    const mql = window.matchMedia(`(min-width: ${tabletMinWidth})`);
    mql.addEventListener?.('change', this.#onChange.bind(this));
    this.#onChange(mql);
  }

  #onChange({ matches }) {
    this.#isMobileViewport = !matches;
    this.#listeners.forEach((listener) => listener(this.#isMobileViewport));
  }

  addListener(listener: MobileViewportChangesListener) {
    this.#listeners.push(listener);
    listener(this.#isMobileViewport);
  }

  removeListener(listener: MobileViewportChangesListener) {
    this.#listeners = this.#listeners.filter((item) => item !== listener);
  }
}

export const mobileViewportChanges = new MobileViewportChanges();

export default {
  IS_DESKTOP,
  IS_MOBILE,
  IS_POINTER,
  IS_TOUCH,
  mobileViewportChanges,
};
