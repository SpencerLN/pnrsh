FROM golang:1.24-bookworm AS builder

WORKDIR /build
COPY . /build
RUN cd /build/cmd && go build -v

FROM debian:bookworm-slim

EXPOSE 8080
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/cmd/cmd /pnrsh
ENTRYPOINT /pnrsh
