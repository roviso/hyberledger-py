package main

import (
//     "bytes"
    "encoding/json"
    "fmt"
    "time"
//     "strconv"
    "fabric-chaincode-go/shim"
	sc "fabric-protos-go/peer"
)

type SmartContract struct{
}

type Person struct {
    Id string `json:"id"`
    Class string `json:"class"`
    Name string `json:"name"`
    Email string `json:"email"`
}

type Art struct {
    Id string `json:"id"`
    Class string `json:"class"`
    Name string `json:"name"`
    Description string `json:"description"`
    Artist string `json:"artist"`
    Owner string `json:"owner"`
    CreatedAt time.Time `json:"created_at"`
}


// Init function
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
    fmt.Println("chaincode Initializing")
    _, args := stub.GetFunctionAndParameters()
    fmt.Println(args)
    var user Person


    if len(args) != 3 {
		return shim.Error("Incorrect argument numbers. Expecting 4")
	}

    user.Id = args[0]
	user.Class = args[1]
	user.Name = args[2]
	user.Email = args[3]
    UserBytes, _ := json.Marshal(user)
    APIstub.PutState(user.Id , UserBytes)

    return shim.Success(nil)
}

func (s *SmartContract) CreateUser(APIstub shim.ChaincodeStubInterface) sc.Response {
    Id := "user-"+utils.RandStringBytes(32) // utils is a custom package. You can write your own too :P
    var user = Person{Class:"Person", Id: Id, Name: args[0], Email: args[1]}
    UserBytes, _ := json.Marshal(user)
    APIstub.PutState(Id, UserBytes)
    return shim.Success(nil)
}

func (s *SmartContract) queryUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting UserID")
    }
    UserBytes, err := APIstub.GetState(args[0])
    if err != nil {
        return shim.Error(err.Error())
    }
    return shim.Success(UserBytes)
}