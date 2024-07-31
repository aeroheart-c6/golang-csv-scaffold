package asset

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"reflect"
)

const (
	recordsBatchSize int = 100
)

type CSVFileName string

func (f CSVFileName) IsValid() bool {
	switch f {
	case
		substationCSVFileName,
		switchboardDNCSVFileName,
		switchboardPanelDNFileName:
		return true
	default:
		return false
	}
}

func (f CSVFileName) String() string {
	return string(f)
}

const (
	substationCSVFileName      CSVFileName = "adwh_elec_substation_%s.csv"
	switchboardDNCSVFileName   CSVFileName = "adwh_elec_dx_swb_%s.csv"
	switchboardPanelDNFileName CSVFileName = "adwh_elec_dx_swb_pnl_%s.csv"
)

type CSVRecord interface {
	SubstationCSV |
		SwitchboardCSV
}

func (i impl) ImportAssets(ctx context.Context) error {
	root := "/var/data/gemini/adwh"
	date := "20230721"

	importers := map[CSVFileName]func(context.Context, *csv.Reader) error{
		substationCSVFileName:    i.ImportSubstations,
		switchboardDNCSVFileName: i.ImportDNSwitchboards,
	}

	for _, pattern := range []CSVFileName{
		substationCSVFileName,
		switchboardDNCSVFileName,
	} {
		log.Println("Reading file:", pattern)

		importer, ok := importers[pattern]
		if !ok {
			return fmt.Errorf("importer not found for %v", pattern)
		}

		file, err := os.Open(path.Join(root, fmt.Sprintf(pattern.String(), date)))
		if err != nil {
			return err
		}

		reader := csv.NewReader(file)
		reader.ReuseRecord = true

		err = importer(ctx, reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseAssetCSV[T CSVRecord](ctx context.Context, reader *csv.Reader, batchSize int) (chan []T, chan error, error) {
	// Read the header as fields
	values, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}
	fields := make([]string, 0, len(values))
	fields = append(fields, values...)

	log.Println("CSV fields:", fields)

	var (
		chanRecords = make(chan []T, 5)
		chanErr     = make(chan error)
	)

	// Read the records
	log.Println("CSV parsing records...")
	go func() {
		var records []T = make([]T, 0, batchSize)

		for {
			values, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				chanErr <- err
				break
			}

			records = append(records, makeAssetRecord[T](values, fields, len(fields)))

			if len(records) >= batchSize && chanRecords != nil {
				chanRecords <- records
			}
		}

		if err != nil && err != io.EOF {
			chanErr <- err
		}

		close(chanRecords)
		close(chanErr)
	}()

	return chanRecords, chanErr, nil
}

func makeAssetRecord[T CSVRecord](values []string, fields []string, fieldsLen int) T {
	mapping := make(map[string]string, fieldsLen)

	// zip the field and values together in a map
	for idx, field := range fields {
		mapping[field] = values[idx]
	}

	var (
		record     T
		recordType reflect.Type  = reflect.TypeOf(record)
		recordVal  reflect.Value = reflect.ValueOf(&record).Elem()
	)

	// the annoying part of assigning to struct using reflection
	for idx := 0; idx < recordType.NumField(); idx++ {
		field := recordType.Field(idx).Tag.Get("csv")
		value := recordVal.Field(idx)

		if !value.CanSet() && value.Kind() != reflect.String {
			continue
		}

		data, ok := mapping[field]
		if !ok {
			continue
		}

		value.SetString(data)
	}

	return record
}
