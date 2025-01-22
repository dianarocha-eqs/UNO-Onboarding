# GO API

## Installation

Follow these steps to set up and run the Go API on your local machine.

### Prerequisites

Make sure you have the following tools installed on your system:

- [Go](https://golang.org/dl/) (v1.23)
- [Docker](https://www.docker.com/get-started) (optional, for running the database locally)
- [Gorm](https://gorm.io/index.html) (ORM for Go, used for database interactions)

### Clone the Repository

```bash
git clone https://github.com/yourusername/UNO-Onboarding.git
git checkout task/IP-241-api-base-structure
cd src/api

```
---

### Set Up the Azure SQL Server and Database

Run the following commands to set up the Azure SQL Server, configure firewall rules, and create the database:

1. **Create the Azure SQL Server**:
   ```bash
   az sql server create --name uno-onboarding --resource-group UNOnboarding --location northeurope --admin-user eqs-digital --admin-password B3N5GA!QpgCX^&@Bt+pXhTY6NcZD7NWw
   ```

2. **Add a Firewall Rule to Allow Your IP Address**:
   Replace `192.168.1.233` with your actual public IP address:
   ```bash
   az sql server firewall-rule create --resource-group UNOnboarding --server uno-onboarding --name AllowMyIP --start-ip-address 192.168.1.233 --end-ip-address 192.168.1.233
   ```

3. **Create the Azure SQL Database**:
   ```bash
   az sql db create \
     --resource-group UNOnboarding \
     --server uno-onboarding \
     --name uno-onboarding \
     --service-objective S1
   ```
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

## API Endpoints

Here are the available endpoints for interacting with the sensor data in sensors/routes.go:

- `GET /api/sensors`  
  Retrieves a list of all sensors.

- `GET /api/sensors/:id`  
  Retrieves a specific sensor by its ID.

- `POST /api/sensors`  
  Creates a new sensor.

- `PUT /api/sensors/:id`  
  Updates an existing sensor by its ID.

- `DELETE /api/sensors/:id`  
  Deletes a specific sensor by its ID.
