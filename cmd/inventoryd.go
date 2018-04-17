// Copyrights 2018 Saurabh Badhwar
// The use of this package is governed by MIT license
// which can be found in the LICENSE file.

package main

import (
	"flag"
	"encoding/json"
	"os"
	"net/http"
	"log"
	inventory "inventory/lib"
)

// Configuration provides a structure for holding the configuration data
// for the inventory service.
type Configuration struct {
	DataStorePath	string
	FlushInterval	uint16
}

var (
	configLookupPath  string
	config			  *Configuration
)

// ConfigurationParser parses the configuration file 
func ConfigurationParser() bool {
	f, ok := os.Open(configLookupPath)
	if ok != nil {
		f.Close()
	} else {
		json.NewDecoder(f).Decode(&config)
		f.Close()
		return true
	}
	return false
}

// flagParser parses the flags from the command line
func flagParser() {
	configPath := flag.String("configFile", "/etc/bolt/inventory.json", "Provide the path where bolt can find its configuration")

	flag.Parse()
	configLookupPath = *configPath
}

func main() {
	flagParser()
	ConfigurationParser()
	api := inventory.APIInit(config.DataStorePath, config.FlushInterval)
	log.SetOutput(os.Stdout)
	log.Fatal(http.ListenAndServe(":8250", api))
}