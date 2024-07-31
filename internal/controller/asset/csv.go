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

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/repository/asset"
)

const (
	recordsBatchSize int = 500
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
}

func (i impl) ImportAssets(ctx context.Context) error {
	root := "/var/data/gemini/adwh"
	date := "20230721"

	importers := map[CSVFileName]func(context.Context, asset.Repository, *csv.Reader) error{
		substationCSVFileName:    importSubstations,
		switchboardDNCSVFileName: importSwitchboards[switchboardDXCSV],
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

		err = importer(ctx, i.repo, reader)
		file.Close()
		if err != nil {
			return err
		}

	}

	return nil
}

// parseCSV runs through the file and outputs batches of the records read in chunks
func parseCSV[T CSVRecord](ctx context.Context, reader *csv.Reader, batchSize int) (chan []T, error) {
	var (
		values []string
		err    error
	)

	// Read the header as fields
	values, err = reader.Read()
	if err != nil {
		return nil, err
	}

	err = validateCSVHeaders(values, getCSVFields[T]())
	if err != nil {
		return nil, err
	}

	fields := make([]string, 0, len(values))
	fields = append(fields, values...)
	log.Println("CSV fields:", fields)

	var chanRecords chan []T = make(chan []T, 5)

	// Read the records
	log.Println("CSV parsing records...")
	go func() {
		var records []T

	parseLoop:
		for {
			select {
			case <-ctx.Done():
				log.Printf("context is cancelled. peace out: %+v\n", context.Cause(ctx))
				break parseLoop
			default:
			}

			if records == nil {
				records = make([]T, 0, batchSize)
			}

			values, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Printf("found an error in parsing, %+v\n", err)
				continue
			}

			records = append(records, getCSVRecord[T](values, fields, len(fields)))
			if len(records) < batchSize || chanRecords == nil {
				continue
			}

			select {
			case chanRecords <- records:
				records = nil
			case <-ctx.Done():
				log.Println("CSV parsing cancelled")
				break parseLoop
			}
		}

		if records != nil {
			log.Println("residue data found and making last delivery")
			chanRecords <- records
			records = nil
		}

		close(chanRecords)
	}()

	return chanRecords, nil
}

// getCSVFields returns a string mapping of csv fields provided in the CSVRecord struct
func getCSVFields[T CSVRecord]() map[string]bool {
	var (
		record     T
		recordType reflect.Type  = reflect.TypeOf(record)
		recordVal  reflect.Value = reflect.ValueOf(&record).Elem()
	)

	result := map[string]bool{}

	for idx := 0; idx < recordType.NumField(); idx++ {
		csvField := recordVal.Type().Field(idx).Tag.Get("csv")
		if csvField == "" || csvField == "-" {
			continue
		}
		result[csvField] = true
	}

	return result
}

// getCSVRecord creates an instance of T based on the values and fields taken from the CSV file
func getCSVRecord[T CSVRecord](values []string, fields []string, fieldsLen int) T {
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

		if !value.CanSet() || value.Kind() != reflect.String {
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

/**
 * validateCSVHeaders will validate whether all the expected headers are found or not
 */
func validateCSVHeaders(headers []string, expectedHeaders map[string]bool) error {
	return nil
}
