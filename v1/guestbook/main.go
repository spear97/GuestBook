/*
Package main implements a simple guestbook application with REST API endpoints.
It provides functionality to manage guestbook entries, including storing and retrieving them.
This application supports two modes: one using Redis as a backend, and another in-memory mode.

The following endpoints are provided:
- GET /lrange/{key}: Retrieves a list by key and returns it as a JSON array.
- GET /rpush/{key}/{value}: Adds a value to a list identified by key.
- GET /info: Returns information about the database backend (Redis or in-memory).
- GET /env: Returns a JSON object with the environment variables of the application.
- GET /hello: Returns a simple greeting message with the hostname of the container.

This application is meant for demonstration purposes and does not include error handling.
*/

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/xyproto/simpleredis/v2"
)

// Declaring global variables for Redis connection pools and in-memory storage
var (
	masterPool *simpleredis.ConnectionPool
	slavePool  *simpleredis.ConnectionPool
	lists      map[string][]string = map[string][]string{} // in-memory storage for lists
)

// Input struct for parsing JSON input
type Input struct {
	InputText string `json:"input_text"`
}

// GetList retrieves a list by key.
// If Redis is configured, it checks the slave pool first, then falls back to the master.
// If no Redis is configured, it retrieves from the in-memory storage.
func GetList(key string) ([]string, error) {
	if slavePool != nil { // Using Redis with slave pool
		list := simpleredis.NewList(slavePool, key)
		if result, err := list.GetAll(); err == nil {
			return result, err
		}
		// If we can't talk to the slave, try using the master instead
	}

	// If the slave doesn't exist or is unavailable, read from the master (if configured)
	if masterPool != nil {
		list := simpleredis.NewList(masterPool, key)
		return list.GetAll()
	}

	// If neither Redis pool exists, use in-memory storage
	return lists[key], nil
}

// AppendToList appends an item to a list identified by key.
// If Redis is configured, it adds the item to the master pool.
// If no Redis is configured, it appends to the in-memory storage.
func AppendToList(item string, key string) ([]string, error) {
	var err error
	items := []string{}

	if masterPool != nil { // Using Redis with master pool
		list := simpleredis.NewList(masterPool, key)
		list.Add(item)
		items, err = list.GetAll()
		if err != nil {
			return nil, err
		}
	} else { // No Redis, using in-memory storage
		items = lists[key]
		items = append(items, item)
		lists[key] = items
	}
	return items, nil
}

// ListRangeHandler handles the GET request to retrieve a list by key.
func ListRangeHandler(rw http.ResponseWriter, req *http.Request) {
	var data []byte

	// Get the list based on the key in the request URL
	items, err := GetList(mux.Vars(req)["key"])
	if err != nil {
		data = []byte("Error getting list: " + err.Error() + "\n")
	} else {
		// Marshal the items into JSON format
		if data, err = json.MarshalIndent(items, "", ""); err != nil {
			data = []byte("Error marshalling list: " + err.Error() + "\n")
		}
	}

	// Write the response
	rw.Write(data)
}

// ListPushHandler handles the GET request to add a value to a list identified by key.
func ListPushHandler(rw http.ResponseWriter, req *http.Request) {
	var data []byte

	// Get the key and value from the request URL
	key := mux.Vars(req)["key"]
	value := mux.Vars(req)["value"]

	// Append the value to the list
	items, err := AppendToList(value, key)

	// Prepare the response data
	if err != nil {
		data = []byte("Error adding to list: " + err.Error() + "\n")
	} else {
		if data, err = json.MarshalIndent(items, "", ""); err != nil {
			data = []byte("Error marshalling list: " + err.Error() + "\n")
		}
	}

	// Write the response
	rw.Write(data)
}

// InfoHandler handles the GET request to get information about the database backend.
func InfoHandler(rw http.ResponseWriter, req *http.Request) {
	info := ""

	// If using Redis, attempt to get information from the master
	if masterPool != nil {
		i, err := masterPool.Get(0).Do("INFO")
		if err != nil {
			info = "Error getting DB info: " + err.Error()
		} else {
			info = string(i.([]byte))
		}
	} else {
		info = "In-memory datastore (not Redis)"
	}

	// Write the response
	rw.Write([]byte(info + "\n"))
}

// EnvHandler handles the GET request to get environment variables of the application.
func EnvHandler(rw http.ResponseWriter, req *http.Request) {
	environment := make(map[string]string)

	// Iterate through the environment variables
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := strings.Join(splits[1:], "=")
		environment[key] = val
	}

	// Marshal the environment variables into JSON format
	data, err := json.MarshalIndent(environment, "", "")
	if err != nil {
		data = []byte("Error marshalling env vars: " + err.Error())
	}

	// Write the response
	rw.Write(data)
}

// HelloHandler handles the GET request to provide a simple greeting message.
func HelloHandler(rw http.ResponseWriter, req *http.Request) {
	// Construct a greeting message including the hostname
	rw.Write([]byte("Hello from guestbook. " +
		"Your app is up! (Hostname: " +
		os.Getenv("HOSTNAME") +
		")\n"))
}

// findRedisURL checks environment variables to determine the Redis URL.
// It supports multiple URL schemes based on the provided environment variables.
func findRedisURL() string {
	host := os.Getenv("REDIS_MASTER_SERVICE_HOST")
	port := os.Getenv("REDIS_MASTER_SERVICE_PORT")
	password := os.Getenv("REDIS_MASTER_SERVICE_PASSWORD")
	masterPort := os.Getenv("REDIS_MASTER_PORT")

	if host != "" && port != "" && password != "" {
		return password + "@" + host + ":" + port
	} else if masterPort != "" {
		return "redis-master:6379"
	}
	return ""
}

// main function initializes the application.
func main() {
	// When using Redis, setup our DB connections
	url := findRedisURL()
	if url != "" {
		masterPool = simpleredis.NewConnectionPoolHost(url)
		defer masterPool.Close() // Close the master pool when main function exits
		slavePool = simpleredis.NewConnectionPoolHost("redis-slave:6379")
		defer slavePool.Close() // Close the slave pool when main function exits
	}

	// Create a new Gorilla mux router
	r := mux.NewRouter()

	// Define endpoints and their corresponding handlers
	r.Path("/lrange/{key}").Methods("GET").HandlerFunc(ListRangeHandler)
	r.Path("/rpush/{key}/{value}").Methods("GET").HandlerFunc(ListPushHandler)
	r.Path("/info").Methods("GET").HandlerFunc(InfoHandler)
	r.Path("/env").Methods("GET").HandlerFunc(EnvHandler)
	r.Path("/hello").Methods("GET").HandlerFunc(HelloHandler)

	// Create a new Negroni middleware with default settings
	n := negroni.Classic()
	n.UseHandler(r) // Set the Gorilla mux router as the handler for Negroni
	n.Run(":3000")  // Start the server on port 3000
}
