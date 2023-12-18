import { ROUTES } from 'config';
import { generateState } from './auth';

export const generateNonce = (): string => generateState();
export const generateVisitorId = (): string => generateState().replace(/-/g, '');

export const redactUrl = (url: string): string => {
  const pattern = new RegExp(ROUTES.SEARCH_QUERY.replace(':query', '([^/?]+)'));
  return url.replace(pattern, (match, p1) => match.replace(p1, 'redacted'));
};
