package entity

import "testing"

func TestRiskLevel_String(t *testing.T) {
	tests := []struct {
		name string
		r    RiskLevel
		want string
	}{
		{"Low", RiskLevelLow, "low"},
		{"Medium", RiskLevelMedium, "medium"},
		{"High", RiskLevelHigh, "high"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.want {
				t.Errorf("RiskLevel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRiskLevel_IsValid(t *testing.T) {
	tests := []struct {
		name string
		r    RiskLevel
		want bool
	}{
		{ "Low", RiskLevelLow, true},
		{ "Medium", RiskLevelMedium, true},
		{ "High", RiskLevelHigh, true},
		{ "Invalid", "invalid", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsValid(); got != tt.want {
				t.Errorf("RiskLevel.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRiskLevel_Emoji(t *testing.T) {
	tests := []struct {
		name string
		r    RiskLevel
		want string
	}{
		{"Low", RiskLevelLow, "🟢"},
		{"Medium", RiskLevelMedium, "🟡"},
		{"High", RiskLevelHigh, "🔴"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Emoji(); got != tt.want {
				t.Errorf("RiskLevel.Emoji() = %v, want %v", got, tt.want)
			}
		})
	}
}