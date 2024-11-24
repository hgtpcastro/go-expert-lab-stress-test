FROM golang:1.22.4 as builder

WORKDIR /app
COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ./go-http-bench ./benchmark.go

FROM debian:buster-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /root/

COPY --from=builder /app/go-http-bench .

ENTRYPOINT ["./go-http-bench"]