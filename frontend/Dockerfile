### BUILD ###

FROM node:18-alpine AS base

WORKDIR /build

COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .
COPY stencil.config.ts .
COPY src/ ./src/
COPY scripts/build-static-html.js ./scripts/build-static-html.js

COPY .env.example .env

RUN npm config set update-notifier false
RUN npm ci

FROM base AS unit-test

RUN npm run test

FROM base AS build

ARG APP_VERSION
RUN ENVIRONMENT=production APP_VERSION="$APP_VERSION" npm run build
RUN npm run build:static

### SERVE ###

FROM nginx:latest AS serve
COPY --from=build /build/www/ /usr/share/nginx/html/
COPY nginx/ /etc/nginx/
