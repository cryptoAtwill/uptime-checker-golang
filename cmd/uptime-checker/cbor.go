package main

import (
	"fmt"
	"io"
	"math"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

func (t *MapState) UnmarshalCBOR(r io.Reader) (err error) {
	return nil
}

func (t *Votes) UnmarshalCBOR(r io.Reader) (err error) {
	*t = Votes{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.lastVote (string) (string)
	{
		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}

		t.lastVote = ChainEpoch(extra)
	}

	// t.votes ([]PeerID) (slice)
	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Miners: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.votes = make([]PeerID, extra)
	}

	for i := 0; i < int(extra); i++ {
		v, err := cbg.ReadString(cr)
		if err != nil {
			return err
		}
		t.votes[i] = v
	}

	return nil
}

func (t *NodeInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = NodeInfo{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 3 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.id
	{
		v, err := cbg.ReadString(cr)
		if err != nil {
			return err
		}
		t.id = PeerID(v)
	}

	// t.creator
	{
		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}

		t.creator = ActorID(extra)
	}

	{
		maj, extra, err = cr.ReadHeader()
		if err != nil {
			return err
		}

		if extra > cbg.MaxLength {
			return fmt.Errorf("t.Miners: array too large (%d)", extra)
		}

		if maj != cbg.MajArray {
			return fmt.Errorf("expected cbor array")
		}

		if extra > 0 {
			t.addresses = make([]MultiAddr, extra)
		}

		for i := 0; i < int(extra); i++ {
			v, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}
			t.addresses[i] = v
		}
	}

	return nil
}
