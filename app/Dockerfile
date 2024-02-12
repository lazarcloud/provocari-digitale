FROM node:alpine AS build

COPY ./app/package.json .
RUN npm install

COPY ./app/ ./
RUN npm run build

FROM busybox:latest AS runtime
COPY --from=build /build .
CMD ["busybox", "httpd", "-f", "-v", "-p", "80"]