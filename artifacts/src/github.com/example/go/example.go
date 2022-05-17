package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	case "queryAllData":
		return s.queryAllData(APIstub)
	case "createData":
		return s.createData(APIstub,args)
	case "queryData":
		return s.queryData(APIstub,args)
	case "queryDataBySurveyId":
		return s.queryDataBySurveyId(APIstub,args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}


func (s *SmartContract) createData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {


var breakPoint int
breakPoint=len(args)/4
dt := time.Now()
var record Record
var survey Survey;
var questions []Question;

record.UserId=args[0]
//record.createdDate=args[1]
record.CreatedDate=dt.Format("01-02-2006 15:04:05")
survey.SurveyId=args[2]
survey.SurveyDescription=args[3]


var question Question
var answer Answer

for index, element := range args {
	if index>=breakPoint*4{
		break
	}

	if index>3{
		if index%4==0{
			question.QuestionId=element
		} 
		if index%4==1{
			question.QuestionDescription=element
		} 
		if index%4==2 {
			answer.AnswerId=element
		} 
		if index%4==3 {
			answer.AnswerDescription=element
			question.Answer=answer
			questions=append(questions,question)
		}
	}

}
survey.Question=questions
record.Survey=survey



recordAsBytes, _ := json.Marshal(record)
APIstub.PutState(args[0], recordAsBytes)


indexName := "status~key"
nameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{record.Survey.SurveyId, args[0]})
if err != nil {
	return shim.Error(err.Error())
}
value := []byte{0x00}
APIstub.PutState(nameIndexKey, value)
return shim.Success(recordAsBytes)

}


func (s *SmartContract) queryAllData(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "DATA1"
	endKey := "DATA999"

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


func (s *SmartContract) queryData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	dataAsBytes, _ := APIstub.GetState(args[0])

	var record Record
	_ = json.Unmarshal(dataAsBytes, record)

	return shim.Success(dataAsBytes)
}


func (S *SmartContract) queryDataBySurveyId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	surveyId := args[0]

	surveyIdResultIterator, err := APIstub.GetStateByPartialCompositeKey("status~key", []string{surveyId})
	if err != nil {
		return shim.Error(err.Error())
	}

	defer surveyIdResultIterator.Close()

	var i int
	var id string

	var records []byte
	bArrayMemberAlreadyWritten := false

	records = append([]byte("["))

	for i = 0; surveyIdResultIterator.HasNext(); i++ {
		responseRange, err := surveyIdResultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		id = compositeKeyParts[1]
		assetAsBytes, err := APIstub.GetState(id)

		if bArrayMemberAlreadyWritten == true {
			newBytes := append([]byte(","), assetAsBytes...)
			records = append(records, newBytes...)

		} else {
			records = append(records, assetAsBytes...)
		}

		fmt.Printf("Found a asset for index : %s asset id : ", objectType, compositeKeyParts[0], compositeKeyParts[1])
		bArrayMemberAlreadyWritten = true

	}

	records = append(records, []byte("]")...)

	return shim.Success(records)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
