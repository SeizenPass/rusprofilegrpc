# Build executable stage
FROM golang
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o grpcClient client/*.go
ENTRYPOINT /app/grpcClient
# Build final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=0 /app/grpcClient .
EXPOSE 8081
ENTRYPOINT ["./grpcClient", "-host=0.0.0.0:8081", "-grpc-server-endpoint=rusprofile-grpc-server:8085"]