
FROM golang:1.19.0-buster AS builder

WORKDIR /app-build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /kube-sync
##
## Deploy
##
FROM gcr.io/distroless/static-debian11
#COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /etc/group /etc/group

WORKDIR /
COPY --from=builder --chown=nonroot /kube-sync /app/kube-sync
USER nonroot

ENTRYPOINT ["/app/kube-sync"]
