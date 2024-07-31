package asset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	csvFields []string = []string{
		"one",
		"two",
		"three",
	}
	csvValues []string = []string{
		"Foo",
		"Bar",
		"Baz",
	}
)

func Test_makeAssetRecord_OK(t *testing.T) {
	defer func() {
		require.Equal(t, recover(), nil)
	}()

	type CSVRecordOK struct {
		One   string `csv:"one"`
		Two   string `csv:"two"`
		Three string `csv:"three"`
	}

	record := makeAssetRecord[CSVRecordOK](csvValues, csvFields, len(csvFields))

	require.Equal(t, record, CSVRecordOK{
		One:   csvValues[0],
		Two:   csvValues[1],
		Three: csvValues[2],
	})
}

func Test_makeAssetRecord_MissingTags(t *testing.T) {
	defer func() {
		require.Equal(t, recover(), nil)
	}()

	type CSVRecordMissingTags struct {
		One   string `csv:"one"`
		Two   string
		Three string `csv:"three"`
	}

	record := makeAssetRecord[CSVRecordMissingTags](csvValues, csvFields, len(csvFields))

	require.Equal(t, record, CSVRecordMissingTags{
		One:   csvValues[0],
		Two:   "",
		Three: csvValues[2],
	})
}

func Test_makeAssetRecord_Private(t *testing.T) {
	defer func() {
		require.Equal(t, recover(), nil)
	}()

	type CSVRecordPrivate struct {
		One   string `csv:"one"`
		two   string
		Three string `csv:"three"`
	}

	record := makeAssetRecord[CSVRecordPrivate](csvValues, csvFields, len(csvFields))

	require.Equal(t, record, CSVRecordPrivate{
		One:   csvValues[0],
		two:   "",
		Three: csvValues[2],
	})
}

func Test_makeAssetRecord_NonString(t *testing.T) {
	defer func() {
		require.Equal(t, recover(), nil)
	}()

	type CSVRecordNonString struct {
		One   string `csv:"one"`
		Two   string `csv:"two"`
		Three int    `csv:"three"`
	}

	record := makeAssetRecord[CSVRecordNonString](csvValues, csvFields, len(csvFields))

	require.Equal(t, record, CSVRecordNonString{
		One:   csvValues[0],
		Two:   csvValues[1],
		Three: 0,
	})
}

func Test_makeAssetRecord_Pointer(t *testing.T) {
	defer func() {
		require.Equal(t, recover(), nil)
	}()

	type CSVRecordPointer struct {
		One   string  `csv:"one"`
		Two   string  `csv:"two"`
		Three *string `csv:"three"`
	}

	record := makeAssetRecord[CSVRecordPointer](csvValues, csvFields, len(csvFields))

	require.Equal(t, record, CSVRecordPointer{
		One:   csvValues[0],
		Two:   csvValues[1],
		Three: nil,
	})
}

func Test_makeAssetRecord_Nested(t *testing.T) {
	defer func() {
		require.Equal(t, recover(), nil)
	}()

	type CSVRecordNested struct {
		One   string `csv:"one"`
		Two   string `csv:"two"`
		Three struct {
			Value string
		} `csv:"three"`
	}

	record := makeAssetRecord[CSVRecordNested](csvValues, csvFields, len(csvFields))

	require.Equal(t, record, CSVRecordNested{
		One:   csvValues[0],
		Two:   csvValues[1],
		Three: struct{ Value string }{},
	})
}
