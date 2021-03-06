package value

import (
	"net/http"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/transaction"
	"github.com/labstack/echo"
)

// getTransactionByIDHandler gets the transaction by id.
func getTransactionByIDHandler(c echo.Context) error {
	txnID, err := transaction.IDFromBase58(c.QueryParam("txnID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetTransactionByIDResponse{Error: err.Error()})
	}

	// get txn by txn id
	cachedTxnMetaObj := valuetransfers.Tangle().TransactionMetadata(txnID)
	defer cachedTxnMetaObj.Release()
	if !cachedTxnMetaObj.Exists() {
		return c.JSON(http.StatusNotFound, GetTransactionByIDResponse{Error: "Transaction not found"})
	}
	cachedTxnObj := valuetransfers.Tangle().Transaction(txnID)
	defer cachedTxnObj.Release()
	if !cachedTxnObj.Exists() {
		return c.JSON(http.StatusNotFound, GetTransactionByIDResponse{Error: "Transaction not found"})
	}
	txn := ParseTransaction(cachedTxnObj.Unwrap())

	txnMeta := cachedTxnMetaObj.Unwrap()
	txnMeta.Preferred()
	return c.JSON(http.StatusOK, GetTransactionByIDResponse{
		Transaction: txn,
		InclusionState: InclusionState{
			Confirmed:   txnMeta.Confirmed(),
			Conflicting: txnMeta.Conflicting(),
			Liked:       txnMeta.Liked(),
			Solid:       txnMeta.Solid(),
			Rejected:    txnMeta.Rejected(),
			Finalized:   txnMeta.Finalized(),
			Preferred:   txnMeta.Preferred(),
		},
	})
}

// GetTransactionByIDResponse is the HTTP response from retrieving transaction.
type GetTransactionByIDResponse struct {
	Transaction    Transaction    `json:"transaction,omitempty"`
	InclusionState InclusionState `json:"inclusion_state,omitempty"`
	Error          string         `json:"error,omitempty"`
}
