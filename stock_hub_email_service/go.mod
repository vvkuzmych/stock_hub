module email_service

go 1.25

require (
	github.com/lib/pq v1.10.9
	stock_hub_trade v0.0.0
)

// Local development: use local stock_hub_trade models
replace stock_hub_trade => ../stock_hub_trade
