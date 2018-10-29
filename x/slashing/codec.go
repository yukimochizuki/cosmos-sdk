package slashing

import (
	"github.com/yukimochizuki/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUnjail{}, "cosmos-sdk/MsgUnjail", nil)
}

var cdcEmpty = codec.New()
