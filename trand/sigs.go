package trand

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"golang.org/x/xerrors"
)

// SigsSign takes in signature type, private key and message. Returns a signature for that message.
// Valid sigTypes are: "secp256k1" and "bls"
func SigsSign(sigType crypto.SigType, privkey []byte, msg []byte) (*crypto.Signature, error) {
	sv, ok := sigs[sigType]
	if !ok {
		return nil, fmt.Errorf("cannot sign message with signature of unsupported type: %v", sigType)
	}

	sb, err := sv.Sign(privkey, msg)
	if err != nil {
		return nil, err
	}
	return &crypto.Signature{
		Type: sigType,
		Data: sb,
	}, nil
}

// SigsVerify verifies signatures
func SigsVerify(sig *crypto.Signature, addr address.Address, msg []byte) error {
	if sig == nil {
		return xerrors.Errorf("signature is nil")
	}

	if addr.Protocol() == address.ID {
		return fmt.Errorf("must resolve ID addresses before using them to verify a signature")
	}

	sv, ok := sigs[sig.Type]
	if !ok {
		return fmt.Errorf("cannot verify signature of unsupported type: %v", sig.Type)
	}

	return sv.Verify(sig.Data, addr, msg)
}

// SigsGenerate generates private key of given type
func SigsGenerate(sigType crypto.SigType) ([]byte, error) {
	sv, ok := sigs[sigType]
	if !ok {
		return nil, fmt.Errorf("cannot generate private key of unsupported type: %v", sigType)
	}

	return sv.GenPrivate()
}

// SigsToPublic converts private key to public key
func SigsToPublic(sigType crypto.SigType, pk []byte) ([]byte, error) {
	sv, ok := sigs[sigType]
	if !ok {
		return nil, fmt.Errorf("cannot generate public key of unsupported type: %v", sigType)
	}

	return sv.ToPublic(pk)
}

// SigShim is used for introducing signature functions
type SigShim interface {
	GenPrivate() ([]byte, error)
	ToPublic(pk []byte) ([]byte, error)
	Sign(pk []byte, msg []byte) ([]byte, error)
	Verify(sig []byte, a address.Address, msg []byte) error
}

var sigs map[crypto.SigType]SigShim

// RegisterSignature should be only used during init
func RegisterSignature(typ crypto.SigType, vs SigShim) {
	if sigs == nil {
		sigs = make(map[crypto.SigType]SigShim)
	}
	sigs[typ] = vs
}
