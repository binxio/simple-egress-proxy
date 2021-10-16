FROM golang:1.17

WORKDIR /app
ADD ./ /app/
RUN CGO_ENABLED=0 go build -ldflags '-w -s' .

FROM scratch

COPY --from=0 /app/simple-egress-proxy /
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt


EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/simple-egress-proxy"]
