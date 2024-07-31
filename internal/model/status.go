package model

type AssetStatus string

func (s AssetStatus) String() string {
	return string(s)
}

func (s AssetStatus) IsValid() bool {
	switch s {
	case
		AssetStatusCommissioned,
		AssetStatusDecommissioned,
		AssetStatusDemolished,
		AssetStatusPlanned:
		return true
	default:
		return false
	}
}

const (
	AssetStatusCommissioned   = "COMMISSIONED"
	AssetStatusDecommissioned = "DECOMMISSIONED"
	AssetStatusDemolished     = "DEMOLISHED"
	AssetStatusPlanned        = "PLANNED"
)
