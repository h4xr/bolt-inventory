// Copyrights 2018 Saurabh Badhwar.
// The use of this package is governed by MIT License
// which can be found in the LICENSE file.

// Package inventory specifies the inventory service for the bolt
// inventory manager service.
package inventory

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

const (
	// HostgroupCapacity defines the max number of hosts that can be
	// created under a single hostgroup.
	HostgroupCapacity = 65000
	// InventoryCapacity specifies the maximum number of hostgroups
	// that can be stored inside a single inventory service.
	InventoryCapacity = 32764
)

// Host defines the structure for storing the data related to
// individual hosts including their names and local variables.
type Host struct {
	// hostname The address through which the host can be reached
	Hostname string
	// fcats The host specific variable
	Facts map[string]string
}

// NewHost defines the initializer for creating a new host
func NewHost(hostname string) *Host {
	return &Host{Hostname: hostname, Facts: make(map[string]string)}
}

// GetHostName returns the hostname of the host
func (h Host) GetHostName() string {
	return h.Hostname
}

// GetHostFacts returns the facts specific to the host
func (h Host) GetHostFacts() map[string]string {
	return h.Facts
}

// SetFact sets a new host fact as defined by the name and value
func (h *Host) SetFact(name string, value string) {
	h.Facts[name] = value
}

// DeleteFact deletes a host local fact from the mapping
func (h *Host) DeleteFact(name string) {
	if _, ok := h.Facts[name]; ok {
		delete(h.Facts, name)
	}
}

// HostGroup defines the structure used for storing the data for
// the hostgroups that are registered individually in the inventory
// service.
type HostGroup struct {
	// name defines the name of the hostgroup through which it can be identified
	Name string
	// hosts defines a slice in which the hosts belonging to a particular hostgroup can be grouped together
	Hosts map[string]*Host
}

// NewHostGroup creates a new hostgroup for the inventory
func NewHostGroup(name string) *HostGroup {
	return &HostGroup{Name: name, Hosts: make(map[string]*Host)}
}

// AddHost adds a new host to the existing hostgroup
func (hg *HostGroup) AddHost(h *Host) {
	hostname := h.GetHostName()
	if _, ok := hg.Hosts[hostname]; !ok {
		hg.Hosts[hostname] = h
	}
}

// DeleteHost removes a host from the Hostgroup
func (hg *HostGroup) DeleteHost(hname string) {
	if _, ok := hg.Hosts[hname]; ok {
		hg.Hosts[hname] = nil
	}
}

// GetHost returns the host object provided the host name
// If the host doesn't exists, nil is returned
func (hg HostGroup) GetHost(hname string) *Host {
	if host, ok := hg.Hosts[hname]; ok {
		return host
	}
	return nil
}

// GetHosts returns all the hosts from a hostgroup
func (hg HostGroup) GetHosts() map[string]*Host {
	return hg.Hosts
}

// GetHostgroupName returns the name of the hostgroup
func (hg HostGroup) GetHostgroupName() string {
	return hg.Name
}

// Inventory struct defines the global service based inventory database
// used to store the information of all the hostgroups and hosts.
// The Inventory struct is used to retrieve all the data that needs to be
// sent back as JSON to the requesting client.
// Since the service can suffer errors at any point of time, this database is
// written to the disk periodically so as to avoid any inconsistency that may
// take place due to unpredicted failure of the code.
type Inventory struct {
	// hostgroups store the created hostgroups along with the data related
	// to the individual hosts inside them.
	Hostgroups map[string]*HostGroup
	// dataStorePath defines the path where the inventory database is created
	// on the disk. Whenever the service starts, it will look for the inventory
	// database at the specified path and try to load the data from it.
	DataStorePath string
	// flushInterval defines the time in milliseconds at which the inventory
	// flush service will write the data to the disk file.
	FlushInterval uint16
	// pendingOps provide the information about how many operations are still
	// pending to be written to the disk. This provides some data into how much
	// data is inventory service storing in its volatile state. This parameter
	// can also be used in future to enhance the inventory data flush service
	// to be more consistent and aggressive in writing the inventory to disk.
	PendingOps uint32

	// A Reader Writer mutex lock to help during the Marshalling of data
	sync.RWMutex

	// inventoryInactive defines a channel which is used to signal the
	// goroutines that we are closing, and they need to exit
	inventoryInactive chan bool
}

// NewInventory creates a new Inventory store to be used by the Inventory
// Service.
func NewInventory(dataStorePath string, flushInterval uint16) *Inventory {
	if !checkDatastorePath(dataStorePath) {
		ok, err := createDatastore(dataStorePath)
		if ok != true {
			log.Fatalf("Unable to create a datastore %s", err)
		}
	}

	inv := Inventory{
		Hostgroups:        make(map[string]*HostGroup),
		DataStorePath:     dataStorePath,
		FlushInterval:     flushInterval,
		PendingOps:        0,
		inventoryInactive: make(chan bool),
	}

	go inv.flushInventoryService()
	return &inv
}

// GetInventory retrieves the inventory from the inventory database
func (inv *Inventory) GetInventory() map[string]*HostGroup {
	return inv.Hostgroups
}
// toJSON converts the current state of the inventory structure to JSON
// representational form which can be written to disk or transmitted back
// to the caller. In case of error, the function returns a nil value.
func (inv *Inventory) toJSON() []byte {
	inv.RLock()
	defer inv.RUnlock()
	invJSON, err := json.Marshal(inv)
	if err != nil {
		log.Fatalf("Unable to encode the data as valid JSON", err)
		return nil
	}
	return invJSON
}

// Save defines a public interface for the inventory structure to write its
// data to the datastore.
func (inv *Inventory) Save() {
	// check if the datastore actually exists or not
	if checkDatastorePath(inv.DataStorePath) == false {
		// we don't have the data store present, try to create one
		_, err := createDatastore(inv.DataStorePath)
		if err != nil {
			log.Fatalf("Unable to create a datastore", err)
		}
	}
	jsonData := inv.toJSON()
	if jsonData == nil {
		log.Fatalf("Unable to convert the data into valid JSON")
	}
	if inv.WriteData(jsonData) == false {
		log.Fatalf("Scheduled inventory save failed. Exiting...")
	}
}

// WriteData writes the binary data to the datastore and
// returns a boolean to indicate if the write was successful
// or not.
func (inv *Inventory) WriteData(data []byte) bool {
	inv.Lock()
	defer inv.Unlock()
	err := ioutil.WriteFile(inv.DataStorePath, data, 0644)
	if err != nil {
		log.Printf("File data write failed", err)
		return false
	}
	return true
}

// NewHostgroup creates a new hostgroup and adds it to the
// inventory. If the hostgroup already exists, the call returns
// without making any changes.
func (inv *Inventory) NewHostgroup(hgname string) {
	if _, ok := inv.Hostgroups[hgname]; !ok {
		hg := NewHostGroup(hgname)
		inv.Hostgroups[hgname] = hg
	}
}

// GetHostgroup retrieves the hostgroup when the name is provided
// if the hostgroup doesn't exists, the call returns a nil
func (inv Inventory) GetHostgroup(hgname string) *HostGroup {
	if hg, ok := inv.Hostgroups[hgname]; ok {
		return hg
	}
	return nil
}

// NewHost creates a new host under the specified hostgroup
// if the hostgroup doesn't exists, then it is created and then
// a new host added to it.
func (inv *Inventory) NewHost(hgname string, hname string) {
	// check if the hostgroup already exists, and create if it doesn't
	inv.NewHostgroup(hgname)
	// retrieve the hostgroup
	hostgroup := inv.GetHostgroup(hgname)
	if hostgroup == nil {
		log.Fatalf("Unable to retireve the hostgroup")
	}
	host := NewHost(hname)
	hostgroup.AddHost(host)
}

// GetHosts returns the list of hosts based in a hostgroup
func (inv Inventory) GetHosts(hgname string) map[string]*Host {
	hostgroup := inv.GetHostgroup(hgname)
	if hostgroup != nil {
		return hostgroup.GetHosts()
	}
	return nil
}

// SetHostFact sets a new fact for the host. If the fact already exists,
// it's value is overwritten
func (inv *Inventory) SetHostFact(hgname string, hname string, fname string, fval string) bool {
	hostgroup := inv.GetHostgroup(hgname)
	if hostgroup != nil {
		host := hostgroup.GetHost(hname)
		if host != nil {
			host.SetFact(fname, fval)
			return true
		}
		return false
	}
	return false
}

func (inv *Inventory) flushInventoryService() {
	log.Printf("Starting the flushInventory service")
	for {
		select {
		case sig := <-inv.inventoryInactive:
			log.Printf("Shutdown signal received. Storing the structures and shutting down")
			if sig == true {
				inv.Save()
				inv.inventoryInactive <- true
				return
			}
		default:
			inv.Save()
			time.Sleep(time.Duration(inv.FlushInterval))
		}
	}
}

// StopInventory signals the inventory service to exit gracefully
func (inv *Inventory) StopInventory() {
	log.Printf("Shutdown request received. Signalling the routines to terminate")
	inv.inventoryInactive <- true
	sig := <-inv.inventoryInactive
	if sig != true {
		log.Fatalf("Error encountered while trying to shutdown the service")
	}
}

// checkDatastorePath validates if a path provided exists on the disk or not
func checkDatastorePath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// createDatastore creates a new datastore if it doesn't exists. In case of err
// the bool parameter returns a false and an error type object which can be
// queried for the error message
func createDatastore(path string) (bool, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}
