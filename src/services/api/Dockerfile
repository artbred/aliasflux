FROM golang:1.20 AS builder

LABEL version="1.0"

RUN mkdir /go/src/aliasflux
COPY . /go/src/aliasflux
WORKDIR /go/src/aliasflux

RUN go mod download
RUN cd src/services/api && CGO_ENABLED=0 GOOS=linux go build

FROM scratch AS production

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl /etc/ssl
COPY --from=builder /go/src/aliasflux .

CMD ["./src/services/api/api"]
