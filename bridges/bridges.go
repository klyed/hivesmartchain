package bridges

import (
	"fmt"
	"reflect"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	client "github.com/klyed/hiverpc-go"
	"github.com/klyed/hiverpc-go/types"
	"github.com/klyed/hiverpc-go/transports/websocket"
	"github.com/klyed/hivesmartchain/bhandlers"
)

type HiveConfig struct {
	Protocol      string
	RemoteAddress string
	KeysDirectory string
}

var (
	//aliasList = []string{"hive.smart.chain", "hive.side.chain", "hsc"}
	opBlock = uint32(0)
	Client  = reflect.Func
)

type HiveBridge interface {
	Startbridge(call string, interrupted bool, reconnect bool) error
}

func HiveBlock() uint32 {
	//if OpBlock != nil {
	return opBlock
	//}
	//return nil
}

func Startbridge(call string, interrupted bool, reconnect bool) error {

	var (
		url = []string{"ws://185.130.44.165:8090"}
		//reconnect = true
		signalCh    = make(chan os.Signal, 1)
		monitorChan = make(chan interface{}, 1)
		t, err    = websocket.NewTransport(url,
			websocket.SetAutoReconnectEnabled(reconnect),
			websocket.SetAutoReconnectMaxDelay(30*time.Second),
			websocket.SetMonitor(monitorChan))
	  	Client, clienterr = client.NewClient(t)
	)
	if call == "start" {
		// Start catching signals.
		//signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

		// Start the connection monitor.
		//monitorChan := make(chan interface{}, 1)
		if reconnect != false {
			go func() {
				for {
					event, ok := <-monitorChan
					if ok {
						log.Println(event)
					}
				}
			}()
		}

		// Instantiate the WebSocket transport.
		log.Printf("HIVEOP: Connecting to HIVE Over WebSockets: (\"%v\")\n", url)
/*
			t, err := websocket.NewTransport(url,
				websocket.SetAutoReconnectEnabled(reconnect),
				websocket.SetAutoReconnectMaxDelay(30*time.Second),
				websocket.SetMonitor(monitorChan))
*/

		//if wserr != nil {
		//	log.Println(wserr)
			//return err
		//}

		// Use the transport to get an RPC client.
		//Client, err := client.NewClient(url)
		if clienterr != nil {
			log.Println(clienterr)
			Client.Close()
			Client, clienterr = client.NewClient(t)
		}
		defer func() {
			if !interrupted {
				Client.Close()
			}
			Client.Close()
		}()

		// Start processing signals.
		/*
			go func() {
				<-signalCh
				fmt.Println()
				log.Println("HIVEOP: Signal received, exiting...")
				signal.Stop(signalCh)
				interrupted = true
				Client.Close()
			}()
		*/

		// Drop the error in case it is a request being interrupted.
		//return Client, nil
	}
	if call == "stop" {
		interrupted = true
		<-signalCh
		fmt.Println()
		log.Println("HIVEOP: Signal received, exiting...")
		signal.Stop(signalCh)
		Client.Close()
		if Client.Close() != nil {
			log.Println(Client.Close())
		}
		//return Client, nil
	}

	// Get config.
	//log.Println("HIVEOP: GetConfig()")
	//config, err := Client.API.GetConfig()
	//if err != nil {
	//	log.Println(err)
	//}

	// Use the last irreversible block number as the initial last block number.
	props, err := Client.Database.GetDynamicGlobalProperties()
	if err != nil {
		log.Println(err)
	}
	//Latest Actual Block
	lastBlock := uint32(props.HeadBlockNumber)
	log.Printf("HIVEOP: LAST BLOCK (last block = %v)\n", lastBlock)
	//Latest "Safe" Block
	//lastBlock := uint32(props.LastIrreversibleBlockNum)

	// Keep processing incoming blocks forever.
	log.Printf("HIVEOP: Starting HIVE Block Processing Bridge (last block = %v)\n", lastBlock)
	for {
		// Get current properties.
		props, err := Client.Database.GetDynamicGlobalProperties()
		if err != nil {
			log.Println(err)
		}
		opBlock = uint32(props.HeadBlockNumber)
		// Process new blocks.
		for opBlock-lastBlock > 0 {
			block, err := Client.Database.GetBlock(lastBlock)
			//log.Println(opBlock)
			if err != nil {
				log.Println(err)
			}

			// Process the transactions.
			log.Println(HiveBlock())
			for _, tx := range block.Transactions {
				for _, operation := range tx.Operations {
					//Uncomment line below to watch all ops
					//fmt.Printf("HIVEOP: operation:\n %v", operation)
					//fmt.Printf("HIVEOP: Block: #%v - Op Type: %v", lastBlock, operation.Type())
					switch op := operation.Data().(type) {
					//case *types.VoteOperation:
					//fmt.Printf("HIVEOP: @\"%v\"voted for @\"%v\"/\"%v\" \n", op.Voter, op.Author, op.Permlink)
					//case *types.CustomOperation:
					//fmt.Printf("HIVEOP: transfer:\n %v \n%v", tx)
					case *types.CustomJSONOperation:
						if op.ID == "HSC" {
							//fmt.Printf("HIVEOP: Block: #%v - custom_json:\n %v", lastBlock, op)
							handler, err := bhandlers.CustomJSON(lastBlock, tx, op)
							fmt.Printf("HIVEOP: custom_json:\n %v", handler)
							if err != nil {
								log.Println(err)
							}
							//return Client, nil
						}
						//return op

					case *types.TransferOperation:
						if op.To == "hive.smart.chain" {
							//fmt.Printf("HIVEOP: Block: #%v - transfer:\n %v", lastBlock, op)
							handler, err := bhandlers.Transfer(lastBlock, tx, op)
							fmt.Printf("HIVEOP: transfer:\n %v", handler)
							if err != nil {
								log.Println(err)
							}
							//return Client, nil
						}

						//case *types.CustomJSONOperation:
						//	fmt.Printf("HIVEOP: custom_json:\n %v", op)

						// Vote operation.
						//case *types.TransferOperation:
						//	fmt.Printf("HIVEOP: transfer:\n %v", op)

						// You can add more cases here, it depends on
						// what operations you actually need to process.

					}
				}
			}

			lastBlock++
			opBlock = lastBlock
		}

		// Sleep for HIVE_BLOCK_INTERVAL seconds before the next iteration.
		time.Sleep(time.Duration(3) * time.Second)
	}
}

//func Stop() {
//	close(Client)
//}

//func HiveBlock() {
//	//fmt Printf("HIVEOP: HIVE Block Height on Side Chain: %v", lastBlock)
//	return lastBlock
//}

/*
package bridges

import (
	"fmt"
	"time"
	"github.com/klyed/hiverpc-go"
	"github.com/klyed/hiverpc-go/transports/websocket"
	"github.com/klyed/hiverpc-go/types"
)

func Run() {
	fmt.Printf("\n\n=========================================\n=========================================\nHIVE BRIDGE CONNECTED - SIDE CHAIN ACTIVATED\n=========================================\n=========================================\n\n")
	// Instantiate the WebSocket transport.
	t, _ := websocket.NewTransport([]string{"ws://185.130.44.165:8090", "true"})
	fmt.Printf("T @%v @%v/%v\n", t)

	// Use the transport to create an RPC client.
	client, _ := rpc.NewClient(t)
	fmt.Printf("CLIENT @%v @%v/%v\n", client)
	defer client.Close()

	// Call "get_config".
	config, _ := client.Database.GetConfig()
	fmt.Printf("CONFIG @%v @%v/%v\n", config)
	// Start processing blocks.
	for {
		// Call "get_dynamic_global_properties".
		props, _ := client.Database.GetDynamicGlobalProperties()
		fmt.Printf("PROPS @%v @%v/%v\n", props)
		lastBlock := uint32(props.HeadBlockNumber)
		//for props.LastIrreversibleBlockNum-lastBlock > 0 {
			// Call "get_block".
			block, _ := client.Database.GetBlock(lastBlock)
			fmt.Printf("BLOCK @%v @%v/%v\n", block)
			// Process the transactions.
			for _, tx := range block.Transactions {
				for _, op := range tx.Operations {
					switch body := op.Data().(type) {
					// Comment operation.
					case *types.CustomJSONOperation:
						//content, _ := client.Database.GetContent(body.Author, body.Permlink)
						fmt.Printf("OPERATION @%v @%v/%v\n", op)
						fmt.Printf("CUSTOM_JSON @%v %v\n", body)

					// Vote operation.
					case *types.TransferOperation:
						fmt.Printf("OPERATION @%v @%v/%v\n", op)
						fmt.Printf("TRANSFER @%v @%v/%v\n", body.To, body.From, body.Amount, body)

						// You can add more cases, it depends on what
						// operations you actually need to process.
					}
				}
			}

		//	lastBlock++
		//}

		time.Sleep(time.Duration(config.HiveBlockInterval) * time.Second)
	}
}
*/
