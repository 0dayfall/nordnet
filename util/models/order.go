package models

type Order struct {
	Accno               int64               `json:"accno"`
	OrderId             int64               `json:"order_id"`
	Price               Amount              `json:"price"`
	Volume              float64             `json:"volume"`
	Tradable            TradableId          `json:"tradable"`
	OpenVolume          float64             `json:"open_volume"`
	TradedVolume        float64             `json:"traded_volume"`
	Side                string              `json:"side"`
	Modified            int64               `json:"modified"`
	Reference           string              `json:"reference"`
	ActivationCondition ActivationCondition `json:"activation_condition"`
	PriceCondition      string              `json:"price_condition"`
	VolumeCondition     string              `json:"volume_condition"`
	Validity            Validity            `json:"validity"`
	ActionState         string              `json:"action_state"`
	OrderState          string              `json:"order_state"`
}

type TradableId struct {
	Identifier string `json:"identifier"`
	MarketId   int64  `json:"market_id"`
}

type ActivationCondition struct {
	Type             string  `json:"type"`
	TrailingValue    float64 `json:"trailing_value"`
	TriggerValue     float64 `json:"trigger_value"`
	TriggerCondition string  `json:"trigger_condition"`
}
