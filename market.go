package kalshi

import (
	"context"
	"fmt"
	"time"
)

// Event is described here:
// https://trading-api.readme.io/reference/getevents.
type Event struct {
	Category          string    `json:"category"`
	EventTicker       string    `json:"event_ticker"`
	MutuallyExclusive bool      `json:"mutually_exclusive"`
	SeriesTicker      string    `json:"series_ticker"`
	StrikeDate        time.Time `json:"strike_date"`
	StrikePeriod      string    `json:"strike_period"`
	SubTitle          string    `json:"sub_title"`
	Title             string    `json:"title"`
}

// EventsResponse is described here:
// https://trading-api.readme.io/reference/getevents.
type EventsResponse struct {
	CursorResponse
	Events []Event `json:"events"`
}

// EventsRequest is described here:
// https://trading-api.readme.io/reference/getevents.
type EventsRequest struct {
	CursorRequest
	// Status is one of "open", "closed", or "settled".
	Status       string `url:"status,omitempty"`
	SeriesTicker string `url:"series_ticker,omitempty"`
}

// Events is described here:
// https://trading-api.readme.io/reference/getevents.
func (c *Client) Events(ctx context.Context, req EventsRequest) (*EventsResponse, error) {
	var resp EventsResponse

	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     "events",
		QueryParams:  req,
		JSONRequest:  nil,
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// EventResponse is described here:
// https://trading-api.readme.io/reference/getevent.
type EventResponse struct {
	Event   Event    `json:"event"`
	Markets []Market `json:"markets"`
}

// Event is described here:
// https://trading-api.readme.io/reference/getevent.
func (c *Client) Event(ctx context.Context, event string) (*EventResponse, error) {
	var resp EventResponse

	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     "events/" + event,
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// MarketsRequest is described here:
// https://trading-api.readme.io/reference/getmarkets.
type MarketsRequest struct {
	CursorRequest
	EventTicker  string `url:"event_ticker,omitempty"`
	SeriesTicker string `url:"series_ticker,omitempty"`
	MaxCloseTs   int    `url:"max_close_ts,omitempty"`
	MinCloseTs   int    `url:"min_close_ts,omitempty"`
	// Status is one of "open", "closed", and "settled"
	Status  string   `url:"status,omitempty"`
	Tickers []string `url:"status,omitempty"`
}

// Market is described here:
// https://trading-api.readme.io/reference/getmarkets.
type Market struct {
	Ticker                   string    `json:"ticker"`
	EventTicker              string    `json:"event_ticker"`
	MarketType               string    `json:"market_type"`
	Title                    string    `json:"title"`
	Subtitle                 string    `json:"subtitle"`
	YesSubTitle              string    `json:"yes_sub_title"`
	NoSubTitle               string    `json:"no_sub_title"`
	CreatedTime              time.Time `json:"created_time"`
	UpdatedTime              time.Time `json:"updated_time"`
	OpenTime                 time.Time `json:"open_time"`
	CloseTime                time.Time `json:"close_time"`
	ExpirationTime           time.Time `json:"expiration_time"`
	LatestExpirationTime     time.Time `json:"latest_expiration_time"`
	ExpectedExpirationTime   time.Time `json:"expected_expiration_time"`
	SettlementTimerSeconds   int       `json:"settlement_timer_seconds"`
	Status                   string    `json:"status"`
	ResponsePriceUnits       string    `json:"response_price_units"`
	YesBid                   Cents     `json:"yes_bid"`
	YesBidDollars            string    `json:"yes_bid_dollars"`
	YesBidSizeFp             string    `json:"yes_bid_size_fp"`
	YesAsk                   Cents     `json:"yes_ask"`
	YesAskDollars            string    `json:"yes_ask_dollars"`
	YesAskSizeFp             string    `json:"yes_ask_size_fp"`
	NoBid                    Cents     `json:"no_bid"`
	NoBidDollars             string    `json:"no_bid_dollars"`
	NoAsk                    Cents     `json:"no_ask"`
	NoAskDollars             string    `json:"no_ask_dollars"`
	LastPrice                Cents     `json:"last_price"`
	LastPriceDollars         string    `json:"last_price_dollars"`
	PreviousYesBid           Cents     `json:"previous_yes_bid"`
	PreviousYesBidDollars    string    `json:"previous_yes_bid_dollars"`
	PreviousYesAsk           Cents     `json:"previous_yes_ask"`
	PreviousYesAskDollars    string    `json:"previous_yes_ask_dollars"`
	PreviousPrice            Cents     `json:"previous_price"`
	PreviousPriceDollars     string    `json:"previous_price_dollars"`
	Volume                   int       `json:"volume"`
	VolumeFp                 string    `json:"volume_fp"`
	Volume24H                int       `json:"volume_24h"`
	Volume24HFp              string    `json:"volume_24h_fp"`
	Liquidity                Cents     `json:"liquidity"`
	LiquidityDollars         string    `json:"liquidity_dollars"`
	OpenInterest             int       `json:"open_interest"`
	OpenInterestFp           string    `json:"open_interest_fp"`
	NotionalValue            Cents     `json:"notional_value"`
	NotionalValueDollars     string    `json:"notional_value_dollars"`
	Result                   string    `json:"result"`
	CanCloseEarly            bool      `json:"can_close_early"`
	FractionalTradingEnabled bool      `json:"fractional_trading_enabled"`
	ExpirationValue          string    `json:"expiration_value"`
	SettlementValue          Cents     `json:"settlement_value"`
	SettlementValueDollars   string    `json:"settlement_value_dollars"`
	SettlementTs             time.Time `json:"settlement_ts"`
	FeeWaiverExpirationTime  time.Time `json:"fee_waiver_expiration_time"`
	EarlyCloseCondition      string    `json:"early_close_condition"`
	TickSize                 Cents     `json:"tick_size"`
	RulesPrimary             string    `json:"rules_primary"`
	RulesSecondary           string    `json:"rules_secondary"`
	PriceLevelStructure      string    `json:"price_level_structure"`
	PriceRanges              []struct {
		Start string `json:"start"`
		End   string `json:"end"`
		Step  string `json:"step"`
	} `json:"price_ranges"`
	StrikeType          string         `json:"strike_type"`
	FloorStrike         float64        `json:"floor_strike,omitempty"`
	CapStrike           float64        `json:"cap_strike,omitempty"`
	FunctionalStrike    string         `json:"functional_strike"`
	CustomStrike        map[string]any `json:"custom_strike"`
	MveCollectionTicker string         `json:"mve_collection_ticker"`
	MveSelectedLegs     []struct {
		EventTicker               string `json:"event_ticker"`
		MarketTicker              string `json:"market_ticker"`
		Side                      string `json:"side"`
		YesSettlementValueDollars string `json:"yes_settlement_value_dollars"`
	} `json:"mve_selected_legs"`
	PrimaryParticipantKey string `json:"primary_participant_key"`
	IsProvisional         bool   `json:"is_provisional"`
}

func (m *Market) YesMidPrice() Cents {
	return (m.YesBid + m.YesAsk) / 2
}

func (m *Market) NoMidPrice() Cents {
	return (m.NoBid + m.NoAsk) / 2
}

func (m *Market) MarketValue(p *MarketPosition) Cents {
	if p == nil {
		return 0
	}

	if p.Position < 0 {
		return Cents(p.AbsPosition()) * m.NoMidPrice()
	}
	return Cents(p.AbsPosition()) * m.YesMidPrice()
}

// EstimateReturn shows the estimated return for an open position.
func (m *Market) EstimateReturn(p *MarketPosition) Cents {
	if p == nil {
		return 0
	}

	posMarketValue := m.MarketValue(p)
	costBasis := p.MarketExposure
	return p.RealizedPnl - p.FeesPaid + (posMarketValue - costBasis)
}

// MarketsResponse is described here:
// https://trading-api.readme.io/reference/getmarkets.
type MarketsResponse struct {
	Markets []Market `json:"markets,omitempty"`
	CursorResponse
}

// Markets is described here:
// https://trading-api.readme.io/reference/getmarkets.
func (c *Client) Markets(
	ctx context.Context,
	req MarketsRequest,
) (*MarketsResponse, error) {
	var resp MarketsResponse

	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     "markets",
		QueryParams:  req,
		JSONRequest:  nil,
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Trade is described here:
// https://trading-api.readme.io/reference/gettrades.
type Trade struct {
	TradeID         string    `json:"trade_id"`
	Ticker          string    `json:"ticker"`
	Price           float64   `json:"price"`
	Count           int       `json:"count"`
	CountFp         string    `json:"count_fp"`
	YesPrice        Cents     `json:"yes_price"`
	NoPrice         Cents     `json:"no_price"`
	YesPriceDollars string    `json:"yes_price_dollars"`
	NoPriceDollars  string    `json:"no_price_dollars"`
	TakerSide       Side      `json:"taker_side"`
	CreatedTime     time.Time `json:"created_time"`
}

// TradesResponse is described here:
// https://trading-api.readme.io/reference/gettrades.
type TradesResponse struct {
	CursorResponse
	Trades []Trade `json:"trades,omitempty"`
}

// TradesRequest is described here:
// https://trading-api.readme.io/reference/gettrades.
type TradesRequest struct {
	CursorRequest
	Ticker string `url:"ticker,omitempty"`
	MinTS  int    `url:"min_ts,omitempty"`
	MaxTS  int    `url:"max_ts,omitempty"`
}

// Trades is described here:
// https://trading-api.readme.io/reference/gettrades.
func (c *Client) Trades(
	ctx context.Context,
	req TradesRequest,
) (*TradesResponse, error) {
	var resp TradesResponse

	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     "markets/trades",
		QueryParams:  req,
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Market is described here:
// https://trading-api.readme.io/reference/getmarket.
func (c *Client) Market(ctx context.Context, ticker string) (*Market, error) {
	var resp struct {
		Market Market `json:"market"`
	}
	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("markets/%s", ticker),
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}
	return &resp.Market, nil
}

// HistoricalMarket is described here:
// https://docs.kalshi.com/api-reference/historical/get-historical-market.
func (c *Client) HistoricalMarket(ctx context.Context, ticker string) (*Market, error) {
	var resp struct {
		Market Market `json:"market"`
	}
	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("historical/markets/%s", ticker),
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}
	return &resp.Market, nil
}

// GetHistoricalMarketCandlesticksRequest is described here:
// https://docs.kalshi.com/api-reference/historical/get-historical-market-candlesticks.
type GetHistoricalMarketCandlesticksRequest struct {
	Ticker         string `url:"-"`
	StartTs        int64  `url:"start_ts,omitempty"`
	EndTs          int64  `url:"end_ts,omitempty"`
	PeriodInterval int    `url:"period_interval,omitempty"`
}

// HistoricalOhlc represents dollar-string-only OHLC values for historical markets.
type HistoricalOhlc struct {
	Open  string `json:"open"`
	Low   string `json:"low"`
	High  string `json:"high"`
	Close string `json:"close"`
}

// HistoricalOhlcExtended represents extended OHLC values including Mean and Previous.
type HistoricalOhlcExtended struct {
	HistoricalOhlc
	Mean     string `json:"mean"`
	Previous string `json:"previous"`
}

// HistoricalCandlestick represents a single historical candlestick data point.
type HistoricalCandlestick struct {
	EndPeriodTs  int64                  `json:"end_period_ts"`
	YesBid       HistoricalOhlc         `json:"yes_bid"`
	YesAsk       HistoricalOhlc         `json:"yes_ask"`
	Price        HistoricalOhlcExtended `json:"price"`
	Volume       string                 `json:"volume"`
	OpenInterest string                 `json:"open_interest"`
}

// GetHistoricalMarketCandlesticksResponse is described here:
// https://docs.kalshi.com/api-reference/historical/get-historical-market-candlesticks.
type GetHistoricalMarketCandlesticksResponse struct {
	Ticker       string                  `json:"ticker"`
	Candlesticks []HistoricalCandlestick `json:"candlesticks"`
}

// GetHistoricalMarketCandlesticks is described here:
// https://docs.kalshi.com/api-reference/historical/get-historical-market-candlesticks.
func (c *Client) GetHistoricalMarketCandlesticks(
	ctx context.Context,
	req GetHistoricalMarketCandlesticksRequest,
) (*GetHistoricalMarketCandlesticksResponse, error) {
	var resp GetHistoricalMarketCandlesticksResponse

	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("historical/markets/%s/candlesticks", req.Ticker),
		QueryParams:  req,
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// MarketOrderBook is described here:
// https://trading-api.readme.io/reference/getmarketorderbook.
func (c *Client) MarketOrderBook(ctx context.Context, ticker string) (*MarketOrderBookResponse, error) {
	var resp MarketOrderBookResponse
	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("markets/%s/orderbook/?depth=100", ticker),
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// SettlementSource represents a settlement source for a series.
type SettlementSource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Series is described here:
// https://trading-api.readme.io/reference/getseries.
type Series struct {
	Ticker                 string             `json:"ticker"`
	Frequency              string             `json:"frequency"`
	Title                  string             `json:"title"`
	Category               string             `json:"category"`
	Tags                   []string           `json:"tags"`
	SettlementSources      []SettlementSource `json:"settlement_sources"`
	ContractURL            string             `json:"contract_url"`
	ContractTermsURL       string             `json:"contract_terms_url"`
	FeeType                string             `json:"fee_type"`
	FeeMultiplier          float64            `json:"fee_multiplier"`
	AdditionalProhibitions []string           `json:"additional_prohibitions"`
	ProductMetadata        map[string]any     `json:"product_metadata"`
	Volume                 int                `json:"volume"`
	VolumeFp               string             `json:"volume_fp"`
}

// ListSeriesRequest is described here:
// https://trading-api.readme.io/reference/listseries.
type ListSeriesRequest struct {
	Category               string `url:"category,omitempty"`
	Tags                   string `url:"tags,omitempty"`
	IncludeProductMetadata bool   `url:"include_product_metadata,omitempty"`
	IncludeVolume          bool   `url:"include_volume,omitempty"`
}

// ListSeriesResponse is described here:
// https://trading-api.readme.io/reference/listseries.
type ListSeriesResponse struct {
	Series []Series `json:"series"`
}

// ListSeries is described here:
// https://trading-api.readme.io/reference/listseries.
func (c *Client) ListSeries(ctx context.Context, req ListSeriesRequest) (*ListSeriesResponse, error) {
	var resp ListSeriesResponse

	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     "series",
		QueryParams:  req,
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetSeriesRequest is described here:
// https://trading-api.readme.io/reference/getseries.
type GetSeriesRequest struct {
	IncludeVolume bool `url:"include_volume,omitempty"`
}

// Series is described here:
// https://trading-api.readme.io/reference/getseries.
func (c *Client) Series(ctx context.Context, seriesTicker string, req GetSeriesRequest) (*Series, error) {
	var resp struct {
		Series Series `json:"series"`
	}
	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("series/%s", seriesTicker),
		QueryParams:  req,
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}
	return &resp.Series, nil
}

// GetMarketCandlesticksRequest is described here:
// https://docs.kalshi.com/api-reference/market/get-market-candlesticks.
type GetMarketCandlesticksRequest struct {
	SeriesTicker             string `url:"-"`
	Ticker                   string `url:"-"`
	StartTs                  int64  `url:"start_ts,omitempty"`
	EndTs                    int64  `url:"end_ts,omitempty"`
	PeriodInterval           int    `url:"period_interval,omitempty"`
	IncludeLatestBeforeStart bool   `url:"include_latest_before_start,omitempty"`
}

// Ohlc represents Open, High, Low, Close values.
type Ohlc struct {
	Open         int64  `json:"open"`
	OpenDollars  string `json:"open_dollars"`
	Low          int64  `json:"low"`
	LowDollars   string `json:"low_dollars"`
	High         int64  `json:"high"`
	HighDollars  string `json:"high_dollars"`
	Close        int64  `json:"close"`
	CloseDollars string `json:"close_dollars"`
}

// OhlcExtended represents extended OHLC values including Mean, Previous, Min, Max.
type OhlcExtended struct {
	Ohlc
	Mean            int64  `json:"mean"`
	MeanDollars     string `json:"mean_dollars"`
	Previous        int64  `json:"previous"`
	PreviousDollars string `json:"previous_dollars"`
	Min             int64  `json:"min"`
	MinDollars      string `json:"min_dollars"`
	Max             int64  `json:"max"`
	MaxDollars      string `json:"max_dollars"`
}

// Candlestick represents a single candlestick data point.
type Candlestick struct {
	EndPeriodTs    int64        `json:"end_period_ts"`
	YesBid         Ohlc         `json:"yes_bid"`
	YesAsk         Ohlc         `json:"yes_ask"`
	Price          OhlcExtended `json:"price"`
	Volume         int64        `json:"volume"`
	VolumeFp       string       `json:"volume_fp"`
	OpenInterest   int64        `json:"open_interest"`
	OpenInterestFp string       `json:"open_interest_fp"`
}

// GetMarketCandlesticksResponse is described here:
// https://docs.kalshi.com/api-reference/market/get-market-candlesticks.
type GetMarketCandlesticksResponse struct {
	Ticker       string        `json:"ticker"`
	Candlesticks []Candlestick `json:"candlesticks"`
}

// GetMarketCandlesticks is described here:
// https://docs.kalshi.com/api-reference/market/get-market-candlesticks.
func (c *Client) GetMarketCandlesticks(
	ctx context.Context,
	req GetMarketCandlesticksRequest,
) (*GetMarketCandlesticksResponse, error) {
	var resp GetMarketCandlesticksResponse

	err := c.request(ctx, request{
		Method:       "GET",
		Endpoint:     fmt.Sprintf("series/%s/markets/%s/candlesticks", req.SeriesTicker, req.Ticker),
		QueryParams:  req,
		JSONResponse: &resp,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
