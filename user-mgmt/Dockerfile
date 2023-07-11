# syntax=docker/dockerfile:1

################################
# Build
################################
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN ./scripts/build.sh

################################
# Deploy
################################
FROM --platform=linux/amd64 debian:bullseye-slim

# Copy built binaries
RUN mkdir -p /app
COPY --from=build /app/bin/users-api /app/
COPY --from=build /app/scripts /app/scripts
COPY --from=build /app/migrations /app/migrations
WORKDIR /app

# Run app
EXPOSE 8080
ENTRYPOINT ["./scripts/entrypoint.sh"]