package core

import (
	"github.com/tonkeeper/opentonapi/pkg/addressbook"
	"github.com/tonkeeper/tongo"
)

type addressBook interface {
	IsWallet(a tongo.AccountID) (bool, error)
	GetAddressInfoByAddress(a tongo.AccountID) (addressbook.KnownAddress, bool) // todo: maybe rewrite to pointer
	GetCollectionInfoByAddress(a tongo.AccountID) (addressbook.KnownCollection, bool)
	GetJettonInfoByAddress(a tongo.AccountID) (addressbook.KnownJetton, bool)
	GetTFPoolInfo(a tongo.AccountID) (addressbook.TFPoolInfo, bool)
	GetKnownJettons() map[tongo.AccountID]addressbook.KnownJetton
	GetKnownCollections() map[tongo.AccountID]addressbook.KnownCollection
	SearchAttachedAccountsByPrefix(prefix string) []addressbook.AttachedAccount
}

// SimpleAddressBook is a simple implementation of the AddressBook interface.
type SimpleAddressBook struct {
	addressBook
	wallets map[tongo.AccountID]bool
}

// NewSimpleAddressBook creates a new instance of SimpleAddressBook.
func NewSimpleAddressBook(wallets map[tongo.AccountID]bool, defaultBook addressBook) *SimpleAddressBook {
	s := &SimpleAddressBook{
		wallets:     wallets,
		addressBook: defaultBook,
	}
	return s
}

// IsWallet checks if the given account ID is a wallet.
func (s *SimpleAddressBook) IsWallet(accountID tongo.AccountID) (bool, error) {
	if isWallet, exists := s.wallets[accountID]; exists {
		return isWallet, nil
	}
	return s.addressBook.IsWallet(accountID)
}
