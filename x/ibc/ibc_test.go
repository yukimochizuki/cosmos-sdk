package ibc

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/yukimochizuki/cosmos-sdk/codec"
	"github.com/yukimochizuki/cosmos-sdk/store"
	sdk "github.com/yukimochizuki/cosmos-sdk/types"
	"github.com/yukimochizuki/cosmos-sdk/x/auth"
	"github.com/yukimochizuki/cosmos-sdk/x/bank"
)

// AccountKeeper(/Keeper) and IBCMapper should use different StoreKey later

func defaultContext(key sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func newAddress() sdk.AccAddress {
	return sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
}

func getCoins(ck bank.Keeper, ctx sdk.Context, addr sdk.AccAddress) (sdk.Coins, sdk.Error) {
	zero := sdk.Coins(nil)
	coins, _, err := ck.AddCoins(ctx, addr, zero)
	return coins, err
}

func makeCodec() *codec.Codec {
	var cdc = codec.New()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(bank.MsgSend{}, "test/ibc/Send", nil)
	cdc.RegisterConcrete(bank.MsgIssue{}, "test/ibc/Issue", nil)
	cdc.RegisterConcrete(IBCTransferMsg{}, "test/ibc/IBCTransferMsg", nil)
	cdc.RegisterConcrete(IBCReceiveMsg{}, "test/ibc/IBCReceiveMsg", nil)

	// Register AppAccount
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/ibc/Account", nil)
	codec.RegisterCrypto(cdc)

	cdc.Seal()

	return cdc
}

func TestIBC(t *testing.T) {
	cdc := makeCodec()

	key := sdk.NewKVStoreKey("ibc")
	ctx := defaultContext(key)

	am := auth.NewAccountKeeper(cdc, key, auth.ProtoBaseAccount)
	ck := bank.NewBaseKeeper(am)

	src := newAddress()
	dest := newAddress()
	chainid := "ibcchain"
	zero := sdk.Coins(nil)
	mycoins := sdk.Coins{sdk.NewInt64Coin("mycoin", 10)}

	coins, _, err := ck.AddCoins(ctx, src, mycoins)
	require.Nil(t, err)
	require.Equal(t, mycoins, coins)

	ibcm := NewMapper(cdc, key, DefaultCodespace)
	h := NewHandler(ibcm, ck)
	packet := IBCPacket{
		SrcAddr:   src,
		DestAddr:  dest,
		Coins:     mycoins,
		SrcChain:  chainid,
		DestChain: chainid,
	}

	store := ctx.KVStore(key)

	var msg sdk.Msg
	var res sdk.Result
	var egl int64
	var igs int64

	egl = ibcm.getEgressLength(store, chainid)
	require.Equal(t, egl, int64(0))

	msg = IBCTransferMsg{
		IBCPacket: packet,
	}
	res = h(ctx, msg)
	require.True(t, res.IsOK())

	coins, err = getCoins(ck, ctx, src)
	require.Nil(t, err)
	require.Equal(t, zero, coins)

	egl = ibcm.getEgressLength(store, chainid)
	require.Equal(t, egl, int64(1))

	igs = ibcm.GetIngressSequence(ctx, chainid)
	require.Equal(t, igs, int64(0))

	msg = IBCReceiveMsg{
		IBCPacket: packet,
		Relayer:   src,
		Sequence:  0,
	}
	res = h(ctx, msg)
	require.True(t, res.IsOK())

	coins, err = getCoins(ck, ctx, dest)
	require.Nil(t, err)
	require.Equal(t, mycoins, coins)

	igs = ibcm.GetIngressSequence(ctx, chainid)
	require.Equal(t, igs, int64(1))

	res = h(ctx, msg)
	require.False(t, res.IsOK())

	igs = ibcm.GetIngressSequence(ctx, chainid)
	require.Equal(t, igs, int64(1))
}
