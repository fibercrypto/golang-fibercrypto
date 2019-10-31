package core

// KeyValueStorage provides read / write access to values given a key
type KeyValueStorage interface {
	GetValue(key string) interface{}
}

// AddressBook provides method to manage a contact database.
type AddressBook interface {
	Open() error
	GetContact(id uint64, password []byte) (Contact, error)
	ListContact(password []byte) ([]Contact, error)
	InsertContact(contact Contact, password []byte) error
	DeleteContact(id uint64) error
	UpdateContact(id uint64, contact Contact, password []byte) error
}

// Contact provides encrypt / decrypt data.
type Contact interface {
	GetID() uint64
	SetID(id uint64)
	GetAddress(pos int64) ReadableAddress
	SetAddress(ReadableAddress)
}

type ReadableAddress interface {
	GetValue() []byte
	SetValue(val []byte)
	GetCoinType() []byte
	SetType(val []byte)
}
