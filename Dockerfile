FROM golang:1.16-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -ldflags="-X 'github.com/daniildulin/minter-validator-switch-off/core.Vs=$PHRASE'" -o ./builds/linux/switcher ./cmd/switch.go

FROM alpine:3.13

COPY --from=builder /app/builds/linux/switcher /usr/bin/switcher
RUN addgroup minteruser && adduser -D -h /minter -G minteruser minteruser
USER minteruser
WORKDIR /minter
CMD ["/usr/bin/switcher"]
