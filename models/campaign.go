package models

type Campaign struct {
	Date               string  `json:"date"`
	Campaign           string  `json:"id"`
	SignUps            int     `json:"signups"`
	Ftd                int     `json:"ftd"`
	Cpa                int     `json:"cpa"`
	Deposits           float64 `json:"deposits"`
	Clicks             int     `json:"clicks"`
	CpaCommission      float64 `json:"cpaCommission"`
	RevShareCommission float64 `json:"revShareCommission"`
	TotalCommission    float64 `json:"totalCommission"`
	Ggr                float64 `json:"ggr"`
	Bet                float64 `json:"bet"`
	Win                float64 `json:"win"`
	Bonus              float64 `json:"bonus"`
	Depo               int     `json:"depo"`
	Withd              float64 `json:"withd"`
	Netrev             float64 `json:"netrev"`
}
