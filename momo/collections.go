package momo

type CollectionsOp struct {
	client *Client
}

func (c *CollectionsOp) RequestToPay(mobile string, amount int64, id, payeeNote, payerMessage, currency string) (string) {
	return "transaction ID"
}

func (c *CollectionsOp) GetTransaction(transactionID string) {
}

func (c *CollectionsOp) GetBalance() {}

func (c *CollectionsOp) IsPayeeActive(accountHolderType, accountHolderID string) {

}

func (c *CollectionsOp) GetToken() string {
	return ""
}
