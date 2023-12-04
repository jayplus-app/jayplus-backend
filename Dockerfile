FROM golang:1.21 AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend cmd/app/main.go

FROM alpine:latest

COPY --from=builder /backend /backend

RUN apk --no-cache add tzdata

CMD ["/backend"]
