package tipselector

import (
	"github.com/iotaledger/hive.go/events"

	"github.com/iotaledger/goshimmer/packages/binary/datastructure"
	"github.com/iotaledger/goshimmer/packages/binary/messagelayer/message"
)

type TipSelector struct {
	tips   *datastructure.RandomMap
	Events Events
}

func New() *TipSelector {
	return &TipSelector{
		tips: datastructure.NewRandomMap(),
		Events: Events{
			TipAdded:   events.NewEvent(transactionIdEvent),
			TipRemoved: events.NewEvent(transactionIdEvent),
		},
	}
}

func (tipSelector *TipSelector) AddTip(transaction *message.Message) {
	transactionId := transaction.Id()
	if tipSelector.tips.Set(transactionId, transactionId) {
		tipSelector.Events.TipAdded.Trigger(transactionId)
	}

	trunkTransactionId := transaction.TrunkId()
	if _, deleted := tipSelector.tips.Delete(trunkTransactionId); deleted {
		tipSelector.Events.TipRemoved.Trigger(trunkTransactionId)
	}

	branchTransactionId := transaction.BranchId()
	if _, deleted := tipSelector.tips.Delete(branchTransactionId); deleted {
		tipSelector.Events.TipRemoved.Trigger(branchTransactionId)
	}
}

func (tipSelector *TipSelector) GetTips() (trunkTransaction, branchTransaction message.Id) {
	tip := tipSelector.tips.RandomEntry()
	if tip == nil {
		trunkTransaction = message.EmptyId
		branchTransaction = message.EmptyId

		return
	}

	branchTransaction = tip.(message.Id)

	if tipSelector.tips.Size() == 1 {
		trunkTransaction = branchTransaction

		return
	}

	trunkTransaction = tipSelector.tips.RandomEntry().(message.Id)
	for trunkTransaction == branchTransaction && tipSelector.tips.Size() > 1 {
		trunkTransaction = tipSelector.tips.RandomEntry().(message.Id)
	}

	return
}

func (tipSelector *TipSelector) GetTipCount() int {
	return tipSelector.tips.Size()
}
