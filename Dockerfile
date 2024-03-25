FROM golang:1.22.1 as builder
ENV GOPROXY='https://proxy.golang.com.cn,direct'
WORKDIR /app
COPY . .
RUN go build

FROM alpine:3.19
WORKDIR /app
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=builder /app/budget_exporter /app/budget_exporter
COPY --from=builder /app/configs/budget_exporter.yaml.example /app/budget_exporter.yaml
EXPOSE 9901
VOLUME /app/data
ENTRYPOINT ["/app/budget_exporter"]
