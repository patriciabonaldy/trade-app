package model

import _ "embed"

const (
	// BTCUSDPair BTC-USD Pair
	BTCUSDPair = "BTC-USD"
	// ETHUSDPair ETH-USD Pair
	ETHUSDPair = "ETH-USD"
	// ETHBTCPair ETH-BTC Pair
	ETHBTCPair = "ETH-BTC"
)

//go:embed header.txt
var Header []byte

// Data store information about trade
type Data struct {
	Price float64
	Size  float64
}

// VWpaData store Volume Weighted Average Price
type VWpaData struct {
	Price float64
	Size  float64
	Vwpa  float64
}

// CalculateVwpa method to get vwpa
func (d *VWpaData) CalculateVwpa() {
	d.Vwpa = (d.Price * d.Size) / d.Size
}
