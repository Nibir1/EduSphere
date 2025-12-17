# server/Dockerfile

FROM golang:1.25-bookworm

WORKDIR /app

# 1. Install System Dependencies
# We include 'curl' to manually download the Tesseract language file
RUN apt-get update && apt-get install -y \
    poppler-utils \
    tesseract-ocr \
    libtesseract-dev \
    libleptonica-dev \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 2. Cache Go Modules
COPY go.mod go.sum ./
RUN go mod download

# 3. Copy Source Code
COPY . .

# 4. Build the Application
RUN go build -o main main.go

# ------------------------------------------------------------------------------
# 5. Runtime Configuration (The Bulletproof Fix)
# ------------------------------------------------------------------------------
# A. Create a clean, known directory for Tesseract data
RUN mkdir -p /usr/share/tessdata

# B. Manually download the English training data (Fast version)
# This bypasses any confusion about where apt-get installed it.
RUN curl -L -o /usr/share/tessdata/eng.traineddata \
    https://github.com/tesseract-ocr/tessdata_fast/raw/main/eng.traineddata

# C. Point TESSDATA_PREFIX to that specific directory
ENV TESSDATA_PREFIX /usr/share/tessdata

EXPOSE 8080

# Start the application
CMD ["./main"]