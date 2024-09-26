# Receipt Uploader

This project is a receipt uploader service for handling image uploads and scaling them for use in user interfaces. The service is built using Go (Golang) and is containerized using Docker. The service supports basic image upload functionality with optional image resizing, and includes a simple permissions system (stretch goal) to ensure users can only access their own receipts.

## Features

- Upload images of receipts.
- Resize images to different resolutions (proportional scaling, not stretched).
- List all uploaded receipts for a user.
- Fetch specific receipts by ID, with optional resizing.
- Built-in unit tests for services, models, and handlers.
- Containerized using Docker for easy deployment.

## Requirements

- Go 1.18 or later
- Docker
- Git

## Project structure

```
.
└── receipt-uploader/
    ├── handlers/                           # Contains HTTP handlers for uploading, fetching, and listing receipts.
    │   ├── receipts_test.go
    │   └── receipts.go
    ├── models/                             # Manages receipt metadata and file storage.
    │   ├── receipt_test.go
    │   └── receipt.go
    ├── services/                           # Contains helper functions for file handling, image processing, and unit tests.
    │   ├── image_service_test.go
    │   ├── image_service.go
    │   ├── storage_test.go
    │   └── storage.go
    ├── testdata/                           # Contains sample data (e.g., test images).
    ├── Dockerfile                          # Dockerfile for containerizing the Go application.
    ├── go.mod                              # Go module dependencies.
    ├── go.sum                              # Checksums for Go modules.
    └── README.md                           # Project documentation.
```

## Getting Started

### Prerequisites

- Install [Go](https://golang.org/dl/) (1.18 or later).
- Docker

### Clone the Repository

```bash
git clone https://github.com/hpkhanhloc/receipt-uploader.git
cd receipt-uploader
```

### Run locally

1. **Install dependencies:**

```
go mod download
```

2. **Run the application:**

```
go run main.go
```

3. **Access the application:**

Open your browser and navigate to http://localhost:8080.

### Run Tests

To run the tests for the application, use:

```
go test ./... -v
```

This will run all tests for handlers, models, and services.

## Using Docker

### Build Docker Image

You can build the Docker image with:

```
docker build -t receipt-uploader .
```

### Build Docker Image

```
docker run -p 8080:8080 receipt-uploader
```

### Access the Application

Once the container is running, you can access the service at http://localhost:8080.

## API Endpoints

### X-User-ID Header

The `X-User-ID` header is used to authenticate the user making the request. Each user is assigned a unique user ID, and this ID must be included in the request headers for all API calls.

- **Purpose**: This header helps the service associate uploaded receipts with the correct user and enforce permissions.
- **Requirement**: The service will return an error if the `X-User-ID` header is missing or if a user tries to access another user’s receipts.
- **Example Usage**:
  - In `curl`:
    ```bash
    curl -H "X-User-ID: user123" -F "file=@receipt.jpg" http://localhost:8080/receipts
    ```
  - In API requests from client-side applications, this header must be included in each API request to identify the user.

### Upload Receipt

- **URL**: `/receipts`
- **Method**: `POST`
- **Headers**: `X-User-ID`
- **Content-Type**: `multipart/form-data`
- **Description**: Upload an image of a receipt. Returns the receipt ID.
- **Example**:
  ```bash
  curl -H "X-User-ID: user123" -F "file=@receipt.jpg" http://localhost:8080/receipts
  ```

### Get Receipt by ID

- **URL**: `/receipts/{receipt_id}`
- **Method**: `GET`
- **Headers**: `X-User-ID`
- **Description**: Fetch a receipt by its ID. Optional query parameters: width, height for resizing.
- **Example**:
  ```bash
  curl -H "X-User-ID: user123" http://localhost:8080/receipts/{receipt_id}?width=200&height=200
  ```

### List User Receipts

- **URL**: `/receipts`
- **Method**: `GET`
- **Headers**: `X-User-ID`
- **Description**: List all receipts for the authenticated user. The X-User-ID header ensures that the user can only retrieve receipts they own.
- **Example**:
  ```bash
  curl -H "X-User-ID: user123" http://localhost:8080/receipts
  ```
