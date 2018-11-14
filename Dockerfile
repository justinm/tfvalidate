FROM golang:1.11 AS build

ARG REPOSITORY="github.com/justinm/tfvalidate"
ENV GOPATH=/build
WORKDIR /build/src/$REPOSITORY/

COPY vendor/ /build/src/$REPOSITORY/vendor/
COPY tfvalidate.go Gopkg.lock Gopkg.toml /build/src/$REPOSITORY/
COPY tfvalidate/ /build/src/$REPOSITORY/tfvalidate/

RUN go build -o /build/bin/tfvalidate


FROM golang:1.11

RUN mkdir /workspace

COPY --from=build /build/bin/tfvalidate /usr/local/bin/tfvalidate

VOLUME ["/workspace"]
WORKDIR /workspace

ENTRYPOINT ["/usr/local/bin/tfvalidate"]

CMD ["-h"]
