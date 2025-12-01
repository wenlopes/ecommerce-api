package out

type Product struct {
	Code       string     `json:"code"`
	Price      float64    `json:"price"`
	Categories []Category `json:"categories"`
	Variants   []Variant  `json:"variants,omitempty"`
}

type Variant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}
