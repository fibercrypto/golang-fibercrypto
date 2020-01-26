package skycoin

import (
	"encoding/hex"
	"errors"
	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/base58"
	"github.com/fibercrypto/fibercryptowallet/src/coin/mocks"
	"github.com/fibercrypto/fibercryptowallet/src/core"
	"github.com/fibercrypto/fibercryptowallet/src/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestNewSkycoinAddress(t *testing.T) {
	type args struct {
		addrStr string
	}
	tests := []struct {
		name    string
		args    args
		want    SkycoinAddress
		wantErr bool
		err     error
	}{
		{
			name: "address1",
			args: args{addrStr: "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ"},
			want: SkycoinAddress{
				isBip32: false,
				address: cipher.Address{
					Version: 0,
					Key: cipher.Ripemd160{0x3b, 0xe2, 0x53, 0x7f, 0x8c, 0x8, 0x93, 0xfd, 0xdc, 0xdd, 0xc8,
						0x78, 0x51, 0x8f, 0x38, 0xea, 0x49, 0x3d, 0x94, 0x9e},
				},
				poolSection: "skycoin",
			}, wantErr: false},
		{
			name: "address2",
			args: args{addrStr: "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
			want: SkycoinAddress{
				isBip32: false,
				address: cipher.Address{
					Version: 0,
					Key: cipher.Ripemd160{0xfd, 0x4a, 0xd3, 0x13, 0xa1, 0xb6, 0x48, 0xd3, 0x7f, 0xa1, 0xac,
						0x5, 0xcf, 0x42, 0xb, 0x3, 0xd8, 0x86, 0xdd, 0x35},
				},
				poolSection: "skycoin",
			}, wantErr: false},
		{
			name:    "empty",
			args:    args{addrStr: ""},
			want:    SkycoinAddress{},
			wantErr: true,
			err:     base58.ErrInvalidString},
		{
			name:    "invalid character",
			args:    args{addrStr: "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			want:    SkycoinAddress{},
			wantErr: true,
			err:     base58.ErrInvalidChar},
		{
			name:    "invalid checksum",
			args:    args{addrStr: "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			want:    SkycoinAddress{},
			wantErr: true,
			err:     cipher.ErrAddressInvalidChecksum},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSkycoinAddress(tt.args.addrStr)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewSkycoinAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == true {
				assert.Equal(t, tt.err, err)
				return
			}
			cAddrs, err := cipher.AddressFromBytes(got.Bytes())
			require.NoError(t, err)
			assert.Equal(t, got.String(), cAddrs.String())
			assert.False(t, got.IsBip32())
			assert.False(t, got.Null())
			assert.NotNil(t, got.Checksum())
			assert.Implements(t, new(core.CryptoAccount), got.GetCryptoAccount())
			got.isBip32 = true
			assert.True(t, got.IsBip32())
			got.isBip32 = false
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSkycoinAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSkycoinAddressIterator(t *testing.T) {
	var got core.AddressIterator
	type args struct {
		addresses []string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "empty", args: args{
			addresses: []string{},
		}},
		{name: "one-address", args: args{
			addresses: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
		}},
		{name: "two-address", args: args{
			addresses: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt", "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "empty" {
				got = NewSkycoinAddressIterator(nil)
				return
			} else {
				var args []core.Address
				for e := range tt.args.addresses {
					addrs, err := NewSkycoinAddress(tt.args.addresses[e])
					require.NoError(t, err)
					args = append(args, &addrs)
				}
				got = NewSkycoinAddressIterator(args)
			}

			for got.Next() {
				assert.Contains(t, tt.args.addresses, got.Value().String())
			}
		})
	}
}

func TestSkycoinAddress_Bytes(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{name: "addrs1", address: "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
		{name: "addrs2", address: "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := NewSkycoinAddress(tt.address)
			require.NoError(t, err)
			assert.Equal(t, tt.address, addr.String())
		})
	}
}

func TestSkycoinAddress_Verify(t *testing.T) {
	addrsFromString := func(s string) core.Address {
		skyAddrs, err := NewSkycoinAddress(s)
		require.NoError(t, err)
		return &skyAddrs
	}
	pubkeyFromString := func(s string) core.PubKey {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		spk, err := skyPubKeyFromBytes(b)
		require.NoError(t, err)
		return spk
	}
	tests := []struct {
		name        string
		addrsString string
		pkHex       string
		wantErr     bool
	}{
		{
			name:        "Ok",
			addrsString: "sFbkd94v4j1Fw3K3xwyWEPNxEP9mw5p99a",
			pkHex:       "034f1e3f2391bd3670151fd4fa3accc6a0273885984404089e5b846871db4c5304",
			wantErr:     false,
		}, {
			name:        "wrong",
			addrsString: "e1c32uD7QFSNPeACJ9iqB6Bjp4K6E8yT1J",
			pkHex:       "034f1e3f2391bd3670151fd4fa3accc6a0273885984404089e5b846871db4c5304",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addrs := addrsFromString(tt.addrsString)
			pk := pubkeyFromString(tt.pkHex)
			if tt.wantErr {
				require.Error(t, addrs.Verify(pk))
				return
			}
			require.NoError(t, addrs.Verify(pk))
		})
	}

	mockPub := new(mocks.PubKey)
	mockPub.On("Bytes").Return([]byte(""))
	skyAddr := new(SkycoinAddress)
	require.Error(t, skyAddr.Verify(mockPub))
}

func Test_skyPubKeyFromBytes(t *testing.T) {
	pubkeyFromHex := func(s string) []byte {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		return b
	}
	tests := []struct {
		name    string
		pubkHex string
		wantErr bool
	}{
		{name: "OK", pubkHex: "034f1e3f2391bd3670151fd4fa3accc6a0273885984404089e5b846871db4c5304", wantErr: false},
		{name: "wrong", pubkHex: "0213c9273d9f944c3d907bfe844090ecef3d9504c88c5165cb690de98125a4e4", wantErr: true},
		{name: "OK2", pubkHex: "0304eb48d7c0b3a915d0f53c6d966f4d9fa75df645e63dfb91d589592790943613", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := skyPubKeyFromBytes(pubkeyFromHex(tt.pubkHex))
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, pubkeyFromHex(tt.pubkHex), got.Bytes())
			assert.Equal(t, tt.pubkHex, got.pubkey.Hex())
			require.NoError(t, got.Verify())
			assert.False(t, got.Null())
			spk, err := toSkycoinPubKey(got)
			require.NoError(t, err)
			assert.Equal(t, spk.Hex(), tt.pubkHex)
			require.NoError(t, spk.Verify())
			assert.False(t, spk.Null())
		})
	}
}

func Test_skySecKeyFromBytes(t *testing.T) {
	pubkeyFromHex := func(s string) []byte {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		return b
	}
	tests := []struct {
		name    string
		pubkHex string
		wantErr bool
	}{
		{name: "OK", pubkHex: "c9135a2b667eb0847fb7ad3d1ae58a1ea2d0c38526c8948b520417dcab618563", wantErr: false},
		{name: "wrong", pubkHex: "0213c9273d9f944c3d907bfe844090ecef3d9504c88c5165cb690de985a4e4", wantErr: true},
		{name: "OK2", pubkHex: "408ea9aef71391071d275f8255bd9b6d22d5d5a22e6ab2bfc54307fb273d468c", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := skySecKeyFromBytes(pubkeyFromHex(tt.pubkHex))
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, pubkeyFromHex(tt.pubkHex), got.Bytes())
			assert.Equal(t, tt.pubkHex, got.seckey.Hex())
			require.NoError(t, got.Verify())
			assert.False(t, got.Null())
			spk, err := toSkycoinSecKey(got)
			require.NoError(t, err)
			assert.Equal(t, spk.Hex(), tt.pubkHex)
			require.NoError(t, spk.Verify())
			assert.False(t, spk.Null())
		})
	}
}

func TestAddressFromString(t *testing.T) {
	type args struct {
		addrs      string
		coinTicket string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "valid1", args: args{
			addrs:      "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt",
			coinTicket: "SKY",
		}, want: "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt", wantErr: false},
		{name: "invalid_ticket", args: args{
			addrs:      "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt",
			coinTicket: "",
		}, want: "", wantErr: true},
		{name: "invalid_address", args: args{
			addrs:      "2LEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt",
			coinTicket: "SKY",
		}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.AddressFromString(tt.args.addrs, tt.args.coinTicket)
			if tt.wantErr {
				require.Error(t, err)
				switch tt.name {
				case "invalid_ticket":
					assert.Equal(t, err, errors.New("coinTicket not match"))
				case "invalid_address":
					assert.Equal(t, err, cipher.ErrAddressInvalidLength)
				}

			} else {
				assert.NotNil(t, got)
				assert.Equal(t, tt.want, got.String())
			}

		})
	}
}

func TestToSkycoinCipherAddress(t *testing.T) {
	dir := "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ"
	key := cipher.Ripemd160{0x3b, 0xe2, 0x53, 0x7f, 0x8c, 0x8, 0x93, 0xfd, 0xdc, 0xdd, 0xc8, 0x78, 0x51, 0x8f, 0x38, 0xea, 0x49, 0x3d, 0x94, 0x9e}
	addr, err := NewSkycoinAddress(dir)
	require.NoError(t, err)
	skyAddr, err1 := addr.ToSkycoinCipherAddress()
	require.NoError(t, err1)
	require.Equal(t, key, skyAddr.Key)
}

func Test_toSkycoinPubKey(t *testing.T) {
	pub := new(SkycoinPubKey)
	pub.pubkey[4], pub.pubkey[2] = 1, 1

	mockPub := new(mocks.PubKey)
	mockPub.On("Bytes").Return([]byte(""))
	tests := []struct {
		name      string
		pub       core.PubKey
		wantError bool
	}{
		{name: "FromSkycoinPubKey", pub: pub, wantError: false},
		{name: "FromMock", pub: mockPub, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := toSkycoinPubKey(tt.pub)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				for i, b := range tt.pub.Bytes() {
					require.Equal(t, b, key[i])
				}
			}
		})
	}
}

func Test_toSkycoinSecKey(t *testing.T) {
	sec := new(SkycoinSecKey)
	sec.seckey[4], sec.seckey[2] = 1, 1

	mockSec := new(mocks.SecKey)
	mockSec.On("Bytes").Return([]byte(""))
	tests := []struct {
		name      string
		sec       core.SecKey
		wantError bool
	}{
		{name: "FromSkycoinSecKey", sec: sec, wantError: false},
		{name: "FromMock", sec: mockSec, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := toSkycoinSecKey(tt.sec)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				for i, b := range tt.sec.Bytes() {
					require.Equal(t, b, key[i])
				}
			}
		})
	}
}
