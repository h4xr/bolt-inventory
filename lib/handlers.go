// Copyrights 2018 Saurabh Badhwar
// The use of this package is governed by MIT License
// which can be found in the LICENSE file.

package inventory

import (
	//"os"
	"encoding/json"
	"fmt"
	"net/http"
	//"github.com/gorilla/mux"
)

var (
	inv *Inventory
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "Pong")
}

func createHostgroup(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	json.NewDecoder(r.Body).Decode(&params)
	if hgname, ok := params["hostgroup"]; ok {
		inv.NewHostgroup(hgname)
		w.WriteHeader(http.StatusCreated)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	return
}

func createHost(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	json.NewDecoder(r.Body).Decode(&params)
	if hgname, ok := params["hostgroup"]; ok {
		if hname, ok := params["hostname"]; ok {
			inv.NewHost(hgname, hname)
			w.WriteHeader(http.StatusCreated)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func setHostFact(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	json.NewDecoder(r.Body).Decode(&params)
	hostgroup, hgok := params["hostgroup"]
	hostname, hok := params["hostname"]
	if !hgok || !hok {
		w.WriteHeader(http.StatusBadRequest)
	}
	delete(params, "hostgroup")
	delete(params, "hostname")
	for f, v := range params {
		ok := inv.SetHostFact(hostgroup, hostname, f, v)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to process the request. Please try again later"))
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func getInventory(w http.ResponseWriter, r *http.Request) {
	outputInvMap := make(map[string]interface{})
	hostInventory := inv.GetInventory()
	outputInvMap["_meta"] = make(map[string]interface{})
	outputInvMap["_meta"].(map[string]interface{})["hostvars"] = make(map[string]interface{})
	for hgname := range hostInventory {
		outputInvMap[hgname] = make(map[string][]string)
		outputInvMap[hgname].(map[string][]string)["hosts"] = make([]string, 0, 65000)
		hosts := hostInventory[hgname].GetHosts()
		for hostname := range hosts {
			// We dynamically create a inventory as per ansible wants it to be
			// this involves explicitly typecasting an interface value to map
			// value and then allocating memory to that.
			outputInvMap[hgname].(map[string][]string)["hosts"] = append(outputInvMap[hgname].(map[string][]string)["hosts"], hostname)
			// Setup host facts in the inventory
			outputInvMap["_meta"].(map[string]interface{})["hostvars"].(map[string]interface{})[hostname] = make(map[string]string)
			// Assign the host facts to the inventory
			outputInvMap["_meta"].(map[string]interface{})["hostvars"].(map[string]interface{})[hostname] = hosts[hostname].GetHostFacts()
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(outputInvMap)
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Pending implementation"))
}
