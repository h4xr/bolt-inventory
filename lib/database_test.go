// Copyrights 2018 Saurabh Badhwar.
// The use of this package is governed by MIT License
// which can be found in the LICENSE file.
package inventory

import "testing"

func TestGetHostName(t *testing.T) {
	hostname := "m1.example.com"
	host := NewHost(hostname)
	if host.GetHostName() != hostname {
		t.Errorf("New Host was created with a wrong hostname, got %s, required %s", host.GetHostName(), hostname)
	}
}

func TestGetFact(t *testing.T) {
	hostname := "m1.example.com"
	host := NewHost(hostname)
	host.SetFact("test", "value")
	if _, ok := host.Facts["test"]; !ok {
		t.Errorf("Unable to set the host fact properly")
	}
}

func TestGetHostFacts(t *testing.T) {
	hostname := "m1.example.com"
	host := NewHost(hostname)
	host.SetFact("test1", "value1")
	host.SetFact("test2", "value2")
	host.SetFact("test3", "value3")
	facts := host.GetHostFacts()
	if len(facts) != 3 {
		t.Errorf("Unable to retrieve the host facts")
	}
}

func TestDeleteFact(t *testing.T) {
	hostname := "m1.example.com"
	host := NewHost(hostname)
	host.SetFact("test1", "value1")
	host.SetFact("test2", "value2")
	host.DeleteFact("test1")
	if _, ok := host.Facts["test1"]; ok {
		t.Errorf("Unable to delete the host fact")
	}
}

func TestNewHostGroup(t *testing.T) {
	hostgroupName := "TestGroup"
	hostgroup := NewHostGroup(hostgroupName)
	if hostgroup.Name != hostgroupName {
		t.Errorf("Unable to create a hostgroup")
	}
}

func TestAddHost(t *testing.T) {
	hostgroupName := "TestGroup"
	hostname := "m1.example.com"
	hostgroup := NewHostGroup(hostgroupName)
	host := NewHost(hostname)
	hostgroup.AddHost(host)
	if _, ok := hostgroup.Hosts[hostname]; !ok {
		t.Errorf("Unable to add a new host to the hostgroup")
	}
}

func TestDeleteHost(t *testing.T) {
	hostgroupName := "TestGroup"
	hostname1 := "m1.example.com"
	hostname2 := "m2.example.com"
	hostgroup := NewHostGroup(hostgroupName)
	host1 := NewHost(hostname1)
	host2 := NewHost(hostname2)
	hostgroup.AddHost(host1)
	hostgroup.AddHost(host2)
	hostgroup.DeleteHost(hostname1)
	if val, ok := hostgroup.Hosts[hostname1]; ok {
		if val != nil {
			t.Errorf("Unable to delete the host from hostgroup %s", val)
		}
	} else {
		t.Errorf("Unexpected error while trying to remove host from hostgroup")
	}
}

func TestGetHost(t *testing.T) {
	hostgroupName := "TestGroup"
	hostname := "m1.example.com"
	hostgroup := NewHostGroup(hostgroupName)
	host := NewHost(hostname)
	hostgroup.AddHost(host)
	h := hostgroup.GetHost(hostname)
	if h != host {
		t.Errorf("Unable to retireve the host from the hostgroup")
	}
}

func TestGetHosts(t *testing.T) {
	hostgroupName := "TestGroup"
	hostname := "m1.example.com"
	hostgroup := NewHostGroup(hostgroupName)
	host := NewHost(hostname)
	hostgroup.AddHost(host)
	hosts := hostgroup.GetHosts()
	if _, ok := hosts[hostname]; !ok {
		t.Errorf("Unable to retrieve the hosts from the hostgroup correctly")
	}
}

func TestGetHostgroupName(t *testing.T) {
	hostgroupName := "TestGroup"
	hostgroup := NewHostGroup(hostgroupName)
	if hostgroupName != hostgroup.GetHostgroupName() {
		t.Errorf("Unale to retrieve the correct hostgroup name")
	}
}

func TestNewInventory(t *testing.T) {
	dataStorePath := "/home/sbadhwar/upstream/data.db"
	var flushInterval uint16 = 5000
	inventory := NewInventory(dataStorePath, flushInterval)
	inventory.StopInventory()
	if inventory == nil {
		t.Errorf("Unable to construct a new inventory")
	}
}

func TestNewHostgroup(t *testing.T) {
	dataStorePath := "/home/sbadhwar/upstream/data.db"
	var flushInterval uint16 = 5000
	inventory := NewInventory(dataStorePath, flushInterval)
	hostgroupName := "TestGroup"
	inventory.NewHostgroup(hostgroupName)
	inventory.StopInventory()
	if _, ok := inventory.Hostgroups[hostgroupName]; !ok {
		t.Errorf("Unable to create a new hostgroup")
	}
}

func TestGetHostgroup(t *testing.T) {
	dataStorePath := "/home/sbadhwar/upstream/data.db"
	var flushInterval uint16 = 5000
	inventory := NewInventory(dataStorePath, flushInterval)
	hostgroupName := "TestGroup"
	inventory.NewHostgroup(hostgroupName)
	inventory.StopInventory()
	if inventory.GetHostgroup(hostgroupName) == nil {
		t.Errorf("Unable to create a new hostgroup")
	}	
}

func TestNewHost(t *testing.T) {
	dataStorePath := "/home/sbadhwar/upstream/data.db"
	var flushInterval uint16 = 5000
	inventory := NewInventory(dataStorePath, flushInterval)
	hostgroupName := "TestGroup"
	hostname := "m1.example.com"
	inventory.NewHostgroup(hostgroupName)
	inventory.NewHost(hostgroupName, hostname)
	hosts := inventory.GetHosts(hostgroupName)
	inventory.StopInventory()
	flag := false
	for key := range hosts {
		if key == hostname {
			flag = true
		}
	}
	if flag == false {
		t.Errorf("Unable to add a new host to the hostgroup")
	}
}

func TestSetHostFact(t *testing.T) {
	dataStorePath := "/home/sbadhwar/upstream/data.db"
	var flushInterval uint16 = 5000
	inventory := NewInventory(dataStorePath, flushInterval)
	hostgroupName := "TestGroup"
	hostname := "m1.example.com"
	inventory.NewHostgroup(hostgroupName)
	inventory.NewHost(hostgroupName, hostname)
	inventory.SetHostFact(hostgroupName, hostname, "testfact", "testval")
	hosts := inventory.GetHosts(hostgroupName)
	inventory.StopInventory()
	flag := false
	for key := range hosts {
		if key == hostname {
			if _, ok := hosts[key].Facts["testfact"]; ok {
				flag = true
			}
		}
	}
	if flag == false {
		t.Errorf("Unable to add a new host to the hostgroup")
	}
}