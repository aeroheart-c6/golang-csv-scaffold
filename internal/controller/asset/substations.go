package asset

import (
	"context"
	"encoding/csv"
	"errors"
	"log"

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/model"
)

type SubstationCSV struct {
	AssetID string `csv:"asset_id"`
	Name    string `csv:"substation_name"`
	Status  string `csv:"asset_status"`
	Network string `csv:"object_type"`
}

func (s SubstationCSV) ConvertStatus() (model.AssetStatus, error) {
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

func (s SubstationCSV) ConvertNetwork() (model.Network, error) {
	network, ok := map[string]model.Network{
		"FNDXSS": model.NetworkDX,
		"FNTXSS": model.NetworkTX,
	}[s.Network]

	if !ok {
		return "", errors.New("found invalid value for network")
	}

	return network, nil
}

// ImportSubstations runs through the CSV file and saves them into the database
func (i impl) ImportSubstations(ctx context.Context, reader *csv.Reader) error {
	chanRecords, err := parseAssetCSV[SubstationCSV](
		ctx,
		reader,
		recordsBatchSize,
	)
	if err != nil {
		return err
	}

	log.Println("importing substations...")
	for records := range chanRecords {
		models := make([]model.Substation, 0, recordsBatchSize)

		for _, record := range records {
			status, err := record.ConvertStatus()
			if err != nil {
				continue
			}

			network, err := record.ConvertNetwork()
			if err != nil {
				continue
			}

			models = append(models, model.Substation{
				AssetID: record.AssetID,
				Name:    record.Name,
				Status:  status,
				Network: network,
			})
		}

		err = i.repo.UpsertSubstations(ctx, models)
		if err != nil {
			return err
		}
	}

	return nil
}
