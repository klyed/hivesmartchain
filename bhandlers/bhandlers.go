package bhandlers

import (
	"fmt"
	//"log"
  "strings"
  //"github.com/klyed/hivesmartchain/execution"
	"github.com/klyed/hiverpc-go/types"
)

//var (
//  response = fmt.Printf("")
//  error = nil
//)

func CustomJSON(block uint32, tx *types.Transaction, op *types.CustomJSONOperation) {
  //fmt.Printf("\n\nHIVE --(custom_json op)-> HSC: Op: %v", op)
  sender := op.RequiredAuths[0]
  json := op.JSON
  jsonParsed := []interface{}{json}
  action := jsonParsed[0]
  method := jsonParsed[1]
  data := []interface{}{jsonParsed[2]}
  fmt.Printf("HIVE --(custom_json)-> HSC: Sender: %v - Action: %v - Method: %v - Data: %v)", sender, action, method, data)
  //execution
  //return response, error
}

func Transfer(block uint32, tx *types.Transaction, op *types.TransferOperation) {
  sender := op.From
  //receiver := op.To
  amount := op.Amount
  amountsplit := strings.Fields(amount)
  value := amountsplit[0]
  coin :=  amountsplit[1]
  //memo := op.Memo
  fmt.Printf("HIVE --(transfer)-> HSC: Sender: %v - Amount: %v or %v %v)", sender, amount, value, coin)
  //return response, error
}
