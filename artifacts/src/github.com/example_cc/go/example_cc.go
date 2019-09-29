

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
    "fmt"
    "strconv"
    "encoding/json"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
    "bytes"
    "time"
)

// Define the Smart Contract structure
type InsuranceChaincode struct {
}



var insuranceIndexStr = "_insid"
// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Insurance struct {
	//DocType string `json:"docType"`
	InsuranceId   string `json:"insuranceId"`
	InsuranceType  string `json:"insuranceType"`
	PolicyAmount  string `json:"policyAmount"`
	PaymentPolicy  string `json:"paymentPolicy"`
    Insurer string `json:"insurer"`
    Status string `json:"status"`
}
type Customer struct {
	//DocType string `json:"docType"`
   // InsuranceId string `json:"insuranceId"`
    CustomerId   string `json:"customerId"`
    //PolicyNumber string `json:"policyNumber"`
    CustomerName  string `json:"customerName"`
    CustomerEmail  string `json:"customerEmail"`
	CustomerAge string `json:"customerAge"`
	CustomerGender  string `json:"customerGender"`
}

type Policydetail struct {
    InsuranceId   string `json:"insuranceId"`
    InsuranceType  string `json:"insuranceType"`
    PolicyAmount  string `json:"policyAmount"`
    PaymentPolicy  string `json:"paymentPolicy"`
    Insurer string `json:"insurer"`
    Status string `json:"status"`
   // InsuranceId string `json:"insuranceId"`
    CustomerId   string `json:"customerId"`
    PolicyNumber string `json:"policyNumber"`
    CustomerName  string `json:"customerName"`
    CustomerEmail  string `json:"customerEmail"`
    CustomerAge string `json:"customerAge"`
    CustomerGender  string `json:"customerGender"`
}

func (t *InsuranceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *InsuranceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createInsurance" {
		return t.createInsurance(stub, args)

	} else if function == "getInsuranceByID" {  
		return t.getInsuranceByID(stub, args)
	} else if function == "getAllInsurance" {
		return t.getAllInsurance(stub, args)
	} else if function == "createCustomer" {
		return t.createCustomer(stub, args)
	} else if function == "insuranceTransfer" {
        return t.insuranceTransfer(stub, args)
    } else if function == "buyInsurance" {
        return t.buyInsurance(stub, args) 
    }else if function == "getPolicytHistory" {
        return t.getPolicytHistory(stub, args) 
    } else if function == "getInsuranceByPolicynumber" {
        return t.getInsuranceByPolicynumber(stub, args)
    }else if function == "getCustomerByID" {
        return t.getCustomerByID(stub, args)
    }
	return shim.Error("Invalid Smart Contract function name.")  
}


//create insurance
func (t *InsuranceChaincode) createInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error


    if len(args) != 5 {
        return shim.Error("Incorrect number of arguments. Expecting 5")
    }

    // ==== Input sanitation ====
    fmt.Println("- start createInsurance")
    if len(args[0]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    }
    if len(args[1]) <= 0 {
        return shim.Error("2nd argument must be a non-empty string")
    }
    if len(args[2]) <= 0 {
        return shim.Error("3rd argument must be a non-empty numeric string")
    }
    if len(args[3]) <= 0 {
        return shim.Error("4th argument must be a non-empty float string")
    }
    if len(args[4]) <= 0 {
        return shim.Error("5th argument must be a non-empty numeric value")
    }
    
    InsuranceId := args[0]
    InsuranceType := args[1]
    //PolicyNumber := args[2]
    PolicyAmount := args[2]
    PaymentPolicy := args[3]
    Insurer := args[4] 
    Status := "New Insurance created"

    // ==== Check if createInsurance already exists ====
    insuranceAsBytes, err := stub.GetState(InsuranceId)
    if err != nil {
        return shim.Error("Failed to get createInsurance: " + err.Error())
    } else if insuranceAsBytes != nil {
        fmt.Println("This createInsurance already exists: " + InsuranceId)
        return shim.Error("This createInsurance already exists: " + InsuranceId)
    }

    // ====marshal to JSON ====
    insurance := &Insurance{InsuranceId,InsuranceType,PolicyAmount,PaymentPolicy,Insurer,Status}
    insuranceJSONasBytes, err := json.Marshal(insurance)
    if err != nil {
        return shim.Error(err.Error())
    }
    // === Save createInsurance to state ===
    err = stub.PutState(InsuranceId, insuranceJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    //get the Insurance index
    insuranceIndexAsBytes, err := stub.GetState(insuranceIndexStr)
    if err != nil {
        return shim.Error(err.Error())
    }
    var insuranceIndex []string
    
    json.Unmarshal(insuranceIndexAsBytes, &insuranceIndex)                          //un stringify it aka JSON.parse()
    fmt.Print("insuranceIndex: ")
    fmt.Println(insuranceIndex)
    //append
    insuranceIndex = append(insuranceIndex, InsuranceId) 
    //add "insuranceId" to index list
    jsonAsBytes, _ := json.Marshal(insuranceIndex)
    //store "insuranceId" of createInsurance
    err = stub.PutState(insuranceIndexStr, jsonAsBytes)                      
    if err != nil {
        return shim.Error(err.Error())
    }

    eventMessage := "{ \"InsuranceId\" : \""+InsuranceId+"\", \"message\" : \"Insurance created succcessfully\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", []byte(eventMessage))
    if err != nil {
        return shim.Error(err.Error())
    } 

    // ==== createInsurance saved and indexed. Return success ====
    fmt.Println("- end createInsurance")
    return shim.Success(nil)
}



func (t *InsuranceChaincode) getInsuranceByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var insuranceId, jsonResp string

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting insuranceId of the Insurance to query")
    }

    insuranceId = args[0]
    valAsbytes, err := stub.GetState(insuranceId) //get the marble from chaincode state
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + insuranceId + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"InsuranceId\": \""+ insuranceId + "\", \"Error\":\"Insurance does not exist.\"}"
        return shim.Error(jsonResp)
    }

    return shim.Success(valAsbytes)
}

func (t *InsuranceChaincode) getAllInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var jsonResp, errResp string
    var insuranceIndex []string
    var err error
    
    fmt.Println("start getAllInsurance")

    insuranceIndexAsBytes, err := stub.GetState(insuranceIndexStr)
    if err != nil {
        return shim.Error(err.Error())
    }
    json.Unmarshal(insuranceIndexAsBytes, &insuranceIndex)                               
    //unstringify it aka JSON.parse()
    fmt.Print("insuranceIndex : ")
    fmt.Println(insuranceIndex)
    fmt.Println("len(insuranceIndex) : ")
    fmt.Println(len(insuranceIndex))
    jsonResp = "["
    for i,insuranceId := range insuranceIndex{
        fmt.Println(strconv.Itoa(i) + " - looking at " + insuranceId + " for all 'InsuranceId'")
        valueAsBytes, err := stub.GetState(insuranceId)
        if err != nil {
            errResp = "{\"Error\":\"Failed to get state for " + insuranceId + "\"}"
            return shim.Error(errResp)
        }
        jsonResp = jsonResp + string(valueAsBytes[:])
        if i < len(insuranceIndex)-1 {
            jsonResp = jsonResp + ","
        }
    }
    jsonResp = jsonResp + "]"
    fmt.Println("jsonResp : " + jsonResp)
    fmt.Println("end getAllInsurance")   
    return shim.Success([]byte(jsonResp))
}

// Customer Creation :

func (t *InsuranceChaincode) createCustomer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error


    if len(args) != 5 {
        return shim.Error("Incorrect number of arguments. Expecting 5")
    }

    // ==== Input sanitation ====
    fmt.Println("- start createCustomer")
    if len(args[0]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    }
    if len(args[1]) <= 0 {
        return shim.Error("2nd argument must be a non-empty string")
    }
    if len(args[2]) <= 0 {
        return shim.Error("3rd argument must be a non-empty numeric string")
    }
    if len(args[3]) <= 0 {
        return shim.Error("4th argument must be a non-empty float string")
    }
    if len(args[4]) <= 0 {
        return shim.Error("5th argument must be a non-empty numeric value")
    }

    CustomerId := args[0]
    CustomerName  := args[1]
    CustomerEmail  := args[2]
    CustomerAge := args[3]
    CustomerGender := args[4]

    // ==== Check if createInsurance already exists ====
    customerAsBytes, err := stub.GetState(CustomerId)
    if err != nil {
        return shim.Error("Failed to get createInsurance: " + err.Error())
    } else if customerAsBytes != nil {
        fmt.Println("This createCustomer already exists: " + CustomerId)
        return shim.Error("This createCustomer already exists: " + CustomerId)
    }

    // ====marshal to JSON ====
    customer := &Customer{CustomerId,CustomerName,CustomerEmail,CustomerAge,CustomerGender}
    customerJSONasBytes, err := json.Marshal(customer)
    if err != nil {
        return shim.Error(err.Error())
    }
    // === Save createICustomer to state ===
    err = stub.PutState(CustomerId, customerJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }


    eventMessage := "{ \"CustomerId\" : \""+CustomerId+"\", \"message\" : \"customer created succcessfully\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", []byte(eventMessage))
    if err != nil {
        return shim.Error(err.Error())
    } 

    // ==== create saved and indexed. Return success ====
    fmt.Println("- end customer created")
    return shim.Success(nil)
}

func (t *InsuranceChaincode) insuranceTransfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

var jsonResp string

    if len(args) != 3 {
        return shim.Error("Incorrect number of arguments. Expecting 3")
    }


        CustomerId := args[0]
        InsuranceId := args[1]
        PolicyNumber := args[2]


    valAsbytes, err := stub.GetState(CustomerId) 

    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + CustomerId + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"CustomerId\": \""+ CustomerId + "\", \"Error\":\"Customer does not exist.\"}"
        return shim.Error(jsonResp)
    }

    tempCustomerStruct1 := Customer{};
    json.Unmarshal(valAsbytes, &tempCustomerStruct1)
    fmt.Print("temp customer index : ")
    fmt.Println(tempCustomerStruct1)

    CustomerName := tempCustomerStruct1.CustomerName
    CustomerEmail  := tempCustomerStruct1.CustomerEmail
    CustomerAge := tempCustomerStruct1.CustomerAge
    CustomerGender := tempCustomerStruct1.CustomerGender

    InsvalAsbytes, err := stub.GetState(InsuranceId)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + InsuranceId + "\"}"
        return shim.Error(jsonResp)
    } else if InsvalAsbytes == nil {
        jsonResp = "{\"CustomerId\": \""+ InsuranceId + "\", \"Error\":\"Insurance does not exist.\"}"
        return shim.Error(jsonResp)
    }
    tempInsuranceStruct := Insurance{};

    json.Unmarshal(InsvalAsbytes, &tempInsuranceStruct)
    fmt.Print("temp insurance index : ")
    fmt.Println(tempInsuranceStruct)
    InsuranceType := tempInsuranceStruct.InsuranceType
    PolicyAmount := tempInsuranceStruct.PolicyAmount
    PaymentPolicy := tempInsuranceStruct.PaymentPolicy
    Insurer := tempInsuranceStruct.Insurer
    Status := "Insurance transfered"

    detail := &Policydetail{InsuranceId,InsuranceType,PolicyAmount,PaymentPolicy,Insurer,Status,CustomerId,PolicyNumber,CustomerName,CustomerEmail,CustomerAge,CustomerGender}
    detailJSONasBytes, err := json.Marshal(detail)
    if err != nil {
        return shim.Error(err.Error())
    }
    // === Save to state ===
    err = stub.PutState(PolicyNumber, detailJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

     err = stub.PutState(InsuranceId, detailJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    return shim.Success(nil)
}

func (t *InsuranceChaincode) getInsuranceByPolicynumber(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var policyNumber, jsonResp string

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting policyNumber of the Insurance to query")
    }

    policyNumber = args[0]
    valAsbytes, err := stub.GetState(policyNumber) //get the marble from chaincode state
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + policyNumber + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"InsuranceId\": \""+ policyNumber + "\", \"Error\":\"Insurance does not exist.\"}"
        return shim.Error(jsonResp)
    }

    return shim.Success(valAsbytes)
}

func (t *InsuranceChaincode) buyInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {

var jsonResp string
        if len(args) != 3 {
        return shim.Error("Incorrect number of arguments. Expecting 3")
    }


        CustomerId := args[0]
        InsuranceId := args[1]
        PolicyNumber := args[2]

    valAsbytes, err := stub.GetState(CustomerId) 

    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + CustomerId + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"CustomerId\": \""+ CustomerId + "\", \"Error\":\"Customer does not exist.\"}"
        return shim.Error(jsonResp)
    }

    tempCustomerStruct1 := Customer{};
    json.Unmarshal(valAsbytes, &tempCustomerStruct1)
    fmt.Print("temp customer index : ")
    fmt.Println(tempCustomerStruct1)


    CustomerName := tempCustomerStruct1.CustomerName
    CustomerEmail  := tempCustomerStruct1.CustomerEmail
    CustomerAge := tempCustomerStruct1.CustomerAge
    CustomerGender := tempCustomerStruct1.CustomerGender

    InsvalAsbytes, err := stub.GetState(InsuranceId)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + InsuranceId + "\"}"
        return shim.Error(jsonResp)
    } else if InsvalAsbytes == nil {
        jsonResp = "{\"CustomerId\": \""+ InsuranceId + "\", \"Error\":\"Insurance does not exist.\"}"
        return shim.Error(jsonResp)
    }
    tempInsuranceStruct := Insurance{};

    json.Unmarshal(InsvalAsbytes, &tempInsuranceStruct)
    fmt.Print("temp insurance index : ")
    fmt.Println(tempInsuranceStruct)
    InsuranceType := tempInsuranceStruct.InsuranceType
    PolicyAmount := tempInsuranceStruct.PolicyAmount
    PaymentPolicy := tempInsuranceStruct.PaymentPolicy
    Insurer := tempInsuranceStruct.Insurer
    Status := "Insurance bought"

    detail := &Policydetail{InsuranceId,InsuranceType,PolicyAmount,PaymentPolicy,Insurer,Status,CustomerId,PolicyNumber,CustomerName,CustomerEmail,CustomerAge,CustomerGender}
    detailJSONasBytes, err := json.Marshal(detail)
    if err != nil {
        return shim.Error(err.Error())
    }
    // === Save createInsurance to state ===
    err = stub.PutState(PolicyNumber, detailJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    err = stub.PutState(InsuranceId, detailJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    return shim.Success(nil)
}

// Get Customer by Id

func (t *InsuranceChaincode) getCustomerByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var customerId, jsonResp string

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting customerId of the Insurance to query")
    }

    customerId = args[0]
    valAsbytes, err := stub.GetState(customerId) 
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + customerId + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"CustomerId\": \""+ customerId + "\", \"Error\":\"Customer does not exist.\"}"
        return shim.Error(jsonResp)
    }

    return shim.Success(valAsbytes)
}

func (t *InsuranceChaincode) getPolicytHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    InsuranceId := args[0]

    fmt.Printf("- start getPolicytHistory: %s\n", InsuranceId)

    resultsIterator, err := stub.GetHistoryForKey(InsuranceId)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    // buffer is a JSON array containing historic values for the policy
    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        response, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        // Add a comma before array members, suppress it for the first array member
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"TxId\":")
        buffer.WriteString("\"")
        buffer.WriteString(response.TxId)
        buffer.WriteString("\"")

        buffer.WriteString(", \"Value\":")
        // if it was a delete operation on given key, then we need to set the
        //corresponding value null. Else, we will write the response.Value
        //as-is (as the Value itself a JSON Demand request)
        if response.IsDelete {
            buffer.WriteString("null")
        } else {
            buffer.WriteString(string(response.Value))
        }

        buffer.WriteString(", \"Timestamp\":")
        buffer.WriteString("\"")
        buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
        buffer.WriteString("\"")

        buffer.WriteString(", \"IsDelete\":")
        buffer.WriteString("\"")
        buffer.WriteString(strconv.FormatBool(response.IsDelete))
        buffer.WriteString("\"")

        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")

    fmt.Printf("- getPolicytHistory returning:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(InsuranceChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
