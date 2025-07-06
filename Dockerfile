# Stage 1: Build React frontend
FROM node:lts AS frontend-builder

# Set working directory
WORKDIR /app/website
COPY website/ .

# Install pnpm and dependencies
RUN npm install -g pnpm
RUN pnpm install --frozen-lockfile

RUN pnpm run build


# Stage 2: Build Go backend
FROM golang:1.24 AS backend-builder

# Set working directory
WORKDIR /app

COPY . .
RUN make build


# Stage 3: Final runtime image
FROM ubuntu:latest

# Install ca-certificates for HTTPS requests
RUN apt-get update && apt-get install -y ca-certificates


# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=backend-builder /app/bin/web .
COPY --from=backend-builder /app/configs ./configs

COPY --from=frontend-builder /app/website/dist ./website/dist

# Expose port
EXPOSE 8080
CMD ["./web"]
