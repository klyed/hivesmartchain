package commands

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/klyed/hivesmartchain/encoding"

	"github.com/howeyc/gopass"
	cli "github.com/jawher/mow.cli"
	"github.com/klyed/hivesmartchain/config"
	"github.com/klyed/hivesmartchain/config/deployment"
	"github.com/klyed/hivesmartchain/crypto"
	"github.com/klyed/hivesmartchain/keys"
)

// Keys runs as either client or server
func Keys(output Output) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		keysHost := cmd.String(cli.StringOpt{
			Name:   "host",
			Desc:   "set the host for talking to the key daemon",
			Value:  keys.DefaultHost,
			EnvVar: "HSC_KEYS_HOST",
		})

		keysPort := cmd.String(cli.StringOpt{
			Name:   "port",
			Desc:   "set the port for key daemon",
			Value:  keys.DefaultPort,
			EnvVar: "HSC_KEYS_PORT",
		})

		grpcKeysClient := func(output Output) keys.KeysClient {
			conn, err := encoding.GRPCDial(*keysHost + ":" + *keysPort)
			if err != nil {
				output.Fatalf("Failed to connect to grpc server: %v", err)
			}
			return keys.NewKeysClient(conn)
		}

		cmd.Command("server", "run keys server", func(cmd *cli.Cmd) {
			keysDir := cmd.StringOpt("dir", "", "specify the location of the directory containing key files")
			badPerm := cmd.BoolOpt("allow-bad-perm", true, "Allow unix key file permissions to be readable other than user")
			configOpt := cmd.StringOpt("c config", "", "Use the specified Hive Smart Chain config file")

			var conf *config.BurrowConfig

			cmd.Before = func() {
				var err error
				conf, err = obtainDefaultConfig(*configOpt, "")
				if err != nil {
					output.Fatalf("Could not obtain config: %v", err)
				}
			}

			cmd.Action = func() {
				conf.Keys.AllowBadFilePermissions = *badPerm //Original was: *badPerm   Changed to true to enable windows compiling

				if *keysDir != "" {
					conf.Keys.KeysDirectory = *keysDir
				}

				server := keys.StandAloneServer(conf.Keys.KeysDirectory, conf.Keys.AllowBadFilePermissions)
				address := fmt.Sprintf("%s:%s", *keysHost, *keysPort)
				listener, err := net.Listen("tcp", address)
				if err != nil {
					output.Fatalf("Could not listen on %s: %v", address, err)
				}
				err = server.Serve(listener)
				if err != nil {
					output.Fatalf("Keys server terminated with error: %v", err)
				}
			}
		})

		cmd.Command("gen", "Generates a key using (insert crypto pkgs used)", func(cmd *cli.Cmd) {
			noPassword := cmd.BoolOpt("n no-password", false, "don't use a password for this key")

			keyType := cmd.StringOpt("t curvetype", "ed25519", "specify the curve type of key to create. Supports 'secp256k1' (ethereum),  'ed25519' (tendermint)")

			keyName := cmd.StringOpt("name", "", "name of key to use")

			cmd.Action = func() {
				curve, err := crypto.CurveTypeFromString(*keyType)
				if err != nil {
					output.Fatalf("Unrecognised curve type %v", *keyType)
				}

				var password string
				if !*noPassword {
					fmt.Printf("Enter Password:")
					pwd, err := gopass.GetPasswdMasked()
					if err != nil {
						os.Exit(1)
					}
					password = string(pwd)
				}

				c := grpcKeysClient(output)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				resp, err := c.GenerateKey(ctx, &keys.GenRequest{Passphrase: password, CurveType: curve.String(), KeyName: *keyName})
				if err != nil {
					output.Fatalf("failed to generate key: %v", err)
				}

				fmt.Printf("%v\n", resp.GetAddress())
			}
		})

		cmd.Command("hash", "hash <some data>", func(cmd *cli.Cmd) {
			hashType := cmd.StringOpt("t type", keys.DefaultHashType, "specify the hash function to use")

			hexByte := cmd.BoolOpt("hex", false, "the input should be hex decoded to bytes first")

			msg := cmd.StringArg("MSG", "", "message to hash")

			cmd.Action = func() {
				var message []byte
				var err error
				if *hexByte {
					message, err = hex.DecodeString(*msg)
					if err != nil {
						output.Fatalf("failed to hex decode message: %v", err)
					}
				} else {
					message = []byte(*msg)
				}

				c := grpcKeysClient(output)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				resp, err := c.Hash(ctx, &keys.HashRequest{Hashtype: *hashType, Message: message})
				if err != nil {
					output.Fatalf("failed to get public key: %v", err)
				}

				fmt.Printf("%v\n", resp.GetHash())
			}
		})

		cmd.Command("export", "Export a key to tendermint format", func(cmd *cli.Cmd) {
			keyName := cmd.StringOpt("name", "", "name of key to use")
			keyAddr := cmd.StringOpt("addr", "", "address of key to use")
			passphrase := cmd.StringOpt("passphrase", "", "passphrase for encrypted key")
			keyTemplate := cmd.StringOpt("t template", deployment.DefaultKeysExportFormat, "template for export key")

			cmd.Action = func() {
				c := grpcKeysClient(output)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				resp, err := c.Export(ctx, &keys.ExportRequest{Passphrase: *passphrase, Name: *keyName, Address: *keyAddr})
				if err != nil {
					output.Fatalf("failed to export key: %v", err)
				}

				addr, err := crypto.AddressFromBytes(resp.GetAddress())
				if err != nil {
					output.Fatalf("failed to convert address: %v", err)
				}

				key := deployment.Key{
					Name:       *keyName,
					CurveType:  resp.GetCurveType(),
					Address:    addr,
					PublicKey:  resp.GetPublickey(),
					PrivateKey: resp.GetPrivatekey(),
				}

				str, err := key.Dump(*keyTemplate)
				if err != nil {
					output.Fatalf("failed to template key: %v", err)
				}

				fmt.Printf("%s\n", str)
			}
		})

		cmd.Command("import", "import <priv key> | /path/to/keyfile | <key json>", func(cmd *cli.Cmd) {
			curveType := cmd.StringOpt("t curvetype", "ed25519", "specify the curve type of key to create. Supports 'secp256k1' (ethereum),  'ed25519' (tendermint)")
			noPassword := cmd.BoolOpt("n no-password", false, "don't use a password for this key")
			key := cmd.StringArg("KEY", "", "private key, filename, or raw json")

			cmd.Action = func() {
				var password string
				if !*noPassword {
					fmt.Printf("Enter Password:")
					pwd, err := gopass.GetPasswdMasked()
					if err != nil {
						os.Exit(1)
					}
					password = string(pwd)
				}

				var privKeyBytes []byte
				fileContents, err := ioutil.ReadFile(*key)
				if err == nil {
					*key = string(fileContents)
				}

				c := grpcKeysClient(output)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				if (*key)[:1] == "{" {
					resp, err := c.ImportJSON(ctx, &keys.ImportJSONRequest{JSON: *key})
					if err != nil {
						output.Fatalf("failed to import json key: %v", err)
					}

					fmt.Printf("%s\n", resp.GetAddress())
				} else {
					privKeyBytes, err = hex.DecodeString(*key)
					if err != nil {
						output.Fatalf("failed to hex decode key: %s", *key)
					}
					resp, err := c.Import(ctx, &keys.ImportRequest{Passphrase: password, KeyBytes: privKeyBytes, CurveType: *curveType})
					if err != nil {
						output.Fatalf("failed to import json key: %v", err)
					}

					fmt.Printf("%s\n", resp.GetAddress())

				}
			}
		})

		cmd.Command("pub", "public key", func(cmd *cli.Cmd) {
			name := cmd.StringOpt("name", "", "name of key to use")
			addr := cmd.StringOpt("addr", "", "address of key to use")

			cmd.Action = func() {
				c := grpcKeysClient(output)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				resp, err := c.PublicKey(ctx, &keys.PubRequest{Name: *name, Address: *addr})
				if err != nil {
					output.Fatalf("failed to get public key: %v", err)
				}

				fmt.Printf("%X\n", resp.GetPublicKey())
			}
		})

		cmd.Command("sign", "sign <some data>", func(cmd *cli.Cmd) {
			name := cmd.StringOpt("name", "", "name of key to use")
			addr := cmd.StringOpt("addr", "", "address of key to use")
			msg := cmd.StringArg("MSG", "", "message to sign")
			passphrase := cmd.StringOpt("passphrase", "", "passphrase for encrypted key")

			cmd.Action = func() {
				message, err := hex.DecodeString(*msg)
				if err != nil {
					output.Fatalf("failed to hex decode message: %v", err)
				}

				c := grpcKeysClient(output)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				resp, err := c.Sign(ctx, &keys.SignRequest{Passphrase: *passphrase, Name: *name, Address: *addr, Message: message})
				if err != nil {
					output.Fatalf("failed to get public key: %v", err)
				}
				fmt.Printf("%X\n", resp.GetSignature().Signature)
			}
		})

		cmd.Command("verify", "verify <some data> <sig> <pubkey>", func(cmd *cli.Cmd) {
			curveTypeOpt := cmd.StringOpt("t curvetype", "ed25519", "specify the curve type of key to create. Supports 'secp256k1' (ethereum),  'ed25519' (tendermint)")

			msg := cmd.StringArg("MSG", "", "hash/message to check")
			sig := cmd.StringArg("SIG", "", "signature")
			pub := cmd.StringArg("PUBLIC", "", "public key")

			cmd.Action = func() {
				message, err := hex.DecodeString(*msg)
				if err != nil {
					output.Fatalf("failed to hex decode message: %v", err)
				}
				curveType, err := crypto.CurveTypeFromString(*curveTypeOpt)
				if err != nil {
					output.Fatalf("invalid curve type: %v", err)
				}

				signatureBytes, err := hex.DecodeString(*sig)
				if err != nil {
					output.Fatalf("failed to hex decode signature: %v", err)
				}

				signature, err := crypto.SignatureFromBytes(signatureBytes, curveType)
				if err != nil {
					output.Fatalf("could not form signature: %v", err)
				}

				publickey, err := hex.DecodeString(*pub)
				if err != nil {
					output.Fatalf("failed to hex decode publickey: %v", err)
				}

				c := grpcKeysClient(output)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_, err = c.Verify(ctx, &keys.VerifyRequest{
					PublicKey: publickey,
					Signature: signature,
					Message:   message,
				})
				if err != nil {
					output.Fatalf("failed to verify: %v", err)
				}
				output.Printf("true\n")
			}
		})

		cmd.Command("name", "add key name to addr", func(cmd *cli.Cmd) {
			keyname := cmd.StringArg("NAME", "", "name of key to use")
			addr := cmd.StringArg("ADDR", "", "address of key to use")

			cmd.Action = func() {
				c := grpcKeysClient(output)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_, err := c.AddName(ctx, &keys.AddNameRequest{Keyname: *keyname, Address: *addr})
				if err != nil {
					output.Fatalf("failed to add name to addr: %v", err)
				}
			}
		})

		cmd.Command("list", "list keys", func(cmd *cli.Cmd) {
			name := cmd.StringOpt("name", "", "name or address of key to use")

			cmd.Action = func() {
				c := grpcKeysClient(output)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				resp, err := c.List(ctx, &keys.ListRequest{KeyName: *name})
				if err != nil {
					output.Fatalf("failed to list key: %v", err)
				}
				printKeys := resp.Key
				if printKeys == nil {
					printKeys = make([]*keys.KeyID, 0)
				}
				bs, err := json.MarshalIndent(printKeys, "", "    ")
				if err != nil {
					output.Fatalf("failed to json encode keys: %v", err)
				}
				fmt.Printf("%s\n", string(bs))
			}
		})

		cmd.Command("rm", "rm key name", func(cmd *cli.Cmd) {
			name := cmd.StringArg("NAME", "", "key to remove")

			cmd.Action = func() {
				c := grpcKeysClient(output)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_, err := c.RemoveName(ctx, &keys.RemoveNameRequest{KeyName: *name})
				if err != nil {
					output.Fatalf("failed to remove key: %v", err)
				}
			}
		})
	}
}
