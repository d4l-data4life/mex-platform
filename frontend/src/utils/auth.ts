import { v4 as uuid } from 'uuid';
import { AUTH_CHALLENGE_VERIFIER_LENGTH, AUTH_CHALLENGE_VERIFIER_MASK } from 'config';

export const generateChallenge = async (): Promise<[verifier: string, challenge: string]> => {
  if (!('crypto' in window)) {
    return [null, null];
  }

  try {
    const indices = crypto.getRandomValues(new Uint8Array(AUTH_CHALLENGE_VERIFIER_LENGTH));
    const scalingFactor = 256 / Math.min(AUTH_CHALLENGE_VERIFIER_MASK.length, 256);
    const verifier = new Array(AUTH_CHALLENGE_VERIFIER_LENGTH)
      .fill(null)
      .map((_, index) => AUTH_CHALLENGE_VERIFIER_MASK[Math.floor(indices[index] / scalingFactor)])
      .join('');

    const challengeDigest = await crypto.subtle.digest(
      'SHA-256',
      new Uint8Array(verifier.length).map((_, index) => verifier.charCodeAt(index))
    );

    const challenge = window
      .btoa(String.fromCharCode.apply(0, Array.from(new Uint8Array(challengeDigest))))
      .replace(/\+/g, '-')
      .replace(/\//g, '_')
      .replace(/=+$/, '');

    return [verifier, challenge];
  } catch (_) {
    return [null, null];
  }
};

export const generateState = (): string => {
  return 'crypto' in window ? (crypto as any).randomUUID?.() ?? uuid() : uuid();
};

export const parseJwtPayload = (jwt: string): any => {
  try {
    return (
      jwt &&
      JSON.parse(
        decodeURIComponent(
          window
            .atob(jwt.split('.')[1].replace(/-/g, '+').replace(/_/g, '/'))
            .split('')
            .map((char) => '%' + `00${char.charCodeAt(0).toString(16)}`.slice(-2))
            .join('')
        )
      )
    );
  } catch (_) {
    return null;
  }
};
