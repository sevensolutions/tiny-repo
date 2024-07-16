# -------------------------------------------------------------
# Build Stage
# -------------------------------------------------------------

FROM golang:1.22-alpine AS builder

# Install git
# Git is required for fetching the dependencies
RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/sevensolutions/tiny-repo/

COPY . .

# Fetch dependencies
RUN go get -d -v

# Build the binary
RUN go build -o /go/bin/artifact-repo

# -------------------------------------------------------------
# Runtime Stage
# -------------------------------------------------------------

FROM scratch

COPY --from=builder /go/bin/artifact-repo /go/bin/artifact-repo

ENTRYPOINT ["/go/bin/artifact-repo"]
