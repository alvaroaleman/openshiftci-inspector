FROM golang:1.16 AS api
COPY . /usr/src/app
WORKDIR /usr/src/app
RUN go run buildtool.go export-api

FROM node:14 AS frontend
COPY --from=api /usr/src/app /usr/src/app
WORKDIR /usr/src/app/frontend
RUN npm install
RUN npm run client
RUN npm run build

FROM golang:1.16 AS backend
COPY --from=frontend /usr/src/app /usr/src/app
WORKDIR /usr/src/app
RUN go build -o app cmd/

FROM alpine
COPY --from=backend /usr/src/app/app /srv/app
ENTRYPOINT ["/srv/app"]
CMD []
USER 1000:1000
EXPOSE 8080
