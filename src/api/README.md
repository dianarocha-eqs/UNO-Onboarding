# GO API

### Prerequisites

Make sure you have the following tools installed on your system:

- [Go](https://golang.org/dl/) (v1.23)
- [Docker](https://www.docker.com/get-started) (optional, for running the database locally)
- [Gorm](https://gorm.io/index.html) (ORM for Go, used for database interactions)

---

## Setup, Build & Run

### 1. **Install Go dependencies**:
   
   Inside the project directory, run the following command to download and install the necessary Go dependencies.

   ```bash
   go mod tidy
   ```

### 2. **Build the API**:

   Run the following command to build the Go application:

   ```bash
   cd cmd/
   go build -o uno-onboarding-api .
   ```

   This will create an executable file named `uno-onboarding-api`.

### 3. **Run the API**:

   After building the application, you can run it using the following command:

   ```bash
   cd cmd/
   ./uno-onboarding-api
   ```

   If you're in the development environment and want to run the application directly from the Go source code, use:

   ```bash
   cd cmd/
   go run .
   ```

   The server will start on port `8080` by default.

---
