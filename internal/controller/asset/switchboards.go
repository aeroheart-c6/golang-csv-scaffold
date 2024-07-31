package asset

type SwitchboardCSV struct {
	AssetID string
	Name    string
}

func (switchboard *SwitchboardCSV) FromMap(mapping map[string]string) error {
	switchboard = &SwitchboardCSV{
		AssetID: mapping["asset_id"],
		Name:    mapping["name"],
	}

	return nil
}
