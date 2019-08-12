FROM golang:1.13-rc

EXPOSE 8080

WORKDIR /tmp/github.com/tarkov-database/rest-api
COPY . .

RUN make bin && \
    mkdir -p /usr/share/tarkov-database/rest-api && \
    mv -t /usr/share/tarkov-database/rest-api apiserver && \
    rm -rf /tmp/github.com/tarkov-database/rest-api

WORKDIR /usr/share/tarkov-database/rest-api

CMD ["/usr/share/tarkov-database/rest-api/apiserver"]
