FROM golang:1.13-rc

EXPOSE 8080

WORKDIR /tmp/github.com/tarkov-database/api
COPY . .

RUN make bin && \
    mkdir -p /usr/share/tarkov-database/api && \
    mv -t /usr/share/tarkov-database/api apiserver view && \
    rm -rf /tmp/github.com/tarkov-database/api

WORKDIR /usr/share/tarkov-database/api

CMD ["/usr/share/tarkov-database/api/apiserver"]
