package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"

)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

type Record struct {
	UserId string `json:"userId"`
	CreatedDate string `json:"createdDate"`
	Survey Survey `json:"survey"`
}
type Survey struct {
	SurveyId string `json:"surveyId"`
	SurveyDescription string `json:"surveyDescription"`
	Question []Question `json:"questions"`
}
type Question struct {
	QuestionId string `json:"questionId"`
	QuestionDescription string `json:"questionDescription"`
	Answer Answer `json:"answer"`
}
type Answer struct {
	AnswerId string `json:"answerId"`
	AnswerDescription string `json:"answerDescription"`

}

// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("fabcar_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "queryAllProduct":
		return s.queryAllProduct(APIstub)
	case "createData":
		return s.createData(APIstub,args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}


func (s *SmartContract) createData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

var record Record

err := json.Unmarshal([]byte(args[1]), &record)
if err != nil {
	 return shim.Success(nil)
}

dt := time.Now()
record.CreatedDate=dt.Format("01-02-2006 15:04:05")


recordAsBytes, _ := json.Marshal(record)
APIstub.PutState(args[0], recordAsBytes)

return shim.Success(recordAsBytes)

}


func (s *SmartContract) queryAllProduct(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "DATA0"
	endKey := "DATA10000"

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
		// Add a comma before array members, suppress it for the first array member
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

	fmt.Printf("- queryAllProduct:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
