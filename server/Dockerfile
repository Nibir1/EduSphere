# ------------------------------------------------------------------------------
# Build Stage
# ------------------------------------------------------------------------------
FROM golang:alpine AS builder

WORKDIR /app

# Install dependencies for CGO (required by gosseract)
# - build-base: contains gcc, g++, etc.
# - tesseract-ocr-dev: contains C headers for compilation
# - leptonica-dev: image processing lib headers
RUN apk add --no-cache build-base tesseract-ocr-dev leptonica-dev pkgconfig

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main main.go

# ------------------------------------------------------------------------------
# Run Stage
# ------------------------------------------------------------------------------
FROM alpine:3.19
WORKDIR /app

# Install Runtime Dependencies
# - tesseract-ocr: required for the compiled app to run OCR
# - poppler-utils: required for 'pdftoppm' command (pdf -> image conversion)
# - curl: for healthchecks
RUN apk add --no-cache curl tesseract-ocr poppler-utils

COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder /app/db/migration ./db/migration

EXPOSE 8080
CMD ["/app/main"]