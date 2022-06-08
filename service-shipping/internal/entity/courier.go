package entity

type Courier struct {
	Name         string       `json:"name"`
	Logo         string       `json:"logo"` //Logo image url.
	ProviderCode ProviderCode `json:"provider_code"`
	CourierCode  string       `json:"courier_code,omitempty"` //CourierCode only use with Shippop (provider_code=1).
	EnableCOD    bool         `json:"enable_cod"`
}
