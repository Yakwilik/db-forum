FROM golang:1.20.4-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app cmd/db-forum/main.go

FROM postgres:13-alpine

# Copy the application from the builder stage
COPY --from=builder /app/app /app/app
COPY --from=builder /app/db /app/db
COPY --from=builder /app/db/db.sql /docker-entrypoint-initdb.d/db.sql
ENV POSTGRES_PASSWORD=mysecretpassword

# Install the `supervisord`
RUN apk add --no-cache supervisor

EXPOSE 5432
EXPOSE 5000
EXPOSE 9001

# Copy supervisor configuration
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
# Run the command
