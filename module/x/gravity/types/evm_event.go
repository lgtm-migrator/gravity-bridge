package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

var (
	_ EVMEvent = &SendToCosmosEvent{}
	_ EVMEvent = &BatchExecutedEvent{}
	_ EVMEvent = &ContractCallExecutedEvent{}
	_ EVMEvent = &ERC20DeployedEvent{}
	_ EVMEvent = &SignerSetTxExecutedEvent{}
)

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m *EVMEventVoteRecord) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var event EVMEvent
	return unpacker.UnpackAny(m.Event, &event)
}

//////////
// Hash //
//////////

func (stce *SendToCosmosEvent) Hash() tmbytes.HexBytes {
	rcv, _ := sdk.AccAddressFromBech32(stce.CosmosReceiver)
	path := bytes.Join(
		[][]byte{
			sdk.Uint64ToBigEndian(stce.EventNonce),
			common.HexToAddress(stce.TokenContract).Bytes(),
			stce.Amount.BigInt().Bytes(),
			common.Hex2Bytes(stce.EVMSender),
			rcv.Bytes(),
			sdk.Uint64ToBigEndian(stce.EVMHeight),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (bee *BatchExecutedEvent) Hash() tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			common.HexToAddress(bee.TokenContract).Bytes(),
			sdk.Uint64ToBigEndian(bee.EventNonce),
			sdk.Uint64ToBigEndian(bee.BatchNonce),
			sdk.Uint64ToBigEndian(bee.EVMHeight),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (ccee *ContractCallExecutedEvent) Hash() tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			sdk.Uint64ToBigEndian(ccee.EventNonce),
			ccee.InvalidationScope,
			sdk.Uint64ToBigEndian(ccee.InvalidationNonce),
			sdk.Uint64ToBigEndian(ccee.EVMHeight),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (e20de *ERC20DeployedEvent) Hash() tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			sdk.Uint64ToBigEndian(e20de.EventNonce),
			[]byte(e20de.CosmosDenom),
			common.HexToAddress(e20de.TokenContract).Bytes(),
			[]byte(e20de.Erc20Name),
			[]byte(e20de.Erc20Symbol),
			sdk.Uint64ToBigEndian(e20de.Erc20Decimals),
			sdk.Uint64ToBigEndian(e20de.EVMHeight),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (sse *SignerSetTxExecutedEvent) Hash() tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			sdk.Uint64ToBigEndian(sse.EventNonce),
			sdk.Uint64ToBigEndian(sse.SignerSetTxNonce),
			sdk.Uint64ToBigEndian(sse.EVMHeight),
			EVMSigners(sse.Members).Hash(),
		},
		[]byte{},
	)
	hash := sha256.Sum256(([]byte(path)))
	return hash[:]
}

//////////////
// Validate //
//////////////

func (stce *SendToCosmosEvent) Validate() error {
	if stce.EventNonce == 0 {
		return fmt.Errorf("event nonce cannot be 0")
	}
	if !common.IsHexAddress(stce.TokenContract) {
		return sdkerrors.Wrap(ErrInvalid, "EVM contract address")
	}
	if stce.Amount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}
	if !common.IsHexAddress(stce.EVMSender) {
		return sdkerrors.Wrap(ErrInvalid, "EVM sender")
	}
	if _, err := sdk.AccAddressFromBech32(stce.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, stce.CosmosReceiver)
	}
	return nil
}

func (bee *BatchExecutedEvent) Validate() error {
	if bee.EventNonce == 0 {
		return fmt.Errorf("event nonce cannot be 0")
	}
	if !common.IsHexAddress(bee.TokenContract) {
		return sdkerrors.Wrap(ErrInvalid, "EVM contract address")
	}
	return nil
}

func (ccee *ContractCallExecutedEvent) Validate() error {
	if ccee.EventNonce == 0 {
		return fmt.Errorf("event nonce cannot be 0")
	}
	return nil
}

func (e20de *ERC20DeployedEvent) Validate() error {
	if e20de.EventNonce == 0 {
		return fmt.Errorf("event nonce cannot be 0")
	}
	if !common.IsHexAddress(e20de.TokenContract) {
		return sdkerrors.Wrap(ErrInvalid, "EVM contract address")
	}
	if err := sdk.ValidateDenom(e20de.CosmosDenom); err != nil {
		return err
	}
	return nil
}

func (sse *SignerSetTxExecutedEvent) Validate() error {
	if sse.EventNonce == 0 {
		return fmt.Errorf("event nonce cannot be 0")
	}
	if sse.Members == nil {
		return fmt.Errorf("members cannot be nil")
	}
	for i, member := range sse.Members {
		if err := member.ValidateBasic(); err != nil {
			return fmt.Errorf("EVM signer %d error: %w", i, err)
		}
	}
	return nil
}

/////////////
// ChainID //
/////////////

func (stce *SendToCosmosEvent) ChainID() uint32 {
	return stce.GetChainId()
}

func (bee *BatchExecutedEvent) ChainID() uint32 {
	return bee.GetChainId()
}

func (ccee *ContractCallExecutedEvent) ChainID() uint32 {
	return ccee.GetChainId()
}

func (e20de *ERC20DeployedEvent) ChainID() uint32 {
	return e20de.GetChainId()
}

func (sse *SignerSetTxExecutedEvent) ChainID() uint32 {
	return sse.GetChainId()
}
