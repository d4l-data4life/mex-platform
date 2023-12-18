const isPreviewModeFromQueryParam = new URLSearchParams(document.location.search).get('preview') === 'true';
let isPreviewModeFromSessionStorage = false;

try {
  isPreviewModeFromSessionStorage = window.sessionStorage?.getItem?.('preview') === 'true';
  if (isPreviewModeFromQueryParam && !isPreviewModeFromSessionStorage && window.sessionStorage?.setItem) {
    sessionStorage.setItem('preview', 'true');
  }
} catch {}

export const IS_PREVIEW_MODE = isPreviewModeFromQueryParam || isPreviewModeFromSessionStorage;
