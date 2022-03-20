package model

// Data store information about trade
type Data struct {
	Price    float64
	Quantity int
}

// VWpaData store Volume Weighted Average Price
type VWpaData struct {
	Price    float64
	Quantity int
	Vwpa     float64
}

// CalculateVwpa method to get vwpa
func (d *VWpaData) CalculateVwpa() {
	d.Vwpa = d.Price / float64(d.Quantity)
}
