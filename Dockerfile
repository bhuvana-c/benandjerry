FROM golang:1.8
WORKDIR /ben-and-jerry
ENV SRC_DIR=/go/src/github.com/zalora/benandjerry/
RUN mkdir /etc/BAJ/
ADD . $SRC_DIR
COPY ./db/migrations /etc/BAJ/migrations
RUN cd $SRC_DIR/handler;go test -v
RUN cd $SRC_DIR;go build -o benandjerry  -ldflags "-X 'main.buildTimestamp=$(date '+%b %d %Y %T')' -X main.commitID=`git describe --match=rAnDom --always --abbrev --dirty`"; cp benandjerry /ben-and-jerry/
ENV BAJ_LOG_FORMAT json
ENV BAJ_LOG_LEVEL info
ENV BAJ_LISTEN_PORT 8080
ENV BAJ_POSTGRES_URL=postgres://postgres:postgres@localhost:5432/ben-and-jerry?sslmode=disable
ENV BAJ_MIGRATION_SCRIPTS_PATH=/etc/BAJ/migrations
EXPOSE 8080
ENTRYPOINT ["./benandjerry"]