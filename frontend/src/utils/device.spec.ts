const ORIGINAL_USERAGENT = global.navigator.userAgent;
const EXAMPLE_IPHONE =
  'Mozilla/5.0 (iPhone; CPU iPhone OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148';
const EXAMPLE_ANDROID =
  'Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.87 Mobile Safari/537.36';
const EXAMPLE_WINDOWS = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:96.0) Gecko/20100101 Firefox/96.0';
const EXAMPLE_MACOS =
  'Mozilla/5.0 (Macintosh; Intel Mac OS X 12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36';
const EXAMPLE_LINUX =
  'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36';

const getDevice = () => {
  jest.resetModules();
  return require('./device');
};

const setUserAgent = (userAgent?: string) => {
  Object.defineProperty(global.navigator, 'userAgent', {
    get() {
      return userAgent || ORIGINAL_USERAGENT;
    },
  });
};

const setTouch = (isTouch: boolean) => {
  Object.defineProperty(global.navigator, 'maxTouchPoints', {
    get() {
      return isTouch ? 6 : 0;
    },
  });

  if (isTouch) {
    window.ontouchstart = jest.fn();
  } else {
    delete window.ontouchstart;
  }
};

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
  })),
});

describe('device util', () => {
  it('detects mobile devices', () => {
    let device;

    setUserAgent(EXAMPLE_ANDROID);
    device = getDevice();
    expect(device.IS_MOBILE).toBe(true);
    expect(device.IS_DESKTOP).toBe(false);

    setUserAgent(EXAMPLE_IPHONE);
    device = getDevice();
    expect(device.IS_MOBILE).toBe(true);
    expect(device.IS_DESKTOP).toBe(false);
  });

  it('detects desktop devices', () => {
    let device;

    setUserAgent(EXAMPLE_WINDOWS);
    device = getDevice();
    expect(device.IS_DESKTOP).toBe(true);
    expect(device.IS_MOBILE).toBe(false);

    setUserAgent(EXAMPLE_MACOS);
    device = getDevice();
    expect(device.IS_DESKTOP).toBe(true);
    expect(device.IS_MOBILE).toBe(false);

    setUserAgent(EXAMPLE_LINUX);
    device = getDevice();
    expect(device.IS_DESKTOP).toBe(true);
    expect(device.IS_MOBILE).toBe(false);
  });

  it('detects unknown user agents as desktop', () => {
    setUserAgent();
    const device = getDevice();
    expect(device.IS_DESKTOP).toBe(true);
    expect(device.IS_MOBILE).toBe(false);
  });

  it('detects touch devices', () => {
    setTouch(true);
    const device = getDevice();
    expect(device.IS_TOUCH).toBe(true);
    expect(device.IS_POINTER).toBe(false);
    expect(window.ontouchstart).not.toHaveBeenCalled();
  });

  it('detects pointer devices', () => {
    setTouch(false);
    const device = getDevice();
    expect(device.IS_POINTER).toBe(true);
    expect(device.IS_TOUCH).toBe(false);
  });

  describe('MobileViewportChanges', () => {
    it('adds a listener for changes between mobile and non-mobile viewport', () => {
      const listenerMock = jest.fn();
      const device = getDevice();
      device.mobileViewportChanges.addListener(listenerMock);
      expect(listenerMock).toHaveBeenCalledTimes(1);
      expect(window.matchMedia).toHaveBeenCalledWith(expect.stringContaining('(min-width:'));
    });
  });
});
