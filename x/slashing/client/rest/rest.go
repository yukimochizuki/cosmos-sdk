package rest

import (
	"github.com/yukimochizuki/cosmos-sdk/client/context"
	"github.com/yukimochizuki/cosmos-sdk/codec"
	"github.com/yukimochizuki/cosmos-sdk/crypto/keys"

	"github.com/gorilla/mux"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	registerQueryRoutes(cliCtx, r, cdc)
	registerTxRoutes(cliCtx, r, cdc, kb)
}
