FROM golang:1.22.0

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on

COPY . .

RUN go mod tidy
RUN go mod download
RUN go build -o ./bin/main .
RUN go test ./...

CMD ["./bin/main"]
