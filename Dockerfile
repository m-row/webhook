# BUILDER ---------------------------------------------------------------------
FROM golang:1.22.1-alpine3.19 AS builder

WORKDIR /build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download && \ 
    apk --update --no-cache add ca-certificates git

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -ldflags="-s -w -X main.commitCount=$(git rev-list --count HEAD) -X main.commitDescribe=$(git describe --always)" \
    -o main ./

# FINAL APP -------------------------------------------------------------------
FROM scratch

WORKDIR /
# COPY --from=migrator ["/migrations", "/migrations"]
COPY --from=builder ["/build/main", "/"]
COPY --from=builder ["/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/ca-certificates.crt"]

EXPOSE 8000

ENTRYPOINT ["/main"]

