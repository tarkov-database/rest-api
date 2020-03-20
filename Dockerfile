FROM golang:1.14.1

LABEL homepage="https://tarkov-database.com"
LABEL repository="https://github.com/tarkov-database/rest-api"
LABEL maintainer="Markus Wiegand <mail@morphy2k.dev>"

EXPOSE 8080

WORKDIR /tmp/github.com/tarkov-database/rest-api
COPY . .

RUN make bin && \
    mkdir -p /usr/share/tarkov-database/rest-api && \
    mv -t /usr/share/tarkov-database/rest-api apiserver && \
    rm -rf /tmp/github.com/tarkov-database/rest-api

WORKDIR /usr/share/tarkov-database/rest-api

CMD ["/usr/share/tarkov-database/rest-api/apiserver"]
