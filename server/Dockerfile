# server/Dockerfile

# Use a single-stage build with Go installed.
# This ensures we have the 'go' command available in the final container
# to run your seed script (db/seed/seed_courses.go).
FROM golang:1.25-bookworm

# Set the working directory inside the container
WORKDIR /app

# ------------------------------------------------------------------------------
# 1. Install System Dependencies
# ------------------------------------------------------------------------------
# We install:
# - poppler-utils: for PDF -> Image conversion (pdftoppm)
# - tesseract-ocr: the OCR engine
# - tesseract-ocr-eng: English language data for OCR
# - libtesseract-dev: headers for compiling Go CGO bindings
# - libleptonica-dev: image processing headers
RUN apt-get update && apt-get install -y \
    poppler-utils \
    tesseract-ocr \
    tesseract-ocr-eng \
    libtesseract-dev \
    libleptonica-dev \
    && rm -rf /var/lib/apt/lists/*

# ------------------------------------------------------------------------------
# 2. Cache Go Modules
# ------------------------------------------------------------------------------
COPY go.mod go.sum ./
RUN go mod download

# ------------------------------------------------------------------------------
# 3. Copy Source Code
# ------------------------------------------------------------------------------
# We copy the entire project directory.
# This ensures 'db/seed/seed_courses.go' is present in the container.
COPY . .

# ------------------------------------------------------------------------------
# 4. Build the Application
# ------------------------------------------------------------------------------
RUN go build -o main main.go

# ------------------------------------------------------------------------------
# 5. Runtime Configuration
# ------------------------------------------------------------------------------
# ⭐ FIX: Point to the generic path for Tesseract 5 (Debian Bookworm)
# The '/4.00/' subdirectory is removed because it no longer exists in this version.
ENV TESSDATA_PREFIX /usr/share/tesseract-ocr/tessdata

EXPOSE 8080

# Start the application
CMD ["./main"]