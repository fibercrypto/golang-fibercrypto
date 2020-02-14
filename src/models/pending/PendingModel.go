package pending

import (
	"github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/models" //callable as skycoin
	"github.com/fibercrypto/fibercryptowallet/src/core"
	local "github.com/fibercrypto/fibercryptowallet/src/main"
	"github.com/fibercrypto/fibercryptowallet/src/util"
	"github.com/fibercrypto/fibercryptowallet/src/util/logging"
	qtCore "github.com/therecipe/qt/core"
)

var logPendingTxn = logging.MustGetLogger("modelsPendingTransaction")

type PendingTransactionList struct {
	qtCore.QObject
	PEX       core.PEX
	WalletEnv core.WalletEnv

	_ func() `constructor:"init"`

	_ []*PendingTransaction            `property:"transactions"`
	_ bool                             `property:"loading"`
	_ func(bool) []*PendingTransaction `slot:"recoverTransactions"`
	_ func()                           `slot:"getAll"`
	_ func()                           `slot:"cleanPendingTxns"`
	_ func()                           `slot:"getMine"`
}

type PendingTransaction struct {
	qtCore.QObject

	_ string            `property:"sky"`
	_ string            `property:"coinHours"`
	_ *qtCore.QDateTime `property:"timeStamp"`
	_ string            `property:"transactionID"`
	_ int               `property:"mine"`
}

func (model *PendingTransactionList) init() {
	model.ConnectGetAll(model.getAll)
	model.ConnectGetMine(model.getMine)
	model.ConnectRecoverTransactions(model.recoverTransactions)
	model.SetLoading(true)
	model.ConnectCleanPendingTxns(model.cleanPendingTxns)

	altManager := local.LoadAltcoinManager()
	walletsEnvs := make([]core.WalletEnv, 0)
	for _, plug := range altManager.ListRegisteredPlugins() {
		walletsEnvs = append(walletsEnvs, plug.LoadWalletEnvs()...)
	}
	model.PEX = &skycoin.SkycoinPEX{}
	model.WalletEnv = walletsEnvs[0]

}

func (model *PendingTransactionList) cleanPendingTxns() {
	model.SetTransactions(make([]*PendingTransaction, 0))
}

func (model *PendingTransactionList) recoverTransactions(mine bool) []*PendingTransaction {
	model.SetLoading(true)
	if mine {
		model.getMine()
	} else {
		model.getAll()
	}
	return model.Transactions()
}

func (model *PendingTransactionList) getAll() {
	logPendingTxn.Info("Getting txn details")

	txns, err := model.PEX.GetTxnPool()
	if err != nil {
		logPendingTxn.WithError(err).Warn("Couldn't get txn pool")
		return
	}
	ptModels := make([]*PendingTransaction, 0)
	for txns.Next() {
		t := txns.Value()
		ptModel := TransactionToPendingTransaction(t)
		ptModel.SetMine(0)
		ptModels = append(ptModels, ptModel)
	}
	model.SetLoading(false)
	model.SetTransactions(ptModels)

}

func (model *PendingTransactionList) getMine() {
	logPendingTxn.Info("Getting txn details")

	wallets := model.WalletEnv.GetWalletSet().ListWallets()
	if wallets == nil {
		logPendingTxn.WithError(nil).Warn("Couldn't load list of wallets")
		return
	}
	ptModels := make([]*PendingTransaction, 0)
	for wallets.Next() {
		wlt := wallets.Value()
		account := wlt.GetCryptoAccount()
		txns, err := account.ListPendingTransactions()
		if err != nil {
			//display an error in qml app when Mine is selected
			logPendingTxn.WithError(err).Warn("Couldn't list pending transactions")
			continue
		}
		for txns.Next() {
			txn := txns.Value()
			if txn.GetStatus() == core.TXN_STATUS_PENDING {
				ptModel := TransactionToPendingTransaction(txn)
				ptModel.SetMine(1)
				ptModels = append(ptModels, ptModel)
			}
		}
	}
	model.SetLoading(false)
	model.SetTransactions(ptModels)
}

func TransactionToPendingTransaction(stxn core.Transaction) *PendingTransaction {
	pt := NewPendingTransaction(nil)
	year, month, day, h, m, s := util.ParseDate(int64(stxn.GetTimestamp()))
	pt.SetTimeStamp(qtCore.NewQDateTime3(qtCore.NewQDate3(year, month, day), qtCore.NewQTime3(h, m, s, 0), qtCore.Qt__LocalTime))
	pt.SetTransactionID(stxn.GetId())
	iterator := skycoin.NewSkycoinTransactionOutputIterator(stxn.GetOutputs())
	sky, coinHours := uint64(0), uint64(0)
	skyErr, coinHoursErr := false, false
	for iterator.Next() {
		output := iterator.Value()
		val, err := output.GetCoins(skycoin.Sky)
		if err != nil {
			logPendingTxn.WithError(err).Warn("Couldn't get Skycoins")
		}
		skyErr = skyErr || (err != nil)
		if !skyErr {
			sky = sky + val
		}
		val, err = output.GetCoins(skycoin.CoinHour)
		if err != nil {
			logPendingTxn.WithError(err).Warn("Couldn't get Coin Hours")
		}
		coinHoursErr = coinHoursErr || (err != nil)
		if !coinHoursErr {
			coinHours = coinHours + val
		}
	}
	skyAccuracy, err := util.AltcoinQuotient(skycoin.Sky)
	if err != nil {
		logPendingTxn.WithError(err).Warn("Couldn't get Skycoins quotient")
	}
	skyErr = skyErr || (err != nil)
	float_sky := ""
	if skyErr {
		float_sky = "NA"
	} else {
		float_sky = util.FormatCoins(sky, skyAccuracy)
	}
	pt.SetSky(float_sky)
	skychAccuracy, err2 := util.AltcoinQuotient(skycoin.CoinHour)
	if err != nil {
		logPendingTxn.WithError(err).Warn("Couldn't get Coin Hours quotient")
	}
	coinHoursErr = coinHoursErr || (err2 != nil)
	uint_ch := ""
	if coinHoursErr {
		uint_ch = "NA"
	} else {
		uint_ch = util.FormatCoins(coinHours, skychAccuracy)
	}
	pt.SetCoinHours(uint_ch)
	return pt
}
