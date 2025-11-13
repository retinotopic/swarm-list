FROM golang:1.24.3-alpine3.21 AS builder
# RUN apk add --no-cache git

WORKDIR /build
COPY --from=cmd_server . .
COPY go.mod go.sum .
COPY /server ./server
RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go

# FROM nixos/nix:latest AS builder
# WORKDIR /build
# COPY --from=cmd_server . .
# COPY go.mod go.sum .
# COPY /server ./server

# RUN nix \
#     --extra-experimental-features "nix-command flakes" \
#     --option filter-syscalls false \
#     build
# RUN mkdir /tmp/nix-store-closure
# RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure



FROM alpine:3.21
WORKDIR /app
RUN apk add --no-cache curl
COPY --from=builder /build/server/db/postgres/migrations /bin/maindir/migrations
# COPY --from=builder /tmp/nix-store-closure /nix/store
# COPY --from=builder /build/result/ /app

COPY --from=builder /main /app
# RUN echo $(ls -a /app)
# RUN echo $(ls -a /app/bin)
# COPY --from=certs /key.pem /cert.pem /bin/maindir
ENTRYPOINT ["/app/main"]
