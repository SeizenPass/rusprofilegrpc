# Build executable stage
FROM golang
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o grpcServer server/*.go
ENTRYPOINT /app/grpcServer
# Build final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=0 /app/grpcServer .
EXPOSE 8085
ENTRYPOINT ["./grpcServer", "-host=0.0.0.0:8085"]