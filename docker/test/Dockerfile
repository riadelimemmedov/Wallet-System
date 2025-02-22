#! Test Dockerfile
FROM golang:1.21-alpine3.18 AS test

#! Add git and basic tools
RUN apk add --no-cache \
    git \
    postgresql-client \
    curl \
    make \
    gcc

#! Set working directory
WORKDIR /app

#! Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

#! Copy the rest of the source code
COPY . .

#! Copy wait-for script
COPY docker/scripts/wait-for.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/wait-for.sh

#! Run tests
CMD ["go", "test", "-v", "-cover", "./..."]