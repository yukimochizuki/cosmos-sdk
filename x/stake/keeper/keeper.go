package keeper

import (
	"github.com/yukimochizuki/cosmos-sdk/codec"
	sdk "github.com/yukimochizuki/cosmos-sdk/types"

	"github.com/yukimochizuki/cosmos-sdk/x/bank"
	"github.com/yukimochizuki/cosmos-sdk/x/params"
	"github.com/yukimochizuki/cosmos-sdk/x/stake/types"
)

// keeper of the stake store
type Keeper struct {
	storeKey   sdk.StoreKey
	storeTKey  sdk.StoreKey
	cdc        *codec.Codec
	bankKeeper bank.Keeper
	hooks      sdk.StakingHooks
	paramstore params.Subspace

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey, ck bank.Keeper, paramstore params.Subspace, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:   key,
		storeTKey:  tkey,
		cdc:        cdc,
		bankKeeper: ck,
		paramstore: paramstore.WithTypeTable(ParamTypeTable()),
		hooks:      nil,
		codespace:  codespace,
	}
	return keeper
}

// Set the validator hooks
func (k Keeper) WithHooks(sh sdk.StakingHooks) Keeper {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}
	k.hooks = sh
	return k
}

//_________________________________________________________________________

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

//_______________________________________________________________________

// load the pool
func (k Keeper) GetPool(ctx sdk.Context) (pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(PoolKey)
	if b == nil {
		panic("stored pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinary(b, &pool)
	return
}

// set the pool
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinary(pool)
	store.Set(PoolKey, b)
}

//_______________________________________________________________________

// Load the last total validator power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(LastTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinary(b, &power)
	return
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinary(power)
	store.Set(LastTotalPowerKey, b)
}

//_______________________________________________________________________

// Load the last validator power.
// Returns zero if the operator was not a validator last block.
func (k Keeper) GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetLastValidatorPowerKey(operator))
	if bz == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinary(bz, &power)
	return
}

// Set the last validator power.
func (k Keeper) SetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(power)
	store.Set(GetLastValidatorPowerKey(operator), bz)
}

// Delete the last validator power.
func (k Keeper) DeleteLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetLastValidatorPowerKey(operator))
}

//__________________________________________________________________________

// get the current in-block validator operation counter
func (k Keeper) GetIntraTxCounter(ctx sdk.Context) int16 {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(IntraTxCounterKey)
	if b == nil {
		return 0
	}
	var counter int16
	k.cdc.MustUnmarshalBinary(b, &counter)
	return counter
}

// set the current in-block validator operation counter
func (k Keeper) SetIntraTxCounter(ctx sdk.Context, counter int16) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinary(counter)
	store.Set(IntraTxCounterKey, bz)
}
