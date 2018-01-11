package main

import (
	"fmt"

	"github.com/inklabsfoundation/inkchain/core/chaincode/shim"
	pb "github.com/inklabsfoundation/inkchain/protos/peer"
	"encoding/json"
)

// Definitions of a service's status
const (
	S_Created = "created"
	S_Available = "available"
	S_Invalid = "invalid"
)

// Prefixes for user and service separately
const (
	UserPrefix	= "USER_"
	ServicePrefix	= "SER_"
)

// Invoke functions definition
const (
	// User-related basic invoke
	RegisterUser 	= "registerUser"
	RemoveUser 		= "removeUser"
	QueryUser		= "queryUser"

	// Service-related invoke
	RegisterService 	= "registerService"
	InvalidateService 	= "invalidateService"	// mark whether the service is validated
	CreateMashup 		= "createMashup"		// utilize services to create a new mashup
	QueryService		= "queryService"
	EditService			= "editService"
	QueryServiceByUser	= "queryServiceByUser"
	QueryServiceByRange	= "queryServiceByRange"

	// User-related reward invoke
	RewardService = "rewardService"

)

// Chaincode for DSES (Decentralized Service Eco-System)
type serviceChaincode struct {
}

// Structure definition for user
type user struct {
	Name 			string	`json:"name"`
	Introduction	string	`json:"introduction"`
	Address 		string 	`json:"address"`
	// There is a one-to-one correspondence between "Name" and "Address"
	// The Address records the user's profit from creating valuable services or mashups.

	Contribution	int		`json:"contribution"`
	// "Contribution" evaluates the user's contribution to the service ecosystem.

	// Benefit of "Contribution":
	// 1. construct a evaluation for every user's contribution on the service ecosystem
	// 2. inspire users to participate in creating new services and mashups

}

// Structure definition for service
// type "service" defines conventional services as well as mashups.
type service struct {
	Name 			string	`json:"name"`
	Type			string  `json:"type"`
	Developer		string	`json:"developer"`		// record the user that developed this service
	Description		string 	`json:"description"`

	CreatedTime		string	`json:"createdTime"`
	UpdatedTime		string	`json:"updatedTime"`

	// Status records the status of a service:
	// created/available/invalid
	Status			string 	`json:"status"`

	// Whether the service is a mashup or not.
	IsMashup		bool 	`json:"isMashup"`

	// if the service is a mashup, "Composited" records the services that it invokes;
	// if the service is not a mashup, "Composited" records the co-occurrence documents of the service
	Composition		map[string]int	`json:"composition"`

	// Benefit of "Composited":
	// 1. Automatically create service co-occurrence documents and store it into the ledger
	// 2. Promote the security and integrality of service data

	// future: people need to pay if they want to use the record information
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(serviceChaincode))
	if err != nil {
		fmt.Printf("Error starting assetChaincode: %s", err)
	}
}

// Init initializes chaincode
// ==================================================================================
func (t *serviceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("assetChaincode Init.")
	return shim.Success([]byte("Init success."))
}

// Invoke func
// ==================================================================================
func (t *serviceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("assetChaincode Invoke.")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	// ********************************************************
	// PART 1: User-related invokes
	case RegisterUser:
		if len(args) != 2 {
			return shim.Error("Incorrect number of arguments. Expecting 2.")
		}
		// args[0]: user name
		return t.registerUser(stub, args)

	case RemoveUser:
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1.")
		}
		// args[0]: user name
		return t.removeUser(stub, args)

	case QueryUser:
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1.")
		}
		// args[0]: user name
		return t.queryUser(stub, args)

	// ********************************************************
	// PART 2: service-related invokes
	case RegisterService:
		if len(args) != 3 {
			return shim.Error("Incorrect number of arguments. Expecting 3.")
		}
		// args[0]: service name
		// args[1]: service type
		// args[2]: service description
		return t.registerService(stub, args)

	case InvalidateService:
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1.")
		}
		// args[0]: service name
		return t.invalidateService(stub, args)

	case QueryService:
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1.")
		}
		// args[0]: service name
		return t.queryService(stub, args)

	// ********************************************************
	// PART 3: user-related reward invokes
	}

	return shim.Error("Invalid invoke function name.")
}

// Invoke func about user
// ==================================================================================

// ==================================
// registerUser: Register a new user
// ==================================
func (t *serviceChaincode) registerUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var new_name string
	var new_intro string
	var new_add string
	var err error

	new_name = args[0]
	new_intro = args[1]

	// Get the user's address automatically through INKchian's GetSender() interface
	new_add, err = stub.GetSender()
	if err != nil {
		return shim.Error("Fail to get the sender's address.")
	}

	// check if user exists
	user_key := UserPrefix + new_name
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	} else if userAsBytes != nil {
		return shim.Error("This user already exists: " + new_name)
	}

	// register user
	user := &user{new_name, new_intro, new_add, 0}
	userJSONasBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(user_key, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("User register success."))
}

// ===================================
// removeUser: Remove an existed user
// ===================================
func (t *serviceChaincode) removeUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var user_name string
	var err error

	user_name = args[0]

	// check if user exists
	user_key := UserPrefix + user_name
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	} else if userAsBytes == nil {
		return shim.Error("This user does not exist: " + user_name)
	}

	err = stub.DelState(user_key)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("User delete success."))
}

// ===================================
// queryUser: Query an existed user
// ===================================
func (t *serviceChaincode) queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var user_name string
	var err error

	user_name = args[0]

	// check if user exists
	user_key := UserPrefix + user_name
	userAsBytes, err := stub.GetState(user_key)
	if err != nil {
		return shim.Error("Fail to get user: " + err.Error())
	} else if userAsBytes == nil {
		return shim.Error("This user does not exist: " + user_name)
	}

	// return user info
	return shim.Success(userAsBytes)
}

// Invoke func about service
// ==================================================================================

// =======================================
// registerService: Register a new service
// =======================================
func (t *serviceChaincode) registerService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var service_name string
	var service_type string
	var service_des  string
	var service_dev  string
	var err error

	service_name = args[0]
	service_type = args[1]
	service_des = args[2]

	// get service developer
	service_dev, err = stub.GetSender()
	if err != nil {
		return shim.Error("Fail to get the sender's address.")
	}

	// check if service exists
	service_key := ServicePrefix + service_name
	serviceAsBytes, err := stub.GetState(service_key)
	if err != nil {
		return shim.Error("Fail to get service: " + err.Error())
	} else if serviceAsBytes != nil {
		return shim.Error("This service already exists: " + service_name)
	}

	// register service
	newS := &service{service_name, service_type, service_dev,
		service_des, "", "", S_Created,
		false, make(map[string]int)}
	serviceJSONasBytes, err := json.Marshal(newS)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(service_key, serviceJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("Service register success."))
}

// =================================================
// invalidateService: Invalidate an existed service
// =================================================
func (t *serviceChaincode) invalidateService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var service_name string
	var err error

	service_name = args[0]

	// STEP 0: check if service exists
	service_key := ServicePrefix + service_name
	serviceAsBytes, err := stub.GetState(service_key)
	if err != nil {
		return shim.Error("Fail to get service: " + err.Error())
	} else if serviceAsBytes == nil {
		return shim.Error("This service does not exists: " + service_name)
	}

	// STEP 1: check whether it is the service's developer's invocation
	var senderAdd string
	senderAdd, err = stub.GetSender()
	if err != nil {
		return shim.Error("Fail to get the sender's address.")
	}

	var serviceJSON service
	err = json.Unmarshal([]byte(serviceAsBytes), &serviceJSON)
	if err != nil {
		return shim.Error("Error unmarshal service bytes.")
	}

	if senderAdd != serviceJSON.Developer {
		return shim.Error("Aurthority err! Not invoke by the service's developer.")
	}

	// STEP 2: invalidate the service and store it.
	// new service, make it invalidated
	new_service := &service{serviceJSON.Name, serviceJSON.Type, serviceJSON.Developer,
							serviceJSON.Description, serviceJSON.CreatedTime, serviceJSON.UpdatedTime,
							S_Invalid, serviceJSON.IsMashup, serviceJSON.Composition}
	// store the new service
	assetJSONasBytes, err := json.Marshal(new_service)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(service_key, assetJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("Invalidate Service success."))
}

// ======================================
// queryService: Query an existed service
// ======================================
func (t *serviceChaincode) queryService(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var service_name string
	var err error

	service_name = args[0]

	// check if service exists
	service_key := ServicePrefix + service_name
	serviceAsBytes, err := stub.GetState(service_key)
	if err != nil {
		return shim.Error("Fail to get service: " + err.Error())
	} else if serviceAsBytes == nil {
		return shim.Error("This service does not exist: " + service_name)
	}

	// return service info
	return shim.Success(serviceAsBytes)
}
