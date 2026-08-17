ARG APP_NAME=server
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY motadata-go-sdk/ motadata-go-sdk/
RUN go mod download
COPY src/ src/
ARG APP_NAME
RUN CGO_ENABLED=0 go build -o ${APP_NAME} ./src/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
ARG APP_NAME
COPY --from=builder /app/${APP_NAME} /usr/local/bin/app
ENTRYPOINT ["app"]
