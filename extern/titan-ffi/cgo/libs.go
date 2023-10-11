package cgo

/*
#cgo LDFLAGS: -L${SRCDIR}/.. -lticrypto
#cgo pkg-config: ${SRCDIR}/../ticrypto.pc
#include "../ticrypto.h"
#include <stdlib.h>
*/
import "C"
