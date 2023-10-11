package trand

import (
	"context"
	"encoding/binary"
	"fmt"
	"titan-vrf/filrpc"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/minio/blake2b-simd"
	"golang.org/x/xerrors"
)

func drawRandomness(rbase []byte, pers DomainSeparationTag, height abi.ChainEpoch, entropy []byte) ([]byte, error) {
	h := blake2b.New256()
	if err := binary.Write(h, binary.BigEndian, int64(pers)); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %w", err)
	}
	VRFDigest := blake2b.Sum256(rbase)
	_, err := h.Write(VRFDigest[:])
	if err != nil {
		return nil, xerrors.Errorf("hashing VRFDigest: %w", err)
	}
	if err := binary.Write(h, binary.BigEndian, height); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %w", err)
	}
	_, err = h.Write(entropy)
	if err != nil {
		return nil, xerrors.Errorf("hashing entropy: %w", err)
	}

	return h.Sum(nil), nil
}

type SignFunc func(context.Context, address.Address, []byte) (*crypto.Signature, error)

func computeVRF(ctx context.Context, privateKey, sigInput []byte) ([]byte, error) {
	sig, err := SigsSign(crypto.SigTypeBLS, privateKey, sigInput)
	if err != nil {
		return nil, err
	}

	if sig.Type != crypto.SigTypeBLS {
		return nil, fmt.Errorf("miner worker address was not a BLS key")
	}

	return sig.Data, nil
}

func verifyVRF(ctx context.Context, worker address.Address, vrfBase, vrfproof []byte) error {
	sig := &crypto.Signature{
		Type: crypto.SigTypeBLS,
		Data: vrfproof,
	}

	if err := SigsVerify(sig, worker, vrfBase); err != nil {
		return xerrors.Errorf("vrf was invalid: %w", err)
	}

	return nil
}

type VRFOut struct {
	Height abi.ChainEpoch
	Proof  []byte
}

type BeaconEntry struct {
	Round uint64
	Data  []byte
}

func (vrf *VRFOut) Sum256() [32]byte {
	return blake2b.Sum256(vrf.Proof)
}

func GenerateVRF(ctx context.Context, pers DomainSeparationTag,
	privateKey []byte, rbase []byte, height abi.ChainEpoch, entropy []byte) (*VRFOut, error) {

	// draw randomness
	randomness, err := drawRandomness(rbase, pers, height, entropy)
	if err != nil {
		return nil, xerrors.Errorf("GenerateVRF drawRandomness failed: %w", err)
	}

	// compute vrf
	vrf, err := computeVRF(ctx, privateKey, randomness)
	if err != nil {
		return nil, xerrors.Errorf("GenerateVRF computeVRF failed: %w", err)
	}

	return &VRFOut{
		Height: height,
		Proof:  vrf,
	}, nil
}

func VerifyVRF(ctx context.Context, worker address.Address,
	pers DomainSeparationTag, rbase []byte, entropy []byte, vrf *VRFOut) error {

	// draw randomness
	randomness, err := drawRandomness(rbase, pers, vrf.Height, entropy)
	if err != nil {
		return xerrors.Errorf("VerifyVRF drawRandomness failed: %w", err)
	}

	return verifyVRF(ctx, worker, randomness, vrf.Proof)
}

func GenerateVRFByTipSet(ctx context.Context, pers DomainSeparationTag,
	privateKey []byte, ts *filrpc.TipSet, entropy []byte) (*VRFOut, error) {

	// use min ticket
	minTicket := ts.MinTicket()
	return GenerateVRF(ctx, pers, privateKey, minTicket.VRFProof, ts.Height(), entropy)
}

func VerifyVRFByTipSet(ctx context.Context, worker address.Address,
	pers DomainSeparationTag, ts *filrpc.TipSet, entropy []byte, vrf *VRFOut) error {
	if ts.Height() != vrf.Height {
		return xerrors.Errorf("VerifyVRFByTipSet tipset height %d != %d(vrf)", ts.Height(), vrf.Height)
	}

	// use min ticket
	minTicket := ts.MinTicket()
	return VerifyVRF(ctx, worker, pers, minTicket.VRFProof, entropy, vrf)
}
