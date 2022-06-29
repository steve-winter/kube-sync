
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /kube-sync

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /kube-sync /kube-sync

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/kube-sync"]
