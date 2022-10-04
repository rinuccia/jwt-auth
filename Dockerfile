# Step 1: Modules caching
FROM golang:1.19-alpine3.16 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.19-alpine3.16 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /build
WORKDIR /build
RUN go build -o main .

# Step 3: Final
FROM alpine:3.16
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY . /app
COPY --from=builder /build/main /app
WORKDIR /app
EXPOSE 9000
CMD ["./main"]