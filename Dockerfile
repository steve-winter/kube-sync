
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN useradd -u 10001 scratchuser

RUN CGO_ENABLED=0 GOOS=linux go build -o /kube-sync

##
## Deploy
##
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/passwd /etc/passwd

WORKDIR /
COPY --from=build /kube-sync /kube-sync

USER scratchuser
ENTRYPOINT ["/kube-sync"]
