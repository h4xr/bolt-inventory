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
#### /create/hostgroup [POST]
Create a new hostgroup.

Parameters to pass in body:

`hostgroup` : The name with which the hostgroup should be created

#### /create/host [POST]
Create a new host in the inventory

Parameters to pass in body:

`hostgroup`: The name of the hostgroup under which the host belongs
`hostname`: The hostname of the host to be added to inventory

#### /create/fact [POST]
Create a new `{host}` local variable whose name is specified by `{key}` and value is specified by `{value}`
If for some reason, the variable is already set due to some previous execution, the value of the variable will be overwritten by the new variable.

Parameters:

`hostname`: The name of the hostgroup to which the host belongs
`hostname`: The hostname for which the facts should be created
`{{ key }}`: `{{ value }}` The key-value pair consisting the fact. These can be multiple in the body

#### /get/inventory [GET]
Retrieve the list of all the hosts under all hostgroups along with their facts

#### /get/hosts
Retrieve all the hosts with their data under the provided hostgroup

Parameters:

`hostgroup`: The name of the hostgroup for which the facts should be retrieved
