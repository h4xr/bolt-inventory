# Bolt Inventory Service
## Introduction
---
Bolt Inventory Service Provides APIs for maintaining an Inventory record, which can be consumed by the different requests through the use of REST API.

The following features are currently supported by the Inventory service:

* Creation of Hostgroups
* Creation of Hosts
* Setting up of host local variables
* Setting up of hostgroup local variables
* Retrieval of the inventory based on the following parameters
    * All hostgroups
    * Filtered by hostgroup
* Deletion of Hostgroups
* Deletion of Hosts

## Public REST APIs
---
#### /inventory/create/{hostgroup}
Create a new hostgroup with the name as specified in `{hostgroup}`

#### /inventory/create/{hostgroup}/{host}
Create a new `{host}` entry under the mentioned `{hostgroup}`

#### /inventory/create/{hostgroup}/{host}/{key}/{value}
Create a new `{host}` local variable whose name is specified by `{key}` and value is specified by `{value}`
If for some reason, the variable is already set due to some previous execution, the value of the variable will be overwritten by the new variable.

#### /inventory/list
Retrieve the list of all the hosts under all hostgroups

#### /inventory/list/{hostgroup}
Retrieve all the hosts with their data under the provided `{hostgroup}`

#### /inventory/delete/{hostgroup}
Delete a `{hostgroup}` as well as all the hosts under it

#### /inventory/delete/{hostgroup}/{host}
Delete a `{host}` under the provided `{hostgroup}`   