package logic

type InventoryDto struct {
}

type InventoryItem struct {
	Title    string `json:"type"`
	Quantity int    `json:"quantity"`
}

type TransactionDto struct {
	Received []*ReceiveTransaction `json:"received"`
	Sent     []*SentTransaction    `json:"sent"`
}

type ReceiveTransaction struct {
	Username string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentTransaction struct {
	Username string `json:"toUser"`
	Amount   int    `json:"amount"`
}

type InfoDto struct {
	Coins        int              `json:"coins"`
	Inventory    []*InventoryItem `json:"inventory"`
	CoinsHistory *TransactionDto  `json:"coinHistory"`
}
