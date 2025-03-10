package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

// SmartContract provides functions for managing DIDs
type SmartContract struct {
	contractapi.Contract
}

// DID represents a Decentralized Identifier
type DID struct {
	PhoneNumber string `json:"phoneNumber"`
	DID         string `json:"did"`
}

// Nominee represents nominee details record
type Nominee struct {
	NomineeID              string `json:"nominee_id"`
	UserDID                string `json:"user_did"`
	NomineeDID             string `json:"nominee_did"`
	FinancialInstitutionID string `json:"financial_institution_id"`
	NomineeType            string `json:"nominee_type"`
	PercentageShare        int    `json:"percentage_share"`
	Status                 string `json:"status"`
	CreationTime           string `json:"creation_time"`
	ExpiryTime             string `json:"expiry_time"`
}

// Consent represents a minimal consent record
type Consent struct {
	ConsentID            string   `json:"consent_id"`
	UserDID              string   `json:"user_did"`
	TPAID                string   `json:"tpa_id"`
	FinancialInstitution string   `json:"financial_institution_id"`
	RequestedData        []string `json:"requested_data"`
	ExpiryTime           string   `json:"expiry_time"`
	Status               string   `json:"status"` // Active, Expired, Revoked
}

// StoreDID saves a DID on the ledger
func (s *SmartContract) StoreDID(ctx contractapi.TransactionContextInterface, phoneNumber string, did string) error {
	exists, err := s.DIDExists(ctx, phoneNumber)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("DID for phone number %s already exists", phoneNumber)
	}

	didRecord := DID{
		PhoneNumber: phoneNumber,
		DID:         did,
	}
	didBytes, err := json.Marshal(didRecord)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(phoneNumber, didBytes)
}

// GetDID retrieves a DID from the ledger
func (s *SmartContract) GetDID(ctx contractapi.TransactionContextInterface, phoneNumber string) (*DID, error) {
	didBytes, err := ctx.GetStub().GetState(phoneNumber)
	if err != nil {
		return nil, err
	}
	if didBytes == nil {
		return nil, fmt.Errorf("DID not found for phone number %s", phoneNumber)
	}

	var did DID
	err = json.Unmarshal(didBytes, &did)
	if err != nil {
		return nil, err
	}

	return &did, nil
}

// DIDExists checks if a DID exists on the ledger
func (s *SmartContract) DIDExists(ctx contractapi.TransactionContextInterface, phoneNumber string) (bool, error) {
	didBytes, err := ctx.GetStub().GetState(phoneNumber)
	if err != nil {
		return false, err
	}
	return didBytes != nil, nil
}

// CreateConsent stores a new consent record on the ledger
func (s *SmartContract) CreateConsent(ctx contractapi.TransactionContextInterface, consentID, userDID, tpaID, financialInstitution string, requestedData []string, expiryTime string) error {

	// Check if the consent already exists
	existingConsent, _ := ctx.GetStub().GetState(consentID)
	if existingConsent != nil {
		return fmt.Errorf("consent with ID %s already exists", consentID)
	}

	// Create a new consent record
	consent := Consent{
		ConsentID:            consentID,
		UserDID:              userDID,
		TPAID:                tpaID,
		FinancialInstitution: financialInstitution,
		RequestedData:        requestedData,
		ExpiryTime:           expiryTime,
		Status:               "Active",
	}

	// Convert to JSON
	consentJSON, err := json.Marshal(consent)
	if err != nil {
		return err
	}

	// Store in blockchain
	return ctx.GetStub().PutState(consentID, consentJSON)
}

// GetConsent retrieves a consent record by its ID
func (s *SmartContract) GetConsent(ctx contractapi.TransactionContextInterface, consentID string) (*Consent, error) {

	// Fetch the consent data
	consentJSON, err := ctx.GetStub().GetState(consentID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from ledger: %v", err)
	}
	if consentJSON == nil {
		return nil, fmt.Errorf("consent with ID %s does not exist", consentID)
	}

	// Convert JSON to struct
	var consent Consent
	err = json.Unmarshal(consentJSON, &consent)
	if err != nil {
		return nil, err
	}

	return &consent, nil
}

// RevokeConsent updates the consent status to "Revoked"
func (s *SmartContract) RevokeConsent(ctx contractapi.TransactionContextInterface, consentID string) error {

	// Get existing consent
	consent, err := s.GetConsent(ctx, consentID)
	if err != nil {
		return err
	}

	// Update status
	consent.Status = "Revoked"

	// Convert back to JSON and update ledger
	updatedConsentJSON, err := json.Marshal(consent)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(consentID, updatedConsentJSON)
}

// CheckConsentValidity verifies if the consent is still valid
func (s *SmartContract) CheckConsentValidity(ctx contractapi.TransactionContextInterface, consentID string) (bool, error) {

	// Get existing consent
	consent, err := s.GetConsent(ctx, consentID)
	if err != nil {
		return false, err
	}

	// Check if consent is active
	if consent.Status != "Active" {
		return false, fmt.Errorf("consent is %s", consent.Status)
	}

	// Check expiry time
	currentTime := time.Now()
	expiryTime, err := time.Parse(time.RFC3339, consent.ExpiryTime)
	if err != nil {
		return false, fmt.Errorf("invalid expiry time format")
	}

	if currentTime.After(expiryTime) {
		return false, fmt.Errorf("consent has expired")
	}

	return true, nil
}

// CreateNominee stores a new nominee in the ledger
func (s *SmartContract) CreateNominee(ctx contractapi.TransactionContextInterface, nomineeID, userDID, nomineeDID, financialInstitutionID, nomineeType string, percentageShare int, expiryTime string) error {
	// Ensure nominee does not already exist
	existingNominee, err := ctx.GetStub().GetState(nomineeID)
	if err != nil {
		return fmt.Errorf("failed to check nominee existence: %v", err)
	}
	if existingNominee != nil {
		return fmt.Errorf("nominee already exists")
	}

	// Set current time as creation time
	creationTime := time.Now().UTC().Format(time.RFC3339)

	// Create nominee object
	nominee := Nominee{
		NomineeID:              nomineeID,
		UserDID:                userDID,
		NomineeDID:             nomineeDID,
		FinancialInstitutionID: financialInstitutionID,
		NomineeType:            nomineeType,
		PercentageShare:        percentageShare,
		Status:                 "Active",
		CreationTime:           creationTime,
		ExpiryTime:             expiryTime,
	}

	// Convert to JSON
	nomineeJSON, err := json.Marshal(nominee)
	if err != nil {
		return fmt.Errorf("failed to marshal nominee JSON: %v", err)
	}

	// Save nominee to blockchain
	return ctx.GetStub().PutState(nomineeID, nomineeJSON)
}

// GetNominee retrieves a nominee from the ledger
func (s *SmartContract) GetNominee(ctx contractapi.TransactionContextInterface, nomineeID string) (*Nominee, error) {
	nomineeJSON, err := ctx.GetStub().GetState(nomineeID)
	if err != nil {
		return nil, fmt.Errorf("failed to read nominee: %v", err)
	}
	if nomineeJSON == nil {
		return nil, fmt.Errorf("nominee not found")
	}

	// Convert JSON to Nominee struct
	var nominee Nominee
	err = json.Unmarshal(nomineeJSON, &nominee)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal nominee JSON: %v", err)
	}

	return &nominee, nil
}

// RevokeNominee updates nominee status to "Revoked"
func (s *SmartContract) RevokeNominee(ctx contractapi.TransactionContextInterface, nomineeID string) error {
	// Get nominee from blockchain
	nominee, err := s.GetNominee(ctx, nomineeID)
	if err != nil {
		return err
	}

	// Update status to "Revoked"
	nominee.Status = "Revoked"

	// Convert to JSON
	updatedNomineeJSON, err := json.Marshal(nominee)
	if err != nil {
		return fmt.Errorf("failed to marshal updated nominee JSON: %v", err)
	}

	// Save updated nominee to blockchain
	return ctx.GetStub().PutState(nomineeID, updatedNomineeJSON)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
