FROM golang:1.8
WORKDIR /ben-and-jerry
ENV SRC_DIR=/go/src/github.com/zalora/benandjerry/
ADD . $SRC_DIR
RUN cd $SRC_DIR;cd handler;go test -v ;go build -o benandjerry  -ldflags "-X 'main.buildTimestamp=$(date '+%b %d %Y %T')' -X main.commitID=`git describe --match=rAnDom --always --abbrev --dirty`"; cp benandjerry /ben-and-jerry/
ENV BAJ_LOG_FORMAT json
ENV BAJ_LOG_LEVEL info
ENV BAJ_LISTEN_PORT 8080
ENV BAJ_POSTGRES_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
ENTRYPOINT ["./benandjerry"]