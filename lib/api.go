// Coyrights 2018 Saurabh Badhwar
// The use of this package is governed by MIT License
// which can be found in the LICENSE file.

package inventory

import (
	"github.com/gorilla/mux"
)

// APIInit initializes the API service using the mux router
// engine and maps the endpoints to the required call handlers
func APIInit(dataStorePath string, flushInterval uint16) *mux.Router {
	// Setup the inventory before we can use the router
	setupInventory(dataStorePath, flushInterval)
	// We are good to go with a new router
	router := mux.NewRouter()
	// Register the handlers here
	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/create/hostgroup", createHostgroup).Methods("POST")
	router.HandleFunc("/create/host", createHost).Methods("POST")
	router.HandleFunc("/create/fact", setHostFact).Methods("POST")
	router.HandleFunc("/get/inventory", getInventory).Methods("GET")
	router.HandleFunc("/get/hosts/{hostgroup}", getHosts).Methods("GET")
	return router
}

// setupInventory initializes the inventory variable which is then
// used by the API to actually run the inventory service
func setupInventory(dataStorePath string, flushInterval uint16) {
	inv = NewInventory(dataStorePath, flushInterval)
}

