package util

import (
	"errors"
	"strconv"

	"github.com/fibercrypto/fibercryptowallet/src/core"
	local "github.com/fibercrypto/fibercryptowallet/src/main"
	"github.com/fibercrypto/fibercryptowallet/src/util/logging"
)

var logUtil = logging.MustGetLogger("FiberCrypto util")

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func GetCoinValue(value string, ticker string) (uint64, error) {
	accuracy, err := AltcoinQuotient(ticker)
	if err != nil {
		return uint64(0), err
	}
	coin, err2 := strconv.ParseFloat(value, 64)
	if err2 != nil {
		return uint64(0), err2
	}
	return uint64(coin * float64(accuracy)), nil
}

func FormatUint64(n uint64) string {
	in := strconv.FormatUint(n, 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

func FormatCoins(n uint64, quotient uint64) string {
	if n == uint64(0) {
		return "0"
	}

	number := strconv.FormatUint(n, 10)
	lenQ := len(strconv.FormatUint(quotient, 10)) - 1
	nFormatted := FormatUint64(n / quotient)
	if lenQ > len(number) {
		return nFormatted
	}
	reminder := number[len(number)-lenQ:]
	reminder = RemoveZeros(reminder)
	if len(reminder) == 0 {
		return nFormatted
	}
	return nFormatted + "." + reminder
}

func RemoveZeros(s string) string {
	index := 0
	temp := 0
	for {
		temp = index
		for _, c := range s[index:] {
			if string(c) != "0" {
				index++
				break
			}
		}
		if temp == index {
			break
		}
	}
	return s[:index]
}

func StringInList(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// AddressFromString returns a core.Address if match with string address.
// If the coinTicket parameter not match with any address type returns 'coinTicket not match' error.
func AddressFromString(addrs, coinTicket string) (core.Address, error) {
	altPlugin, ok := local.LoadAltcoinManager().LookupAltcoinPlugin(coinTicket)
	if !ok {
		return nil, errors.New("coinTicket not match")
	}
	return altPlugin.AddressFromString(addrs)
}
