#!/bin/bash

# test_image_converter.sh
# Tests all endpoints of the imageConverter API with image.jpg

# Configuration
SERVER_URL="http://localhost:8080"
INPUT_IMAGE="image.jpg"
OUTPUT_DIR="test_outputs"

# Ensure input image exists
if [ ! -f "$INPUT_IMAGE" ]; then
    echo "Error: Input image $INPUT_IMAGE not found"
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Function to check file size and validity
check_output() {
    local output_file=$1
    if [ -f "$output_file" ]; then
        size=$(stat -f%z "$output_file" 2>/dev/null || stat -c%s "$output_file")
        if [ "$size" -eq 0 ]; then
            echo "Error: Output file $output_file is empty (0 bytes)"
        else
            echo "Success: Output file $output_file created (size: $size bytes)"
            file "$output_file"
        fi
    else
        echo "Error: Output file $output_file was not created"
    fi
}

# Test /convert endpoint (JPEG output)
echo "Testing /convert endpoint..."
curl -v -X POST "$SERVER_URL/convert" \
    -F "file=@$INPUT_IMAGE" \
    -F "output_format=jpeg" \
    -F "quality=80" \
    -o "$OUTPUT_DIR/converted_image.jpg"
check_output "$OUTPUT_DIR/converted_image.jpg"
echo

# Test /square-crop endpoint
echo "Testing /square-crop endpoint..."
curl -v -X POST "$SERVER_URL/square-crop" \
    -F "file=@$INPUT_IMAGE" \
    -o "$OUTPUT_DIR/square_cropped_image.jpg"
check_output "$OUTPUT_DIR/square_cropped_image.jpg"
echo

# Test /fit-to-square endpoint (PNG output)
echo "Testing /fit-to-square endpoint..."
curl -v -X POST "$SERVER_URL/fit-to-square" \
    -F "file=@$INPUT_IMAGE" \
    -F "output_format=png" \
    -F "quality=80" \
    -o "$OUTPUT_DIR/square_fitted_image.png"
check_output "$OUTPUT_DIR/square_fitted_image.png"
echo

# Test /invert endpoint
echo "Testing /invert endpoint..."
curl -v -X POST "$SERVER_URL/invert" \
    -F "file=@$INPUT_IMAGE" \
    -o "$OUTPUT_DIR/inverted_image.jpg"
check_output "$OUTPUT_DIR/inverted_image.jpg"
echo

# Test /apply-filter endpoint (grayscale filter)
echo "Testing /apply-filter endpoint (grayscale)..."
curl -v -X POST "$SERVER_URL/apply-filter" \
    -F "file=@$INPUT_IMAGE" \
    -F "filter_name=grayscale" \
    -F "intensity=20" \
    -o "$OUTPUT_DIR/grayscale_image.jpg"
check_output "$OUTPUT_DIR/grayscale_image.jpg"
echo

# Test /apply-filter endpoint (blur filter)
echo "Testing /apply-filter endpoint (blur)..."
curl -v -X POST "$SERVER_URL/apply-filter" \
    -F "file=@$INPUT_IMAGE" \
    -F "filter_name=blur" \
    -F "intensity=5" \
    -o "$OUTPUT_DIR/blurred_image.jpg"
check_output "$OUTPUT_DIR/blurred_image.jpg"
echo

echo "All tests completed. Check $OUTPUT_DIR for output images."
