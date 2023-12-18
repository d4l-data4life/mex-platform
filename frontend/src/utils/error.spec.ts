jest.mock('stencil-router-v2');

import { catchRetryableAction } from './error';

describe('error util', () => {
  describe('catchRetryableAction()', () => {
    let dispatchEventSpy;

    beforeEach(() => {
      jest.clearAllMocks();
      jest.spyOn(console, 'error').mockImplementation(() => {});
      dispatchEventSpy = jest.spyOn(window, 'dispatchEvent');
    });

    it('runs the given (async) action', async () => {
      const spy = jest.fn();
      const functionA = jest.fn(() => {});
      const functionB = async () => {
        await new Promise((resolve) => setTimeout(resolve, 1));
        spy();
      };

      await catchRetryableAction(functionA);
      expect(functionA).toHaveBeenCalled();

      await catchRetryableAction(functionB);
      expect(spy).toHaveBeenCalled();
    });

    it('dispatches custom event when action throws error', async () => {
      const error = new Error('some error');

      await expect(
        catchRetryableAction(() => {
          throw error;
        })
      ).resolves.not.toThrow();

      expect(dispatchEventSpy).toHaveBeenCalledWith(
        expect.objectContaining({
          detail: expect.objectContaining({
            error,
            retryHandler: expect.objectContaining({}),
          }),
        })
      );
    });
  });
});
