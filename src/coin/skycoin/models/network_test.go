package skycoin

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/readable"
)

func TestSkycoinPEXGetTxnPool(t *testing.T) {
	global_mock.On("PendingTransactionsVerbose").Return(
		[]readable.UnconfirmedTransactionVerbose{
			readable.UnconfirmedTransactionVerbose{
				Transaction: readable.BlockTransactionVerbose{
					Out: []readable.TransactionOutput{
						readable.TransactionOutput{
							Coins: "1",
							Hours: 2000,
						},
						readable.TransactionOutput{
							Coins: "1",
							Hours: 2000,
						},
					},
				},
			},
			readable.UnconfirmedTransactionVerbose{
				Transaction: readable.BlockTransactionVerbose{
					Out: []readable.TransactionOutput{
						readable.TransactionOutput{
							Coins: "1",
							Hours: 2000,
						},
						readable.TransactionOutput{
							Coins: "1",
							Hours: 2000,
						},
					},
				},
			},
		}, nil)

	pex := &SkycoinPEX{poolSection: PoolSection}
	txns, err2 := pex.GetTxnPool()
	require.Nil(t, err2)

	for txns.Next() {
		iter := NewSkycoinTransactionOutputIterator(txns.Value().GetOutputs())
		for iter.Next() {
			output := iter.Value()
			val, err3 := output.GetCoins(Sky)
			require.Nil(t, err3)
			require.Equal(t, val, uint64(1000000))
			val, err3 = output.GetCoins(CoinHour)
			require.Nil(t, err3)
			require.Equal(t, val, uint64(2000))
		}
	}
}
