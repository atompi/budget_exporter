FROM golang:1.22.1 as builder
ENV GOPROXY='https://proxy.golang.com.cn,direct'
WORKDIR /app
COPY . .
RUN go build

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/budget_exporter /app/budget_exporter
COPY --from=builder /app/configs/budget_exporter.yaml.example /app/budget_exporter.yaml
EXPOSE 9901
VOLUME /app/data
ENTRYPOINT ["/app/budget_exporter"]
