# Buses Live Location

This project is a live bus tracking system that displays real-time bus locations on a map using Leaflet and WebSockets. The backend is built with Go using the Gin framework.

## Features

- Real-time location updates via WebSockets.
- Display bus locations on a responsive Leaflet map.
- REST endpoints for sending and fetching bus location data.

## Installation

### Prerequisites

- Go
- Docker (optional)

### Running with Go

1. Clone the repository.
2. Download the Go modules:

   ```
   go mod download
   ```

3. Run the application:

   ```
   go run .
   ```

4. Open your browser at [http://localhost:8080](http://localhost:8080).

### Running with Docker

1. Build the Docker image:

   ```
   docker build -t buses-live-location .
   ```

2. Run the container:

   ```
   docker run -p 8080:8080 buses-live-location
   ```

## API Endpoints

- **POST /bus**: Submit a bus location.
  - Example payload:

      ```json
      {
          "latitude": "34.796921",
          "longitude": "48.485994",
          "bus_id": 1
      }
      ```

- **GET /bus**: Retrieve all bus locations.
- **WebSocket Endpoints**:
  - **/clientws**: For clients to receive real-time updates.
  - **/busws**: For buses to send location updates in real-time.

## Testing

API test example files are available under the `api-test` directory:

- `sendLocation.http` for POSTing a bus location.
- `get-buses-location.http` for fetching all bus locations.

## Project Structure

- **main.go**: Application entry point.
- **handlers.go**: Contains API and WebSocket handlers.
- **public/**: Contains frontend files such as HTML, JavaScript, and assets.
- **go.mod/go.sum**: Go module files.
