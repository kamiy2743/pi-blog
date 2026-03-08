FROM node:25-bookworm-slim AS deps

WORKDIR /app

COPY ./frontend/package.json ./
COPY ./frontend/package-lock.json ./
RUN npm ci

FROM node:25-bookworm-slim AS builder

WORKDIR /app

COPY ./frontend/ ./
COPY --from=deps /app/node_modules/ ./node_modules/
RUN npm run build:client

FROM nginx:1.27-alpine

COPY ./nginx/blog.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist/client/ /usr/share/nginx/html/dist/client/

USER nginx:nginx
