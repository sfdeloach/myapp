FROM golang:1.25.5 AS build

WORKDIR /server

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Now copy the source code
COPY . .

# Builds the application as a statically linked binary,
# to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


FROM alpine:3.23.2

WORKDIR /server

# Add packages and create non-root user
RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && addgroup -g 1000 appuser \
    && adduser -D -u 1000 -G appuser appuser \
    && chown -R appuser:appuser /server

# Copy binary and assets
COPY --from=build /server/app .
COPY --chown=appuser:appuser ./public ./public
COPY --chown=appuser:appuser ./views ./views

RUN chmod +x /server/app

# Switch to non-root user
USER appuser

EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./app"]
