FROM golang:1.22-alpine AS builder

RUN apk update && apk add --no-cache git

ENV GO111MODULE=on

WORKDIR /

COPY . .

RUN go get -d -v

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o /app

RUN mkdir /new_tmp && chmod 777 /new_tmp

# Second-stage using an image without anything! :-)
FROM scratch

COPY --from=builder /new_tmp /tmp
COPY --from=builder /app /app

ENTRYPOINT ["/app"]