package state

import (
	"github.com/klyed/hivesmartchain/encoding"

	"github.com/klyed/hivesmartchain/execution/proposal"
	"github.com/klyed/hivesmartchain/txs/payload"
)

var _ proposal.IterableReader = &State{}

func (s *ReadState) GetProposal(proposalHash []byte) (*payload.Ballot, error) {
	tree, err := s.Forest.Reader(keys.Proposal.Prefix())
	if err != nil {
		return nil, err
	}
	bs, err := tree.Get(keys.Proposal.KeyNoPrefix(proposalHash))
	if err != nil {
		return nil, err
	} else if len(bs) == 0 {
		return nil, nil
	}

	ballot := new(payload.Ballot)
	err = encoding.Decode(bs, ballot)
	if err != nil {
		return nil, err
	}
	return ballot, nil
}

func (ws *writeState) UpdateProposal(proposalHash []byte, p *payload.Ballot) error {
	tree, err := ws.forest.Writer(keys.Proposal.Prefix())
	if err != nil {
		return err
	}
	bs, err := encoding.Encode(p)
	if err != nil {
		return err
	}

	tree.Set(keys.Proposal.KeyNoPrefix(proposalHash), bs)
	return nil
}

func (ws *writeState) RemoveProposal(proposalHash []byte) error {
	tree, err := ws.forest.Writer(keys.Proposal.Prefix())
	if err != nil {
		return err
	}
	tree.Delete(keys.Proposal.KeyNoPrefix(proposalHash))
	return nil
}

func (s *ReadState) IterateProposals(consumer func(proposalHash []byte, proposal *payload.Ballot) error) error {
	tree, err := s.Forest.Reader(keys.Proposal.Prefix())
	if err != nil {
		return err
	}
	return tree.Iterate(nil, nil, true, func(key []byte, value []byte) error {
		ballot := new(payload.Ballot)
		err := encoding.Decode(value, ballot)
		if err != nil {
			return err
		}
		return consumer(key, ballot)
	})
}
