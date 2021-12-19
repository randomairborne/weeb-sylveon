###
# Builder to compile our golang code
###
FROM golang:alpine AS builder

WORKDIR /build
COPY . .

RUN go build -o absol -v github.com/randomairborne/eevee/core

###
# Now generate our smaller image
###
FROM alpine

COPY --from=builder /build/absol /go/bin/absol

ENTRYPOINT ["/go/bin/absol"]
CMD ["alert", "cleaner", "factoids", "log", "twitch", "hjt", "search", "mcping"]
