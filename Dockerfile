###
# Builder to compile our golang code
###
FROM golang:alpine AS builder

WORKDIR /build
COPY . .

RUN go build -o eevee -v github.com/randomairborne/eevee/core

###
# Now generate our smaller image
###
FROM alpine

COPY --from=builder /build/eevee /go/bin/eevee

ENTRYPOINT ["/go/bin/absol"]
