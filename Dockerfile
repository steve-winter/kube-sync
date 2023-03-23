
FROM golang:1.20.2-buster@sha256:5f9d227a93a4453df44827dbf44acef84e0b0c3abf39470d8a69b28bbd56733f AS builder

WORKDIR /app-build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /kube-sync
##
## Deploy
##
FROM gcr.io/distroless/static-debian11@sha256:72924583773eeeb9a6200e9f6dbfd95a27fbf25d39bfe7062c46d2654628f007
#COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /etc/group /etc/group

WORKDIR /
COPY --from=builder --chown=nonroot /kube-sync /app/kube-sync
USER nonroot

ENTRYPOINT ["/app/kube-sync"]
