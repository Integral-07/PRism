package entity

type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

func (r RiskLevel) String() string {
	return string(r)
}

func (r RiskLevel) IsValid() bool {
	switch r {
	case RiskLevelLow, RiskLevelMedium, RiskLevelHigh:
		return true
	default:
		return false
	}
}

func (r RiskLevel) Emoji() string {
	return map[RiskLevel]string{
		RiskLevelLow:    "🟢",
		RiskLevelMedium: "🟡",
		RiskLevelHigh:   "🔴",
	}[r]
}