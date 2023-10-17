
FROM golang:1.20.2-buster@sha256:57dbdd5c8fe24e357b15a4ed045b0b1607a99a1465b5101304ea39e11547be27 AS builder

WORKDIR /app-build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /kube-sync
##
## Deploy
##
FROM gcr.io/distroless/static-debian11@sha256:6706c73aae2afaa8201d63cc3dda48753c09bcd6c300762251065c0f7e602b25
#COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /etc/group /etc/group

WORKDIR /
COPY --from=builder --chown=nonroot /kube-sync /app/kube-sync
USER nonroot

ENTRYPOINT ["/app/kube-sync"]
