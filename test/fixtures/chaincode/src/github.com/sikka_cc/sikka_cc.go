package main

import (
    "fmt"
    "log"
    "encoding/json"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
    contractapi.Contract
}


// Asset describes basic details of what makes up a simple asset
type Vendor struct {
    ID              string `json:"ID"`
    FirstName       string `json:"FirstName"`
    LastName        string `json:"LastName"`
    Phone_number    string  `json:"Phone_number"`
    SikkaOwned       int    `json:"SikkaOwned"`
    District        string `json:"District"`
    Local_level     string  `json:"Local_level"`
    Ward            string `json:"Ward"`
    Tole            string `json:"Tole"`
}


type Beneficiary struct {
    ID              string `json:"ID"`
    FirstName       string `json:"FirstName"`
    LastName        string `json:"LastName"`
    Phone_number    string  `json:"Phone_number"`
    SikkaOwned       int    `json:"SikkaOwned"`
    District        string `json:"District"`
    Local_level     string  `json:"Local_level"`
    Ward            string `json:"Ward"`
    Tole            string `json:"Tole"`
}


type Asset struct {
    ID             string `json:"ID"`
    Color          string `json:"color"`
    Size           int    `json:"size"`
    Owner          string `json:"owner"`
    AppraisedValue int    `json:"appraisedValue"`
}


// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    beneficiaries := []Beneficiary{
        {ID: "NP202106B001BKTBM", FirstName: "Ravi", LastName:"Prajapati",Phone_number:"9860073000",SikkaOwned: 500, District: "BKT", Local_level: "BM", Ward:"06",Tole:"Gun" },
        {ID: "NP202103B002KTMKM", FirstName: "Ankit", LastName:"Shakya",Phone_number:"9843219933",SikkaOwned: 700, District: "KTM", Local_level: "KM", Ward:"03",Tole:"Sanka" },
    }

    vendors := []Vendor{
        {ID: "NP202109V001BKTBM", FirstName: "Ram", LastName:"Shrestha",Phone_number:"9860321523",SikkaOwned: 5, District: "BKT", Local_level: "BM", Ward:"09",Tole:"Don" },
        {ID: "NP202103V002KTMKM", FirstName: "Anjila", LastName:"Gurung",Phone_number:"9841232152",SikkaOwned: 33, District: "KTM", Local_level: "KM", Ward:"03",Tole:"NayaBajar" },
    }

    for _, beneficiary := range beneficiaries {
        beneficiaryJSON, err := json.Marshal(beneficiary)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(beneficiary.ID, beneficiaryJSON)
        if err != nil {
            return fmt.Errorf("failed to put to world state. %v", err)
        }
    }

    for _, vendor := range vendors {
        vendorJSON, err := json.Marshal(vendor)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(vendor.ID, vendorJSON)
        if err != nil {
            return fmt.Errorf("failed to put to world state. %v", err)
        }
    }

    return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("the asset %s already exists color: %s, size: %s , owner: %s, appraisedvalue: %s",id, color, size, owner, appraisedValue)
    }

    asset := Asset{
        ID:             id,
        Color:          color,
        Size:           size,
        Owner:          owner,
        AppraisedValue: appraisedValue,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    err =  ctx.GetStub().PutState(id, assetJSON)
    if err != nil {
            return fmt.Errorf("failed to put to world state. %v", err)
        }

    return assetJSON
}


// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
        return nil, fmt.Errorf("the asset %s does not exist", id)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
        return nil, err
    }

    return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", id)
    }

    // overwriting original asset with new asset
    asset := Asset{
        ID:             id,
        Color:          color,
        Size:           size,
        Owner:          owner,
        AppraisedValue: appraisedValue,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", id)
    }
    fmt.Println("deleting asset : %s", id)

    return ctx.GetStub().DelState(id)
}


// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return assetJSON != nil, nil
}


// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
    asset, err := s.ReadAsset(ctx, id)
    if err != nil {
        return err
    }

    asset.Owner = newOwner
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllBeneficiary(ctx contractapi.TransactionContextInterface) ([]*Beneficiary, error) {
// range query with empty string for startKey and endKey does an
// open-ended query of all assets in the chaincode namespace.
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var beneficiaries []*Beneficiary
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var beneficiary Beneficiary
        err = json.Unmarshal(queryResponse.Value, &beneficiary)
        if err != nil {
            return nil, err
        }
        beneficiaries = append(beneficiaries, &beneficiary)
    }

    return beneficiaries, nil
}


func main() {
    assetChaincode, err := contractapi.NewChaincode(&SmartContract{})

    if err != nil {
        log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
    }

    if err := assetChaincode.Start(); err != nil {
        log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
    }
}