# Use a base image for the app
FROM golang:1.19 AS builder

# Set the working directory
WORKDIR /src

# Copy the app files into the container
COPY . .

# Build the app
RUN go build -o ./email-collector

# Use a minimal base image for the app container
FROM ubuntu:latest

# Install ca-certificates package
RUN apt-get update && apt-get install -y ca-certificates

# Set the working directory
WORKDIR /app

# Copy the app binary from the build container
COPY --from=builder /src/email-collector ./

# Create a directory for the database and set permissions
RUN mkdir ./db && chmod 777 ./db

# Postmark environment variables to send emails
ENV POSTMARK_TOKEN "your-token"
ENV POSTMARK_FROM "from@yourdomain.com"
ENV POSTMARK_TEMPLATE_ALIAS "template-alias"

# Mount a volume for the database
VOLUME ./db

# Expose the port that the app listens on
EXPOSE 8080

# Start the app
ENTRYPOINT ["./email-collector"]
