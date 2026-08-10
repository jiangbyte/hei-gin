#FROM node:22-alpine AS build
FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/node:22-alpine AS build

ENV NODE_OPTIONS="--dns-result-order=ipv4first"

WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .
# Prefer .env.production (empty VITE_API_URL) so browser uses same-origin /api via nginx.
RUN pnpm exec vite build --mode production

FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/nginx:1.27-alpine

COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx/default.conf.template /etc/nginx/templates/default.conf.template

ENV BACKEND_URL="http://127.0.0.1:8000" \
    CLIENT_MAX_BODY_SIZE="10m" \
    CONTENT_SECURITY_POLICY="default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'" \
    HSTS_HEADER=""

EXPOSE 80
