FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o /out/api ./api

FROM alpine:3.20
RUN adduser -D -u 10001 appuser
USER appuser
COPY --from=build /out/api /usr/local/bin/api
COPY --from=build /app/migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/api"]
