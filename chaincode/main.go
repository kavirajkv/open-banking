package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SimpleChaincode struct {
	contractapi.Contract
}

// InitLedger - Initialize some data on the ledger
func (s *SimpleChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Create an initial nominee
	nominee := Nominee{ID: "1", Name: "John Doe", Age: 30, Status: "Active"}
	// Save to ledger (marshalling the nominee into bytes)
	nomineeAsBytes, err := json.Marshal(nominee)
	if err != nil {
		return fmt.Errorf("failed to marshal nominee: %v", err)
	}
	err = ctx.GetStub().PutState(nominee.ID, nomineeAsBytes)
	if err != nil {
		return fmt.Errorf("failed to create initial nominee: %v", err)
	}
	return nil
}

// AddNominee - Adds a nominee to the ledger
func (s *SimpleChaincode) AddNominee(ctx contractapi.TransactionContextInterface, id string, name string, age int, status string) error {
	// Create a nominee struct
	nominee := Nominee{ID: id, Name: name, Age: age, Status: status}
	// Marshal the nominee into bytes
	nomineeAsBytes, err := json.Marshal(nominee)
	if err != nil {
		return fmt.Errorf("failed to marshal nominee: %v", err)
	}
	// Save nominee to the ledger
	err = ctx.GetStub().PutState(nominee.ID, nomineeAsBytes)
	if err != nil {
		return fmt.Errorf("failed to add nominee: %v", err)
	}
	return nil
}

// QueryNominee - Queries the nominee information from the ledger
func (s *SimpleChaincode) QueryNominee(ctx contractapi.TransactionContextInterface, id string) (Nominee, error) {
	// Retrieve nominee from the ledger
	nomineeAsBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return Nominee{}, fmt.Errorf("failed to read from world state: %v", err)
	}
	if nomineeAsBytes == nil {
		return Nominee{}, fmt.Errorf("nominee with ID %s does not exist", id)
	}

	var nominee Nominee
	// Unmarshal nominee data
	err = json.Unmarshal(nomineeAsBytes, &nominee)
	if err != nil {
		return Nominee{}, fmt.Errorf("failed to unmarshal nominee: %v", err)
	}
	return nominee, nil
}

// Nominee struct
type Nominee struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Status string `json:"status"`
}

func main() {
	// Initialize the chaincode
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err.Error())
		return
	}

	// Start the chaincode
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
