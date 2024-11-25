package core

import (
	"math/big"

	"github.com/shopspring/decimal"
	"github.com/tonkeeper/tongo"
)

type JettonWallet struct {
	// Address of a jetton wallet.
	Address      tongo.AccountID
	Balance      decimal.Decimal
	OwnerAddress *tongo.AccountID
	// JettonAddress of a jetton master.
	JettonAddress tongo.AccountID
	Lock          *JettonWalletLockData
	Extensions    []string
}

type JettonHolder struct {
	JettonAddress tongo.AccountID
	Address       tongo.AccountID
	Owner         tongo.AccountID
	Balance       decimal.Decimal
}

type JettonMaster struct {
	// Address of a jetton master.
	Address     tongo.AccountID
	TotalSupply big.Int
	Mintable    bool
	Admin       *tongo.AccountID
}

type JettonWalletLockData struct {
	FullBalance decimal.Decimal
	UnlockTime  int64
}

type JettonsAdditionalInfo struct {
	JettonWallets map[tongo.AccountID]JettonWallet
}

func (info *JettonsAdditionalInfo) JettonWallet(jettonWallet tongo.AccountID) (JettonWallet, bool) {
	if info.JettonWallets == nil {
		return JettonWallet{}, false
	}
	value, ok := info.JettonWallets[jettonWallet]
	return value, ok
}

func (info *JettonsAdditionalInfo) SetJettonWallet(jettonWallet tongo.AccountID, value JettonWallet) {
	if info.JettonWallets == nil {
		info.JettonWallets = make(map[tongo.AccountID]JettonWallet)
	}
	info.JettonWallets[jettonWallet] = value
}
