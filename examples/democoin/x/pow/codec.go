package pow

import (
	"github.com/yukimochizuki/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgMine{}, "pow/Mine", nil)
}
