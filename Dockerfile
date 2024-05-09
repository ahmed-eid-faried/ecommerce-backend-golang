FROM golang

# Install Redis
RUN apt-get update && \
    apt-get install -y redis-server && \
    apt-get clean

# Install PostgreSQL client
RUN apt-get update && \
    apt-get install -y postgresql-client && \
    apt-get clean

# Set up working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download Go modules
RUN go mod download

# Copy the sample config file
COPY ./pkg/config/config.sample.yaml ./pkg/config/config.yaml

# Build the Go application
RUN go build -o /app/goshop ./cmd/api

# Expose port 8888
EXPOSE 8888

# Start Redis and PostgreSQL services
CMD service redis-server start && \
    service postgresql start && \
    /app/goshop
