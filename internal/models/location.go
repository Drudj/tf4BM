package models

// Location представляет локацию (дата-центр)
type Location struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Country     string `json:"country"`
	City        string `json:"city"`
	Description string `json:"description,omitempty"`
	Available   bool   `json:"available"`
}

// LocationsResponse представляет ответ со списком локаций
type LocationsResponse struct {
	Locations []Location `json:"locations"`
}
