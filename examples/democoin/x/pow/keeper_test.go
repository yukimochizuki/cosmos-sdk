package pow

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/yukimochizuki/cosmos-sdk/codec"
	"github.com/yukimochizuki/cosmos-sdk/store"
	sdk "github.com/yukimochizuki/cosmos-sdk/types"
	auth "github.com/yukimochizuki/cosmos-sdk/x/auth"
	bank "github.com/yukimochizuki/cosmos-sdk/x/bank"
)

// possibly share this kind of setup functionality between module testsuites?
func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return ms, capKey
}

func TestPowKeeperGetSet(t *testing.T) {
	ms, capKey := setupMultiStore()
	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	am := auth.NewAccountKeeper(cdc, capKey, auth.ProtoBaseAccount)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	config := NewConfig("pow", int64(1))
	ck := bank.NewBaseKeeper(am)
	keeper := NewKeeper(capKey, config, ck, DefaultCodespace)

	err := InitGenesis(ctx, keeper, Genesis{uint64(1), uint64(0)})
	require.Nil(t, err)

	genesis := WriteGenesis(ctx, keeper)
	require.Nil(t, err)
	require.Equal(t, genesis, Genesis{uint64(1), uint64(0)})

	res, err := keeper.GetLastDifficulty(ctx)
	require.Nil(t, err)
	require.Equal(t, res, uint64(1))

	keeper.SetLastDifficulty(ctx, 2)

	res, err = keeper.GetLastDifficulty(ctx)
	require.Nil(t, err)
	require.Equal(t, res, uint64(2))
}
