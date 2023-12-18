# syntax = docker/dockerfile:1-experimental

FROM node:18-alpine AS pairgen

WORKDIR /app/tools/pairgen
COPY . .

RUN npm install && npm run build

CMD [ "npm", "start" ]
