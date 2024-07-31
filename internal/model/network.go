package model

// Network represents the power grid network types. Primarily used for populating `network_type` fields in the PG
type Network string

// Network enum values
const (
	NetworkDX Network = "DX"
	NetworkTX Network = "TX"
)

// String returns the string representation of this enum value
func (network Network) String() string {
	return string(network)
}

// IsValid checks if given value matches one of the enum values declared
func (network Network) IsValid() bool {
	switch network {
	case
		NetworkDX,
		NetworkTX:
		return true
	default:
		return false
	}
}
