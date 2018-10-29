package cli

import (
	"github.com/yukimochizuki/cosmos-sdk/client/context"
	"github.com/yukimochizuki/cosmos-sdk/client/utils"
	"github.com/yukimochizuki/cosmos-sdk/codec"
	sdk "github.com/yukimochizuki/cosmos-sdk/types"
	authcmd "github.com/yukimochizuki/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/yukimochizuki/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/yukimochizuki/cosmos-sdk/x/slashing"

	"github.com/spf13/cobra"
)

// GetCmdUnjail implements the create unjail validator command.
func GetCmdUnjail(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unjail",
		Args:  cobra.NoArgs,
		Short: "unjail validator previously jailed for downtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			valAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := slashing.NewMsgUnjail(sdk.ValAddress(valAddr))
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg}, false)
			}
			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
