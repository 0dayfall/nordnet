package models

type Account struct {
	Accno         int64  `json:"accno"`
	Type          string `json:"type"`
	Default       bool   `json:"default"`
	Alias         string `json:"alias"`
	Blocked       bool   `json:"blocked"`
	BlockedReason string `json:"blocked_reason"`
}

type AccountInfo struct {
	AccountCurrency            string `json:"account_currency"`
	AccountCredit              Amount `json:"account_credit"`
	AccountSum                 Amount `json:"account_sum"`
	Collateral                 Amount `json:"collateral"`
	CreditAccountSum           Amount `json:"credit_account_sum"`
	ForwardSum                 Amount `json:"forward_sum"`
	FutureSum                  Amount `json:"future_sum"`
	UnrealizedFutureProfitLoss Amount `json:"unrealized_future_profit_loss"`
	FullMarketvalue            Amount `json:"full_marketvalue"`
	Interest                   Amount `json:"interest"`
	IntradayCredit             Amount `json:"intraday_credit"`
	LoanLimit                  Amount `json:"loan_limit"`
	OwnCapital                 Amount `json:"own_capital"`
	OwnCapitalMorning          Amount `json:"own_capital_morning"`
	PawnValue                  Amount `json:"pawn_value"`
	TradingPower               Amount `json:"trading_power"`
}
