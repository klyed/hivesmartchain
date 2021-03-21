package bridges

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/klyed/hiverpc-go"
	"github.com/klyed/hiverpc-go/transports/websocket"
	"github.com/klyed/hiverpc-go/types"
)

func Run() {
	if err := Startbridge(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func Startbridge() (err error) {
	// Process flags.
	//flagAddress := []string{"rpc_endpoint", "ws://185.130.44.165:8090", "steemd RPC endpoint address"}
	//flagReconnect := flag.Bool("reconnect", false, "enable auto-reconnect mode")
	//flag.Parse()

	var (
		url       = []string{"ws://185.130.44.165:8090"}
		reconnect = true
	)

	// Start catching signals.
	var interrupted bool
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Drop the error in case it is a request being interrupted.
	defer func() {
		if err == websocket.ErrClosing && interrupted {
			err = nil
		}
	}()

	// Start the connection monitor.
	monitorChan := make(chan interface{}, 1)
	if reconnect {
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
	log.Printf("---> Dial(\"%v\")\n", url)
	t, err := websocket.NewTransport(url,
		websocket.SetAutoReconnectEnabled(reconnect),
		websocket.SetAutoReconnectMaxDelay(30*time.Second),
		websocket.SetMonitor(monitorChan))
	if err != nil {
		return err
	}

	// Use the transport to get an RPC client.
	client, err := rpc.NewClient(t)
	if err != nil {
		return err
	}
	defer func() {
		if !interrupted {
			client.Close()
		}
	}()

	// Start processing signals.
	go func() {
		<-signalCh
		fmt.Println()
		log.Println("Signal received, exiting...")
		signal.Stop(signalCh)
		interrupted = true
		client.Close()
	}()

	// Get config.
	log.Println("---> GetConfig()")
	config, err := client.Database.GetConfig()
	if err != nil {
		return err
	}

	// Use the last irreversible block number as the initial last block number.
	props, err := client.Database.GetDynamicGlobalProperties()
	if err != nil {
		return err
	}
	lastBlock := props.LastIrreversibleBlockNum

	// Keep processing incoming blocks forever.
	log.Printf("---> Entering the block processing loop (last block = %v)\n", lastBlock)
	for {
		// Get current properties.
		props, err := client.Database.GetDynamicGlobalProperties()
		if err != nil {
			return err
		}

		// Process new blocks.
		for props.LastIrreversibleBlockNum-lastBlock > 0 {
			block, err := client.Database.GetBlock(lastBlock)
			if err != nil {
				return err
			}

			// Process the transactions.
			for _, tx := range block.Transactions {
				for _, operation := range tx.Operations {
					switch op := operation.Data().(type) {
					case *types.VoteOperation:
						fmt.Printf("@\"%v\"voted for @\"%v\"/\"%v\"\n", op.Voter, op.Author, op.Permlink)

					case *types.CustomJSONOperation:
						//content, _ := client.Database.GetContent(body.Author, body.Permlink)
						fmt.Printf("OPERATION \"%v\"\n", op)

					// Vote operation.
					case *types.TransferOperation:
						fmt.Printf("OPERATION \"%v\"\n", op)

						// You can add more cases, it depends on what
						// operations you actually need to process.
						// You can add more cases here, it depends on
						// what operations you actually need to process.
					}
				}
			}

			lastBlock++
		}

		// Sleep for HIVE_BLOCK_INTERVAL seconds before the next iteration.
		time.Sleep(time.Duration(config.HiveBlockInterval) * time.Second)
	}
}


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
