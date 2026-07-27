FROM node:20-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
# Vite bakes VITE_* vars into the built JS at build time — the browser, not
# this container, is what calls the API, so this must be a URL reachable
# from the user's machine (see .env).
ARG VITE_GRAPHQL_URL=http://localhost:8080/graphql
ENV VITE_GRAPHQL_URL=$VITE_GRAPHQL_URL
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
