package skycoin

import (
	skyparams "github.com/SkycoinProject/skycoin/src/params"
	"github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/params"
)

var (
	SkycoinMainNetParams = params.SkyFiberParams{
		Distribution: skyparams.MainNetDistribution,
	}
)

const (
	SkycoinTicker              = params.SkycoinTicker
	SkycoinName                = params.SkycoinName
	SkycoinFamily              = params.SkycoinFamily
	SkycoinDescription         = params.SkycoinDescription
	CoinHoursTicker            = params.CoinHoursTicker
	CoinHoursName              = params.CoinHoursName
	CoinHoursDescription       = params.CoinHoursDescription
	CalculatedHoursTicker      = params.CalculatedHoursTicker
	CalculatedHoursName        = params.CalculatedHoursName
	CalculatedHoursDescription = params.CalculatedHoursDescription
)
