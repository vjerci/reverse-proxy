FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /app/jsonendpoint ./cmd/jsonendpoint/jsonendpoint.go


FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/jsonendpoint /app/jsonendpoint
RUN chmod 755 /app/jsonendpoint

EXPOSE 8000
CMD [ "/app/jsonendpoint" ]


