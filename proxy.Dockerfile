FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /app/proxy ./cmd/proxy/proxy.go


FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/proxy /app/proxy
RUN chmod 755 /app/proxy

EXPOSE 8000
CMD [ "/app/proxy" ]


