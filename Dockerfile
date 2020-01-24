FROM golang:1.13.5-alpine

WORKDIR /app

COPY ./ /app

RUN apk add --no-cache make gcc musl-dev linux-headers git

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --exclude-dir=.git --build="go build -o ./builds/linux/switch ./cmd/switch.go" --command=./builds/linux/switch
