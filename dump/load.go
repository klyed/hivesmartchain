package dump

import (
	"crypto/sha256"
	bin "encoding/binary"
	"io"

	"github.com/klyed/hivesmartchain/acm"
	"github.com/klyed/hivesmartchain/acm/acmstate"
	"github.com/klyed/hivesmartchain/binary"
	"github.com/klyed/hivesmartchain/execution/exec"
	"github.com/klyed/hivesmartchain/execution/state"
	"github.com/klyed/hivesmartchain/txs/payload"
)

// Load a dump into state
func Load(source Source, st *state.State) error {
	_, _, err := st.Update(func(s state.Updatable) error {
		txs := make([]*exec.TxExecution, 0)

		var tx *exec.TxExecution

		for {
			row, err := source.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			if row.Account != nil {
				if row.Account.Address != acm.GlobalPermissionsAddress {
					for _, m := range row.Account.ContractMeta {
						metahash := acmstate.GetMetadataHash(m.Metadata)
						err = s.SetMetadata(metahash, m.Metadata)
						if err != nil {
							return err
						}
						m.MetadataHash = metahash.Bytes()
						m.Metadata = ""
					}
					err := s.UpdateAccount(row.Account)
					if err != nil {
						return err
					}
				}
			}

			if row.AccountStorage != nil {
				for _, storage := range row.AccountStorage.Storage {
					err := s.SetStorage(row.AccountStorage.Address, storage.Key, storage.Value)
					if err != nil {
						return err
					}
				}
			}

			if row.Name != nil {
				err := s.UpdateName(row.Name)
				if err != nil {
					return err
				}
			}

			if row.EVMEvent != nil {
				if tx != nil && row.Height != tx.Height {
					txs = append(txs, tx)
					tx = nil
				}
				if tx == nil {
					tx = &exec.TxExecution{
						TxHeader: &exec.TxHeader{
							TxHash: dumpTxHash(row.EVMEvent.ChainID, row.Height),
							TxType: payload.TypeCall,
							Origin: &exec.Origin{
								ChainID: row.EVMEvent.ChainID,
								Height:  row.Height,
								Time:    row.EVMEvent.Time,
								Index:   row.EVMEvent.Index,
							},
						},
					}
				}

				tx.Events = append(tx.Events, &exec.Event{
					Header: &exec.Header{
						TxType:    payload.TypeCall,
						EventType: exec.TypeLog,
						Height:    row.Height,
					},
					Log: row.EVMEvent.Event,
				})
			}
		}

		if tx != nil {
			txs = append(txs, tx)
		}

		return s.AddBlock(&exec.BlockExecution{
			Height:       0,
			TxExecutions: txs,
		})
	})
	return err
}

// Provides a psuedo-hash for the singular 'dump tx' that is generated by a restore
func dumpTxHash(chainID string, lastHeight uint64) binary.HexBytes {
	hasher := sha256.New()
	hasher.Write([]byte(chainID))
	bs := make([]byte, 8)
	bin.BigEndian.PutUint64(bs, lastHeight)
	hasher.Write(bs)
	return hasher.Sum(nil)
}
