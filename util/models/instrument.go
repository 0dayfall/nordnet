package models

type Instrument struct {
	InstrumentId        int64            `json:"instrument_id"`
	Tradables           []Tradable       `json:"tradables"`
	Currency            string           `json:"currency"`
	InstrumentGroupType string           `json:"instrument_group_type"`
	InstrumentType      string           `json:"instrument_type"`
	Multiplier          float64          `json:"multiplier"`
	Symbol              string           `json:"symbol"`
	IsinCode            string           `json:"isin_code"`
	MarketView          string           `json:"market_view"`
	StrikePrice         float64          `json:"strike_price"`
	NumberOfSecurities  float64          `json:"number_of_securities"`
	ProspectusUrl       string           `json:"prospectus_url"`
	ExpirationDate      string           `json:"expiration_date"`
	Name                string           `json:"name"`
	Sector              string           `json:"sector"`
	SectorGroup         string           `json:"sector_group"`
	Underlyings         []UnderlyingInfo `json:"underlyings"`
}

type Tradable struct {
	TradableId
	TickSizeId   int64   `json:"tick_size_id"`
	LotSize      float64 `json:"lot_size"`
	DisplayOrder int64   `json:"display_order"`
}
