package asset

import (
	"context"
	"encoding/csv"
	"errors"
	"log"

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/model"
)

type SwitchboardCSV struct {
	SubstationID string `csv:"parent_asset_id"`
	AssetID      string `csv:"asset_id"`
	Name         string `csv:"switchboard_name"`
	Status       string `csv:"asset_equipment_status"`
}

func (s SwitchboardCSV) ConvertStatus() (model.AssetStatus, error) {
	status, ok := map[string]model.AssetStatus{
		"C":  model.AssetStatusCommissioned,
		"DC": model.AssetStatusDecommissioned,
		"DM": model.AssetStatusDemolished,
		"P":  model.AssetStatusPlanned,
	}[s.Status]

	if !ok {
		return "", errors.New("found invalid value for status")
	}

	return status, nil
}

func (i impl) ImportDNSwitchboards(ctx context.Context, reader *csv.Reader) error {
	chanRecords, err := parseAssetCSV[SwitchboardCSV](
		ctx,
		reader,
		recordsBatchSize,
	)
	if err != nil {
		return err
	}

	log.Println("importing distribution switchboards...")
	for records := range chanRecords {
		var (
			models      = make([]model.Switchboard, 0, len(records))
			substations = make(map[string]model.Substation, 0)
		)

		for _, record := range records {
			status, err := record.ConvertStatus()
			if err != nil {
				continue
			}

			substation, ok := substations[record.SubstationID]
			if !ok {
				substation, err = i.repo.GetSubstation(ctx, record.SubstationID)
				if err != nil {
					continue
				}
				substations[record.SubstationID] = substation
			}

			models = append(models, model.Switchboard{
				AssetID:           record.AssetID,
				SubstationID:      substation.ID,
				Name:              record.Name,
				Status:            status,
				Network:           model.NetworkDX,
				SubstationAssetID: record.SubstationID,
			})
		}

		log.Println("CSV saving records...")
		err := i.repo.UpsertSwitchboards(ctx, models)
		if err != nil {
			return err
		}
	}

	return nil
}
