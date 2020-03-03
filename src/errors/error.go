package errors

import (
	"errors"
)

var (
	// ErrInvalidPoolSection no section for name found in the pool
	ErrInvalidPoolSection = errors.New("Invalid Section")
	// ErrObjectPoolUndeflow no objects can be allocated from the pool
	ErrObjectPoolUndeflow = errors.New("There is not available objects")
	// ErrInvalidLogLevelName invalid log level
	ErrInvalidLogLevelName = errors.New("Could not convert string to log level")
	// ErrInvalidAltcoinTicker ticker string not bound to any altcoin plugin
	ErrInvalidAltcoinTicker = errors.New("Invalid ticker")
	// ErrUnsupportedSigner signer not supported
	ErrUnsupportedSigner = errors.New("Unsupported signer")
	// ErrInvalidTxn invalid transaction
	ErrInvalidTxn = errors.New("Invalid transaction")
	// ErrInvalidOptions invalid options
	ErrInvalidOptions = errors.New("Invalid options")
	// ErrIntegerInputsRequired Input IDs must be integer values
	ErrIntegerInputsRequired = errors.New("Value error: Transaction output references must be integer for signing")
	// ErrBlockNotSet could not find reference to block
	ErrBlockNotSet = errors.New("Block not set or nil")
	// ErrInvalidTransactionNoBlocks unknown number of blocks for unconfirmed transaction
	ErrInvalidTransactionNoBlocks = errors.New("Invalid unconfirmed transaction. Unknown number of blocks")
	// ErrInvalidUnconfirmedTxn invalid unconfirmed transaction
	ErrInvalidUnconfirmedTxn = errors.New("Invalid unconfirmed transaction")
	// ErrInvalidNetworkType invalid network type
	ErrInvalidNetworkType = errors.New("Invalid netType")
	// ErrInvalidID invalid ID
	ErrInvalidID = errors.New("Invalid Id")
	// ErrNotFound target item not found in collection
	ErrNotFound = errors.New("Item not found in collection")
	// ErrParseTxID invalid string value for transaction hash ID
	ErrParseTxID = errors.New("Error parsing transaction hash")
	// ErrParseSHA256 invalid SHA256 hash
	ErrParseSHA256 = errors.New("Error parsing SHA256 hash")
	// ErrParseTxnFee invalid string value for transaction fee
	ErrParseTxnFee = errors.New("Error parsing transaction fee")
	// ErrParseTxnCoins transaction coins can not be parsed
	ErrParseTxnCoins = errors.New("Error parsing transaction coins")
	// ErrInvalidAddressString could not decode base58 address
	ErrInvalidAddressString = errors.New("Error decoding base58 address")
	// ErrTxnSignFailure signing strategy reported an error whie signing transaction
	ErrTxnSignFailure = errors.New("Transaction signing failed for txn")
	// ErrNoMoreElements no more elements in a sequence or succession
	ErrNoMoreElements = errors.New("No more elements in the sequence")
	// ErrUnexpectedUxOutAddress unexpected address
	ErrUnexpectedUxOutAddress = errors.New("Unexpected address")
	// ErrInvalidPoolObjectType clients in the pool do not match expected type
	ErrInvalidPoolObjectType = errors.New("There is not propers client in pool")
	// ErrInvalidWalletEntropy entropy complexity does not match supported values
	ErrInvalidWalletEntropy = errors.New("Entropy must be 128 or 256")
	// ErrInvalidValue invalid value was supplied in to function
	ErrInvalidValue = errors.New("Value errors")
	// ErrInvalidTypeAssertion the current value have not the expected type
	ErrInvalidTypeAssertion = errors.New("Invalid type assertion")
	// ErrWalletCantSign wallet can not sign transactions
	ErrWalletCantSign = errors.New("Wallet does not support transaction signing")
	// ErrNotImplemented feature not implemented
	ErrNotImplemented = errors.New("Feature not implemented")
	// ErrWltFromAddrNotFound trying to find a wallet from the first address in the deterministic sequence failed
	// wallet not found
	ErrWltFromAddrNotFound = errors.New("Error getting wallet from addr, no wallet found")
	// ErrInvalidIndex invalid index detected
	// ErrNotImplemented feature not implemented
	ErrNotImplemented = errors.New("Feature not implemented")
	// ErrKeyNotFound unknown key
	ErrKeyNotFound = errors.New("Key not found")
	ErrInvalidIndex = errors.New("Invalid index")
	// ErrHwUnexpected unexpected error with sky-hw
	ErrHwUnexpected = errors.New("Unexpected error with device")
	// ErrHwSignTransactionCanceled cancelled by user in the physical device
	ErrHwSignTransactionCanceled = errors.New("Sign transaction with hardware wallet has been canceled")
	// ErrNilValue object should not be null
	ErrNilValue = errors.New("Object should not be null")
)
