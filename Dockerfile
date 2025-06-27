FROM golang:1.23-alpine AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /runbox

FROM alpine:latest
WORKDIR /
COPY --from=build-stage /runbox /runbox
RUN apk add --no-cache docker-cli
EXPOSE 8080
ENTRYPOINT ["/runbox"]
