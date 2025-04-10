# Use a debian base image (which includes glibc)
FROM debian:latest

# Install necessary dependencies including libc6 and CA certificates
RUN apt-get update && apt-get install -y \
  libc6 \
  ca-certificates

# Set a working directory inside the container
WORKDIR /app

# Copy everything
COPY . .

# Copy the Go executable into the container
COPY cmd/countries-dashboard-service/main .

# Copy your Firestore service account key file into the container (optional if needed)
# COPY path/to/your/service_account_key.json /app/service_account_key.json

# Make the Go executable the entry point
CMD ["./main"]
