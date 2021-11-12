# syntax=docker/dockerfile:1

FROM golang:1.17.3 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY html/ ./html
COPY pkg/ ./pkg
COPY main.go ./

RUN go build -o /kubevis

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /kubevis /kubevis

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/kubevis"]
