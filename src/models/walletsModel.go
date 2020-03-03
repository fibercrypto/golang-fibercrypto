package models

import (
	hardware "github.com/fibercrypto/fibercryptowallet/src/contrib/skywallet"
	"github.com/fibercrypto/fibercryptowallet/src/util/logging"
	"github.com/fibercrypto/skywallet-go/src/integration/proxy"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/qml"
	skyWallet "github.com/fibercrypto/skywallet-go/src/skywallet"
	wlcore "github.com/fibercrypto/fibercryptowallet/src/main"
	fccore "github.com/fibercrypto/fibercryptowallet/src/core"
	"time"
)

const (
	Name = int(core.Qt__UserRole) + iota + 1
	EncryptionEnabled
	Sky
	CoinHours
	FileName
	Expand
	HasHardwareWallet
)

var logWalletsModel = logging.MustGetLogger("Wallets Model")
var dev skyWallet.Devicer
var hadHwConnected = false
var hwConnectedOn []int

type WalletModel struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*QWallet               `property:"wallets"`

	_             func(*QWallet)                                                                   `slot:"addWallet"`
	_             func(row int, name string, encryptionEnabled bool, sky string, coinHours string) `slot:"editWallet"`
	_             func(row int)                                                                    `slot:"removeWallet"`
	_             func([]*QWallet)                                                                 `slot:"loadModel"`
	_             func([]*QWallet)                                                                 `slot:"updateModel"`
	_ func()                                                                            `slot:"sniffHw"`
	_             func(string)                                                                     `slot:"changeExpanded"`
	_             int                                                                              `property:"count"`
	receivChannel chan *updateWalletInfo
	walletByName  map[string]*QWallet
}

type QWallet struct {
	core.QObject
	_ string `property:"name"`
	_ int    `property:"encryptionEnabled"`
	_ string `property:"sky"`
	_ string `property:"coinHours"`
	_ string `property:"fileName"`
	_ bool   `property:"expand"`
	_ bool    `property:"hasHardwareWallet"`
}

func (walletModel *WalletModel) init() {
	logWalletsModel.Info("Initialize Wallet model")
	dev = proxy.NewSequencer(skyWallet.NewDevice(skyWallet.DeviceTypeUSB), true, func() string{
		return "not implemented"
	})
	walletModel.SetRoles(map[int]*core.QByteArray{
		Name:              core.NewQByteArray2("name", -1),
		EncryptionEnabled: core.NewQByteArray2("encryptionEnabled", -1),
		Sky:               core.NewQByteArray2("sky", -1),
		CoinHours:         core.NewQByteArray2("coinHours", -1),
		FileName:          core.NewQByteArray2("fileName", -1),
		Expand:            core.NewQByteArray2("expand", -1),
		HasHardwareWallet: core.NewQByteArray2("hasHardwareWallet", -1),
	})
	qml.QQmlEngine_SetObjectOwnership(walletModel, qml.QQmlEngine__CppOwnership)
	walletModel.ConnectData(walletModel.data)
	walletModel.ConnectSetData(walletModel.setData)
	walletModel.ConnectRowCount(walletModel.rowCount)
	walletModel.ConnectColumnCount(walletModel.columnCount)
	walletModel.ConnectRoleNames(walletModel.roleNames)

	walletModel.ConnectAddWallet(walletModel.addWallet)
	walletModel.ConnectEditWallet(walletModel.editWallet)
	walletModel.ConnectRemoveWallet(walletModel.removeWallet)
	walletModel.ConnectLoadModel(walletModel.loadModel)
	walletModel.ConnectUpdateModel(walletModel.updateModel)
	walletModel.ConnectChangeExpanded(walletModel.changeExpanded)
	walletModel.receivChannel = walletManager.suscribe()
	walletModel.walletByName = make(map[string]*QWallet, 0)
	go func() {
		for {
			wi := <-walletModel.receivChannel
			if wi.isNew {
				//walletModel.addWallet(wi.wallet)
			} else {
				encrypted := false
				if wi.wallet.EncryptionEnabled() == 1 {
					encrypted = true
				}
				walletModel.editWallet(wi.row, wi.wallet.Name(), encrypted, wi.wallet.Sky(), wi.wallet.CoinHours())
	walletModel.ConnectSniffHw(walletModel.sniffHw)
		}
	}()
}

// attachHwAsSigner add a hw as signer
func attachHwAsSigner(wlt fccore.Wallet, dev skyWallet.Devicer) error {
	hw := hardware.NewSkyWallet(wlt, dev)
	am := wlcore.LoadAltcoinManager()
	if err := am.AttachSignService(hw); err != nil {
		logSignersModel.Errorln("error registering hardware wallet as signer")
		return err
	}
	return nil
}

// sniffHw notify the model about available hardware wallet device if any
func (walletModel *WalletModel) sniffHw() {
	checkForDerivationType := func(dt string) {
		addr, err := hardware.HwFirstAddr(dev, dt)
		if err == nil {
			wlt, err := walletManager.WalletEnv.LookupWallet(addr)
			if err != nil {
				logSignersModel.Warnln("can not find a wallet matching the hardware one")
				// FIXME handle this scenario with a wallet registration.
				return
			}
			err = attachHwAsSigner(wlt, dev)
			if err != nil {
				logSignersModel.WithError(err).Errorln("unable to attach signer")
				return
			}
			hadHwConnected = true
			walletModel.updateWallet(wlt.GetId())
		} else {
			if hadHwConnected {
				hadHwConnected = false
				hwConnectedOn = []int{}
				beginIndex := walletModel.Index(0, 0, core.NewQModelIndex())
				endIndex := walletModel.Index(walletModel.rowCount(core.NewQModelIndex())-1, 0, core.NewQModelIndex())
				walletModel.DataChanged(beginIndex, endIndex, []int{HasHardwareWallet})
				logSignersModel.WithError(err).Info("connection to hardware wallet was lose")
			}
		}
	}
	go func() {
		for {
			hwConnectedOn = []int{}
			checkForDerivationType(skyWallet.WalletTypeDeterministic)
			//checkForDerivationType(skyWallet.WalletTypeBip44)
			time.Sleep(time.Millisecond * 500)
		}
	}()
}

func (walletModel *WalletModel) updateWallet(fn string) {
	index := &core.QModelIndex{}
	for row := 0; row < walletModel.rowCount(core.NewQModelIndex()); row++ {
		index = walletModel.Index(row, 0, core.NewQModelIndex())
		fileName := walletModel.data(index, FileName)
		if  fileName.ToString() == fn {
			hwConnectedOn = append(hwConnectedOn, row)
			walletModel.DataChanged(index, index, []int{HasHardwareWallet})
			break
		}
	}
			}
		}
	}()

}

func (walletModel *WalletModel) changeExpanded(id string) {
	w := walletModel.walletByName[id]
	w.SetExpand(!w.IsExpand())
}

func (walletModel *WalletModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(walletModel.Wallets()) {
		return core.NewQVariant()
	}

	var w = walletModel.Wallets()[index.Row()]

	switch role {
	case Name:
		{
			return core.NewQVariant1(w.Name())
		}

	case EncryptionEnabled:
		{
			return core.NewQVariant1(w.EncryptionEnabled())
		}

	case Sky:
		{
			return core.NewQVariant1(w.Sky())
		}

	case CoinHours:
		{
			return core.NewQVariant1(w.CoinHours())
		}
	case FileName:
		{
			return core.NewQVariant1(w.FileName())
		}
	case HasHardwareWallet:
		{
			valInSlice := func() bool {
				for idx := range hwConnectedOn {
					if hwConnectedOn[idx] == index.Row() && hadHwConnected {
						return true
					}
				}
				return false
			}
			// FIXME: consider a double checking here instead of hadHwConnected
			// be careful this can have a big performance impact
			return core.NewQVariant1(valInSlice())
		}
	case Expand:
		{
			return core.NewQVariant1(w.IsExpand())
		}
	default:
		{
			return core.NewQVariant()
		}
	}
}

func (walletModel *WalletModel) setData(index *core.QModelIndex, value *core.QVariant, role int) bool {

	if !index.IsValid() {
		return false
	}

	if index.Row() >= len(walletModel.Wallets()) {
		return false
	}

	var w = walletModel.Wallets()[index.Row()]

	switch role {
	case Name:
		{
			w.SetName(value.ToString())
		}
	case EncryptionEnabled:
		{
			w.SetEncryptionEnabled(value.ToInt(nil))
		}
	case Sky:
		{
			w.SetSky(value.ToString())
		}
	case CoinHours:
		{
			w.SetCoinHours(value.ToString())
		}
	case FileName:
		{
			w.SetFileName(value.ToString())
		}
	case Expand:
		{
			w.SetExpand(value.ToBool())
		}
	default:
		{
			return false
		}
	}

	walletModel.DataChanged(index, index, []int{role})
	return true
}

func (walletModel *WalletModel) rowCount(parent *core.QModelIndex) int {
	return len(walletModel.Wallets())
}

func (walletModel *WalletModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (walletModel *WalletModel) roleNames() map[int]*core.QByteArray {
	return walletModel.Roles()
}

func (walletModel *WalletModel) addWallet(w *QWallet) {
	logWalletsModel.Info("Add Wallet")
	if w.Pointer() == nil {
		return
	}
	walletModel.walletByName[w.FileName()] = w
	walletModel.BeginInsertRows(core.NewQModelIndex(), len(walletModel.Wallets()), len(walletModel.Wallets()))
	qml.QQmlEngine_SetObjectOwnership(w, qml.QQmlEngine__CppOwnership)
	walletModel.SetWallets(append(walletModel.Wallets(), w))
	walletModel.SetCount(walletModel.Count() + 1)
	walletModel.EndInsertRows()
}

func (walletModel *WalletModel) editWallet(row int, name string, encrypted bool, sky string, coinHours string) {
	logWalletsModel.Info("Edit Wallet")
	pIndex := walletModel.Index(0, 0, core.NewQModelIndex())
	lIndex := walletModel.Index(len(walletModel.Wallets())-1, 0, core.NewQModelIndex())
	w := walletModel.Wallets()[row]
	w.SetName(name)
	if encrypted {
		w.SetEncryptionEnabled(1)
	} else {
		w.SetEncryptionEnabled(0)
	}
	w.SetSky(sky)
	w.SetCoinHours(coinHours)
	walletModel.DataChanged(pIndex, lIndex, []int{Name, EncryptionEnabled, Sky, CoinHours})
}

func (walletModel *WalletModel) removeWallet(row int) {
	logWalletsModel.Info("Remove wallets for index")
	walletModel.BeginRemoveRows(core.NewQModelIndex(), row, row)
	walletModel.SetWallets(append(walletModel.Wallets()[:row], walletModel.Wallets()[row+1:]...))
	walletModel.SetCount(walletModel.Count() - 1)
	walletModel.EndRemoveRows()
}

func (walletModel *WalletModel) updateModel(wallets []*QWallet) {
	for i, wlt := range wallets {
		walletModel.editWallet(i, wlt.Name(), wlt.EncryptionEnabled() == 1, wlt.Sky(), wlt.CoinHours())
	}
}

func (walletModel *WalletModel) loadModel(wallets []*QWallet) {
	logWalletsModel.Info("Loading wallets")
	for _, wlt := range wallets {
		qml.QQmlEngine_SetObjectOwnership(wlt, qml.QQmlEngine__CppOwnership)
	}
	for _, w := range wallets {
		walletModel.walletByName[w.FileName()] = w
	}
	walletModel.BeginResetModel()
	walletModel.SetWallets(wallets)

	walletModel.EndResetModel()
	walletModel.SetCount(len(walletModel.Wallets()))
}
