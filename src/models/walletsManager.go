package models

import (
	"sync"

	"github.com/fibercrypto/fibercryptowallet/src/coin/skycoin"
	"github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/config"

	"github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/params"

	"github.com/fibercrypto/fibercryptowallet/src/util"

	"github.com/therecipe/qt/qml"

	"time"

	sky "github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/models"
	"github.com/fibercrypto/fibercryptowallet/src/core"
	local "github.com/fibercrypto/fibercryptowallet/src/main"
	"github.com/fibercrypto/fibercryptowallet/src/util/logging"
	qtCore "github.com/therecipe/qt/core"
)

var logWalletManager = logging.MustGetLogger("modelsWalletManager")

var once sync.Once
var walletManager *WalletManager

type updateWalletInfo struct {
	isNew  bool
	row    int
	wallet *QWallet
}
type WalletManager struct {
	qtCore.QObject
	WalletEnv                core.WalletEnv
	SeedGenerator            core.SeedGenerator
	wallets                  []*QWallet
	addresseseByWallets      map[string][]*QAddress
	addressesAndWalletsMutex sync.Mutex
	markedAddress            map[string]int
	outputsByAddress         map[string][]*QOutput
	outputsByAddressMutex    sync.Mutex
	altManager               core.AltcoinManager
	signer                   core.BlockchainSignService
	transactionAPI           core.BlockchainTransactionAPI
	walletsIterator          core.WalletIterator
	updaterChannel           chan *updateWalletInfo
	timerUpdate              chan time.Duration

	_ func()                                                                                                                           `slot:"updateWalletEnvs"`
	_ func(wltId, address string)                                                                                                      `slot:"updateOutputs"`
	_ func(string)                                                                                                                     `slot:"updateAddresses"`
	_ func()                                                                                                                           `slot:"updateWallets"`
	_ func()                                                                                                                           `slot:"updateAll"`
	_ func()                                                                                                                           `constructor:"init"`
	_ func(seed string, label string, walletType string, password string, scanN int) *QWallet                                          `slot:"createEncryptedWallet"`
	_ func(seed string, label string, walletType string, scanN int) *QWallet                                                           `slot:"createUnencryptedWallet"`
	_ func(entropy int) string                                                                                                         `slot:"getNewSeed"`
	_ func(seed string) int                                                                                                            `slot:"verifySeed"`
	_ func(id string, n int, password string)                                                                                          `slot:"newWalletAddress"`
	_ func(id string, password string) int                                                                                             `slot:"encryptWallet"`
	_ func(id string, password string) int                                                                                             `slot:"decryptWallet"`
	_ func() []*QWallet                                                                                                                `slot:"getWallets"`
	_ func(id string) []*QAddress                                                                                                      `slot:"getAddresses"`
	_ func(wltIds, addresses []string, source string, pwd interface{}, index []int, qTxn *QTransaction) *QTransaction                  `slot:"signTxn"`
	_ func(wltId string, destinationAddress string, amount string) *QTransaction                                                       `slot:"sendTo"`
	_ func(id, label string) *QWallet                                                                                                  `slot:"editWallet"`
	_ func(wltId, address string) []*QOutput                                                                                           `slot:"getOutputs"`
	_ func(txn *QTransaction) bool                                                                                                     `slot:"broadcastTxn"`
	_ func(wltIds, from, addrTo, skyTo, coinHoursTo []string, change string, automaticCoinHours bool, burnFactor string) *QTransaction `slot:"sendFromAddresses"`
	_ func(wltIds, outs, addrTo, skyTo, coinHoursTo []string, change string, automaticCoinHours bool, burnFactor string) *QTransaction `slot:"sendFromOutputs"`
	_ func() []*QAddress                                                                                                               `slot:"getAllAddresses"`
	_ func(wltId string) []*QOutput                                                                                                    `slot:"getOutputsFromWallet"`
	_ func() string                                                                                                                    `slot:"getDefaultWalletType"`
	_ func(wltIds, addresses []string, source string, bridgeForPassword *QBridge, index []int, qTxn *QTransaction)                     `slot:"signAndBroadcastTxnAsync"`
	_ func() []string                                                                                                                  `slot:"getAvailableWalletTypes"`
	_ func(address string, value int)                                                                                                  `slot:"editMarkAddress"`
	_ func(address string) int                                                                                                         `slot:"markFieldOfAddress"`
}

func (walletM *WalletManager) init() {
	logWalletManager.Info("Initializing WalletManager")
	logWalletManager.Debug("Starting WalletManager")
	once.Do(func() {
		qml.QQmlEngine_SetObjectOwnership(walletM, qml.QQmlEngine__CppOwnership)
		walletM.ConnectEditWallet(walletM.editWallet)
		walletM.ConnectCreateEncryptedWallet(walletM.createEncryptedWallet)
		walletM.ConnectCreateUnencryptedWallet(walletM.createUnencryptedWallet)
		walletM.ConnectGetNewSeed(walletM.getNewSeed)
		walletM.ConnectVerifySeed(walletM.verifySeed)
		walletM.ConnectNewWalletAddress(walletM.newWalletAddress)
		walletM.ConnectEncryptWallet(walletM.encryptWallet)
		walletM.ConnectDecryptWallet(walletM.decryptWallet)
		walletM.ConnectGetWallets(walletM.getWallets)
		walletM.ConnectGetAddresses(walletM.getAddresses)
		walletM.ConnectSendTo(walletM.sendTo)
		walletM.ConnectSignTxn(walletM.signTxn)
		walletM.ConnectGetOutputs(walletM.getOutputs)
		walletM.ConnectSendFromAddresses(walletM.sendFromAddresses)
		walletM.ConnectSendFromOutputs(walletM.sendFromOutputs)
		walletM.ConnectBroadcastTxn(walletM.broadcastTxn)
		walletM.ConnectGetAllAddresses(walletM.getAllAddresses)
		walletM.ConnectGetOutputsFromWallet(walletM.getOutputsFromWallet)
		walletM.ConnectUpdateWalletEnvs(walletM.updateWalletEnvs)
		walletM.ConnectUpdateWallets(walletM.updateWallets)
		walletM.ConnectUpdateAll(walletM.updateAll)
		walletM.ConnectUpdateAddresses(walletM.updateAddresses)
		walletM.ConnectUpdateOutputs(walletM.updateOutputs)
		walletM.ConnectSignAndBroadcastTxnAsync(walletM.signAndBroadcastTxnAsync)
		walletM.ConnectGetDefaultWalletType(walletM.getDefaultWalletType)
		walletM.ConnectGetAvailableWalletTypes(walletM.getAvailableWalletTypes)
		walletM.ConnectEditMarkAddress(walletM.editMarkAddress)
		walletM.ConnectMarkFieldOfAddress(walletM.markFieldOfAddress)
		walletM.addresseseByWallets = make(map[string][]*QAddress, 0)
		walletM.outputsByAddress = make(map[string][]*QOutput, 0)
		walletM.SeedGenerator = new(sky.SeedService)
		walletManager = walletM
		walletM.markedAddress = make(map[string]int)

	})
	walletM.altManager = local.LoadAltcoinManager()
	walletM.updateTransactionAPI()
	walletM.updateSigner()
	walletM.updateWalletEnvs()
	walletM.timerUpdate = make(chan time.Duration)
	for walletM.WalletEnv == nil {
		walletM.updateWalletEnvs()
	}
	walletM = walletManager
	walletM.updaterChannel = make(chan *updateWalletInfo)

	qWallets := make([]*QWallet, 0)
	logWalletManager.Debug("Getting iterator")
	it := walletM.WalletEnv.GetWalletSet().ListWallets()
	logWalletManager.Debug("Obtained iterator")
	if it == nil {
		logWalletManager.WithError(nil).Warn("Couldn't get a wallet iterator")
		return
	}
	for it.Next() {
		qWallet := NewQWallet(nil)
		qml.QQmlEngine_SetObjectOwnership(qWallet, qml.QQmlEngine__CppOwnership)
		qWallet.SetName(it.Value().GetLabel())
		qWallet.SetExpand(false)
		qWallet.SetFileName(it.Value().GetId())
		qWallet.SetEncryptionEnabled(0)

		qWallet.SetSky("N/A")
		qWallet.SetCoinHours("N/A")
		qWallets = append(qWallets, qWallet)
		walletM.initWalletAddresses(it.Value().GetId())
	}
	logWalletManager.Debug("Finish wallets")
	walletM.wallets = qWallets
	go func() {
		updateTime := config.GetDataUpdateTime()
		logWalletManager.Debug("Update time is :=> ", time.Duration(updateTime)*time.Second)
		uptimeTicker := time.NewTicker(time.Duration(updateTime) * time.Second)

		for {
			select {
			case <-uptimeTicker.C:
				go walletM.updateWallets()
				walletManager = walletM
				break
			case t := <-walletM.timerUpdate:
				uptimeTicker = time.NewTicker(t)
				break
			}
			<-uptimeTicker.C
			go walletM.updateWallets()
			walletManager = walletM
		}
	}()
}

func (walletM *WalletManager) suscribe() chan *updateWalletInfo {
	return walletM.updaterChannel
}
func (walletManager *WalletManager) editMarkAddress(address string, value int) {
	walletManager.markedAddress[address] = value
}

func (walletM *WalletManager) markFieldOfAddress(address string) int {
	val, ok := walletM.markedAddress[address]
	if !ok {
		val = 0
	}
	return val
}

func (walletM *WalletManager) updateAll() {
	logWalletManager.Debug("Updating Wallet manager")
	walletM.altManager = local.LoadAltcoinManager()
	walletM.updateTransactionAPI()
	walletM.updateSigner()
	walletM.updateWalletEnvs()
	skycoin.UpdateAltcoin()
	updateTime := config.GetDataUpdateTime()
	logWalletManager.Debug("Update time is :=> ", time.Duration(updateTime)*time.Second)
	walletM.timerUpdate <- time.Duration(updateTime) * time.Second
}

func GetWalletEnv() core.WalletEnv {
	logWalletManager.Info("Getting WalletEnv")
	wm := GetWalletManager()
	if wm == nil {
		return nil
	}

	return wm.WalletEnv
}

func GetWalletManager() *WalletManager {
	return walletManager
}

func (walletM *WalletManager) getDefaultWalletType() string {
	return walletM.WalletEnv.GetWalletSet().DefaultWalletType()
}

func (walletM *WalletManager) getAvailableWalletTypes() []string {
	return walletM.WalletEnv.GetWalletSet().SupportedWalletTypes()
}

func (walletM *WalletManager) updateSigner() {

	logWalletManager.Info("Updating Signers")
	signers := make([]core.BlockchainSignService, 0)

	for _, plug := range walletM.altManager.ListRegisteredPlugins() {
		sing, err := plug.LoadSignService()
		if err != nil {
			logWalletManager.WithError(err).Errorf("Error loading signer from %s plugin", plug.GetName())
		}
		signers = append(signers, sing)
	}

	walletM.signer = signers[0]
}

func (walletM *WalletManager) updateTransactionAPI() {
	logWalletManager.Info("Updating TransactionAPI")
	txnAPIS := make([]core.BlockchainTransactionAPI, 0)

	for _, plug := range walletM.altManager.ListRegisteredPlugins() {
		txnAPI, err := plug.LoadTransactionAPI("MainNet")
		if err != nil {
			logWalletManager.WithError(err).Errorf("Error loading transaction API from %s plugin", plug.GetName())
		}
		txnAPIS = append(txnAPIS, txnAPI)
	}

	walletM.transactionAPI = txnAPIS[0]
}
func (walletM *WalletManager) updateWalletEnvs() {
	logWalletManager.Info("Updating WalletEnvs")
	walletsEnvs := make([]core.WalletEnv, 0)

	for _, plug := range walletM.altManager.ListRegisteredPlugins() {
		walletsEnvs = append(walletsEnvs, plug.LoadWalletEnvs()...)
	}
	if len(walletsEnvs) == 0 {
		logWalletManager.Error("Error loading wallet envs")
		return
	}
	walletM.WalletEnv = walletsEnvs[0]
}

func (walletM *WalletManager) initWalletAddresses(wltId string) {
	logWalletManager.Info("Updating Addresses")
	wlt := walletM.WalletEnv.GetWalletSet().GetWallet(wltId)
	qAddresses := make([]*QAddress, 0)
	it, err := wlt.GetLoadedAddresses()
	if err != nil {
		logWalletManager.WithError(err).Warn("Couldn't loaded addresses")
		return
	}
	for it.Next() {
		addr := it.Value()
		qAddress := NewQAddress(nil)
		qml.QQmlEngine_SetObjectOwnership(qAddress, qml.QQmlEngine__CppOwnership)
		qAddress.SetAddress(addr.String())
		qAddress.SetMarked(0)
		qAddress.SetWallet(wlt.GetLabel())
		qAddress.SetWalletId(wlt.GetId())
		qAddress.SetAddressSky("N/A")
		qAddress.SetAddressCoinHours("N/A")
		qml.QQmlEngine_SetObjectOwnership(qAddress, qml.QQmlEngine__CppOwnership)

		qAddresses = append(qAddresses, qAddress)
		go walletM.getOutputs(wltId, addr.String())
	}

	walletM.addresseseByWallets[wltId] = qAddresses

}

func (walletM *WalletManager) updateAddresses(wltId string) {
	logWalletManager.Info("Updating Addresses")
	wlt := walletM.WalletEnv.GetWalletSet().GetWallet(wltId)
	if wlt == nil {
		logWalletManager.Warn("Couldn't load wallet")
		return
	}
	qAddresses := make([]*QAddress, 0)
	it, err := wlt.GetLoadedAddresses()
	if err != nil {
		logWalletManager.WithError(err).Warn("Couldn't loaded addresses")
		return
	}
	for it.Next() {
		addr := it.Value()
		qAddress := NewQAddress(nil)
		qml.QQmlEngine_SetObjectOwnership(qAddress, qml.QQmlEngine__CppOwnership)
		qAddress.SetAddress(addr.String())
		qAddress.SetMarked(0)
		qAddress.SetWallet(wlt.GetLabel())
		qAddress.SetWalletId(wlt.GetId())
		skyFl, err := addr.GetCryptoAccount().GetBalance(params.SkycoinTicker)
		if err != nil {
			qAddress.SetAddressSky("N/A")
			qAddress.SetAddressCoinHours("N/A")
			logWalletManager.WithError(err).Warn("Couldn't load " + params.SkycoinTicker + " Balance")
			qAddresses = append(qAddresses, qAddress)
			continue
		}
		accuracy, err := util.AltcoinQuotient(params.SkycoinTicker)
		if err != nil {
			qAddress.SetAddressSky("N/A")
			qAddress.SetAddressCoinHours("N/A")
			logWalletManager.WithError(err).Warn("Couldn't load " + params.SkycoinTicker + " quotient")
			qAddresses = append(qAddresses, qAddress)
			continue
		}
		qAddress.SetAddressSky(util.FormatCoins(skyFl, accuracy))
		coinH, err := addr.GetCryptoAccount().GetBalance(params.CoinHoursTicker)
		if err != nil {
			qAddress.SetAddressCoinHours("N/A")
			logWalletManager.WithError(err).Warn("Couldn't load " + params.CoinHoursTicker + " Balance")
			qAddresses = append(qAddresses, qAddress)
			continue
		}
		accuracy, err = util.AltcoinQuotient(params.CoinHoursTicker)
		if err != nil {
			qAddress.SetAddressCoinHours("N/A")
			logWalletManager.WithError(err).Warn("Couldn't load " + params.CoinHoursTicker + " quotient")
			qAddresses = append(qAddresses, qAddress)
			continue
		}
		qAddress.SetAddressCoinHours(util.FormatCoins(coinH, accuracy))
		qml.QQmlEngine_SetObjectOwnership(qAddress, qml.QQmlEngine__CppOwnership)

		qAddresses = append(qAddresses, qAddress)

	}
	walletM.addressesAndWalletsMutex.Lock()
	walletM.addresseseByWallets[wltId] = qAddresses
	walletM.addressesAndWalletsMutex.Unlock()
}

func (walletM *WalletManager) updateOutputs(wltId, address string) {
	outs := make([]*QOutput, 0)
	addressIterator, err := walletM.WalletEnv.GetWalletSet().GetWallet(wltId).GetLoadedAddresses()
	if err != nil {
		logWalletManager.WithError(err).Warn("Couldn't get an addresses iterator")
		walletM.outputsByAddressMutex.Lock()
		walletM.outputsByAddress[address] = outs
		walletM.outputsByAddressMutex.Unlock()
		return
	}
	var addr core.Address
	for addressIterator.Next() {
		if addressIterator.Value().String() == address {
			addr = addressIterator.Value()
			break
		}
	}
	if addr == nil {
		logWalletManager.WithError(err).Warn("Couldn't get address")
		walletM.outputsByAddressMutex.Lock()
		walletM.outputsByAddress[address] = outs
		walletM.outputsByAddressMutex.Unlock()
		return
	}
	outsIter, err := addr.GetCryptoAccount().ScanUnspentOutputs()
	if err != nil {
		logWalletManager.WithError(err).Warn("Couldn't scan unspent outputs")
		walletM.outputsByAddressMutex.Lock()
		walletM.outputsByAddress[address] = outs
		walletM.outputsByAddressMutex.Unlock()
		return
	}
	for outsIter.Next() {
		qout := NewQOutput(nil)
		qml.QQmlEngine_SetObjectOwnership(qout, qml.QQmlEngine__CppOwnership)
		qout.SetOutputID(outsIter.Value().GetId())
		skyV, err := outsIter.Value().GetCoins(sky.Sky)
		if err != nil {
			qout.SetAddressSky("N/A")
			logWalletManager.WithError(err).Warn("Couldn't get " + sky.Sky + " coins")
			continue
		}
		quotient, err := util.AltcoinQuotient(sky.Sky)
		if err != nil {
			qout.SetAddressSky("N/A")
			logWalletManager.WithError(err).Warn("Couldn't get " + sky.Sky + " quotient")
			continue
		}
		sSky := util.FormatCoins(skyV, quotient)
		qout.SetAddressSky(sSky)
		ch, err := outsIter.Value().GetCoins(sky.CoinHour)
		if err != nil {
			qout.SetAddressCoinHours("N/A")
			logWalletManager.WithError(err).Warn("Couldn't get " + sky.CoinHour + " coins")
			continue
		}
		quotient, err = util.AltcoinQuotient(sky.CoinHour)
		if err != nil {
			qout.SetAddressCoinHours("N/A")
			logWalletManager.WithError(err).Warn("Couldn't get " + sky.Sky + " quotient")
			continue
		}
		sCh := util.FormatCoins(ch, quotient)
		qout.SetAddressCoinHours(sCh)
		qout.SetAddressOwner(addr.String())
		qout.SetWalletOwner(wltId)
		outs = append(outs, qout)
	}
	walletM.outputsByAddressMutex.Lock()
	walletM.outputsByAddress[address] = outs
	walletM.outputsByAddressMutex.Unlock()
}

func (walletM *WalletManager) updateWallets() {

	it := walletM.WalletEnv.GetWalletSet().ListWallets()
	if it == nil {
		logWalletManager.WithError(nil).Warn("Couldn't get a wallet iterator")
		return
	}
	for it.Next() {

		go walletM.updateAddresses(it.Value().GetId())

		encrypted, err := walletM.WalletEnv.GetStorage().IsEncrypted(it.Value().GetId())
		if err != nil {
			logWalletManager.WithError(err).Warn("Couldn't get wallet by id")
			continue
		}
		qw := fromWalletToQWallet(it.Value(), encrypted, false)
		founded, changed := false, false
		row := 0
		for i := range walletM.wallets {
			if walletM.wallets[i].FileName() == qw.FileName() {
				row = i
				founded = true
				if !isEqual(walletM.wallets[i], qw) && (qw.Sky() != "N/A" || qw.CoinHours() != "N/A") {
					qw.SetExpand(walletM.wallets[i].IsExpand())
					walletM.wallets[i] = qw
					changed = true
				}
				break
			}
		}
		if !founded {
			changed = true
			walletM.wallets = append(walletM.wallets, qw)
		}

		if changed {
			wi := &updateWalletInfo{
				isNew:  !founded,
				row:    row,
				wallet: qw,
			}
			walletManager.updaterChannel <- wi
		}
	}
}

func (walletM *WalletManager) getAllAddresses() []*QAddress {
	logWalletManager.Info("Getting all addresses")
	qAddresses := make([]*QAddress, 0)
	wltIter := walletM.WalletEnv.GetWalletSet().ListWallets()
	if wltIter == nil {
		logWalletManager.Warn("Error getting the list of wallets")
		return nil
	}
	for wltIter.Next() {
		qAddresses = append(qAddresses, walletM.getAddresses(wltIter.Value().GetId())...)
	}
	return qAddresses
}

func isEqual(a, b *QWallet) bool {
	return a.Name() == b.Name() && a.Sky() == b.Sky() && a.CoinHours() == b.CoinHours() && a.EncryptionEnabled() == b.EncryptionEnabled()
}
func (walletM *WalletManager) broadcastTxn(txn *QTransaction) bool {
	logWalletManager.Info("Broadcasting transaction")
	altManager := local.LoadAltcoinManager()
	plug, _ := altManager.LookupAltcoinPlugin(params.SkycoinTicker)
	pex, err := plug.LoadPEX("MainNet")
	if err != nil {
		logWalletManager.WithError(err).Warn("Error loading PEX")
		return false
	}
	isSigned, err := txn.txn.IsFullySigned()
	if err != nil {
		logWalletManager.WithError(err).Warn("Error checking if transaction if fully signed")
		return false
	}
	if !isSigned {
		logWalletManager.Warn("Transaction is not fully signed")
		return false
	}
	err = pex.BroadcastTxn(txn.txn)
	if err != nil {
		logWalletManager.WithError(err).Warn("Error broadcasting transaction")
		return false
	}
	logWalletManager.Info("Transaction Injected")
	return true
}

func (walletM *WalletManager) sendFromOutputs(wltIds []string, from, addrTo, skyTo, coinHoursTo []string, change string, automaticCoinHours bool, burnFactor string) *QTransaction {
	logWalletManager.Info("Creating transaction")
	wltCache := make(map[string]core.Wallet, 0)
	wlts := make([]core.Wallet, 0)
	for _, wltId := range wltIds {
		var wlt core.Wallet
		wlt, exist := wltCache[wltId]
		if !exist {
			wlt = walletM.WalletEnv.GetWalletSet().GetWallet(wltId)
			if wlt == nil {
				logWalletManager.Warn("Couldn't load wallet to create transaction")
				return nil
			}
			wltCache[wltId] = wlt
		}
		wlts = append(wlts, wlt)
	}

	outputsFrom := make([]core.TransactionOutput, 0)
	for _, outAddr := range from {
		out := util.NewGenericOutput(nil, outAddr)
		outputsFrom = append(outputsFrom, &out)
	}
	outputsTo := make([]core.TransactionOutput, 0)
	for i := 0; i < len(addrTo); i++ {
		ch := ""
		if !automaticCoinHours {
			ch = coinHoursTo[i]
		}
		addr := util.NewGenericAddress(addrTo[i])
		out := util.NewGenericOutput(&addr, "")
		// FIXME: Remove explicit reference to Skycoin
		err := out.PushCoins(sky.Sky, skyTo[i])
		if err != nil {
			logWalletManager.WithError(err).Warn("Error parsing value for %s", sky.Sky)
			return nil
		}
		// FIXME: Remove explicit reference to Skycoin
		err = out.PushCoins(sky.CoinHour, ch)
		if err != nil {
			logWalletManager.WithError(err).Warn("Error parsing value for %s", sky.Sky)
			return nil
		}
		outputsTo = append(outputsTo, &out)
	}
	changeAddr := util.NewGenericAddress(change)
	opt := util.NewKeyValueMap()
	opt.SetValue("BurnFactor", burnFactor)
	if automaticCoinHours {
		opt.SetValue("CoinHoursSelectionType", "auto")
	} else {
		opt.SetValue("CoinHoursSelectionType", "manual")
	}
	var txn core.Transaction
	var err error
	if len(wltCache) > 1 {
		walletsOutputs := make([]core.WalletOutput, 0)
		for i, wlt := range wlts {
			walletsOutputs = append(walletsOutputs, &util.SimpleWalletOutput{
				Wallet: wlt,
				UxOut:  outputsFrom[i],
			})
		}
		txn, err = walletM.transactionAPI.Spend(walletsOutputs, outputsTo, &changeAddr, opt)
	} else {
		txn, err = wlts[0].Spend(outputsFrom, outputsTo, &changeAddr, opt)
	}

	if err != nil {
		logWalletManager.WithError(err).Info("Error creating transaction")
		return nil
	}

	qTransaction, err := NewQTransactionFromTransaction(txn)
	if err != nil {
		logWalletManager.WithError(err).Info("Error converting transaction")
		return nil
	}
	return qTransaction
}
func (walletM *WalletManager) sendFromAddresses(wltIds []string, from, addrTo, skyTo, coinHoursTo []string, change string, automaticCoinHours bool, burnFactor string) *QTransaction {
	wltCache := make(map[string]core.Wallet, 0)
	wlts := make([]core.Wallet, 0)
	for _, wltId := range wltIds {
		var wlt core.Wallet
		wlt, exist := wltCache[wltId]
		if !exist {
			wlt = walletM.WalletEnv.GetWalletSet().GetWallet(wltId)
			if wlt == nil {
				logWalletManager.Warn("Couldn't load wallet to create transaction")
				return nil
			}
			wltCache[wltId] = wlt
		}
		wlts = append(wlts, wlt)
	}

	addrsFrom := make([]core.Address, 0)
	for _, addr := range from {

		addrsFrom = append(addrsFrom, &util.GenericAddress{addr})
	}
	outputsTo := make([]core.TransactionOutput, 0)
	for i := 0; i < len(addrTo); i++ {
		ch := ""
		if !automaticCoinHours {
			ch = coinHoursTo[i]
		}
		addr := util.NewGenericAddress(addrTo[i])
		out := util.NewGenericOutput(&addr, "")
		// FIXME: Remove explicit reference to Skycoin
		err := out.PushCoins(sky.Sky, skyTo[i])
		if err != nil {
			logWalletManager.WithError(err).Warnf("Error parsing value for %s", sky.Sky)
			return nil
		}
		// FIXME: Remove explicit reference to Skycoin
		err = out.PushCoins(sky.CoinHour, ch)
		if err != nil {
			logWalletManager.WithError(err).Warnf("Error parsing value for %s", sky.Sky)
			return nil
		}
		outputsTo = append(outputsTo, &out)
	}
	changeAddr := &util.GenericAddress{change}

	opt := util.NewKeyValueMap()
	opt.SetValue("BurnFactor", burnFactor)
	if automaticCoinHours {
		opt.SetValue("CoinHoursSelectionType", "auto")
	} else {
		opt.SetValue("CoinHoursSelectionType", "manual")
	}
	var txn core.Transaction
	var err error
	if len(wltCache) > 1 {
		walletsAddresses := make([]core.WalletAddress, 0)
		for i, wlt := range wlts {
			walletsAddresses = append(walletsAddresses, &util.SimpleWalletAddress{
				Wallet: wlt,
				UxOut:  addrsFrom[i],
			})
		}
		txn, err = walletM.transactionAPI.SendFromAddress(walletsAddresses, outputsTo, changeAddr, opt)
	} else {
		txn, err = wlts[0].SendFromAddress(addrsFrom, outputsTo, changeAddr, opt)
	}

	if err != nil {
		logWalletManager.WithError(err).Info("Error creating transaction")
		return nil
	}

	qtxn, err := NewQTransactionFromTransaction(txn)
	if err != nil {
		logWalletManager.WithError(err).Info("Error converting transaction")
		return nil
	}
	logWalletManager.Info("Transaction created")
	return qtxn

}

func (walletM *WalletManager) getOutputs(wltId, address string) []*QOutput {
	outs, ok := walletM.outputsByAddress[address]
	if !ok {
		walletM.updateOutputs(wltId, address)
	}
	return outs
}

func (walletM *WalletManager) getOutputsFromWallet(wltId string) []*QOutput {
	logWalletManager.Info("Getting Outputs from wallet by Id")
	outs := make([]*QOutput, 0)
	addrIter, err := walletM.WalletEnv.GetWalletSet().GetWallet(wltId).GetLoadedAddresses()
	if err != nil {
		logWalletManager.WithError(err).Warn("Couldn't load addresses iterator")
		return nil
	}
	for addrIter.Next() {
		outs = append(outs, walletM.getOutputs(wltId, addrIter.Value().String())...)
	}
	logWalletManager.Info("Loaded all outputs")
	return outs
}

func (walletM *WalletManager) sendTo(wltId, destinationAddress, amount string) *QTransaction {
	logWalletManager.Info("Creating Transaction")
	wlt := walletM.WalletEnv.GetWalletSet().GetWallet(wltId)
	addr := util.NewGenericAddress(destinationAddress)
	opt := util.NewKeyValueMap()
	opt.SetValue("BurnFactor", "0.5")
	opt.SetValue("CoinHoursSelectionType", "auto")
	if wlt == nil {
		logWalletManager.Warn("Couldn't load wallet to create transaction")
		return nil
	}
	txOut := util.NewGenericOutput(&addr, "")
	// FIXME: Remove explicit reference to Skycoin
	err := txOut.PushCoins(sky.Sky, amount)
	if err != nil {
		logWalletManager.WithError(err).Warn("Error parsing value for %s", sky.Sky)
		return nil
	}
	txn, err := wlt.Transfer(&txOut, opt)
	if err != nil {
		logWalletManager.WithError(err).Warn("Couldn't create transaction")
		return nil
	}
	qTxn, err := NewQTransactionFromTransaction(txn)
	if err != nil {
		logWalletManager.WithError(err).Warn("Couldn't convert transaction")
		return nil
	}
	logWalletManager.Info("Transaction created")
	return qTxn

}

func (walletM *WalletManager) signTxn(wltIds, address []string, source string, tmpPwd interface{}, index []int, qTxn *QTransaction) *QTransaction {
	pwd, isPwdReader := tmpPwd.(core.PasswordReader)
	if !isPwdReader {
		return nil
	}
	logWalletManager.Info("Signig transaction")

	if len(wltIds) != len(address) {
		logWalletManager.Error("Wallets and addresses provided are incorrect")
		return nil
	}
	wltCache := make(map[string]core.Wallet)
	wltByAddr := make(map[string]core.Wallet)
	wlts := make([]core.Wallet, 0)

	for i, wltId := range wltIds {
		var wlt core.Wallet
		wlt, exist := wltCache[wltId]
		if !exist {
			wlt = walletM.WalletEnv.GetWalletSet().GetWallet(wltId)
			if wlt == nil {
				logWalletManager.Warn("Couldn't load wallet to Sign transaction")
				return nil
			}
			wltCache[wltId] = wlt
		}
		wltByAddr[address[i]] = wlt
		wlts = append(wlts, wlt)
	}

	var txn core.Transaction
	var err error

	if len(wltCache) > 1 {
		signDescriptors := make([]core.InputSignDescriptor, 0)
		for _, in := range qTxn.txn.GetInputs() {
			sd := core.InputSignDescriptor{
				InputIndex: in.GetId(),
				SignerID:   core.UID(source),
				Wallet:     wltByAddr[in.GetSpentOutput().GetAddress().String()],
			}
			signDescriptors = append(signDescriptors, sd)
		}
		txn, err = walletM.signer.Sign(qTxn.txn, signDescriptors, pwd)
	} else {
		signer, err2 := util.LookupSignServiceForWallet(wlts[0], core.UID(source))
		if err2 != nil {
			logWalletManager.WithError(err).Warnf("No signer %s for wallet %v", source, wlts[0])
			return nil
		}
		txn, err = wlts[0].Sign(qTxn.txn, signer, pwd, nil)
	}

	if err != nil {
		logWalletManager.WithError(err).Warn("Error signing txn")
		return nil
	}

	qTxn, err = NewQTransactionFromTransaction(txn)
	if err != nil {
		logWalletManager.WithError(err).Warn("Error converting transaction")
		return nil
	}

	return qTxn

}
func (walletM *WalletManager) signAndBroadcastTxnAsync(wltIds, addresses []string, source string, bridgeForPassword *QBridge, index []int, qTxn *QTransaction) {
	channel := make(chan *QTransaction)
	go func() {
		var pwd core.PasswordReader = func(message string, ctx core.KeyValueStore) (string, error) {
			bridgeForPassword.BeginUse()
			defer bridgeForPassword.EndUse()
			bridgeForPassword.lock()
			suffix := ""
			v := ctx.GetValue(core.StrWalletLabel)
			if v == nil {
				v = ctx.GetValue(core.StrWalletName)
			}
			if v != nil {
				if str, isStr := v.(string); isStr {
					suffix = " for " + str
				}
			}
			bridgeForPassword.GetPassword(message + suffix)
			bridgeForPassword.lock()
			pass := bridgeForPassword.getResult()
			bridgeForPassword.unlock()
			return pass, nil
		}

		channel <- walletM.signTxn(wltIds, addresses, source, pwd, index, qTxn)
	}()

	go func() {
		txn := <-channel
		if txn != nil {
			walletM.broadcastTxn(txn)
		}
	}()
}

func (walletM *WalletManager) createEncryptedWallet(seed, label, wltType, password string, scanN int) *QWallet {
	logWalletManager.Info("Creating encrypted wallet")
	pwd := util.ConstantPassword(password)
	// NOTE: No easy way to get plain passwords in memory
	password = ""
	wlt, err := walletM.WalletEnv.GetWalletSet().CreateWallet(label, seed, wltType, true, pwd, scanN)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't create encrypted wallet")
		return nil
	}

	logWalletManager.Info("Created encrypted wallet")
	qWallet := fromWalletToQWallet(wlt, true, false)
	walletM.wallets = append(walletM.wallets, qWallet)
	return qWallet
}

func (walletM *WalletManager) createUnencryptedWallet(seed, label, wltType string, scanN int) *QWallet {
	logWalletManager.Info("Creating encrypted wallet")
	pwd := util.EmptyPassword
	wlt, err := walletM.WalletEnv.GetWalletSet().CreateWallet(label, seed, wltType, false, pwd, scanN)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't create unencrypted wallet")
		return nil
	}
	logWalletManager.Info("Created unencrypted wallet")

	qWallet := fromWalletToQWallet(wlt, true, false)
	walletM.wallets = append(walletM.wallets, qWallet)
	return qWallet

}

func (walletM *WalletManager) getNewSeed(entropy int) string {
	logWalletManager.Info("Getting new seed")
	seed, err := walletM.SeedGenerator.GenerateMnemonic(entropy)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't get new seed")
		return ""
	}
	logWalletManager.Info("Generated new seed")
	return seed
}

func (walletM *WalletManager) verifySeed(seed string) int {
	logWalletManager.Info("Verifying seed")
	ok, err := walletM.SeedGenerator.VerifyMnemonic(seed)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't verify seed")
		return 0
	}
	logWalletManager.Info("Verified seed")
	if ok {
		return 1
	}
	return 0

}

func (walletM *WalletManager) encryptWallet(id, password string) int {
	logWalletManager.Info("Encrypting wallet")
	pwd := util.ConstantPassword(password)
	// NOTE: No easy way to get plain passwords in memory
	password = ""
	walletM.WalletEnv.GetStorage().Encrypt(id, pwd)
	ret, err := walletM.WalletEnv.GetStorage().IsEncrypted(id)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't create encrypted wallets")
	}
	logWalletManager.Info("Wallet encrypted")
	if ret {
		return 1
	}
	return 0
}

func (walletM *WalletManager) decryptWallet(id, password string) int {
	logWalletManager.Info("Decrypt wallet")
	pwd := util.ConstantPassword(password)
	// NOTE: No easy way to get plain passwords in memory
	password = ""
	walletM.WalletEnv.GetStorage().Decrypt(id, pwd)
	ret, err := walletM.WalletEnv.GetStorage().IsEncrypted(id)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't decrypt wallet")
	}
	logWalletManager.Info("Wallet decrypted")
	if ret {
		return 1
	}
	return 0
}

func (walletM *WalletManager) newWalletAddress(id string, n int, password string) {
	logWalletManager.Info("Creating new wallet addresses")
	wlt := walletM.WalletEnv.GetWalletSet().GetWallet(id)
	pwd := util.ConstantPassword(password)
	// NOTE: No easy way to get plain passwords in memory
	password = ""
	wltEntriesLen := 0
	it, err := wlt.GetLoadedAddresses()
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't load addresses")
		return
	}
	for it.Next() {
		wltEntriesLen++
	}
	wlt.GenAddresses(core.AccountAddress, uint32(wltEntriesLen), uint32(n), pwd)
	go walletM.updateAddresses(id)
	logWalletManager.Info("New addresses created")
}

func (walletM *WalletManager) getWallets() []*QWallet {
	if walletM.wallets == nil {
		return make([]*QWallet, 0)
	}
	return walletM.wallets
}

func (walletM *WalletManager) editWallet(id, label string) *QWallet {
	logWalletManager.Info("Editing wallet")

	wlt := walletM.WalletEnv.GetWalletSet().GetWallet(id)
	wlt.SetLabel(label)
	wlt = walletM.WalletEnv.GetWalletSet().GetWallet(id)
	encrypted, err := walletM.WalletEnv.GetStorage().IsEncrypted(wlt.GetId())
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't edit wallet")
		return nil
	}
	qWallet := fromWalletToQWallet(wlt, encrypted, false)
	logWalletManager.Info("Wallet edited")
	return qWallet
}

func (walletM *WalletManager) getAddresses(Id string) []*QAddress {
	addrs, ok := walletM.addresseseByWallets[Id]
	if !ok {
		walletM.updateAddresses(Id)
		addrs = walletM.addresseseByWallets[Id]
	}
	return addrs
}

func fromWalletToQWallet(wlt core.Wallet, isEncrypted, withoutBalance bool) *QWallet {

	qWallet := NewQWallet(nil)
	qml.QQmlEngine_SetObjectOwnership(qWallet, qml.QQmlEngine__CppOwnership)
	qWallet.SetName(wlt.GetLabel())
	qWallet.SetExpand(false)

	qWallet.SetFileName(wlt.GetId())

	qWallet.SetEncryptionEnabled(0)
	if isEncrypted {
		qWallet.SetEncryptionEnabled(1)
	}

	if withoutBalance {
		qWallet.SetSky("N/A")
		qWallet.SetCoinHours("N/A")
		logWalletManager.Info("Passing over default wallet, without balance")
		return qWallet
	}

	bl, err := wlt.GetCryptoAccount().GetBalance(sky.Sky)
	if err != nil {
		qWallet.SetSky("N/A")
		qWallet.SetCoinHours("N/A")
		logWalletManager.WithError(err).Error("Couldn't get Skycoin balance")
		return qWallet
	}

	accuracy, err := util.AltcoinQuotient(params.SkycoinTicker)
	if err != nil {
		qWallet.SetSky("N/A")
		qWallet.SetCoinHours("N/A")
		logWalletManager.WithError(err).Error("Couldn't get Skycoin Altcoin quotient")
		return qWallet
	}

	qWallet.SetSky(util.FormatCoins(bl, accuracy))

	bl, err = wlt.GetCryptoAccount().GetBalance(sky.CoinHoursTicker)
	if err != nil {
		qWallet.SetCoinHours("N/A")
		logWalletManager.WithError(err).Error("Couldn't get Coin Hours balance")
		return qWallet
	}
	accuracy, err = util.AltcoinQuotient(params.CoinHoursTicker)
	if err != nil {
		qWallet.SetCoinHours("N/A")
		logWalletManager.WithError(err).Error("Couldn't get Coin Hours Altcoin quotient")
		panic(nil)
		return qWallet
	}
	qWallet.SetCoinHours(util.FormatCoins(bl, accuracy))
	return qWallet
}
