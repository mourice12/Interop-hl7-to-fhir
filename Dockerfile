#Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

#Copy go mod files
COPY . .

#Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

#Final stage - tiny image
FROM alpine:latest

WORKDIR /

#Copy binary from builder
COPY --from=builder /server /server
#Expose Port
EXPOSE 8000

#Run the server
CMD ["/server"]