package ledgerstate

type LedgerState struct {
	mainReality  *Reality
	subRealities map[string]*Reality
}

func NewLedgerState() *LedgerState {
	return &LedgerState{
		mainReality:  NewReality(),
		subRealities: make(map[string]*Reality),
	}
}

func (ledgerState *LedgerState) AddAddress(address *Address) *LedgerState {
	ledgerState.mainReality.SetAddress(address)

	return ledgerState
}

func (ledgerState *LedgerState) GetUnspentTransferOutputs(address AddressHash, includedSubRealities ...TransferHash) (result TransferOutputs) {
	result = make(TransferOutputs)

	for _, reality := range ledgerState.getRealities(includedSubRealities...) {
		if address := reality.GetAddress(address); address.Exists() {
			for transferHash, coloredBalance := range address.GetUnspentTransferOutputs() {
				result[transferHash] = coloredBalance
			}
		}
	}

	return
}

func (ledgerState *LedgerState) getRealities(includedSubRealities ...TransferHash) (realities []*Reality) {
	realities = append(realities, ledgerState.mainReality)

	for _, subRealityToInclude := range includedSubRealities {
		if subReality := ledgerState.subRealities[subRealityToInclude]; subReality.Exists() {
			realities = append(realities, subReality)
		}
	}

	return
}
