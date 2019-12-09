package skycoin

import (
	"github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/config"
	"github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/params"
	"github.com/fibercrypto/fibercryptowallet/src/core"
	"github.com/fibercrypto/fibercryptowallet/src/errors"
	//local "github.com/fibercrypto/fibercryptowallet/src/main"
)

type SkyFiberPlugin struct {
	Params params.SkyFiberParams
}

func (p *SkyFiberPlugin) ListSupportedAltcoins() []core.AltcoinMetadata {
	return []core.AltcoinMetadata{
		core.AltcoinMetadata{
			Name:     SkycoinName,
			Ticker:   SkycoinTicker,
			Family:   SkycoinFamily,
			HasBip44: false,
			Accuracy: 6,
		},
		core.AltcoinMetadata{
			Name:     CoinHoursName,
			Ticker:   CoinHoursTicker,
			Family:   CoinHoursFamily,
			HasBip44: false,
			Accuracy: 0,
		},
		core.AltcoinMetadata{
			Name:     CalculatedHoursName,
			Ticker:   CalculatedHoursTicker,
			Family:   CalculatedHoursFamily,
			HasBip44: false,
			Accuracy: 0,
		},
	}
}

func (p *SkyFiberPlugin) ListSupportedFamilies() []string {
	return []string{SkycoinFamily}
}

func (p *SkyFiberPlugin) RegisterTo(manager core.AltcoinManager) {
	for _, info := range p.ListSupportedAltcoins() {
		manager.RegisterAltcoin(info, p)
	}
}

func (p *SkyFiberPlugin) GetName() string {
	return "SkyFiber"
}

func (p *SkyFiberPlugin) GetDescription() string {
	return "FiberCrypto wallet connector for Skycoin and SkyFiber altcoins"
}

func (p *SkyFiberPlugin) LoadWalletEnvs() []core.WalletEnv {

	wltSources, err := config.GetWalletSources()
	if err != nil {
		return nil
	}
	wltEnvs := make([]core.WalletEnv, 0)
	for _, wltS := range wltSources {

		tp := wltS.Tp
		source := wltS.Source
		var wltEnv core.WalletEnv
		if tp == string(config.LocalWallet) {
			wltEnv = &WalletDirectory{WalletDir: source}
		} else if tp == string(config.RemoteWallet) {
			wltEnv = NewWalletNode(source)
		}
		wltEnvs = append(wltEnvs, wltEnv)
	}

	return wltEnvs
}

func (p *SkyFiberPlugin) LoadPEX(netType string) (core.PEX, error) {
	var poolSection string
	if netType == "MainNet" {
		poolSection = PoolSection
	} else {
		return nil, errors.ErrInvalidNetworkType
	}
	return NewSkycoinPEX(poolSection), nil

}

func NewSkyFiberPlugin(params params.SkyFiberParams) core.AltcoinPlugin {
	return &SkyFiberPlugin{
		Params: params,
	}
}
