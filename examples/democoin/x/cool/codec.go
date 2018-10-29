package cool

import (
	"github.com/yukimochizuki/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgQuiz{}, "cool/Quiz", nil)
	cdc.RegisterConcrete(MsgSetTrend{}, "cool/SetTrend", nil)
}
