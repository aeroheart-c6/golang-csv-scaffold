package asset

import (
	"context"
	"encoding/csv"
	"errors"
	"log"
	"time"

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
	chanRecords, chanErr, err := parseAssetCSV[SubstationCSV](
		ctx,
		reader,
		recordsBatchSize,
	)
	if err != nil {
		return err
	}

	for records := range chanRecords {
		substations := make([]model.Substation, 0, recordsBatchSize)

		for _, record := range records {
			status, err := record.ConvertStatus()
			if err != nil {
				continue
			}

			network, err := record.ConvertNetwork()
			if err != nil {
				continue
			}

			substations = append(substations, model.Substation{
				ID:        0,
				AssetID:   record.AssetID,
				Name:      record.Name,
				Status:    status,
				Network:   network,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}

		log.Println("CSV saving records...")
		i.repo.UpsertSubstations(ctx, substations)
		break
	}

	if len(chanErr) == 0 {
		return nil
	}

	return <-chanErr
}