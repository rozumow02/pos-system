FROM node:22-alpine AS builder

WORKDIR /app

COPY frontend/package.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

FROM node:22-alpine

WORKDIR /app

ENV HOST=0.0.0.0
ENV PORT=3000
ENV NUXT_API_SERVER_BASE=http://backend:8080/api

COPY --from=builder /app/.output /app/.output

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
