# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine3.12 AS build-env

# Copy the local package files to the container's workspace.

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN apk add --no-cache git
RUN apk add make

ADD ./atlas.com/wrg /atlas.com/wrg
WORKDIR /atlas.com/wrg

RUN go build -o /server

RUN make swagger

FROM alpine:3.12

# Port 8080 belongs to our application
EXPOSE 8080

RUN apk add --no-cache libc6-compat

WORKDIR /

COPY --from=build-env /server /
COPY --from=build-env /atlas.com/wrg/swagger.yaml /
COPY /atlas.com/wrg/config.yaml /

CMD ["/server"]