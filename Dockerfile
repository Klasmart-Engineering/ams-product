# Use alpine image for the actual deployment
FROM alpine:latest

# Create app directory
WORKDIR /usr/src/app

# Copy the app
COPY ./bin/handler bin/
COPY configs configs

# Add missing certificates
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Bind the app port
EXPOSE 8044

# Start the app
ENTRYPOINT [ "bin/handler" ]
