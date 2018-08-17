FROM golang:1.10.3-alpine3.7 as builder

# install git (required by dep ensure)
RUN apk add git

# Download & install the dep
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/github.com/sinnott74/goblogserver
EXPOSE 8000
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only
COPY . ./
RUN CGO_ENABLED=0 go build


FROM scratch
WORKDIR /go/
EXPOSE 8000
ENV POSTGRES_URL=postgres://Sinnott@host.docker.internal:5432/pwadb?sslmode=disable&timezone=UTC
COPY --from=builder /go/src/github.com/sinnott74/goblogserver .
CMD ["./goblogserver"]