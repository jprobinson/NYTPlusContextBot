FROM alpine:3.3

ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD bot /bot

CMD ["/bot"]
