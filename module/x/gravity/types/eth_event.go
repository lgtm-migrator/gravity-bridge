package types

import (
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// EthereumEvent represents a event on ethereum state
type EthereumEvent interface {
	proto.Message
	// All Ethereum event that we relay from the bridge contract and into the module
	// have a nonce that is monotonically increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what event goes to what nonce means someone is lying.
	// GetEventNonce() uint64
	// The block height that the evented event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all events. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	// GetBlockHeight() uint64
	// Which type of event this is
	GetType() string
	ValidateBasic() error
	EventHash() []byte
}

var (
	_ EthereumEvent = &DepositEvent{}
	_ EthereumEvent = &WithdrawEvent{}
	_ EthereumEvent = &CosmosERC20DeployedEvent{}
	_ EthereumEvent = &LogicCallExecutedEvent{}
)

const (
	TypeMsgWithdrawEvent = "withdraw_event"
	TypeMsgDepositEvent  = "deposit_event"
)

// GetType returns the type of the event
func (e DepositEvent) GetType() string {
	return "deposit"
}

// ValidateBasic performs stateless checks
func (e DepositEvent) ValidateBasic() error {
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if !e.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}
	if err := ValidateEthAddress(e.EthereumSender); err != nil {
		return sdkerrors.Wrap(err, "ethereum sender")
	}
	if _, err := sdk.AccAddressFromBech32(e.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.CosmosReceiver)
	}
	return nil
}

const ()

// EventHash implements BridgeDeposit.EventHash
func (e DepositEvent) EventHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/", e.TokenContract, e.EthereumSender, e.CosmosReceiver)
	return tmhash.Sum([]byte(path))
}

// GetType returns the event type
func (e WithdrawEvent) GetType() string {
	return "withdraw"
}

// ValidateBasic performs stateless checks
func (e WithdrawEvent) ValidateBasic() error {
	if e.BatchNonce == 0 {
		return fmt.Errorf("batch_nonce == 0")
	}
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	return nil
}

// EventHash implements WithdrawBatch.EventHash
func (e WithdrawEvent) EventHash() []byte {
	path := fmt.Sprintf("%s/%d/", e.TokenContract, e.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// EthereumEvent implementation for CosmosERC20DeployedEvent
// ======================================================

// GetType returns the type of the event
func (e CosmosERC20DeployedEvent) GetType() string {
	return "cosmos_erc20_deployed"
}

// ValidateBasic performs stateless checks
func (e CosmosERC20DeployedEvent) ValidateBasic() error {
	if err := ValidateEthAddress(e.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if err := sdk.ValidateDenom(e.CosmosDenom); err != nil {
		return err
	}
	if strings.TrimSpace(e.Name) == "" {
		return fmt.Errorf("token name cannot be blank")
	}
	if strings.TrimSpace(e.Symbol) == "" {
		return fmt.Errorf("token symbol cannot be blank")
	}
	return nil
}

// EventHash implements BridgeDeposit.EventHash
func (e CosmosERC20DeployedEvent) EventHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/%s/%d/", e.CosmosDenom, e.TokenContract, e.Name, e.Symbol, e.Decimals)
	return tmhash.Sum([]byte(path))
}

// EthereumEvent implementation for LogicCallExecutedEvent
// ======================================================

// GetType returns the type of the event
func (e LogicCallExecutedEvent) GetType() string {
	return "logic_call_executed"
}

// ValidateBasic performs stateless checks
func (e LogicCallExecutedEvent) ValidateBasic() error {
	if len(e.InvalidationId) == 0 {
		return fmt.Errorf("invalidation id cannot be blank")
	}
	if e.InvalidationNonce == 0 {
		return fmt.Errorf("invalidation nonce cannot be 0")
	}
	return nil
}

// EventHash implements BridgeDeposit.EventHash
func (e LogicCallExecutedEvent) EventHash() []byte {
	path := fmt.Sprintf("%s/%d/", e.InvalidationId, e.InvalidationNonce)
	return tmhash.Sum([]byte(path))
}
