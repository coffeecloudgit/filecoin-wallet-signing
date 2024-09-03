package actors

import (
	"bytes"
	addr "github.com/filecoin-project/go-address"
	builtin14 "github.com/filecoin-project/go-state-types/builtin"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/exitcode"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/coffeecloudgit/filecoin-wallet-signing/chain/types"
)

var (
	InitActorAddr = builtin14.InitActorAddr
	MethodsInit   = builtin14.MethodsInit
)

func SerializeParams(i cbg.CBORMarshaler) ([]byte, types.ActorError) {
	buf := new(bytes.Buffer)
	if err := i.MarshalCBOR(buf); err != nil {
		// TODO: shouldnt this be a fatal error?
		return nil, types.Absorb(err, exitcode.ErrSerialization, "failed to encode parameter")
	}
	return buf.Bytes(), nil
}

type ConstructorParams struct {
	NetworkName string
}

type ExecParams struct {
	CodeCID           cid.Cid `checked:"true"` // invalid CIDs won't get committed to the state tree
	ConstructorParams []byte
}

type ExecReturn struct {
	IDAddress     addr.Address // The canonical ID-based address for the actor.
	RobustAddress addr.Address // A more expensive but re-org-safe address for the newly created actor.
}
