package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SimpleAsset struct {
}

func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {

	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("Incorrect arguments.\n Arguments passed is %v\n", args))
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s\n", args[0]))
	}

	return shim.Success(nil)
}

func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, args := stub.GetFunctionAndParameters()
	var result string
	var err error

	if fn == "set" {
		result, err = set(stub, args)
	} else {
		result, err = get(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success([]byte(result))
	}
}

func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("incorrect arguments")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}

	return args[1], nil
}

func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("incorrect arguments")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to set error %v got %v", args[0], err)
	} else if value == nil {
		return "", fmt.Errorf("Asset not found %v", args[0])
	} else {
		return string(value), nil
	}
}

func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
