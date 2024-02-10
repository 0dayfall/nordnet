package models

type LedgerInformation struct {
	TotalAccIntDeb  Amount   `json:"total_acc_int_deb"`
	TotalAccIntCred Amount   `json:"total_acc_int_cred"`
	Total           Amount   `json:"total"`
	Ledgers         []Ledger `json:"ledgers"`
}

type Ledger struct {
	Currency      string `json:"currency"`
	AccountSum    Amount `json:"account_sum"`
	AccountSumAcc Amount `json:"account_sum_acc"`
	AccIntDeb     Amount `json:"acc_int_deb"`
	AccIntCred    Amount `json:"acc_int_cred"`
	ExchangeRate  Amount `json:"exchange_rate"`
}
