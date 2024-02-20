FROM golang:1.22.0 AS builder

COPY . /go/src/github.com/opa-/nvidia-smi-rest
WORKDIR /go/src/github.com/opa-/nvidia-smi-rest

RUN go build -o /go/bin/nvidia-smi-rest


FROM nvidia/cuda:12.3.1-base-ubuntu22.04

WORKDIR /app
COPY --from=builder /go/bin/nvidia-smi-rest /app/nvidia-smi-rest

EXPOSE 7777
ENTRYPOINT ["/app/nvidia-smi-rest"]
