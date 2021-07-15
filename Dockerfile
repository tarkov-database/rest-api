FROM golang:1.16.6 as build-env

WORKDIR /tmp/github.com/tarkov-database/rest-api
COPY . .

RUN make bin && \
    mkdir -p /usr/share/tarkov-database/rest-api && \
    mv -t /usr/share/tarkov-database/rest-api apiserver

FROM gcr.io/distroless/base-debian10

LABEL homepage="https://tarkov-database.com"
LABEL repository="https://github.com/tarkov-database/rest-api"
LABEL maintainer="Markus Wiegand <mail@morphy2k.dev>"

LABEL org.opencontainers.image.source="https://github.com/tarkov-database/rest-api"

COPY --from=build-env /usr/share/tarkov-database/rest-api /

EXPOSE 8080

CMD ["/apiserver"]
