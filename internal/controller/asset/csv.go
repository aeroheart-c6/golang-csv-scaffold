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
		log.Println("preparing importer")
		importer, ok := importers[pattern]
		if !ok {
			return fmt.Errorf("importer not found for %v", pattern)
		}

		log.Println("reading file:", pattern)
		file, err := os.Open(path.Join(root, fmt.Sprintf(pattern.String(), date)))
		if err != nil {
			return err
		}

		reader := csv.NewReader(file)
		reader.ReuseRecord = true

		log.Println("importing file:", pattern)
		err = importer(ctx, i.repo, reader)
		file.Close()
		if err != nil {
			return err
		}

	}

	return nil
}

// parseCSV reads a chunk of data and returns it. Relies on `reader` to remember the location of the file
func parseCSV[T CSVRecord](ctx context.Context, reader *csv.Reader, batchSize int) ([]T, error) {
	var (
		fieldsRaw map[string]bool
		fields    []string
	)

	// get the required headers
	fieldsRaw = getCSVFields[T]()
	fields = make([]string, 0, len(fieldsRaw))
	for field := range fieldsRaw {
		fields = append(fields, field)
	}
	log.Println("fields:", fields)

	// Read the records
	log.Println("parsing data...")
	var (
		records []T
		err     error
	)
parseLoop:
	for loop := 0; loop < batchSize; loop++ {
		var row []string

		select {
		case <-ctx.Done():
			log.Printf("context is cancelled. peace out: %+v\n", context.Cause(ctx))
			break parseLoop
		default:
		}

		if records == nil {
			records = make([]T, 0, batchSize)
		}

		row, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("found an error in parsing, %+v\n", err)
			continue
		}

		records = append(records, getCSVRecord[T](row, fields, len(fields)))
	}

	return records, err
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
	/**
	 * Why not do the read here?
	 *
	 * It seems weird that function should require that the line it reads are the headers of the CSV data. Like this
	 * particular action is beyond its scope
	 */
	return nil
}
