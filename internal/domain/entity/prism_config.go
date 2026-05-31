package entity

type PrismConfig struct {
	Output OutputConfig `yaml:"output"`
	Custom string       `yaml:"custom"`
}

type OutputConfig struct {
	Triage  TriageConfig  `yaml:"triage"`
	Support SupportConfig `yaml:"support"`
}

type TriageConfig struct {
	PriorityScore       bool `yaml:"priority_score"`
	RiskLevel           bool `yaml:"risk_level"`
	EstimatedReviewTime bool `yaml:"estimated_review_time"`
	FilePriorityList    bool `yaml:"file_priority_list"`
}

type SupportConfig struct {
	ReviewFocus     bool `yaml:"review_focus"`
	BreakingChanges bool `yaml:"breaking_changes"`
	CoverageDrop    bool `yaml:"coverage_drop"`
}

func DefaultPrismConfig() PrismConfig {
	return PrismConfig{
		Output: OutputConfig{
			Triage: TriageConfig{
				PriorityScore:       true,
				RiskLevel:           true,
				EstimatedReviewTime: true,
				FilePriorityList:    true,
			},
			Support: SupportConfig{
				ReviewFocus:     true,
				BreakingChanges: true,
				CoverageDrop:    true,
			},
		},
	}
}
