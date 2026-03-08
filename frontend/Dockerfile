FROM node:25-bookworm-slim AS deps

WORKDIR /app

COPY ./package.json ./
COPY ./package-lock.json ./
RUN npm ci

FROM node:25-bookworm-slim AS builder

WORKDIR /app

COPY ./ ./
COPY --from=deps /app/node_modules/ ./node_modules/
RUN npm run build:server

FROM node:25-bookworm-slim

WORKDIR /app

COPY --from=builder /app/dist/server/ ./dist/server/
COPY --from=builder /app/node_modules/ ./node_modules/

USER node:node

CMD ["node", "dist/server/entry-server.js"]
