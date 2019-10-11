package ledgerstate

type Reality struct {
	addresses map[AddressHash]*Address
}

func NewReality() *Reality {
	return &Reality{
		addresses: make(map[AddressHash]*Address),
	}
}

func (reality *Reality) SetAddress(address *Address) *Reality {
	reality.addresses[address.GetHash()] = address

	return reality
}

func (reality *Reality) GetAddress(address AddressHash) *Address {
	return reality.addresses[address]
}

func (reality *Reality) Exists() bool {
	return reality != nil
}
