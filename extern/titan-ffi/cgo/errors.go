package cgo

// #cgo linux LDFLAGS: ${SRCDIR}/../libticrypto.a -Wl,-unresolved-symbols=ignore-all
// #cgo darwin LDFLAGS: ${SRCDIR}/../libticrypto.a -Wl,-undefined,dynamic_lookup
// #cgo pkg-config: ${SRCDIR}/../ticrypto.pc
// #include "../ticrypto.h"
import "C"
