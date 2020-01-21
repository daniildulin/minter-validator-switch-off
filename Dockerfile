FROM golang:alpine

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --exclude-dir=.git --build="go build -o ./builds/linux/switch ./cmd/switch.go" --command=./builds/linux/switch
