package bhandlers

import (
	"fmt"
	//"log"
  "strings"
	"github.com/klyed/hiverpc-go/types"
)

func CustomJSON(block uint32, tx *types.Transaction, op *types.CustomJSONOperation) (int, error) {
  sender := op.RequiredAuths[0]
  json := op.JSON
  action := []interface{}{json}

  return fmt.Printf("HIVE --(custom_json)-> HSC: Sender: %v - Action: %v - Method: %v)\n", sender, action)
}

func Transfer(block uint32, tx *types.Transaction, op *types.TransferOperation) (int, error) {
  sender := op.From
  //receiver := op.To
  amount := op.Amount
  amountsplit := strings.Fields(amount)
  value := amountsplit[0]
  coin :=  amountsplit[1]
  //memo := op.Memo

  return fmt.Printf("HIVE --(transfer)-> HSC: Sender: %v - Amount: %v or %v %v)\n", sender, amount, value, coin)
}
