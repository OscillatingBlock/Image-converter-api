

# ImageConverter API Documentation

This document provides details for frontend developers integrating with the `imageConverter` API, a Go-based REST API for image processing. The API supports image conversion, cropping, fitting to square, color inversion, and applying filters (grayscale, blur). It runs on `http://localhost:8080` by default.

## Base URL

```
http://localhost:8080
```

## Endpoints

All endpoints use `POST` requests with `multipart/form-data` content type, requiring an image file (`file`) and additional parameters as form fields. Responses are either binary image data (on success) or JSON error messages (on failure).

### 1. Convert Image (`/convert`)

Converts an image to a specified format (JPEG, PNG, WebP) with adjustable quality.

- **Method**: `POST`
- **Path**: `/convert`
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `file` (file, required): Image file (JPEG, PNG, WebP).
  - `output_format` (string, required): Output format (`jpeg`, `png`, `webp`).
  - `quality` (integer, optional): Quality (1–100, default: 80).
- **Response**:
  - **Success**: `200 OK`, binary image data, `Content-Type: image/<output_format>`.
  - **Error**: `400 Bad Request` or `500 Internal Server Error`, JSON: `{"error": "message"}`.
- **Example Image**: [test_outputs/converted_image.jpg](test_outputs/converted_image.jpg)

**Example Request (JavaScript)**:
```javascript
const formData = new FormData();
formData.append('file', document.querySelector('input[type="file"]').files[0]);
formData.append('output_format', 'jpeg');
formData.append('quality', '80');

fetch('http://localhost:8080/convert', {
  method: 'POST',
  body: formData
})
  .then(response => {
    if (!response.ok) throw new Error(response.statusText);
    return response.blob();
  })
  .then(blob => {
    const url = URL.createObjectURL(blob);
    const img = new Image();
    img.src = url;
    document.body.appendChild(img);
  })
  .catch(error => console.error('Error:', error));
```

### 2. Square Crop (`/square-crop`)

Crops an image to a square by trimming edges.

- **Method**: `POST`
- **Path**: `/square-crop`
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `file` (file, required): Image file (JPEG, PNG, WebP).
- **Response**:
  - **Success**: `200 OK`, binary JPEG data, `Content-Type: image/jpeg`.
  - **Error**: `400 Bad Request` or `500 Internal Server Error`, JSON: `{"error": "message"}`.
- **Example Image**: [test_outputs/square_cropped_image.jpg](test_outputs/square_cropped_image.jpg)

**Example Request (JavaScript)**:
```javascript
const formData = new FormData();
formData.append('file', document.querySelector('input[type="file"]').files[0]);

fetch('http://localhost:8080/square-crop', {
  method: 'POST',
  body: formData
})
  .then(response => {
    if (!response.ok) throw new Error(response.statusText);
    return response.blob();
  })
  .then(blob => {
    const url = URL.createObjectURL(blob);
    const img = new Image();
    img.src = url;
    document.body.appendChild(img);
  })
  .catch(error => console.error('Error:', error));
```

### 3. Fit to Square (`/fit-to-square`)

Places an image on a square canvas with padding (white background for JPEG/WebP, transparent for PNG).

- **Method**: `POST`
- **Path**: `/fit-to-square`
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `file` (file, required): Image file (JPEG, PNG, WebP).
  - `output_format` (string, required): Output format (`jpeg`, `png`, `webp`).
  - `quality` (integer, optional): Quality (1–100, default: 80).
- **Response**:
  - **Success**: `200 OK`, binary image data, `Content-Type: image/<output_format>`.
  - **Error**: `400 Bad Request` or `500 Internal Server Error`, JSON: `{"error": "message"}`.
- **Example Image**: [test_outputs/square_fitted_image.png](test_outputs/square_fitted_image.png)

**Example Request (JavaScript)**:
```javascript
const formData = new FormData();
formData.append('file', document.querySelector('input[type="file"]').files[0]);
formData.append('output_format', 'png');
formData.append('quality', '80');

fetch('http://localhost:8080/fit-to-square', {
  method: 'POST',
  body: formData
})
  .then(response => {
    if (!response.ok) throw new Error(response.statusText);
    return response.blob();
  })
  .then(blob => {
    const url = URL.createObjectURL(blob);
    const img = new Image();
    img.src = url;
    document.body.appendChild(img);
  })
  .catch(error => console.error('Error:', error));
```

### 4. Invert Colors (`/invert`)

Inverts the colors of an image.

- **Method**: `POST`
- **Path**: `/invert`
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `file` (file, required): Image file (JPEG, PNG, WebP).
- **Response**:
  - **Success**: `200 OK`, binary JPEG data, `Content-Type: image/jpeg`.
  - **Error**: `400 Bad Request` or `500 Internal Server Error`, JSON: `{"error": "message"}`.
- **Example Image**: [test_outputs/inverted_image.jpg](test_outputs/inverted_image.jpg)

**Example Request (JavaScript)**:
```javascript
const formData = new FormData();
formData.append('file', document.querySelector('input[type="file"]').files[0]);

fetch('http://localhost:8080/invert', {
  method: 'POST',
  body: formData
})
  .then(response => {
    if (!response.ok) throw new Error(response.statusText);
    return response.blob();
  })
  .then(blob => {
    const url = URL.createObjectURL(blob);
    const img = new Image();
    img.src = url;
    document.body.appendChild(img);
  })
  .catch(error => console.error('Error:', error));
```

### 5. Apply Filter (`/apply-filter`)

Applies a grayscale or blur filter with adjustable intensity.

- **Method**: `POST`
- **Path**: `/apply-filter`
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `file` (file, required): Image file (JPEG, PNG, WebP).
  - `filter_name` (string, required): Filter type (`grayscale`, `blur`).
  - `intensity` (integer, optional): Filter intensity (default: 10, recommended 0–100).
- **Response**:
  - **Success**: `200 OK`, binary JPEG data, `Content-Type: image/jpeg`.
  - **Error**: `400 Bad Request` or `500 Internal Server Error`, JSON: `{"error": "message"}`.
- **Example Images**:
  - Grayscale: [test_outputs/grayscale_image.jpg](test_outputs/grayscale_image.jpg)
  - Blur: [test_outputs/blurred_image.jpg](test_outputs/blurred_image.jpg)

**Example Request (JavaScript)**:
```javascript
const formData = new FormData();
formData.append('file', document.querySelector('input[type="file"]').files[0]);
formData.append('filter_name', 'grayscale');
formData.append('intensity', '20');

fetch('http://localhost:8080/apply-filter', {
  method: 'POST',
  body: formData
})
  .then(response => {
    if (!response.ok) throw new Error(response.statusText);
    return response.blob();
  })
  .then(blob => {
    const url = URL.createObjectURL(blob);
    const img = new Image();
    img.src = url;
    document.body.appendChild(img);
  })
  .catch(error => console.error('Error:', error));
```

## Error Handling

- **400 Bad Request**: Invalid input (e.g., missing `file`, invalid `output_format`).
  ```json
  {"error": "Output format is required"}
  ```
- **500 Internal Server Error**: Server-side issues (e.g., image decoding failure).
  ```json
  {"error": "Error while converting image"}
  ```

**JavaScript Error Handling**:
```javascript
fetch('http://localhost:8080/convert', { method: 'POST', body: formData })
  .then(response => {
    if (!response.ok) {
      return response.json().then(err => { throw new Error(err.error); });
    }
    return response.blob();
  })
  .catch(error => alert(`Error: ${error.message}`));
```

## Notes for Frontend Developers

- **File Input**: Use `<input type="file">` to let users select images. Ensure files are valid (JPEG, PNG, WebP) and under 10MB (configurable server-side).
- **Output Display**: Successful responses are image blobs. Use `URL.createObjectURL` to display or download them.
- **CORS**: If the frontend is hosted separately, ensure the server supports CORS (add `middleware.CORS()` in `main.go` if needed).
- **Progress Indicators**: For large images, show a loading indicator, as processing may take a few seconds.
- **Testing**: Use the images in `test_outputs/` to verify expected results. Run `test_image_converter.sh` to generate them.

## Setup for Testing

1. **Start the Server**:
   ```bash
   cd imageConverter
   go run main.go
   ```

2. **Test Images**:
   - Example input: `image.jpg`
   - Outputs: `test_outputs/` (e.g., `converted_image.jpg`, `grayscale_image.jpg`)

3. **Dependencies**:
   - Backend: Go 1.16+, `github.com/labstack/echo/v4`, `github.com/disintegration/imaging`, `golang.org/x/image/webp`.

## Contact

For issues or feature requests, contact via github or email (@aayusharma270506@gmail.com).
