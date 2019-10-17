

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
type LogisticsChaincode struct {
}



var productIndexStr = "_prid"
var poIndexStr = "_poid"
// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Product struct {

    ProductID   string `json:"productID"`
    SellerName  string `json:"sellerName"`
    SellerID  string `json:"sellerID"`
	ItemName string `json:"itemName"`
    Price  string `json:"price"`
    Status  string `json:"status"`
}
type PO struct {
	//DocType string `json:"docType"`
    ProductID   string `json:"productID"`
    PoID   string `json:"poID"`
    SellerName  string `json:"sellerName"`
    SellerID  string `json:"sellerID"`
    BuyerName  string `json:"buyerName"`
    BuyerID  string `json:"buyerID"`
    ItemName string `json:"itemName"`
    NoOfItem string `json:"noOfItem"`
    Price  string `json:"price"`
    TotalPrice string `json:"totalPrice"`
    OrderDate string `json:"orderDate"`
    DeliveryAddress string `json:"deliveryAddress"`
    DeliveryDate string `json:"deliveryAddress"`
    Status string `json:"status"`
}



func (t *LogisticsChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


func (t *LogisticsChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createProduct" {
		return t.createProduct(stub, args)

	} else if function == "getPOByID" {  
		return t.getPOByID(stub, args)
	} else if function == "getAllProduct" {
		return t.getAllProduct(stub, args)
	} else if function == "createPO" {
		return t.createPO(stub, args)
	} else if function == "POinTransit" {
        return t.POinTransit(stub, args)
    } else if function == "buyerAccepted" {
        return t.buyerAccepted(stub, args) 
    }else if function == "buyerRejected" {
        return t.buyerRejected(stub, args) 
    } else if function == "getAllPO" {
        return t.getAllPO(stub, args)
    }
	return shim.Error("Invalid Smart Contract function name.")  
}


//create Product
func (t *LogisticsChaincode) createProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error


    if len(args) != 5 {
        return shim.Error("Incorrect number of arguments. Expecting 5")
    }

    // ==== Input sanitation ====
    fmt.Println("- start createProduct")
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
    
    ProductID := args[0]
    SellerName := args[1]
    //PolicyNumber := args[2]
    SellerID := args[2]
    ItemName := args[3]
    Price := args[4] 
    Status := "With Seller"

    // ==== Check if createProduct already exists ====
    productAsBytes, err := stub.GetState(ProductID)
    if err != nil {
        return shim.Error("Failed to get createProduct: " + err.Error())
    } else if productAsBytes != nil {
        fmt.Println("This createProduct already exists: " + ProductID)
        return shim.Error("This createProduct already exists: " + ProductID)
    }

    // ====marshal to JSON ====
    product := &Product{ProductID,SellerName,SellerID,ItemName,Price,Status}
    productJSONasBytes, err := json.Marshal(product)
    if err != nil {
        return shim.Error(err.Error())
    }
    // === Save createProduct to state ===
    err = stub.PutState(ProductID, productJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    //get the product index
    productIndexAsBytes, err := stub.GetState(productIndexStr)
    if err != nil {
        return shim.Error(err.Error())
    }
    var productIndex []string
    
    json.Unmarshal(productIndexAsBytes, &productIndex)                          //un stringify it aka JSON.parse()
    fmt.Print("productIndex: ")
    fmt.Println(productIndex)
    //append
    productIndex = append(productIndex, ProductID) 
    //add "ProductID" to index list
    jsonAsBytes, _ := json.Marshal(productIndex)
    //store "ProductID" of createProduct
    err = stub.PutState(productIndexStr, jsonAsBytes)                      
    if err != nil {
        return shim.Error(err.Error())
    }

    eventMessage := "{ \"ProductID\" : \""+ProductID+"\", \"message\" : \"Product created succcessfully\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", []byte(eventMessage))
    if err != nil {
        return shim.Error(err.Error())
    } 

    // ==== createProduct saved and indexed. Return success ====
    fmt.Println("- end createProduct")
    return shim.Success(nil)
}

// Create PO 
func (t *LogisticsChaincode) createPO(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error


    if len(args) != 13 {
        return shim.Error("Incorrect number of arguments. Expecting 5")
    }

    // ==== Input sanitation ====
    fmt.Println("- start createPO")
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
    if len(args[5]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    }
    if len(args[6]) <= 0 {
        return shim.Error("2nd argument must be a non-empty string")
    }
    if len(args[7]) <= 0 {
        return shim.Error("3rd argument must be a non-empty numeric string")
    }
    if len(args[8]) <= 0 {
        return shim.Error("4th argument must be a non-empty float string")
    }
    if len(args[9]) <= 0 {
        return shim.Error("5th argument must be a non-empty numeric value")
    }
    if len(args[10]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    }
    if len(args[11]) <= 0 {
        return shim.Error("2nd argument must be a non-empty string")
    }
    if len(args[12]) <= 0 {
        return shim.Error("3rd argument must be a non-empty numeric string")
    }

    ProductID := args[0]
    PoID := args[1]
    SellerName := args[2]
    SellerID := args[3]
    BuyerName := args[4]
    BuyerID := args[5]
    ItemName := args[6]
    NoOfItem := args[7]
    Price := args[8]
    TotalPrice := args[9]
    OrderDate := args[10]
    DeliveryAddress := args[11]
    DeliveryDate := args[12]
    Status := "Po Created"

    // ==== Check if createPO already exists ====
    poAsBytes, err := stub.GetState(PoID)
    if err != nil {
        return shim.Error("Failed to get PO: " + err.Error())
    } else if poAsBytes != nil {
        fmt.Println("This PO already exists: " + PoID)
        return shim.Error("This createProduct already exists: " + PoID)
    }

    // ====marshal to JSON ====
    po := &PO{ProductID,PoID,SellerName,SellerID,BuyerName,BuyerID,ItemName,NoOfItem,Price,TotalPrice,OrderDate,DeliveryAddress,DeliveryDate,Status}
    poJSONasBytes, err := json.Marshal(po)
    if err != nil {
        return shim.Error(err.Error())
    }
    // === Save createPO to state ===
    err = stub.PutState(PoID, poJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    //get the product index
    poIndexAsBytes, err := stub.GetState(poIndexStr)
    if err != nil {
        return shim.Error(err.Error())
    }
    var poIndex []string
    
    json.Unmarshal(poIndexAsBytes, &poIndex)                          //un stringify it aka JSON.parse()
    fmt.Print("productIndex: ")
    fmt.Println(productIndex)
    //append
    poIndex = append(poIndex, PoID) 
    //add "PoID" to index list
    jsonAsBytes, _ := json.Marshal(poIndex)
    //store "PoID" of createPO
    err = stub.PutState(poIndexStr, jsonAsBytes)                      
    if err != nil {
        return shim.Error(err.Error())
    }

    eventMessage := "{ \"PoID\" : \""+PoID+"\", \"message\" : \"PO created succcessfully\", \"code\" : \"200\"}"
    err = stub.SetEvent("evtsender", []byte(eventMessage))
    if err != nil {
        return shim.Error(err.Error())
    } 

    // ==== createPO saved and indexed. Return success ====
    fmt.Println("- end createPO")
    return shim.Success(nil)
}

//get all Product
func (t *LogisticsChaincode) getAllProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var jsonResp, errResp string
    var productIndex []string
    var err error
    
    fmt.Println("start getAllProduct")

    productIndexAsBytes, err := stub.GetState(productIndexStr)
    if err != nil {
        return shim.Error(err.Error())
    }
    json.Unmarshal(productIndexAsBytes, &productIndex)                               
    //unstringify it aka JSON.parse()
    fmt.Print("productIndex : ")
    fmt.Println(productIndex)
    fmt.Println("len(productIndex) : ")
    fmt.Println(len(productIndex))
    jsonResp = "["
    for i,productID := range productIndex{
        fmt.Println(strconv.Itoa(i) + " - looking at " + productID + " for all 'productID'")
        valueAsBytes, err := stub.GetState(productID)
        if err != nil {
            errResp = "{\"Error\":\"Failed to get state for " + productID + "\"}"
            return shim.Error(errResp)
        }
        jsonResp = jsonResp + string(valueAsBytes[:])
        if i < len(productIndex)-1 {
            jsonResp = jsonResp + ","
        }
    }
    jsonResp = jsonResp + "]"
    fmt.Println("jsonResp : " + jsonResp)
    fmt.Println("end getAllProduct")   
    return shim.Success([]byte(jsonResp))
}

//get all PO
func (t *LogisticsChaincode) getAllPO(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var jsonResp, errResp string
    var poIndex []string
    var err error
    
    fmt.Println("start getAllPO")

    poIndexAsBytes, err := stub.GetState(poIndexStr)
    if err != nil {
        return shim.Error(err.Error())
    }
    json.Unmarshal(poIndexAsBytes, &poIndex)                               
    //unstringify it aka JSON.parse()
    fmt.Print("poIndex : ")
    fmt.Println(poIndex)
    fmt.Println("len(poIndex) : ")
    fmt.Println(len(poIndex))
    jsonResp = "["
    for i,poID := range poIndex{
        fmt.Println(strconv.Itoa(i) + " - looking at " + poID + " for all 'productID'")
        valueAsBytes, err := stub.GetState(poID)
        if err != nil {
            errResp = "{\"Error\":\"Failed to get state for " + poID + "\"}"
            return shim.Error(errResp)
        }
        jsonResp = jsonResp + string(valueAsBytes[:])
        if i < len(poIndex)-1 {
            jsonResp = jsonResp + ","
        }
    }
    jsonResp = jsonResp + "]"
    fmt.Println("jsonResp : " + jsonResp)
    fmt.Println("end getAllPO")   
    return shim.Success([]byte(jsonResp))
}

//po in-transit
func (t *LogisticsChaincode) POinTransit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

var jsonResp string

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 3")
    }


    PoID := args[0]

    valAsbytes, err := stub.GetState(PoID) 

    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + PoID + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"PoID\": \""+ PoID + "\", \"Error\":\"PO does not exist.\"}"
        return shim.Error(jsonResp)
    }

    tempPOStruct := PO{};
    json.Unmarshal(valAsbytes, &tempPOStruct)
    fmt.Print("temp PO index : ")
    fmt.Println(tempPOStruct)

    SellerName := tempPOStruct.SellerName
    SellerID  := tempPOStruct.SellerID
    BuyerName := tempPOStruct.BuyerName
    BuyerID := tempPOStruct.BuyerID
    ItemName := tempPOStruct.ItemName
    NoOfItem  := tempPOStruct.NoOfItem
    Price := tempPOStruct.Price
    TotalPrice := tempPOStruct.TotalPrice
    OrderDate := tempPOStruct.OrderDate
    DeliveryAddress  := tempPOStruct.DeliveryAddress
    DeliveryDate := tempPOStruct.DeliveryDate
    Status := "In-Transit"


    detail := &PO{ProductID,PoID,SellerName,SellerID,BuyerName,BuyerID,ItemName,NoOfItem,Price,TotalPrice,OrderDate,DeliveryAddress,DeliveryDate,Status}
    detailJSONasBytes, err := json.Marshal(detail)
    if err != nil {
        return shim.Error(err.Error())
    }
    // === Save to state ===
    err = stub.PutState(PoID, detailJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    return shim.Success(nil)
}

//po Accepted By Buyer
func (t *LogisticsChaincode) POinTransit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    var jsonResp string
    
        if len(args) != 1 {
            return shim.Error("Incorrect number of arguments. Expecting 3")
        }
    
    
        PoID := args[0]
    
        valAsbytes, err := stub.GetState(PoID) 
    
        if err != nil {
            jsonResp = "{\"Error\":\"Failed to get state for " + PoID + "\"}"
            return shim.Error(jsonResp)
        } else if valAsbytes == nil {
            jsonResp = "{\"PoID\": \""+ PoID + "\", \"Error\":\"PO does not exist.\"}"
            return shim.Error(jsonResp)
        }
    
        tempPOStruct := PO{};
        json.Unmarshal(valAsbytes, &tempPOStruct)
        fmt.Print("temp PO index : ")
        fmt.Println(tempPOStruct)
    
        SellerName := tempPOStruct.SellerName
        SellerID  := tempPOStruct.SellerID
        BuyerName := tempPOStruct.BuyerName
        BuyerID := tempPOStruct.BuyerID
        ItemName := tempPOStruct.ItemName
        NoOfItem  := tempPOStruct.NoOfItem
        Price := tempPOStruct.Price
        TotalPrice := tempPOStruct.TotalPrice
        OrderDate := tempPOStruct.OrderDate
        DeliveryAddress  := tempPOStruct.DeliveryAddress
        DeliveryDate := tempPOStruct.DeliveryDate
        Status := "In-Transit"
    
    
        detail := &PO{ProductID,PoID,SellerName,SellerID,BuyerName,BuyerID,ItemName,NoOfItem,Price,TotalPrice,OrderDate,DeliveryAddress,DeliveryDate,Status}
        detailJSONasBytes, err := json.Marshal(detail)
        if err != nil {
            return shim.Error(err.Error())
        }
        // === Save to state ===
        err = stub.PutState(PoID, detailJSONasBytes)
        if err != nil {
            return shim.Error(err.Error())
        }
    
        return shim.Success(nil)
    }

//po Rejected By Buyer
func (t *LogisticsChaincode) POinTransit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    var jsonResp string
    
        if len(args) != 1 {
            return shim.Error("Incorrect number of arguments. Expecting 3")
        }
    
    
        PoID := args[0]
    
        valAsbytes, err := stub.GetState(PoID) 
    
        if err != nil {
            jsonResp = "{\"Error\":\"Failed to get state for " + PoID + "\"}"
            return shim.Error(jsonResp)
        } else if valAsbytes == nil {
            jsonResp = "{\"PoID\": \""+ PoID + "\", \"Error\":\"PO does not exist.\"}"
            return shim.Error(jsonResp)
        }
    
        tempPOStruct := PO{};
        json.Unmarshal(valAsbytes, &tempPOStruct)
        fmt.Print("temp PO index : ")
        fmt.Println(tempPOStruct)
    
        SellerName := tempPOStruct.SellerName
        SellerID  := tempPOStruct.SellerID
        BuyerName := tempPOStruct.BuyerName
        BuyerID := tempPOStruct.BuyerID
        ItemName := tempPOStruct.ItemName
        NoOfItem  := tempPOStruct.NoOfItem
        Price := tempPOStruct.Price
        TotalPrice := tempPOStruct.TotalPrice
        OrderDate := tempPOStruct.OrderDate
        DeliveryAddress  := tempPOStruct.DeliveryAddress
        DeliveryDate := tempPOStruct.DeliveryDate
        Status := "In-Transit"
    
    
        detail := &PO{ProductID,PoID,SellerName,SellerID,BuyerName,BuyerID,ItemName,NoOfItem,Price,TotalPrice,OrderDate,DeliveryAddress,DeliveryDate,Status}
        detailJSONasBytes, err := json.Marshal(detail)
        if err != nil {
            return shim.Error(err.Error())
        }
        // === Save to state ===
        err = stub.PutState(PoID, detailJSONasBytes)
        if err != nil {
            return shim.Error(err.Error())
        }
    
        return shim.Success(nil)
    }

// Get PO by Id

func (t *LogisticsChaincode) getPOByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var poID, jsonResp string

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting poID of the PO to query")
    }

    poID = args[0]
    valAsbytes, err := stub.GetState(poID) 
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + poID + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"poID\": \""+ poID + "\", \"Error\":\"PO does not exist.\"}"
        return shim.Error(jsonResp)
    }

    return shim.Success(valAsbytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(LogisticsChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}