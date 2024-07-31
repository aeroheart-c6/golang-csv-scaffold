package asset

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log"

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/model"
	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/repository/asset"
)

type switchboardCSV interface {
	toModel() (model.Switchboard, error)
	parentAssetID() string
}
type switchboardDXCSV struct {
	SubstationID string `csv:"parent_asset_id"`
	AssetID      string `csv:"asset_id"`
	Name         string `csv:"switchboard_name"`
	Status       string `csv:"asset_equipment_status"`
}

func (s switchboardDXCSV) parentAssetID() string {
	return s.SubstationID
}

func (s switchboardDXCSV) convertStatus() (model.AssetStatus, error) {
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

func (s switchboardDXCSV) toModel() (model.Switchboard, error) {
	status, err := s.convertStatus()
	if err != nil {
		return model.Switchboard{}, errors.New("switchboard status is invalid")
	}

	return model.Switchboard{
		AssetID:           s.AssetID,
		Name:              s.Name,
		Status:            status,
		Network:           model.NetworkDX,
		SubstationAssetID: s.SubstationID,
	}, nil
}

func importSwitchboards[T switchboardCSV](ctx context.Context, repo asset.Repository, reader *csv.Reader) error {
	var (
		values []string
		err    error
	)

	// Read the header as fields
	values, err = reader.Read()
	if err != nil {
		return err
	}

	err = validateCSVHeaders(values, getCSVFields[T]())
	if err != nil {
		return err
	}

	log.Println("importing switchboards...")
	for {
		var (
			records  []T
			err      error
			errParse error
		)

		records, errParse = parseCSV[T](ctx, reader, recordsBatchSize)

		var (
			models      = make([]model.Switchboard, 0, len(records))
			substations = make(map[string]model.Substation, len(records))
		)
		for idx, record := range records {
			var (
				substationID string = record.parentAssetID()
				substation   model.Substation
				switchboard  model.Switchboard
				err          error
				ok           bool
			)

			// check if substation exists
			substation, ok = substations[substationID]
			if !ok {
				substation, err = repo.GetSubstation(ctx, substationID)
				if err != nil {
					log.Printf("unable to find substation [%s]. skipping.", substationID)
					continue
				}

				substations[substation.AssetID] = substation
			}

			switchboard, err = record.toModel()
			if err != nil {
				log.Printf("skipping row (%d) %v", idx, err)
				continue
			}

			switchboard.SubstationID = substation.ID

			models = append(models, switchboard)
		}

		log.Println("saving records...")
		err = repo.UpsertSwitchboards(ctx, models)
		if err != nil {
			return err
		}

		if errors.Is(errParse, io.EOF) {
			break
		}
	}

	return nil
}
