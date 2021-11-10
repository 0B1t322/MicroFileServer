FROM golang:1.16 AS builder

WORKDIR /go/src/MFS/

ENV CGO_ENABLED=0
COPY go.* ./
RUN go mod download

COPY . .
RUN GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/MFS/main .

CMD ./main