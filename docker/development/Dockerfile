#! Development Dockerfile
FROM golang:1.21-alpine3.18 AS development

#! Add git and build essentials for CGO
RUN apk add --no-cache \
    git \
    gcc \
    musl-dev \
    postgresql-client \
    curl \
    make \
    build-base

#! Set working directory
WORKDIR /app

#! Install specific version of air for hot reload
RUN go install github.com/cosmtrek/air@v1.40.4

#! Install specific version of golang-migrate
RUN CGO_ENABLED=1 go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

#! Copy wait-for script first
COPY docker/scripts/wait-for.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/wait-for.sh

#! Copy go mod files
COPY go.mod go.sum ./

#! Download dependencies
RUN go mod download

#! Copy the rest of the source code
COPY . .

#! Expose port
EXPOSE 8080

#! Use air for hot reload
CMD ["air"]