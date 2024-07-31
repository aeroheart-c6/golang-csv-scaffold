package asset

import (
	"context"
	"encoding/csv"
	"errors"
	"log"

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/model"
	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/repository/asset"
)

type substationCSV struct {
	AssetID string `csv:"asset_id"`
	Name    string `csv:"substation_name"`
	Status  string `csv:"asset_status"`
	Network string `csv:"object_type"`
}

func (s substationCSV) convertStatus() (model.AssetStatus, error) {
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

func (s substationCSV) convertNetwork() (model.Network, error) {
	network, ok := map[string]model.Network{
		"FNDXSS": model.NetworkDX,
		"FNTXSS": model.NetworkTX,
	}[s.Network]

	if !ok {
		return "", errors.New("found invalid value for network")
	}

	return network, nil
}

func (s substationCSV) toModel() (model.Substation, error) {
	status, err := s.convertStatus()
	if err != nil {
		return model.Substation{}, errors.New("substation status is invalid")
	}

	network, err := s.convertNetwork()
	if err != nil {
		return model.Substation{}, errors.New("network is invalid")
	}

	return model.Substation{
		AssetID: s.AssetID,
		Name:    s.Name,
		Status:  status,
		Network: network,
	}, nil
}

// ImportSubstations runs through the CSV file and saves them into the database
func importSubstations(ctx context.Context, repo asset.Repository, reader *csv.Reader) error {
	var (
		chanRecords chan []substationCSV
		err         error
	)

	/*
	 *	cancelCtx, cancelFn := context.WithCancelCause(ctx)
	 *	go func() {
	 *		time.Sleep(2 * time.Second)
	 *		cancelFn(errors.New("giving up on ingestion"))
	 *	}()
	 */

	chanRecords, err = parseCSV[substationCSV](
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

		for idx, record := range records {
			var (
				m   model.Substation
				err error
			)

			// translate from raw CSV record to model
			m, err = record.toModel()
			if err != nil {
				log.Printf("skipping row (%d) %v", idx, err)
				continue
			}

			models = append(models, m)
		}

		err = repo.UpsertSubstations(ctx, models)
		if err != nil {
			return err
		}
	}

	return nil
}
