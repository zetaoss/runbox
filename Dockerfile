FROM golang:1.21-alpine AS build-stage
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -C pkg -o /runbox

FROM alpine:latest
WORKDIR /
COPY --from=build-stage /runbox /runbox
RUN apk add --no-cache docker-cli
EXPOSE 8080
ENTRYPOINT ["/runbox"]
