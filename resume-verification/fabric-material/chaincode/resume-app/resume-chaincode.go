// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
 */

package main

/* Imports  
* 4 utility libraries for handling bytes, reading and writing JSON, 
formatting, and string manipulation  
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts  
*/ 
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Resume structure, with 4 properties.  
Structure tags are used by encoding/json library
*/
type Resume struct {
	User_Id string `json:"user_id"`
	Timestamp string `json:"timestamp"`
	Resume  string `json:"resume"`
	Resume_Hash  string `json:"resume_hash"`
}

/*
 * The Init method *
 called when the Smart Contract "resume-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function 
 -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "resume-chaincode"
 The app also specifies the specific smart contract function to call with args
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryResume" {
		return s.queryResume(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordResume" {
		return s.recordResume(APIstub, args)
	} else if function == "queryAllResumes" {
		return s.queryAllResumes(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryResume method *
Used to view the records of one particular resume
It takes one argument -- the key for the resume in question
 */
func (s *SmartContract) queryResume(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	resumeAsBytes, _ := APIstub.GetState(args[0])
	if resumeAsBytes == nil {
		return shim.Error("Could not locate resume")
	}
	return shim.Success(resumeAsBytes)
}

/*
 * The initLedger method *
Will add test data (10 resume catches)to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	resume := []Resume{
		Resume{User_Id: "CUCK", Resume: "67.0006, -70.5476", Timestamp: "1504054225", Resume_Hash: "Miriam"},
		Resume{User_Id: "M83T", Resume: "91.2395, -49.4594", Timestamp: "1504057825", Resume_Hash: "Dave"},
		Resume{User_Id: "T012", Resume: "58.0148, 59.01391", Timestamp: "1493517025", Resume_Hash: "Igor"},
	}

	i := 0
	for i < len(resume) {
		fmt.Println("i is ", i)
		resumeAsBytes, _ := json.Marshal(resume[i])
		APIstub.PutState(strconv.Itoa(i+1), resumeAsBytes)
		fmt.Println("Added", resume[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The recordResume method *
Fisherman like Sarah would use to record each of her resume catches. 
This method takes in five arguments (attributes to be saved in the ledger). 
 */
func (s *SmartContract) recordResume(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var resume = Resume{ User_Id: args[1], Resume: args[2], Timestamp: args[3], Resume_Hash: args[4] }

	resumeAsBytes, _ := json.Marshal(resume)
	err := APIstub.PutState(args[0], resumeAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record resume catch: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllResumes method *
allows for assessing all the records added to the ledger(all resume catches)
This method does not take any arguments. Returns JSON string containing results. 
 */
func (s *SmartContract) queryAllResumes(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllResumes:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}