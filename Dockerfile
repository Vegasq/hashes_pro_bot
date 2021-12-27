FROM golang:1.17 AS builder
WORKDIR /root/
RUN apt-get update
RUN git clone https://github.com/Vegasq/hashes_pro_bot.git
RUN cd hashes_pro_bot && CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"'

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /root/hashes_pro_bot/hpbot /
CMD ["/hpbot"]
