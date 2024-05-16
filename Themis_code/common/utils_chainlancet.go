package common

import (
	"github.com/ethereum/go-ethereum/themis"
)

var GlobalThemisObj *themis.SManager

var Opcodeinfo themis.SOpcodecommon

var Exsit_flag bool

var Transfer bool

var CallGas uint64

var AccountSAS bool
var StorageSAS bool

var MaxMemSize uint64
