/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Developed by IT People Corporation 2017
Some portions of the code were developed by IT People as part
of its set of development accelerators
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of hard-coding.

import (
	"encoding/json"
	"errors"
	"fmt"
	"itpProcUtils"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

////////////////////////////////////////////////////////////////////////////
// The SupplierContract is the agreement between the Buyer and the Supplier
// It is received as a JSON File from the Buyer once both Supplier and Buyer have
// agreed. THIS IS Sprint 1 version of the struct
////////////////////////////////////////////////////////////////////////////

type SupplierContract struct {
	ObjectType           string `json:"doctype"`
	SOWNum               string `json:"sownum"`
	PONumber             string `json:"ponumber"`
	VendorID             string `json:"vendorid"`
	VendorDescription    string `json:"vendorDescription"`
	PODate               string `json:"podate"`
	PaymentTerm          string `json:"paymentTerm"`
	PaymentDays          string `json:"paymentDays"`
	ItemLine1            string `json:"itemLine1"`
	ItemLine2            string `json:"itemLine2"`
	ItemLine3            string `json:"itemLine3"`
	POQuantityLine1      string `json:"poquantityline1"`
	OTQuantityLine2      string `json:"otquantityline2"`
	ExpenseQuantityLine1 string `json:"expenseQuantityLine1"`
	UnitPriceLine1       string `json:"unitPriceLine1"`
	UnitPriceLine2       string `json:"unitPriceLine2"`
	UnitPriceLine3       string `json:"unitPriceLine3"`
	UOMLine1             string `json:"uomline1"`
	UOMLine2             string `json:"uomline2"`
	UOMLine3             string `json:"uomline3"`
	Currency             string `json:"currency"`
	ShortTextLine1       string `json:"shortTextLine1"`
	ShortTextLine2       string `json:"shortTextLine2"`
	ShortTextLine3       string `json:"shortTextLine3"`
	StartDate            string `json:"startDate"`
	EndDate              string `json:"endDate"`
	AgreementNumber      string `json:"agreementNumber"`
	Street               string `json:"street"`
	Postcode             string `json:"postcode"`
	City                 string `json:"city"`
	Country              string `json:"country"`
	Region               string `json:"region"`
	TimeZone             string `json:"timeZone"`
	RequestContactName   string `json:"requestContactName"`
	CumulativeQuantities string `json:"cumulativeQuantities"`
	Ordered              string `json:"ordered"`
	Delivered            string `json:"delivered"`
	Invoiced             string `json:"invoiced"`
}

////////////////////////////////////////////////////////////////////////////
// The ContractorInfo is the details of the Contractor who will work on the
// project against this PO
// THIS IS Sprint 2 version of the struct
////////////////////////////////////////////////////////////////////////////
type ContractorInfo struct {
	ObjectType          string `json:"doctype"`
	PONumber            string `json:"ponumber"`     // Key
	VendorID            string `json:"vendorid"`     // Key
	ContractorID        string `json:"contractorid"` // Key
	Email               string `json:"email"`
	ContractorName      string `json:"contractorname"`
	SupplierName        string `json:"suppliername"`
	Manager             string `json:"manager"`
	BudgetHours         string `json:"budgetedhours"`
	BudgetedRate        string `json:"budgetedrate"`
	OTHours             string `json:"othours"`
	OTRate              string `json:"otrate"`
	Expenses            string `json:"expenses"`
	LocationWorked      string `json:"locationworked"`
	Activity            string `json:"activity"`
	ActivityDescription string `json:"activitydescription"`
	BillCode            string `json:"billcode"`
	StartDate           string `json:"startdate"`
	EndDate             string `json:"enddate"`
	WorkItem            string `json:"workitem"`
	WorkItemDescription string `json:"workitemdescription"`
	SamEmail            string `json:"samemail"`
	PmEmail             string `json:"pmemail"`
	SamAction           string `json:"samaction"`
	PmAction            string `json:"pmAction"`
	SamActionTime       string `json:"samActionTime"`
	PmActionTime        string `json:"pmactionTime"`
}

///////////////////////////////////////////////////////////////////////////////////////////
// Refinement post Sprint 3
// Mohan
///////////////////////////////////////////////////////////////////////////////////////////
type CtrPODetails struct {
	ItemLine1            string `json:"itemLine1"`            // Item Line 1
	ItemLine2            string `json:"itemLine2"`            // Item Line 2
	ItemLine3            string `json:"itemLine3"`            // Item Line 3
	POQuantityLine1      string `json:"poquantityline1"`      // PO Quantity (Line1)
	OTQuantityLine2      string `json:"otquantityline2"`      // OT Quantity (Line2)
	ExpenseQuantityLine1 string `json:"expenseQuantityLine3"` // Expense Quantity (Line3)
	UnitPriceLine1       string `json:"unitPriceLine1"`       // Nett Unit Price (Line 1)
	UnitPriceLine2       string `json:"unitPriceLine2"`       // Nett Unit Price (Line 2)
	UnitPriceLine3       string `json:"unitPriceLine3"`       // Nett Unit Price (Line 3)
	UOMLine1             string `json:"uomline1"`             // UOM (Line 1)
	UOMLine2             string `json:"uomline2"`             // UOM (Line 2)
	UOMLine3             string `json:"uomline3"`             // UOM (Line 3)
	Currency             string `json:"currency"`             // Currency
	ShortTextLine1       string `json:"shortTextLine1"`       // Short Text (Line 1)
	ShortTextLine2       string `json:"shortTextLine2"`       // Short Text (Line 2)
	ShortTextLine3       string `json:"shortTextLine3"`       // Short Text (Line 3)
}

///////////////////////////////////////////////////////////////////////////
// The Workflow Object has to be refined for Sprint 4
///////////////////////////////////////////////////////////////////////////
type WorkflowInfo struct {
	ObjectType     string `json:"docType"`
	PONumber       string `json:"ponumber"`     // Key
	FunctionName   string `json:"functionname"` //Key
	Action         string `json:"action"`
	ActionType     string `json:"actiontype"`
	ActionEndpoint string `json:"actionendpoint"`
}

///////////////////////////////////////////////////////////////////////////
// The TimeSheetRules  Object has to be refined for Sprint 4
// Currently, ruiles are minimal and can be read off the ContractorInfo
///////////////////////////////////////////////////////////////////////////
type TimeSheetRules struct {
	MaxHours       string   `json:"maxhours"`         // Cannot find this - Added by Archana
	MaxOTHours     string   `json:"maxothours"`       // Cannot find this - Added by Archana
	WorkHrsPerDay  string   `json:"workhoursperday"`  // Available in ContractorInfo
	WorkHrsPerWeek string   `json:"workhoursperweek"` // Available in ContractorInfo
	OTHourRule     []string `json:"othourrule"`       // More Details Required
	WorkWeekStart  string   `json:"workweekstart"`    // For Example Sat -> Fri
	PayRollStart   string   `json:"payrollstart"`     // start of payroll cycle eg: Sa for Saturday
}

////////////////////////////////////////////////////////////////////
// ObjectType : TimeSheet
// The UI Will display the segment of the timesheet for the period
// Use ProjectID = "DEFAULT" until we know where it is stored
// TimeSheet -  Keys  Year, PONUM, BUYERCTRID, SupplierCtrID
////////////////////////////////////////////////////////////////////
type TimeSheet struct {
        ObjectType     string                     `json:"doctype"` // TIMESHEET
        Year           string                     `json:"year"`
        PONumber       string                     `json:"ponum"`
        BuyerCtrID     string                     `json:"buyerctrid"` // Key
        SupplierCtrID  string                     `json:"supplierctrid"`
        Assignment     string                     `json:"assignment"`
        ProjectID      string                     `json:"projectid"`
        ContractorName string                     `json:"contractorname"`
        LocationWorked string                     `json:"locationworked"`
        BuyerID        string                     `json:"buyerid"`
        SupplierID     string                     `json:"supplierid"`
        TimeEntries    map[string]TimeEntryForDay `json:"timesheetentry"`
}

// Task ID has not been defined
type TimeEntryForDay struct {
        MonthDay      string `json:"monthday"`
        Activity      string `json:"activity"`
        Description   string `json:"comments"`
        WorkedHours   map [string] string   `json:"workedhours"` // taskid, hours
        TravelHours   string `json:"travelhours"`
        OTHours       string `json:"othours"`
        VacationHours string `json:"vacationhours"` //Worked + Vacation = 8hrs, OThrs = 0
        Expense       string `json:"expense"`
}

type Approval struct {
	ObjectType          string `json:"doctype"`    // APPROVAL
	PONumber            string `json:"ponum"`      // Key
	BuyerCtrID          string `json:"buyerctrid"` // Key
	SupplierCtrID       string `json:"supplierctrid"`
	ApprovalWeekBegin   string `json:"approvalweekbegin"` // The Saturday in case of INFINITE
	ApprovalStatus      string `json:"approvalstatus"`    // submitted, approved or rejected
	BuyerApprover       string `json:"buyerapprover"`
	BuyerApprovalStatus string `json:"buyerapprovalstatus"` //SAM or PM ID
	BuyerApprovaldate   string `json:"approvaltimestamp"`
	BuyerSignature      []byte `json:"buyersignature"`
	SupApprover         string `json:"supapprover"`
	SupApprovalStatus   string `json:"supapprovalstatus"` //SAM or PM ID
	SupApprovaldate     string `json:"supapprovaltimestamp"`
	SupSignature        []byte `json:"supsignature"`
}

////////////////////////////////////////////////////////////////////////
// Read the Timesheet Object and Create view for the week requested
// The view is a slice of the full timesheet
// It is returned at request whenever the Contractor wants to enter
// Time. The TimesheetView is a transient Object
// This Object is not persisted
////////////////////////////////////////////////////////////////////////
type WeeklyTimeSheetView struct {
    ObjectType     string                     `json:"doctype"`
	TSDate         string                     `json:"tsdate"`
	PONumber       string                     `json:"ponum"`
	BuyerCtrID     string                     `json:"buyerctrid"` // Key
	SupplierCtrID  string                     `json:"supplierctrid"`
	Assignment     string                     `json:"assignment"`
	ProjectID      string                     `json:"projectid"`
	ContractorName string                     `json:"contractorname"`
	LocationWorked string                     `json:"locationworked"`
	BuyerID        string                     `json:"buyerid"`
	SupplierID     string                     `json:"supplierid"`
	Status         string                     `json:"tsstatus"`
	TimeEntries    map[string]TimeEntryForDay `json:"timesheetentry"`
}

// ProcurementChaincode example simple Chaincode implementation
type ProcurementChaincode struct {
}

//////////////////////////////////////////////////////////
// JSON To args[] - return a map of the JSON string
//////////////////////////////////////////////////////////
func JSONtoArgs(Avalbytes []byte) (map[string]interface{}, error) {

	var data map[string]interface{}

	if err := json.Unmarshal(Avalbytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}

//////////////////////////////////////////////////////////
// Variation of the above - return value from a JSON string
//////////////////////////////////////////////////////////

func GetKeyValue(Avalbytes []byte, key string) string {
	var dat map[string]interface{}
	if err := json.Unmarshal(Avalbytes, &dat); err != nil {
		panic(err)
	}

	val := dat[key].(string)
	return val
}

//////////////////////////////////////////////////////////
// Converts JSON String to an SupplierContract Object
//////////////////////////////////////////////////////////
func JSONtoSupplierContract(data []byte) (SupplierContract, error) {
	fmt.Println("JSONtoSupplierContract Init")
	sc := SupplierContract{}
	err := json.Unmarshal([]byte(data), &sc)
	if err != nil {
		fmt.Println("Unmarshal failed : ", err)
	}
	fmt.Println("**** Supplier Contract ****", sc)
	return sc, err
}

//////////////////////////////////////////////////////////
// Converts an SupplierContract Object to a JSON String
//////////////////////////////////////////////////////////
func SupplierContracttoJSON(ar SupplierContract) ([]byte, error) {

	ajson, err := json.Marshal(ar)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ajson, nil
}

//////////////////////////////////////////////////////////
// Converts JSON String to an ContractorInfo Object
//////////////////////////////////////////////////////////
func JSONtoCtrInfo(data []byte) (ContractorInfo, error) {

	ci := ContractorInfo{}
	err := json.Unmarshal([]byte(data), &ci)
	if err != nil {
		fmt.Println("Unmarshal failed : ", err)
	}

	return ci, err
}

//////////////////////////////////////////////////////////
// Converts an Contractor information Object to a JSON String
//////////////////////////////////////////////////////////
func CtrInfotoJSON(ci ContractorInfo) ([]byte, error) {

	ajson, err := json.Marshal(ci)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ajson, nil
}

//////////////////////////////////////////////////////////
// Converts JSON String to an WorkfkowInfo Object
//////////////////////////////////////////////////////////
func JSONtoWorkflowInfo(data []byte) (WorkflowInfo, error) {

	wfi := WorkflowInfo{}
	err := json.Unmarshal([]byte(data), &wfi)
	if err != nil {
		fmt.Println("Unmarshal failed : ", err)
	}

	return wfi, err
}

//////////////////////////////////////////////////////////
// Converts TimeSheet information Object to a JSON String
//////////////////////////////////////////////////////////
func TimeSheettoJSON(ts TimeSheet) ([]byte, error) {

	ajson, err := json.Marshal(ts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ajson, nil
}

//////////////////////////////////////////////////////////
// Converts JSON String to an TimeSheet Object
//////////////////////////////////////////////////////////
func JSONtoTimeSheet(data []byte) (TimeSheet, error) {

	ts := TimeSheet{}
	err := json.Unmarshal([]byte(data), &ts)
	if err != nil {
		fmt.Println("Unmarshal failed : ", err)
	}

	return ts, err
}

//////////////////////////////////////////////////////////
// Converts TimeSheetView information Object to a JSON String
//////////////////////////////////////////////////////////
func TimeSheetViewtoJSON(ts WeeklyTimeSheetView) ([]byte, error) {

	ajson, err := json.Marshal(ts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ajson, nil
}

//////////////////////////////////////////////////////////
// Converts JSON String to a TimeSheetView Object
//////////////////////////////////////////////////////////
func JSONtoTimeSheetView(data []byte) (WeeklyTimeSheetView, error) {

	ts := WeeklyTimeSheetView{}
	err := json.Unmarshal([]byte(data), &ts)
	if err != nil {
		fmt.Println("Unmarshal failed : ", err)
	}

	return ts, err
}

//////////////////////////////////////////////////////////
// Converts JSON String to a Approval Object
//////////////////////////////////////////////////////////
func JSONtoApproval(data []byte) (Approval, error) {

	ci := Approval{}
	err := json.Unmarshal([]byte(data), &ci)
	if err != nil {
		fmt.Println("Unmarshal failed : ", err)
	}

	return ci, err
}

//////////////////////////////////////////////////////////
// Converts Approval Object to a JSON String
//////////////////////////////////////////////////////////
func ApprovaltoJSON(ci Approval) ([]byte, error) {

	ajson, err := json.Marshal(ci)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ajson, nil
}

////////////////////////////////////////////
//TSWEEKBEGIN()
// Provides the offset from the week beginning
// provided the functions knows date
// and the week beginning day symbol
// dayIndex("Sa", "Th")
////////////////////////////////////////////
func dayIndex(d1 string, d2 string) int {
	xtr := map[string][]string{
		"Mo": {"Su", "Sa", "Fr", "We", "Th", "Tu"},
		"Tu": {"Mo", "Su", "Sa", "Fr", "Th", "We"},
		"We": {"Tu", "Mo", "Su", "Sa", "Fr", "Th"},
		"Th": {"We", "Tu", "Mo", "Su", "Sa", "Fr"},
		"Fr": {"Th", "We", "Tu", "Mo", "Su", "Sa"},
		"Sa": {"Fr", "Th", "We", "Tu", "Mo", "Su"},
		"Su": {"Sa", "Fr", "Th", "We", "Tu", "Mo"},
	}

	row := xtr[d1]
	// fmt.Println(row)

	for i, each := range row {
		if each == d2 {
			return i
		}
	}

	return -1
}

////////////////////////////////////////////////////////////////
// Determine if a year is a leap year
////////////////////////////////////////////////////////////////
func IsLeapYear(y int) bool {

	// convert int to Time - use the last day of the year, which is 31st December
	year := time.Date(y, time.December, 31, 0, 0, 0, 0, time.Local)
	days := year.YearDay()

	if days > 365 {
		return true
	} else {
		return false
	}
}

func getDaysInMonth(year int, month int) int {

	switch month {
	case
		1, 3, 5, 7, 8, 10, 12:
		return 31
	case 2:
		if IsLeapYear(year) == true {
			return 29
		} else {
			return 28
		}

	}
	return 30
}

func TSWeekBegin(d1 string, dt string) (string, string) {

	// format string date to date
	t, _ := time.Parse("2006-01-02", dt)
	x := fmt.Sprintf("%s", t.Weekday())
	z := dayIndex(x[:2], d1) + 1
	_, week := t.ISOWeek()
	weekString := fmt.Sprintf("%02d", week)
	newDate := t.AddDate(0, 0, -1*z)
	return newDate.Format("20060102"), weekString
}


func NextCalendarDay(currentDate string) string {
	t, _ := time.Parse("20060102", currentDate)
	t = t.AddDate(0, 0, 1)
	return t.Format("20060102")
}

func TSWeekEnd(currentDate string) string {
	t, _ := time.Parse("20060102", currentDate)
	t = t.AddDate(0, 0, 6)
	return t.Format("20060102")
}

func CompareDate(day1, day2 string) bool {
	t1, _ := time.Parse("20060102", day1)
	t2, _ := time.Parse("20060102", day2)
	return t1.Before(t2) || t1.Equal(t2)
}

//////////////////////////////////////////////////////////////////////////////////
// Initialize Chaincode at Deploy Time
//////////////////////////////////////////////////////////////////////////////////

func (t *ProcurementChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ProcurementChaincode Init")
	// TODO - Update Roles, Responsibilities, Permissions
	// TODO - Update permissions

	return shim.Success(nil)
}

//////////////////////////////////////////////////////////////////////////////////
// Invoke Chaincode functions as requested by the Invoke Function
// In fabric 1.0 both Invoke and Query Requests are handled by Invoke
//////////////////////////////////////////////////////////////////////////////////

func (t *ProcurementChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ProcurementChaincode Invoke")
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Function: ", function, "Args: ", args)
	if function == "RecordSupplierContractInfo" {
		// Initialize the SupplierContract
		return t.RecordSupplierContractInfo(stub, args)
	} else if function == "RecordContractorInfo" {
		// Recird Contractor Info
		return t.RecordContractorInfo(stub, args)
	} else if function == "UpdateContractorInfo" {
		// Update Supplier Contract
		return t.UpdateContractorInfo(stub, args)
	} else if function == "QuerySupplierContract" {
		// Retrieve Supplier Contract
		return t.QuerySupplierContract(stub, args)
	} else if function == "QuerySupplierContractList" {
		// Retrieve Supplier Contract
		return t.QuerySupplierContractList(stub, args)
	} else if function == "QueryContractorInfo" {
		// Retrieve Supplier Contractor info
		return t.QueryContractorInfo(stub, args)
	} else if function == "RecordWorkflowInfo" {
		// Record Notification or Workflow info
		return t.RecordWorkFlowInfo(stub, args)
	} else if function == "QueryWorkflowList" {
		// Retreive the list of Workflow info
		return t.QueryWorkFlowList(stub, args[0])
	} else if function == "GetWeeklyTimeSheet" {
        // GetWeeklyTimeSheet
     	return t.GetWeeklyTimeSheet(stub, args)
     } else if function == "ReceiveWeeklyTimesheet" {
        //   ReceiveWeeklyTimesheet
        return t.ReceiveWeeklyTimesheet(stub, args)
    } else if function == "QueryTimeSheet" {
       //   QueryTimeSheet
       return t.QueryTimeSheet(stub, args)
    } else {
	  fmt.Println("Invoke(): Too bad ---- > Function not found : ", function)
	}

	return shim.Error("Invoke(): Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

////////////////////////////////////////////////////////////////////////////////////////
// Records the SupplierContractInfo
// Keys PONumber and VendorID
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) RecordSupplierContractInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//objectID := "SUPPLIERCONTRACT"
	var err error
	var Avalbytes []byte

	if len(args) != 37 {
		return shim.Error("RecordSupplierContractInfo(): Incorrect number of arguments. Expecting 37")
	}

	// Query the Ledger to ensure Record does not exist
	fmt.Println("RecordSupplierContract : ", args)
	Avalbytes, err = itpProcUtils.QueryObject(stub, "SUPPLIERCONTRACT", args[:3])
	if err != nil {
		fmt.Println("RecordSupplierContractInfo() : Failed ------ ", args)
		return shim.Error("RecordSupplierContractInfo(): Failed to get state")
	}

	if Avalbytes == nil {
		fmt.Println(" args **: ", args)
		Avalbytes, _ := t.RecordSupplierContract(stub, args)
		fmt.Println("**** Avalbytes **: ", Avalbytes)
		me, _ := JSONtoSupplierContract(Avalbytes)
		fmt.Println("RecordSupplierContract **: ", me)
	} else {
		return shim.Error("SupplierContract already exist")
	}

	return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////////////////
// Query the Supplier Contractor Records
//
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) QuerySupplierContractList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	queryString := fmt.Sprintf("{\"selector\":{\"doctype\":\"SUPPLIERCONTRACT\"}}")

	queryResults, err := itpProcUtils.GetQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

////////////////////////////////////////////////////////////////////////////////////////
// Query the Supplier Contractor Records
//
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) QueryWorkFlowList(stub shim.ChaincodeStubInterface, ponumber string) pb.Response {
	//queryString := fmt.Sprintf("{\"selector\":{\"doctype\":\"WORKFLOW\",\"ponumber\":\"%s\"}}", ponumber)
	queryString := fmt.Sprintf("{\"selector\":{\"doctype\":\"WORKFLOW\"}}")
	queryResults, err := itpProcUtils.GetQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)

}

////////////////////////////////////////////////////////////////////////////////////////
// Updates the CONTRACTORINFO
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) UpdateContractorInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	objectType := "CONTRACTORINFO"
	var err error
	var Avalbytes []byte

	if len(args) != 26 {
		return shim.Error("UpdateContractorInfo() : Incorrect number of arguments. Expecting 26")
	}

	newObject := &ContractorInfo{objectType, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15], args[16], args[17], args[18], args[19], args[20], args[21], args[22], args[23], args[24], args[25]}

	fmt.Println("*** docType ****  ", newObject)
	Avalbytes, err = itpProcUtils.QueryObject(stub, "CONTRACTORINFO", args[:2])

	if Avalbytes == nil {
		fmt.Println("Failed to update,Contractor does not exist")
		return shim.Error("Failed to update,Contractor does not exist")
	}

	if err != nil {
		fmt.Println("UpdateContractorinfo() : Failed ****  ", args)
		return shim.Error("Failed to get state")
	}

	oldObject := ContractorInfo{}
	err = json.Unmarshal(Avalbytes, &oldObject) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("oldObject **: ", oldObject)
	if newObject.PONumber != "" {
		oldObject.PONumber = newObject.PONumber
	}
	if newObject.ContractorID != "" {
		oldObject.ContractorID = newObject.ContractorID
	}
	if newObject.VendorID != "" {
		oldObject.VendorID = newObject.VendorID
	}
	if newObject.Email != "" {
		oldObject.Email = newObject.Email
	}

	if newObject.ContractorName != "" {
		oldObject.ContractorName = newObject.ContractorName
	}
	if newObject.SupplierName != "" {
		oldObject.SupplierName = newObject.SupplierName
	}
	if newObject.Manager != "" {
		oldObject.Manager = newObject.Manager
	}
	if newObject.BudgetHours != "" {
		oldObject.BudgetHours = newObject.BudgetHours
	}
	if newObject.BudgetedRate != "" {
		oldObject.BudgetedRate = newObject.BudgetedRate
	}
	if newObject.OTHours != "" {
		oldObject.OTHours = newObject.OTHours
	}
	if newObject.OTRate != "" {
		oldObject.OTRate = newObject.OTRate
	}
	if newObject.Expenses != "" {
		oldObject.Expenses = newObject.Expenses
	}
	if newObject.LocationWorked != "" {
		oldObject.LocationWorked = newObject.LocationWorked
	}
	if newObject.Activity != "" {
		oldObject.Activity = newObject.Activity
	}
	if newObject.ActivityDescription != "" {
		oldObject.ActivityDescription = newObject.ActivityDescription
	}
	if newObject.BillCode != "" {
		oldObject.BillCode = newObject.BillCode
	}
	if newObject.StartDate != "" {
		oldObject.StartDate = newObject.StartDate
	}
	if newObject.EndDate != "" {
		oldObject.EndDate = newObject.EndDate
	}
	if newObject.WorkItem != "" {
		oldObject.WorkItem = newObject.WorkItem
	}
	if newObject.WorkItemDescription != "" {
		oldObject.WorkItemDescription = newObject.WorkItemDescription
	}
	if newObject.SamAction != "" {
		oldObject.SamAction = newObject.SamAction
	}
	if newObject.PmActionTime != "" {
		oldObject.PmActionTime = newObject.PmActionTime
	}

	if newObject.SamActionTime != "" {
		oldObject.SamActionTime = newObject.SamActionTime
	}

	if newObject.PmAction != "" {
		oldObject.PmAction = newObject.PmAction
	}
	if newObject.SamEmail != "" {
		oldObject.SamEmail = newObject.SamEmail
	}
	if newObject.PmEmail != "" {
		oldObject.PmEmail = newObject.PmEmail
	}

	scJSONBytes, _ := json.Marshal(oldObject)
	fmt.Println("RecordContractor() : *** scJSONBytes **** ", scJSONBytes)
	fmt.Println("RecordContractor() : CONTRACTORINFO", args[0], args[1])
	err = itpProcUtils.UpdateObject(stub, "CONTRACTORINFO", args[:2], scJSONBytes)

	if err != nil {
		fmt.Println("UpdateContractorinfo() : Failed ****** ", args)
		return shim.Error("UpdateContractorinfo : Failed to get state")
	}

	//TODO avoid updates each time
	workflowArgs := []string{oldObject.PONumber, "RecordContractorInfo", "Notification", "Email", oldObject.SamEmail}
	t.RecordWorkFlowInfo(stub, workflowArgs)

	workflowArgs = []string{oldObject.PONumber, "RecordContractorInfo", "Notification", "Email", oldObject.PmEmail}
	t.RecordWorkFlowInfo(stub, workflowArgs)

	result := t.QueryWorkFlowList(stub, oldObject.PONumber)
	return result
}

////////////////////////////////////////////////////////////////////////////////////////
// Query the Supplier Contractor Record
// In the Ledger args: { PO012345", "INFINITE"}
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) QuerySupplierContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte

	// In future , it should be > 1 and ,= mo_of_keys for object
	if len(args) != 3 {
		return shim.Error("QuerySupplierContract(): Incorrect number of arguments. Expecting 3")
	}

	// itpProcUtils.Query for the information about the PARTY Alice
	// Get the current abalance of the Party from the ledger
	// If nil, then party can be initialized
	fmt.Println("QuerySupplierContract() : Arguments for itpProcUtils.Query: ", args[0], " : ", args[1])
	Avalbytes, err = itpProcUtils.QueryObject(stub, "SUPPLIERCONTRACT", args)
	fmt.Println("QuerySupplierContract() : **** Supplier Contract  Avalbytes ****", Avalbytes)
	me, _ := JSONtoSupplierContract(Avalbytes)
	fmt.Println("QuerySupplierContract() : **** Supplier Contract  me ****", me)

	if err != nil {
		return shim.Error("QuerySupplierContract() : Failed to query object successfully")
	}

	if Avalbytes == nil {
		return shim.Error("QuerySupplierContract() : Supplier Contract Code not found " + args[1])
	}

	fmt.Println("QuerySupplierContract() : Returning results")
	fmt.Println(me)

	return shim.Success(Avalbytes)

}

////////////////////////////////////////////////////////////////////////////////////////
// Query the Supplier Contractoror Record
// In the Ledger args: { "PO123", "INF02345"}
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) QueryContractorInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte

	// In future , it should be > 1 and ,= mo_of_keys for object
	if len(args) != 2 {
		return shim.Error("QueryContractorInfo() : Incorrect number of arguments. Expecting 2")
	}

	// itpProcUtils.Query for the information about the PARTY Alice
	// Get the current abalance of the Party from the ledger
	// If nil, then party can be initialized
	fmt.Println("QueryContractorInfo() : Arguments for itpProcUtils.Query: ", "CONTRACTORINFO : ", args[0], " : ", args[1])
	Avalbytes, err = itpProcUtils.QueryObject(stub, "CONTRACTORINFO", args)
	me, _ := JSONtoSupplierContract(Avalbytes)
	if err != nil {
		return shim.Error("QueryContractorInfo() : Failed to query object successfully")
	}

	if Avalbytes == nil {
		return shim.Error("QueryContractorInfo() : Supplier Contractor Code not found " + args[1])
	}

	fmt.Println(me)
	return shim.Success(Avalbytes)
}

// ========================================================================
// MAIN
// ========================================================================
func main() {
	fmt.Println("Procurement MVP ver S3.4 2017-04-04 - Mohan ")
	err := shim.Start(new(ProcurementChaincode))
	if err != nil {
		fmt.Printf("Proc_MVP: main(): Error starting Procurement chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init Party - create a new Party, store into chaincode state as a JSON record
// ============================================================================================================================
func (t *ProcurementChaincode) RecordSupplierContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 37 {
		return nil, errors.New("RecordSupplierContract() : Incorrect number of arguments. Expecting 37")
	}

	// Create Party (example Alice or Bob)  from inputs
	fmt.Println("SUPPLIERCONTRACT", args[0], args[1])
	objectType := "SUPPLIERCONTRACT"
	sc := &SupplierContract{objectType, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15], args[16], args[17], args[18], args[19], args[20], args[21], args[22], args[23], args[26], args[25], args[26], args[27], args[28], args[29], args[30], args[31], args[32], args[33], args[34], args[35], args[36]}
	scJSONBytes, _ := json.Marshal(sc)

	err := itpProcUtils.UpdateObject(stub, "SUPPLIERCONTRACT", args[:3], scJSONBytes)
	if err != nil {
		fmt.Println("RecordSupplierContract() : Error inserting Object into LedgerState %s", err)
	}
	fmt.Println("RecordSupplierContract() : ", sc)

	fmt.Println("- end init RecordContract")
	return scJSONBytes, err
}

// ============================================================================================================================
// Init ContractorInfo - create a new ContractorInfo, store into chaincode state as a JSON record
// ============================================================================================================================
func (t *ProcurementChaincode) RecordContractor(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 26 {
		return nil, errors.New("RecordContractor() : Incorrect number of arguments. Expecting 26")
	}
	objectType := "CONTRACTORINFO"

	// Create Party (example Alice or Bob)  from inputs
	fmt.Println("CONTRACTORINFO", args[0], args[1])
	sc := &ContractorInfo{objectType, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15], args[16], args[17], args[18], args[19], args[20], args[21], args[22], args[23], args[24], args[25]}
	fmt.Println("RecordContractor() : *** sc **** ", sc)
	scJSONBytes, _ := json.Marshal(sc)
	fmt.Println("RecordContractor() : *** scJSONBytes **** ", scJSONBytes)
	fmt.Println("RecordContractor() : CONTRACTORINFO", args[0], args[1])
	err := itpProcUtils.UpdateObject(stub, "CONTRACTORINFO", args[:2], scJSONBytes)
	if err != nil {
		fmt.Println("RecordContractor() : Error inserting Object into LedgerState %s", err)
	}
	fmt.Println("RecordContractor() : ", sc)

	fmt.Println("- end init RecordContractor()")
	return scJSONBytes, err
}

////////////////////////////////////////////////////////////////////////////////////////
// Initializes the Contractor Info
// Keys PONumber, ContractorID
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) RecordContractorInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	objectType := "CONTRACTORINFO"
	var err error
	var Avalbytes []byte

	if len(args) != 26 {
		return shim.Error("RecordContractorInfo() : Incorrect number of arguments. Expecting 26")
	}

	// Query the Ledger to ensure Record does not exist
	fmt.Println("RecordContractorInfo() : ", args)
	Avalbytes, err = itpProcUtils.QueryObject(stub, "CONTRACTORINFO", args[:2])
	if err != nil {
		fmt.Println("RecordContractorInfo() : Failed ------ ", args)
		return shim.Error("Failed to get state")
	}

	if Avalbytes == nil {
		fmt.Println("RecordContractorInfo() * args *: ", args)
		Avalbytes, _ := t.RecordContractor(stub, args)
		me, _ := JSONtoCtrInfo(Avalbytes)
		fmt.Println("RecordContractorInfo() : me ------ ", me)
	} else {
		return shim.Error("Contractor already exist")
	}

	sc := &ContractorInfo{objectType, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15], args[16], args[17], args[18], args[19], args[20], args[21], args[22], args[23], args[24], args[25]}
	err = json.Unmarshal(Avalbytes, &sc)

	workflowArgs := []string{sc.PONumber, "RecordContractorInfo", "Notification", "Email", sc.SamEmail}
	t.RecordWorkFlowInfo(stub, workflowArgs)
	workflowArgs = []string{sc.PONumber, "RecordContractorInfo", "Notification", "Email", sc.PmEmail}
	t.RecordWorkFlowInfo(stub, workflowArgs)

	result := t.QueryWorkFlowList(stub, sc.PONumber)

	return result
}

////////////////////////////////////////////////////////////////////////////////////////
// Initializes the Contractor Info
// Keys PONumber, ContractorID
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) RecordWorkFlowInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 5 {
		return shim.Error("RecordContractor() : Incorrect number of arguments. Expecting 5")
	}
	objectType := "WORKFLOW"

	// Create Party (example Alice or Bob)  from inputs
	fmt.Println("WORKFLOW", args[0], args[1])
	sc := &WorkflowInfo{objectType, args[0], args[1], args[2], args[3], args[4]}
	fmt.Println("RecordContractor() : *** sc **** ", sc)
	scJSONBytes, _ := json.Marshal(sc)
	fmt.Println("RecordContractor() : *** scJSONBytes **** ", scJSONBytes)
	fmt.Println("RecordContractor() : WORKFLOW", args[0], args[1], args[2], args[3], args[4])
	err := itpProcUtils.UpdateObject(stub, "WORKFLOW", args[:5], scJSONBytes)
	if err != nil {
		fmt.Println("RecordContractor() : Error inserting Object into LedgerState %s", err)
	}
	fmt.Println("RecordContractor() : ", sc)

	fmt.Println("- end init RecordContractor()")

	return shim.Success(scJSONBytes)
}

////////////////////////////////////////////////////////////////////////////////////////
// Updates the WORKFLOW
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) UpdateWorkFlowInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//objectID := "WORKFLOW"
	var err error
	var Avalbytes []byte

	if len(args) != 22 {
		return shim.Error("UpdateWorkFlowInfo() : Incorrect number of arguments. Expecting 22")
	}

	Avalbytes, err = itpProcUtils.QueryObject(stub, "WORKFLOW", args[:5])

	if Avalbytes == nil {
		fmt.Println("Failed to update,workflow does not exist")
		return shim.Error("Failed to update,workflow does not exist")
	}

	if err != nil {
		fmt.Println("UpdateWorkFlowInfo() : Failed ****  ", args)
		return shim.Error("Failed to get state")
	}

	//me, err := JSONtoWorkflowInfo(Avalbytes)
	//fmt.Println("UpdateWorkFlowInfo() **: ", me)
	err = itpProcUtils.UpdateObject(stub, "WORKFLOW", args[:5], Avalbytes)

	if err != nil {
		fmt.Println("UpdateWorkFlowInfo() : Failed ****** ", args)
		return shim.Error("UpdateWorkFlowInfo : Failed to get state")
	}
	//fmt.Println("UpdateWorkFlowInfo() **: ", me)

	return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////////////////
// Query the Supplier Contractor Record
// In the Ledger args: { PO012345", "INFINITE"}
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) QueryWorkFlowInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte

	// In future , it should be > 1 and ,= mo_of_keys for object
	if len(args) != 2 {
		return shim.Error("QuerySupplierContract(): Incorrect number of arguments. Expecting 2")
	}

	// itpProcUtils.Query for the information about the PARTY Alice
	// Get the current abalance of the Party from the ledger
	// If nil, then party can be initialized
	fmt.Println("QuerySupplierContract() : Arguments for itpProcUtils.Query: ", args[0], " : ", args[1])
	Avalbytes, err = itpProcUtils.QueryObject(stub, "WORKFLOW", args)
	fmt.Println("WORKFLOW : **** Supplier Contract  Avalbytes ****", Avalbytes)
	me, _ := JSONtoWorkflowInfo(Avalbytes)
	fmt.Println("QuerySupplierContract() : **** Supplier Contract  me ****", me)

	if err != nil {
		return shim.Error("QuerySupplierContract() : Failed to query object successfully")
	}

	if Avalbytes == nil {
		return shim.Error("QuerySupplierContract() : Supplier Contract Code not found " + args[1])
	}

	fmt.Println("QuerySupplierContract() : Returning results")
	fmt.Println(me)

	return shim.Success(Avalbytes)
}

/////////////////////////////////////////////////////////////////////////////////
// TimeSheet Management Functions
// 1. Create a Brand New Empty Timesheet - Done
// 2. Record the Timesheet (Save it)
// 3. Return a weekly Timesheet (TimesheetView
// 4. Update TimesheetView into the Actual Timesheet
// 5. Query a Timesheet
// 6. Save or Submit a Timesheet - Performed by ReceiveWeeklyTimeSheet
// 7. Record and Query Approval Object
/////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Create New TimeSheet if none exists
// Keys YYYY, PONumber, BuyerCtrID
////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) CreateNewTimeSheet(stub shim.ChaincodeStubInterface, timeSheetView WeeklyTimeSheetView ) ([]byte, error) {
	var te TimeEntryForDay
	var tsEntry TimeSheet // A TimeSheet Record
	mapE := make(map[string]TimeEntryForDay)
   	objectID := "TIMESHEET"

	fmt.Printf("**** Called CreateNewTimeSheet with ******* ", timeSheetView)

	keys := []string{timeSheetView.TSDate[:4],timeSheetView.PONumber,timeSheetView.BuyerCtrID,timeSheetView.SupplierCtrID}

	fmt.Printf("**** QueryObject CreateNewTimeSheet with ******* ", keys)

    Avalbytes, err := itpProcUtils.QueryObject(stub, objectID, keys)
    if err!= nil {
        return nil, err
    }

	if Avalbytes == nil {
        	year, _ := strconv.Atoi(timeSheetView.TSDate[:4])
        	var NumOfDays int

        	for month := 1; month <= 12; month++ {
        		monthString := fmt.Sprintf("%02d", month)
        		if IsLeapYear(year) && month == 2 {
        			NumOfDays = 29
        		} else {
        			NumOfDays = getDaysInMonth(year,month)
        		}
        		for days := 1; days <= NumOfDays; days++ {
        			daysString := fmt.Sprintf("%02d", days)
        			date := monthString + daysString
        			te.MonthDay = date
        			mapE[date] = te
        		}
        	}

         	tsEntry = TimeSheet{objectID,timeSheetView.TSDate[:4], timeSheetView.PONumber, timeSheetView.BuyerCtrID,
        	timeSheetView.SupplierCtrID,timeSheetView.Assignment,timeSheetView.ProjectID,
        	timeSheetView.ContractorName, timeSheetView.LocationWorked, timeSheetView.BuyerID,timeSheetView.SupplierID,mapE}
        	fmt.Printf("TimeSheet Entry : ", tsEntry.PONumber, "\n")
        	Avalbytes, err = TimeSheettoJSON(tsEntry)

        	err = itpProcUtils.UpdateObject(stub, objectID, keys, Avalbytes)

   	}

	return Avalbytes, err

}


////////////////////////////////////////////////////////////////////////////////////////
// Query the Supplier Contractor Records
//
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) QueryTimeSheet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	queryString := fmt.Sprintf("{\"selector\":{\"doctype\":\"TIMESHEET\"}}")

	queryResults, err := itpProcUtils.GetQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)

}

////////////////////////////////////////////////////////////////////////////////////////
// Get Weekly Time Sheet
// The TimeSheet is queried in full  with the keys
// args: YYYY, PONumber, BuyerCtrID, Date or Weekbegin
// ObjectType TIMESHEET, Keys: YYYY, PONumber, BuyerCtrID
// Then Extract the Week from the TimeEntries
////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) GetWeeklyTimeSheet(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var Avalbytes []byte
	ObjectID := "TIMESHEET"
    tsview, err := JSONtoTimeSheetView([]byte(args[0]));
   	// In future , it should be > 1 and = no_of_keys for object
	if err != nil {
	 return shim.Error("GetWeeklyTimeSheet() : Incorrect number of arguments. Expecting 4 - YYYY, PONumber, BuyerCtrID, Date or Weekbegin")
	}

   fmt.Println("ObjectID ", ObjectID)
    Avalbytes, err = t.CreateNewTimeSheet(stub, tsview)

    if err != nil {
       fmt.Println("Error ", err)
       return shim.Error("Error While Creating / Retrieving time sheet")
    }

	queryDate := tsview.TSDate
	weekStartDate, _ := TSWeekBegin("Sa", queryDate)
	weekEndDate := TSWeekEnd(weekStartDate)
	fmt.Println("weekEndDate ", weekEndDate)

	me, err := JSONtoTimeSheet(Avalbytes)
	if err != nil {
	   fmt.Println("Error ", err)
       return shim.Error("Error While Creating / Retrieving time sheet")
    }
	// yearweek := me.Year + weekNumber

	mapE := make(map[string]TimeEntryForDay)

	day := weekStartDate
   fmt.Println(" *** me *******  ", me)
   fmt.Println(" *** me.TimeEntries *******  ", me.TimeEntries)
    for  CompareDate(day, weekEndDate) {
		mapE[day[4:8]] = me.TimeEntries[day[4:8]]
		day = NextCalendarDay(day)
	}

     fmt.Println(" *** mapE *******  ", mapE)

	///////////////////////////////////////////////////////////////////////////////////////////////
	// Is Timesheet Locked -- Check if there is an approval object for that period in the ledger
	// The Client and UI should lock the sheet if the status is "SUBMIT"
	///////////////////////////////////////////////////////////////////////////////////////////////
	tsStatus := "SAVE"
	_, err = itpProcUtils.QueryObject(stub, "APPROVAL", []string{tsview.PONumber, tsview.BuyerCtrID, tsview.SupplierCtrID, weekStartDate})
	if err != nil {
		tsStatus = "SUBMIT"
	}

	///////////////////////////////////////////////////////////////////////////////////////////////
	// Generate the WeeklyTimeSheetView and return back to client
	///////////////////////////////////////////////////////////////////////////////////////////////
	timeSheetView := WeeklyTimeSheetView{ObjectID, me.Year, me.PONumber, me.BuyerCtrID, me.Assignment, me.SupplierCtrID, me.ProjectID, me.ContractorName, me.LocationWorked, me.BuyerID,me.SupplierID, tsStatus, mapE}

	timeSheetViewBytes, _ := TimeSheetViewtoJSON(timeSheetView)

	return shim.Success(timeSheetViewBytes)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Assume that we receive 2 arguments ["SAVE|SUBMIT", "a JSON of the WeeklyTimesheetView"]
// Unmarshall the object
// Retrieve original Timesheet record for the contractor with the Weeklysubmission
// Update the Timesheet Record
// Save the Timesheet record back into the Ledger
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (t *ProcurementChaincode) ReceiveWeeklyTimesheet(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var Avalbytes []byte


	// In future , it should be > 1 and = no_of_keys for object
	/*
	if len(args) != 1 {
		return shim.Error("ReceiveWeeklyTimesheet() : Incorrect number of arguments. Expecting 1 TimeSheetViewJSON")
	}*/

	// We assume that the TimeSheet Submission has 1 arg and that is a WeeklyTimeSheetView JSON
	// Convert the arg to a TimesheetView Object
	fmt.Println("ReceiveWeeklyTimesheet() : Arguments for itpProcUtils.Query: TIMESHEET : ", args[0])
	wtsView, err := JSONtoTimeSheetView([]byte(args[0]))
	if err != nil {
		return shim.Error("ReceiveWeeklyTimesheet() : Failed to convert arg[0] to a TimesheetView object")
	}

	// Query and Retrieve the Full Timesheet for the Contractor
	keys := []string{wtsView.TSDate[:4], wtsView.PONumber, wtsView.BuyerCtrID, wtsView.SupplierCtrID}

	Avalbytes, err = itpProcUtils.QueryObject(stub, "TIMESHEET", keys)
	if err != nil {
		return shim.Error("ReceiveWeeklyTimesheet() : Failed to query Timesheet object successfully")
	}

	if Avalbytes == nil {
		return shim.Error("ReceiveWeeklyTimesheet() : Timesheet object does not exist for Update")
	}

	me, _ := JSONtoTimeSheet(Avalbytes)

	// Iterate through TimeSheetView TimeEntries and Update Timesheet TimeEntries

	//i := 0
	var tsWeek string
	for k, v := range wtsView.TimeEntries {
		me.TimeEntries[k] = v
		/*if i == 0 {
			//tsWeek = wtsView + k
			tsWeek = k
			i = 1
		}*/
	}

	// Update the Timesheet object
	wts, _ := json.Marshal(me)
	err = itpProcUtils.UpdateObject(stub, "TIMESHEET", []string{me.Year, me.PONumber, me.BuyerCtrID}, wts)
	if err != nil {
		fmt.Println("ReceiveWeeklyTimesheet() : Error inserting Object into LedgerState %s", err)
		return shim.Error("ReceiveWeeklyTimesheet() : Timesheet object Update failed")
	}
     me.TimeEntries = wtsView.TimeEntries
     wts, _ = json.Marshal(me)
	// Check if the Contractor Submitted it for SAVE or SUBMIT for approval
	if wtsView.Status == "SAVE" {
		return shim.Success(wts)
	}

	// Create Approval Object
	// The Approvers have to be included in the Object
	appObj := Approval{"APPROVAL", wtsView.PONumber, wtsView.BuyerCtrID, wtsView.SupplierCtrID, tsWeek, "SUBMIT", "<Adam Gordan>adam.gordon@itpeoplecorp.com", "NA", "NA", []byte("NA"), "<srinivas>srinivas.rachakonda@itpeoplecorp.com", "NA", "NA", []byte("NA")}
	approObjBytes, _ := ApprovaltoJSON(appObj)
	err = itpProcUtils.UpdateObject(stub, "APPROVAL", []string{me.Year, me.PONumber, me.BuyerCtrID, tsWeek}, wts)
	if err != nil {
		fmt.Println("ReceiveWeeklyTimesheet() : Error inserting Approval Object into LedgerState %s", err)
		return shim.Error("ReceiveWeeklyTimesheet() : Approval object Update failed")
	}
	return shim.Success(approObjBytes)

}
