module titan-vrf

go 1.20

require (
	github.com/filecoin-project/go-state-types v0.13.0
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2
)

require (
	github.com/TongTianTaiShi/titan-ffi v0.0.0-20231009094233-2ff4a5177a2b
	github.com/filecoin-project/go-address v1.1.0
	github.com/ipfs/go-block-format v0.2.0 // indirect
	github.com/ipfs/go-cid v0.4.1
	github.com/ipfs/go-ipfs-util v0.0.3 // indirect
	github.com/ipfs/go-ipld-cbor v0.1.0 // indirect
	github.com/ipfs/go-ipld-format v0.6.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/multiformats/go-multihash v0.2.3 // indirect
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/polydawn/refmt v0.89.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/whyrusleeping/cbor-gen v0.0.0-20230923211252-36a87e1ba72f
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	lukechampine.com/blake3 v1.2.1 // indirect
)

replace github.com/TongTianTaiShi/titan-ffi => ./extern/titan-ffi
