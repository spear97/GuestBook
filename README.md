# Guestbook Application

## Overview

This Go application implements a simple guestbook with REST API endpoints. It allows users to manage guestbook entries, including storing and retrieving them. The application supports two modes: one using Redis as a backend and another using in-memory storage.

### Endpoints
- `GET /lrange/{key}`: Retrieves a list by key and returns it as a JSON array.
- `GET /rpush/{key}/{value}`: Adds a value to a list identified by key.
- `GET /info`: Returns information about the database backend (Redis or in-memory).
- `GET /env`: Returns a JSON object with the environment variables of the application.
- `ET /hello`: Returns a simple greeting message with the hostname of the container.

### Functionality
- `GetList`: Retrieves a list by key. If Redis is configured, it checks the slave pool first and then falls back to the master. If no Redis is configured, it retrieves from the in-memory storage.
- `AppendToList`: Appends an item to a list identified by key. If Redis is configured, it adds the item to the master pool. If no Redis is configured, it appends to the in-memory storage.

## How to Use

### Running the Application
1. Ensure you have `Go` installed on your machine.
2. Clone this repository.
3. Set up your Redis environment variables if using Redis as the backend.
4. Run the following commands:
   ```shell
     go build
    ./guestbook
   ```

### API Endpoints
- `GET /lrange/{key}`: Retrieve a list by key.
- `GET /rpush/{key}/{value}`: Add a value to a list identified by key.
- `GET /info`: Get information about the database backend.
- `GET /env`: Get environment variables of the application.
- `GET /hello`: Get a simple greeting message.

### Dependencies
- [github.com/codegangsta/negroni](https://github.com/urfave/negroni): Negroni middleware for HTTP requests.
- [github.com/gorilla/mux](https://github.com/gorilla/mux): Gorilla mux router for URL routing.
- [github.com/xyproto/simpleredis/v2](https://github.com/xyproto/simpleredis/tree/main/v2): Simpleredis library for Redis connection pooling.

### Environment Variables
- `REDIS_MASTER_SERVICE_HOST`: Redis master service host.
- `REDIS_MASTER_SERVICE_PORT`: Redis master service port.
- `REDIS_MASTER_SERVICE_PASSWORD`: Redis master service password.
- `REDIS_MASTER_PORT`: Redis master port.
