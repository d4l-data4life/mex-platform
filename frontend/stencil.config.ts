import { Config } from '@stencil/core';
import { postcss } from '@stencil/postcss';
import postcssCustomMedia from 'postcss-custom-media';
import visualizer from 'rollup-plugin-visualizer';
import dotenv from 'dotenv';

dotenv.config();

const env = process.env;

// https://stenciljs.com/docs/config

const REQUIRED_ENV_VARS = ['API_URL', 'CONFIG_STATIC_URL'];
REQUIRED_ENV_VARS.forEach((key) => {
  if (typeof env[key] === 'undefined') {
    throw new Error(`${key} env var missing`);
  }
});

export const config: Config = {
  globalStyle: 'src/global/app.css',
  globalScript: 'src/global/app.ts',
  taskQueue: 'async',
  buildEs5: 'prod' as boolean | 'prod',
  hashedFileNameLength: 8,
  hashFileNames: true,
  outputTargets: [
    {
      type: 'www',
      serviceWorker: null,
      baseUrl: '/',
    },
  ],
  plugins: [
    ...(env.npm_lifecycle_event === 'analyze' ? [visualizer()] : []),
    postcss({
      plugins: [postcssCustomMedia({ preserve: false })],
    }),
  ],
  testing: {
    browserArgs: ['--no-sandbox', '--disable-setuid-sandbox'],
    moduleDirectories: ['node_modules', 'src'],
  },
  devServer: {
    port: parseInt(env.PORT ?? '4000', 10),
  },
  env: {
    PORT: env.PORT,
    API_URL: env.API_URL,
    AUTH_PROVIDER: env.AUTH_PROVIDER,
    AUTH_CLIENT_ID: env.AUTH_CLIENT_ID,
    AUTH_AUTHORIZE_URI: env.AUTH_AUTHORIZE_URI,
    AUTH_TOKEN_URI: env.AUTH_TOKEN_URI,
    AUTH_LOGOUT_URI: env.AUTH_LOGOUT_URI,
    CONFIG_STATIC_URL: env.CONFIG_STATIC_URL,
    CONFIG_CMS_URL: env.CONFIG_CMS_URL,
    DYNAMIC_ENV_URL: env.DYNAMIC_ENV_URL,
    APP_VERSION: env.APP_VERSION ?? 'unknown build',
  },
};
