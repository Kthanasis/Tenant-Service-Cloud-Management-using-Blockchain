
  
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


// SmartContract provides functions for managing a struct 
type SmartContract struct {
	contractapi.Contract
}


//---------------------------Helping Functions----------------------------------------
//Function to check if a revoker is in the Revoker list 
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

//-----------------------------End of Helping Function-----------------------------





//----------------------Structs for Tenants, Services and Delegations-------------------------------


// Tenant describes basic details of what makes up a tenant
type Tenant struct {
	Pck			string   `json:"pck"` 		
	Name        string 	 `json:"name"`
	Email       string 	 `json:"email"`
	Phone       string   `json:"phone"`
	Registered 	bool     `json:"registered"`	//false if not, true if Registered
	Type        string   `json:"type"`		    //T for tenants
}


//Service describes basic details of what makes up a service
type Service struct {
	Pck					string  `json:"pck"`
	Name				string  `json:"name"`
	Registered			bool    `json:"registered"` //true if registered false if not 
	Type 				string 	`json:"Type"`       //S for Services
}


//Delegation describes basic details of what makes up a Delegation
type Delegation struct {

	Pck					string	 `json:"pck"`	               //	
	Grandor				string   `json:"grandor"`              //			 
	Recipient  			string 	 `json:"recipient"`			   //
	Subdel       		uint8 	 `json:"subdel"`			   //
	Issue 				uint64   `json:"issue"`				   //
	Expiry 				uint64   `json:"expiry"`			   //
	Suspended			bool	 `json:"suspended"`			   //false if not, true if suspended
	Revoked 			bool	 `json:"revoked"`			   //false if not, true if revoked
	Revokers 			[]string `json:"revokers"`      	   //list of tenants & services who can revoke the delegation
	DelegationChain		[]string `json:"delegationchain"`	   //
	Type 				string 	 `json:"Type"`				   //D is for Delegation	
}


type SubDelegation struct {
	Pck					string	 	`json:"pck"`	               //
	Grandor				string   	`json:"grandor"`               //		 
	Recipient  			string 	 	`json:"recipient"`			   //
	Subdel       		uint8 	 	`json:"subdel"`			  	   //
	Issue 				uint64   	`json:"issue"`				   //
	Expiry 				uint64   	`json:"expiry"`				   //
	Suspended			bool	 	`json:"suspended"`			   //false if not, true if suspended
	Revoked 			bool	 	`json:"revoked"`			   //false if not, true if revoked
	Revokers 			[]string 	`json:"revokers"`      		   //list of tenants & services who can revoke the subdelegation
	DelegationChain		[]string 	`json:"delegationchain"`       //
	Type 				string 		`json:"Type"`				   //SD is for SubDelegation 
}

//----------------------------------------End Of Structs------------------------------------------- 



// InitLedger adds a base set of tenants and services to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	tenants := []Tenant{
		Tenant{Pck: "T1", Name: "Tenant One",   Email: "t1@mail.com", Phone: "1111111111", Registered: true, Type: "T"},
		Tenant{Pck: "T2", Name: "Tenant Two",   Email: "t2@mail.com", Phone: "2222222222", Registered: true, Type: "T"},
		Tenant{Pck: "T3", Name: "Tenant Three", Email: "t3@mail.com", Phone: "3333333333", Registered: true, Type: "T"},
		Tenant{Pck: "T4", Name: "Tenant Four",  Email: "t4@mail.com", Phone: "4444444444", Registered: true, Type: "T"},
		Tenant{Pck: "T5", Name: "Tenant Five",  Email: "t5@mail.com", Phone: "5555555555", Registered: true, Type: "T"},
		Tenant{Pck: "T6", Name: "Tenant Six",   Email: "t6@mail.com", Phone: "6666666666", Registered: true, Type: "T"},
		Tenant{Pck: "T7", Name: "Tenant Seven", Email: "t7@mail.com", Phone: "7777777777", Registered: true, Type: "T"},
		Tenant{Pck: "T8", Name: "Tenant Eight", Email: "t8@mail.com", Phone: "8888888888", Registered: true, Type: "T"},
	}


	//We save the data to the world State based on their Pck
	for _, tenant := range tenants {
		tenantAsBytes, _ := json.Marshal(tenant)
		err := ctx.GetStub().PutState(tenant.Pck, tenantAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}
	
	services := []Service{
		Service{Pck: "S1", Name: "Service One",   Registered: true, Type: "S"},
		Service{Pck: "S2", Name: "Service Two",   Registered: true, Type: "S"},
		Service{Pck: "S3", Name: "Service Three", Registered: true, Type: "S"},
	}


	//We save the data to the world State based on their Pck
	for _, service := range services {
		serviceAsBytes, _ := json.Marshal(service)
		err := ctx.GetStub().PutState(service.Pck, serviceAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}





//***************************************************************************************************
//**																							   **
//**							The following section Manages the SubDelegations                   ** 
//**                                                                                               **
//***************************************************************************************************


//-------------------------------------SubDelegation Management--------------------------------------
//in this section there are the basic functions to manage SubDelegations 



//RegisterSubDelegation adds a new SubDelegation to the world state with given details
func (s *SmartContract) RegisterSubDelegation(ctx contractapi.TransactionContextInterface, pck string, exdelegation string, recipient string, subdel string, issue string, expiry string) error {
	//getting the previous in chain delegation info 
	delegation, err := s.IsDelegation(ctx, exdelegation)
	if err != nil {
		return err
	}
	
	//turning string to uint8
	tempsubdel, err := strconv.ParseUint(subdel, 10, 8) 
	if err != nil {
		fmt.Printf("Error")
		}
	subdel1 :=uint8(tempsubdel)

	//turning string to uint64
	issue1, err := strconv.ParseUint(issue, 10, 64)
	if 	err != nil {
		fmt.Printf("Error")
		}

	expiry1, err := strconv.ParseUint(expiry, 10, 64)
	if err != nil {
		fmt.Printf("Error")
		}
	
	
	//checking if the recipient is a Tenant, Services are NOT allowed to be sudelegated
	tenant, err := s.IsTenant(ctx, recipient)
	if err != nil {
		return err
	}
	if tenant.Type != "T"{
		return fmt.Errorf("Recipient of a SubDelegation must be a Tenant")
	}
	
	if tenant.Registered == false {
		return fmt.Errorf("Cannot create the Subdelegation with a tenant that is already destroyed")
	}
	
	//ckecking for selfsubdelegations
	if delegation.Recipient == recipient {
		return fmt.Errorf("Cannot Self-SubDelegate")
	}

	//checking the subdel 
	if delegation.Subdel - 1 < subdel1 {
		return fmt.Errorf("Subdel must be smaller")
	}
	
	//checking issue parameter, we accept equals
	if delegation.Issue > issue1 {
		return fmt.Errorf("Issue must be post the issue of the previous delegation")
	}
	
	if issue1 > expiry1 {
		return fmt.Errorf("Delegation Issue and Expiry times should be checked again, Issue cannot be after the expiry")
	}
	
	//checking expiry parameter, we accept equals
	if delegation.Expiry < expiry1  {
		return fmt.Errorf("expiry must be prior the expiry of the previous delegation")
	}
	
	//checking if the previous delegation is suspended
	if delegation.Suspended == true  {
		return fmt.Errorf("Cannot create SubDelegation because Delegation has been Suspended")
	}
	
	//checking if the previous delegation is revoked 
	if delegation.Revoked == true  {
		return fmt.Errorf("Cannot create SubDelegation because Delegation has been Revoked")
	}
	
	
	//checking if it is already expired upon creation
	expiryint64, _ := strconv.ParseUint(expiry, 10, 64)
	timenow := uint64(time.Now().Unix())
	if timenow > expiryint64 {
		return fmt.Errorf("Cannot create a SubDelegation which is already Expired")
	}
	
	//checking if a previous delegation has been revoked or suspended, if yes we do not create the subdelegation
	for _, x := range delegation.DelegationChain{
		temp, err := s.IsDelegation(ctx, x)
		if err != nil {
			return err
		}
		if temp.Suspended == true {
			return fmt.Errorf("Cannot create a SubDelegation because a previous Delegation has been Suspended")
		} else if temp.Revoked == true {
			return fmt.Errorf("Cannot create a SubDelegation because a previous Delegation has been Revoked")
		}
	}
	
	
	//updating the subdel field on the previous delegation 
	if delegation.Type == "D" {
		delegation.Subdel = delegation.Subdel - 1 - subdel1
		delegationAsBytes, _ := json.Marshal(delegation)
		ctx.GetStub().PutState(delegation.Pck, delegationAsBytes)
	} else {
		tempdelegation, err := s.IsSubDelegation(ctx, exdelegation)
		if err != nil {
			return err
		}
		tempdelegation.Subdel = tempdelegation.Subdel - 1 - subdel1
		tempdelegationAsBytes, _ := json.Marshal(tempdelegation)
		ctx.GetStub().PutState(tempdelegation.Pck, tempdelegationAsBytes)
	}
	
	//creating the Revokers list 
	var finalrevokers []string
	finalrevokers = append(delegation.Revokers,recipient)

	//creating the delegationchain 
	var tempdelegationchain []string 
	tempdelegationchain = append(delegation.DelegationChain,pck)


	//we pass the data to the struct and then we save to the world state 
	subdelegation := SubDelegation{
		Pck:				pck,	 	
		Grandor:			delegation.Recipient,   			 
		Recipient:  		recipient, 	 
		Subdel:       		subdel1, 	 
		Issue: 				issue1,   
		Expiry: 			expiry1,   
		Suspended:			false,
		Revoked:			false,
		Revokers: 			finalrevokers, 
		DelegationChain:	tempdelegationchain,
		Type: 				"SD",
	}

	subdelegationAsBytes, _ := json.Marshal(subdelegation)

	return ctx.GetStub().PutState(pck, subdelegationAsBytes)
}


//function to update the world state with the new Suspended status, true when the delegation has been suspended
func (s *SmartContract) SuspendSubDelegation(ctx contractapi.TransactionContextInterface, pck string) error {
	//we pull from the world state the data for the delegation
	subdelegation, err := s.IsSubDelegation(ctx, pck)

	if err != nil {
		return err
	}

	//we update the Suspended field 
	subdelegation.Suspended = true
	
	//we update the subdel field to the privious delegation 
	previousdelegation := subdelegation.DelegationChain[len(subdelegation.DelegationChain)-2]
	temppreviousdelegation, err := s.IsDelegation(ctx, previousdelegation)
		if err != nil {
			return err
		}
	
	
	//updating the subdel field on the previous delegation 
	if temppreviousdelegation.Type == "D" {
		newsubdel:=temppreviousdelegation.Subdel + 1 + subdelegation.Subdel
		temppreviousdelegation.Subdel = newsubdel 
		temppreviousdelegationAsBytes, _ := json.Marshal(temppreviousdelegation)
		ctx.GetStub().PutState(temppreviousdelegation.Pck, temppreviousdelegationAsBytes)
	} else {
		tempdelegation, err := s.IsSubDelegation(ctx, previousdelegation)
		if err != nil {
			return err
		}
		tempdelegation.Subdel = tempdelegation.Subdel + 1 + subdelegation.Subdel
		tempdelegationAsBytes, _ := json.Marshal(tempdelegation)
		ctx.GetStub().PutState(tempdelegation.Pck, tempdelegationAsBytes)
	}
	
	
	//we store back to the world state
	subdelegationAsBytes, _ := json.Marshal(subdelegation)

	return ctx.GetStub().PutState(pck, subdelegationAsBytes)
}


//function to update the world state with the new Revoke status, true when the delegation has been revoked
func (s *SmartContract) RevokeSubDelegation(ctx contractapi.TransactionContextInterface, pck string, revoker string ) error {
	//we pull from the world state the data for the delegation
	subdelegation, err := s.IsSubDelegation(ctx, pck)

	if err != nil {
		return err
	}

	if stringInSlice(revoker, subdelegation.Revokers){
			//we update the Revoked field 
			subdelegation.Revoked = true
			} else {
				return fmt.Errorf("%s is not an authorized Revoker", revoker)
		}
		
	//we update the subdel field to the privious delegation 
	previousdelegation := subdelegation.DelegationChain[len(subdelegation.DelegationChain)-2]
	temppreviousdelegation, err := s.IsDelegation(ctx, previousdelegation)
		if err != nil {
			return err
		}
	
	
	//updating the subdel field on the previous delegation 
	if temppreviousdelegation.Type == "D" {
		temppreviousdelegation.Subdel = temppreviousdelegation.Subdel + 1 + subdelegation.Subdel
		temppreviousdelegationAsBytes, _ := json.Marshal(temppreviousdelegation)
		ctx.GetStub().PutState(temppreviousdelegation.Pck, temppreviousdelegationAsBytes)
	} else {
		tempdelegation, err := s.IsSubDelegation(ctx, previousdelegation)
		if err != nil {
			return err
		}
		tempdelegation.Subdel = tempdelegation.Subdel + 1 + subdelegation.Subdel
		tempdelegationAsBytes, _ := json.Marshal(tempdelegation)
		ctx.GetStub().PutState(tempdelegation.Pck, tempdelegationAsBytes)
	}
		
		

	//we store back to the world state
	subdelegationAsBytes, _ := json.Marshal(subdelegation)

	return ctx.GetStub().PutState(pck, subdelegationAsBytes)
}


//IsSuspended checks if the Delegation is Suspended based on the delegation.Suspended
func (s *SmartContract) IsSubSuspended(ctx contractapi.TransactionContextInterface, pck string) bool {
	//we pull from the world state the data for the delegation
	subdelegation, _ := s.IsSubDelegation(ctx, pck)
	
	//checking if a previous delegation has been suspended
	for _, x := range subdelegation.DelegationChain{
		temp, _ := s.IsDelegation(ctx, x)
		if temp.Suspended == true {
			return true
		} 
	}
	
	return false 
}


//IsRevoked checks if the Delegation is Revoked based on the delegation.Revoked 
func (s *SmartContract) IsSubRevoked(ctx contractapi.TransactionContextInterface, pck string) bool {
	//we pull from the world state the data for the delegation
	subdelegation, _ := s.IsSubDelegation(ctx, pck)

	//checking if a previous delegation has been revoked
	for _, x := range subdelegation.DelegationChain{
		temp, _ := s.IsDelegation(ctx, x)
		if temp.Revoked == true {
			return true
		} 
	}
	
	return false 
}


//Isvalid checks if the Delegation is valid based on the delegation.Expiry and delegation.Issue timestamp
func (s *SmartContract) IsSubValid(ctx contractapi.TransactionContextInterface, pck string) bool {
	//we pull from the world state the data for the delegation
	subdelegation, _ := s.IsSubDelegation(ctx, pck)

	//we pull the current system time to check if it surpasses the Expired field of the delegation
	timenow := uint64(time.Now().Unix())

	if subdelegation.Expiry > timenow && subdelegation.Issue < timenow {
		//checking if a previous delegation has been revoked or suspended
		for _, x := range subdelegation.DelegationChain{
			temp, _ := s.IsDelegation(ctx, x)
			if temp.Suspended == true {
				return false
			} else if temp.Revoked == true {
				return false
			}
		}
		return true 
	}	else {
		return false 
	}
}	
	
//IsExpired checks if the Delegation has expired based on the delegation.Expiry timestamp
func (s *SmartContract) IsSubExpired(ctx contractapi.TransactionContextInterface, pck string) (bool,error) {
	//we pull from the world state the data for the delegation
	subdelegation, _ := s.IsSubDelegation(ctx, pck)

	//we pull the current system time to check if it surpasses the Expired field of the delegation
	timenow := uint64(time.Now().Unix())

	//we check if this delegation is Expired and we return true or false 
	if subdelegation.Expiry > timenow {
		return false, nil 
	} else {
		if subdelegation.Suspended == false && subdelegation.Revoked == false {
			//we update the Suspended field 
			subdelegation.Suspended = true
	
			//we update the subdel field to the privious delegation 
			previousdelegation := subdelegation.DelegationChain[len(subdelegation.DelegationChain)-2]
			temppreviousdelegation, _ := s.IsDelegation(ctx, previousdelegation)
		
	
	
			//updating the subdel field on the previous delegation 
			if temppreviousdelegation.Type == "D" {
				newsubdel:=temppreviousdelegation.Subdel + 1 + subdelegation.Subdel
				temppreviousdelegation.Subdel = newsubdel 
				temppreviousdelegationAsBytes, _ := json.Marshal(temppreviousdelegation)
				ctx.GetStub().PutState(temppreviousdelegation.Pck, temppreviousdelegationAsBytes)
			} else {
				tempdelegation, _ := s.IsSubDelegation(ctx, previousdelegation)
				tempdelegation.Subdel = tempdelegation.Subdel + 1 + subdelegation.Subdel
				tempdelegationAsBytes, _ := json.Marshal(tempdelegation)
				ctx.GetStub().PutState(tempdelegation.Pck, tempdelegationAsBytes)
			}
	
	
			//we store back to the world state
			subdelegationAsBytes, _ := json.Marshal(subdelegation)
			ctx.GetStub().PutState(pck, subdelegationAsBytes)
		}
				
		return true, nil 
	}
}


//for the charging we use the ChargingDel function
//to check if a subdelegation is expired we can use the IsExpired function  

//IsSubDelegation returns the subdelegation stored in the world state with given Pck (Key)
func (s *SmartContract)IsSubDelegation(ctx contractapi.TransactionContextInterface, pck string) (*SubDelegation, error) {
	//we pull from the world state the data for the subdelegation
	subdelegationAsBytes, err := ctx.GetStub().GetState(pck)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if subdelegationAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", pck)
	}

	subdelegation := new(SubDelegation)
	_ = json.Unmarshal(subdelegationAsBytes, subdelegation)

	return subdelegation, nil
}


//--------------------------------------End Of SubDelegation Management------------------------------




//***************************************************************************************************
//**																							   **
//**							The following section Manages the Delegations                      ** 
//**                                                                                               **
//***************************************************************************************************


//----------------------------------Delegation Management--------------------------------------------
//in this section there are the basic functions to manage Delegations  

//RegisterDelegation adds a new Delegation to the world state with given details
func (s *SmartContract) RegisterDelegation(ctx contractapi.TransactionContextInterface, pck string, grandor string, recipient string, subdel string, issue string, expiry string) error {
	
	//getting the service data from the world state 
	service, err := s.IsService(ctx, grandor)
	
	if err != nil {
		return err
	}

	recipientcheck,err := s.IsService(ctx, recipient)
	
	if err != nil {
		return err
	}
	
	if grandor == recipient {
		return fmt.Errorf("Cannot Self-Delegate")
	}
	
	//if issue > expiry {
	//	return fmt.Errorf("Delegation Issue and Expiry times should be checked again")
	//}
	
	expiryint64, _ := strconv.ParseUint(expiry, 10, 64)
	timenow := uint64(time.Now().Unix())
	if timenow > expiryint64 {
		return fmt.Errorf("Cannot create a Delegation which is already Expired")
	}
	
	//if everything above is ok we procceed with the checking and creation of the Delegation 
	if service.Type == "S" && service.Registered == true  && recipientcheck.Registered == true {
		
	//turning string to uint8
	tempsubdel, err := strconv.ParseUint(subdel, 10, 8) 

	if err != nil {
		fmt.Printf("Error")
		}

	subdel1 :=uint8(tempsubdel)

	//turning string to uint64
	issue1, err := strconv.ParseUint(issue, 10, 64)

	if 	err != nil {
		fmt.Printf("Error")
		}

	expiry1, err := strconv.ParseUint(expiry, 10, 64)

	if err != nil {
		fmt.Printf("Error")
		}
		
	if issue1 > expiry1 {
		return fmt.Errorf("Delegation Issue and Expiry times should be checked again, Issue cannot be after the expiry")
	}

	var finalrevokers []string
	finalrevokers = append(finalrevokers,grandor)
	finalrevokers = append(finalrevokers,recipient)

	var tempdelegationchain []string
	tempdelegationchain = append(tempdelegationchain,pck)

	//we pass the data to the struct and then we save to the world state 
	delegation := Delegation{
		Pck:				pck,	 	
		Grandor:			grandor,   			 
		Recipient:  		recipient, 	 
		Subdel:       		subdel1, 	 
		Issue: 				issue1,   
		Expiry: 			expiry1,   
		Suspended:			false,
		Revoked:			false,
		Revokers: 			finalrevokers,
		DelegationChain:	tempdelegationchain,
		Type:				"D",
	}

	delegationAsBytes, _ := json.Marshal(delegation)

	return ctx.GetStub().PutState(pck, delegationAsBytes)
	
	} else {
		return fmt.Errorf("Grandor Or Reciepient Error")
		 
	}
}


//function to update the world state with the new Suspended status, true when the delegation has been suspended
func (s *SmartContract) SuspendDelegation(ctx contractapi.TransactionContextInterface, pck string) error {
	//we pull from the world state the data for the delegation
	delegation, err := s.IsDelegation(ctx, pck)

	if err != nil {
		return err
	}

	//we update the Suspended field 
	delegation.Suspended = true

	//we store back to the world state
	delegationAsBytes, _ := json.Marshal(delegation)

	return ctx.GetStub().PutState(pck, delegationAsBytes)
}


//function to update the world state with the new Revoke status, true when the delegation has been revoked
func (s *SmartContract) RevokeDelegation(ctx contractapi.TransactionContextInterface, pck string, revoker string ) error {
	//we pull from the world state the data for the delegation
	delegation, err := s.IsDelegation(ctx, pck)

	if err != nil {
		return err
	}

	if stringInSlice(revoker, delegation.Revokers){
			//we update the Revoked field 
			delegation.Revoked = true
			} else {
				return fmt.Errorf("%s is not an authorized Revoker", revoker)
		}

	//we store back to the world state
	delegationAsBytes, _ := json.Marshal(delegation)

	return ctx.GetStub().PutState(pck, delegationAsBytes)
}


//IsSuspended checks if the Delegation is Suspended based on the delegation.Suspended
func (s *SmartContract) IsSuspended(ctx contractapi.TransactionContextInterface, pck string) (bool,error) {
	//we pull from the world state the data for the delegation
	delegation, _ := s.IsDelegation(ctx, pck)

	//we check if this delegation is Suspended and we return true or false 
	if delegation.Suspended == true {
		return true, nil 
	} else {
		return false, nil 
	}
}


//IsRevoked checks if the Delegation is Revoked based on the delegation.Revoked 
func (s *SmartContract) IsRevoked(ctx contractapi.TransactionContextInterface, pck string) (bool,error) {
	//we pull from the world state the data for the delegation
	delegation, _ := s.IsDelegation(ctx, pck)

	//we check if this delegation is Revoked and we return true or false 
	if delegation.Revoked == true {
		return true, nil 
	} else {
		return false, nil 
	}
}


//Isvalid checks if the Delegation is valid based on the delegation.Expiry and delegation.Issue timestamp
func (s *SmartContract) IsValid(ctx contractapi.TransactionContextInterface, pck string) (bool,error) {
	//we pull from the world state the data for the delegation
	delegation, _ := s.IsDelegation(ctx, pck)

	//we pull the current system time to check if it surpasses the Expired field of the delegation
	timenow := uint64(time.Now().Unix())

	//we check if this delegation is Valid and we return true or false 
	if delegation.Expiry > timenow && delegation.Issue < timenow && delegation.Suspended == false && delegation.Revoked == false{
		return true, nil 
	} else {
		return false, nil 
	}
}


//IsExpired checks if the Delegation has expired based on the delegation.Expiry timestamp
func (s *SmartContract) IsExpired(ctx contractapi.TransactionContextInterface, pck string) (bool,error) {
	//we pull from the world state the data for the delegation
	delegation, _ := s.IsDelegation(ctx, pck)

	//we pull the current system time to check if it surpasses the Expired field of the delegation
	timenow := uint64(time.Now().Unix())

	//we check if this delegation is Expired and we return true or false 
	if delegation.Expiry > timenow {
		return false, nil 
	} else {
		return true, nil 
	}
}


//function to charge the Delegations
func (s *SmartContract) ChargingDel(ctx contractapi.TransactionContextInterface, pck string, ncores string ) (uint64, error) {
	//we pull from the world state the data for the subdelegation
	delegation, err:= s.IsDelegation(ctx, pck)
	if err != nil {
		return 0, fmt.Errorf("There is no such Delegation")
	}

	var chargetime  uint64
	var hours	    uint64
	var costperhour uint64
	var singlecost  uint64
	var totalcost   uint64

	//turning string to uint64 for issue 
	ncores1, _ := strconv.ParseUint(ncores, 10, 64)

	//we pull the current system time to check if it surpasses the Expired field of the delegation and if yes the delegation has started to charge 
	timenow := uint64(time.Now().Unix())

	if delegation.Issue > timenow {
		return 0, nil
	} else if delegation.Issue < timenow && delegation.Expiry > timenow {
		chargetime = timenow - delegation.Issue
		hours = chargetime/3600
		costperhour = 2
		singlecost = hours * costperhour
		totalcost = singlecost * ncores1
		return totalcost, nil
	} else {
		chargetime = delegation.Expiry - delegation.Issue
		hours = chargetime/3600
		costperhour = 2
		singlecost = hours * costperhour
		totalcost = singlecost * ncores1
		return totalcost, nil
	}
}


//IsDelegation returns the delegation stored in the world state with given Pck (Key)
func (s *SmartContract)IsDelegation(ctx contractapi.TransactionContextInterface, pck string) (*Delegation, error) {
	//we pull from the world state the data for the delegation
	delegationAsBytes, err := ctx.GetStub().GetState(pck)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if delegationAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", pck)
	}

	delegation := new(Delegation)
	_ = json.Unmarshal(delegationAsBytes, delegation)

	return delegation, nil
}


//--------------------------------------End Of Delegation Management---------------------------------





//***************************************************************************************************
//**																							   **
//**							The following section Manages the Service Accounts                 ** 
//**                                                                                               **
//***************************************************************************************************


//----------------------------------Service Account Management-------------------------------------
//in this section there are the basic functions to manage Services 

//Register_Service creates a service and adds its info and the tenant owner of it in the world state
func (s *SmartContract) Register_Service(ctx contractapi.TransactionContextInterface, pck string, name string ) error {
	//matching the data given with the service fields 
	service := Service{
		Pck: 	    		  pck,
		Name: 	   			  name,
		Registered:			  true,
		Type:				  "S",
	}

	//storing to the world state 
	serviceAsBytes, _ := json.Marshal(service)

	return ctx.GetStub().PutState(pck, serviceAsBytes)
}


//UnRegister_Service updates the Registered field of a service in the world state  
func (s *SmartContract) UnRegister_Service(ctx contractapi.TransactionContextInterface, pck string) error {
	//getting the service data from the world state 
	service, err := s.IsService(ctx, pck)

	if err != nil {
		return err
	}

	//updating the Registered field of the service 
	service.Registered = false

	//storing the updated data back to the world state 
	serviceAsBytes, _ := json.Marshal(service)

	return ctx.GetStub().PutState(pck, serviceAsBytes)
}


//IsService returns the service stored in the world state with given Pck (Key)
func (s *SmartContract)IsService(ctx contractapi.TransactionContextInterface, pck string) (*Service, error) {
	//geting the service data from the world state based on the pck 
	serviceAsBytes, err := ctx.GetStub().GetState(pck)

	if err != nil {
		//return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
		return nil, fmt.Errorf("Failed to read from world state. No such Service Exists")
	}

	if serviceAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", pck)
	}

	service := new(Service)
	_ = json.Unmarshal(serviceAsBytes, service)

	return service, nil
}


//--------------------------------------End Of Service Account Management----------------------------





//***************************************************************************************************
//**																							   **
//**							The following section Manages the Tenant Accounts                  ** 
//**                                                                                               **
//***************************************************************************************************


//-------------------------------------Tenant Account Management-------------------------------------
//in this section there are the basic function to manage (Enroll, Update, Destroy, and Search) a Tenant 


//Enroll adds a new tenant to the world state with given details
func (s *SmartContract) Enroll(ctx contractapi.TransactionContextInterface, pck string, name string, email string, phone string) error {
	//matching the given data to the tenant fields 
	tenant := Tenant{
		Pck: 	   pck,
		Name: 	   name,
		Email:     email,
		Phone: 	   phone,
		Registered: true,
		Type:      "T",
	}

	//storing to the world state based on the pck 
	tenantAsBytes, _ := json.Marshal(tenant)
	return ctx.GetStub().PutState(pck, tenantAsBytes)
}


//Update function updates the info of a tenant with new info in world state 
func (s *SmartContract) Update(ctx contractapi.TransactionContextInterface, tenantNumber string, newName string, newEmail string, newPhone string) error {
	//getting the data from the world state 
	tenant, err := s.IsTenant(ctx, tenantNumber)
	
	if err != nil {
		return err
	}
		
	//updating tenant info, all fields 
	tenant.Name = newName
	tenant.Email = newEmail
	tenant.Phone = newPhone

	//storing back to the world state the updated info 
	tenantAsBytes, _ := json.Marshal(tenant)

	return ctx.GetStub().PutState(tenantNumber, tenantAsBytes)
}


//DestroyTenant updated the Registered field of a tenant with given Pck (Key) in world state 
func (s *SmartContract) DestroyTenant(ctx contractapi.TransactionContextInterface, pck string) error {
	//getting the data from the world state based on the pck
	tenant, err := s.IsTenant(ctx, pck)

	if err != nil {
		return err
	}

	//updating the Registered field for the specific tenant 
	tenant.Registered = false 

	//storing the updated data back to the world state 
	tenantAsBytes, _ := json.Marshal(tenant)

	return ctx.GetStub().PutState(pck, tenantAsBytes)
}


//IsTenant returns the tenant stored in the world state with given Pck (Key)
func (s *SmartContract)IsTenant(ctx contractapi.TransactionContextInterface, pck string) (*Tenant, error) {
	//getting the data from world state based on the pck 
	tenantAsBytes, err := ctx.GetStub().GetState(pck)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if tenantAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", pck)
	}

	tenant := new(Tenant)
	_ = json.Unmarshal(tenantAsBytes, tenant)

	return tenant, nil
}


//---------------------------------------End Of Tenant Account Management----------------------------





//***************************************************************************************************
//**																							   **
//**						      	The End of The Management Functions 		                   ** 
//**                                                                                               **
//***************************************************************************************************





//---------------------------------------Main Func---------------------------------------------------

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create Saranyu chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Saranyu chaincode: %s", err.Error())
	}
}
//---------------------------------------END OF CODE--------------------------------------------------


//***************************************************************************************************
//**																							   **
//**						      	Code Created By Athanasios G Kostas  		                   ** 
//**                                                                                               **
//***************************************************************************************************

