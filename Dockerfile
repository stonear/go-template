# syntax=docker/dockerfile:1

# Build the application from source
FROM stonear/golang:latest AS build-stage

WORKDIR /go/src

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN make gen
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o app -tags=viper_bind_struct

# Deploy the application binary into a lean image
FROM alpine:latest

RUN apk add -U tzdata
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

WORKDIR /

COPY --from=build-stage /go/src .

EXPOSE 8080

# Create a group and user
RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

# Tell docker that all future commands should run as the nonroot user
USER nonroot:nonroot

ENTRYPOINT ["./app"]
