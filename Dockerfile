FROM golang:1.11 AS build

ARG REPOSITORY="github.com/justinm/tfvalidate"
ENV GOPATH=/build
WORKDIR /build/src/$REPOSITORY/tfvalidate

RUN mkdir /build/bin \
    && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . /build/src/$REPOSITORY/tfvalidate

RUN /build/bin/dep ensure
RUN go build -o /build/bin/tfvalidate


FROM golang:1.11

RUN mkdir /workspace

COPY --from=build /build/bin/tfvalidate /usr/local/bin/tfvalidate

VOLUME ["/workspace"]
WORKDIR /workspace

ENTRYPOINT ["/usr/local/bin/tfvalidate"]

CMD ["-h"]
