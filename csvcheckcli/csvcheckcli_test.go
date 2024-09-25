package csvcheckcli_test

import (
	"csvcheckcli/csvcheckcli"
	"encoding/csv"
	"fmt"
	"strings"
	"testing"

	"github.com/BrianWeiHaoMa/csvcheck"
	"github.com/stretchr/testify/assert"
)

func Get2DArrayFromCsvString(csvString string) [][]csvcheck.StringHashable {
	reader := csv.NewReader(strings.NewReader(csvString))
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	res := make([][]csvcheck.StringHashable, len(records))
	for i, row := range records {
		res[i] = make([]csvcheck.StringHashable, len(row))
		for j, cell := range row {
			res[i][j] = csvcheck.BasicStringHashable(cell)
		}
	}
	return res
}

type userInputSolid struct {
	inputDir            string
	files               []string
	method              string
	function            string
	keepIndex           bool
	outputDir           string
	addTimestamp        bool
	columnsToUse        []string
	columnsToIgnore     []string
	autoAlign           bool
	useCommonColumns    bool
	ColumnsToKeep       []string
	ColumnsToDelete     []string
	ColumnsArrangement1 []string
	ColumnsArrangement2 []string
}

func (o userInputSolid) getUserInput() csvcheckcli.UserInput {
	return csvcheckcli.UserInput{
		InputDir:            &o.inputDir,
		Files:               &o.files,
		Method:              &o.method,
		Function:            &o.function,
		KeepIndex:           &o.keepIndex,
		OutputDir:           &o.outputDir,
		AddTimestamp:        &o.addTimestamp,
		ColumnsToUse:        &o.columnsToUse,
		ColumnsToIgnore:     &o.columnsToIgnore,
		AutoAlign:           &o.autoAlign,
		UseCommonColumns:    &o.useCommonColumns,
		ColumnsToKeep:       &o.ColumnsToKeep,
		ColumnsToDelete:     &o.ColumnsToDelete,
		ColumnsArrangement1: &o.ColumnsArrangement1,
		ColumnsArrangement2: &o.ColumnsArrangement2,
	}
}

func TestParseUserInputProperAndImproperInputs(t *testing.T) {
	for i, data := range []struct {
		input       csvcheckcli.UserInput
		expectError bool
	}{
		{
			input: userInputSolid{
				inputDir:         "/path/to/input/dir",
				files:            []string{"file1.csv", "file2.csv"},
				method:           csvcheckcli.MethodStringMatch,
				function:         csvcheckcli.FunctionStringCommon,
				keepIndex:        true,
				outputDir:        "/path/to/output/dir",
				addTimestamp:     true,
				columnsToUse:     []string{"column1", "column2"},
				columnsToIgnore:  nil,
				autoAlign:        true,
				useCommonColumns: false,
			}.getUserInput(),
			expectError: false,
		},
		{
			input: userInputSolid{
				inputDir:         "/path/to/input/dir",
				files:            []string{"file1.csv", "file2.csv"},
				method:           csvcheckcli.MethodStringMatch,
				function:         csvcheckcli.FunctionStringCommon,
				keepIndex:        false,
				outputDir:        "/path/to/output/dir",
				addTimestamp:     true,
				columnsToUse:     nil,
				columnsToIgnore:  nil,
				autoAlign:        true,
				useCommonColumns: true,
			}.getUserInput(),
			expectError: false,
		},
		{
			input: userInputSolid{
				inputDir:         "/path/to/input/dir",
				files:            []string{"file2.csv"},
				method:           csvcheckcli.MethodStringMatch,
				function:         csvcheckcli.FunctionStringCommon,
				keepIndex:        false,
				outputDir:        "/path/to/output/dir",
				addTimestamp:     true,
				columnsToUse:     nil,
				columnsToIgnore:  nil,
				autoAlign:        true,
				useCommonColumns: true,
			}.getUserInput(),
			expectError: true,
		},
		{
			input: userInputSolid{
				inputDir:         "/path/to/input/dir",
				files:            []string{"file2.csv", "file2.csv", "file2.csv"},
				method:           csvcheckcli.MethodStringMatch,
				function:         csvcheckcli.FunctionStringCommon,
				keepIndex:        false,
				outputDir:        "/path/to/output/dir",
				addTimestamp:     true,
				columnsToUse:     nil,
				columnsToIgnore:  nil,
				autoAlign:        true,
				useCommonColumns: true,
			}.getUserInput(),
			expectError: true,
		},
		{
			input: userInputSolid{
				inputDir:         "/path/to/input/dir",
				files:            []string{"file2.csv", "file2.csv"},
				method:           csvcheckcli.MethodStringMatch,
				function:         csvcheckcli.FunctionStringCommon,
				keepIndex:        false,
				outputDir:        "/path/to/output/dir",
				addTimestamp:     true,
				columnsToUse:     []string{"column1", "column2"},
				columnsToIgnore:  nil,
				autoAlign:        true,
				useCommonColumns: true,
			}.getUserInput(),
			expectError: true,
		},
		{
			input: userInputSolid{
				inputDir:         "/path/to/input/dir",
				files:            []string{"file2.csv", "file2.csv"},
				method:           csvcheckcli.MethodStringMatch,
				function:         csvcheckcli.FunctionStringCommon,
				keepIndex:        false,
				outputDir:        "/path/to/output/dir",
				addTimestamp:     true,
				columnsToUse:     []string{"column1", "column2"},
				columnsToIgnore:  []string{"column1", "column2"},
				autoAlign:        true,
				useCommonColumns: false,
			}.getUserInput(),
			expectError: true,
		},
		{
			input: userInputSolid{
				inputDir:         "/path/to/input/dir",
				files:            []string{"file2.csv", "file2.csv"},
				method:           "subtract",
				function:         csvcheckcli.FunctionStringCommon,
				keepIndex:        false,
				outputDir:        "/path/to/output/dir",
				addTimestamp:     true,
				columnsToUse:     nil,
				columnsToIgnore:  nil,
				autoAlign:        true,
				useCommonColumns: true,
			}.getUserInput(),
			expectError: true,
		},
		{
			input: userInputSolid{
				inputDir:         "/path/to/input/dir",
				files:            []string{"file2.csv", "file2.csv"},
				method:           csvcheckcli.MethodStringMatch,
				function:         "flip",
				keepIndex:        false,
				outputDir:        "/path/to/output/dir",
				addTimestamp:     true,
				columnsToUse:     nil,
				columnsToIgnore:  nil,
				autoAlign:        true,
				useCommonColumns: true,
			}.getUserInput(),
			expectError: true,
		},
		{
			input: userInputSolid{
				inputDir:            "/path/to/input/dir",
				files:               []string{"file2.csv", "file2.csv"},
				method:              csvcheckcli.MethodStringMatch,
				function:            "common",
				keepIndex:           false,
				outputDir:           "/path/to/output/dir",
				addTimestamp:        true,
				columnsToUse:        []string{"column1", "column2"},
				columnsToIgnore:     nil,
				autoAlign:           true,
				useCommonColumns:    false,
				ColumnsToKeep:       []string{"column1"},
				ColumnsToDelete:     nil,
				ColumnsArrangement1: nil,
				ColumnsArrangement2: []string{"column1"},
			}.getUserInput(),
			expectError: false,
		},
		{
			input: userInputSolid{
				inputDir:            "/path/to/input/dir",
				files:               []string{"file2.csv", "file2.csv"},
				method:              csvcheckcli.MethodStringMatch,
				function:            "common",
				keepIndex:           false,
				outputDir:           "/path/to/output/dir",
				addTimestamp:        true,
				columnsToUse:        []string{"column1", "column2"},
				columnsToIgnore:     nil,
				autoAlign:           true,
				useCommonColumns:    false,
				ColumnsToKeep:       []string{"column1"},
				ColumnsToDelete:     []string{"column1"},
				ColumnsArrangement1: nil,
				ColumnsArrangement2: []string{"column1"},
			}.getUserInput(),
			expectError: true,
		},
		{
			input: userInputSolid{
				inputDir:            "/path/to/input/dir",
				files:               []string{"file2.csv", "file2.csv"},
				method:              csvcheckcli.MethodStringMatch,
				function:            "common",
				keepIndex:           false,
				outputDir:           "/path/to/output/dir",
				addTimestamp:        true,
				columnsToUse:        nil,
				columnsToIgnore:     nil,
				autoAlign:           true,
				useCommonColumns:    true,
				ColumnsToKeep:       nil,
				ColumnsToDelete:     []string{"column1"},
				ColumnsArrangement1: []string{"column1"},
				ColumnsArrangement2: []string{"column1"},
			}.getUserInput(),
			expectError: false,
		},
		{
			input: userInputSolid{
				inputDir:            "/path/to/input/dir",
				files:               []string{"file2.csv", "file2.csv"},
				method:              csvcheckcli.MethodStringMatch,
				function:            "",
				keepIndex:           false,
				outputDir:           "/path/to/output/dir",
				addTimestamp:        true,
				columnsToUse:        nil,
				columnsToIgnore:     nil,
				autoAlign:           true,
				useCommonColumns:    true,
				ColumnsToKeep:       nil,
				ColumnsToDelete:     []string{"column1"},
				ColumnsArrangement1: []string{"column1"},
				ColumnsArrangement2: []string{"column1"},
			}.getUserInput(),
			expectError: true,
		},
	} {
		indexString := fmt.Sprintf("Test case index: %d", i)
		if data.expectError {
			_, err := csvcheckcli.ParseUserInput(&data.input)
			assert.NotNil(t, err, indexString)
		} else {
			_, err := csvcheckcli.ParseUserInput(&data.input)
			assert.Nil(t, err, indexString)
		}
	}
}

func TestGetResArraysCommonMatchKeepIndex(t *testing.T) {
	input := userInputSolid{
		inputDir:         "/path/to/input/dir",
		files:            []string{"file2.csv", "file2.csv"},
		method:           csvcheckcli.MethodStringMatch,
		function:         csvcheckcli.FunctionStringCommon,
		keepIndex:        true,
		outputDir:        "/path/to/output/dir",
		addTimestamp:     true,
		columnsToUse:     nil,
		columnsToIgnore:  nil,
		autoAlign:        false,
		useCommonColumns: true,
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
a,b,c
7,8,9
1,2,3
1,2,3
4,5,6
7,8,9
7,8,9
10,10,10
`)
	res1, res2, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	expected1 := Get2DArrayFromCsvString(fmt.Sprintf(`
a,b,c,%s
1,2,3,1
4,5,6,2
7,8,9,3
7,8,9,4
`, csvcheckcli.IndexColumnName))

	expected2 := Get2DArrayFromCsvString(fmt.Sprintf(`
a,b,c,%s
7,8,9,1
1,2,3,2
4,5,6,4
7,8,9,5
`, csvcheckcli.IndexColumnName))

	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)
	assert.Equal(t, expected2, res2)
}

func TestGetResArraysDifferentSetKeepIndex(t *testing.T) {
	input := userInputSolid{
		inputDir:         "/path/to/input/dir",
		files:            []string{"file2.csv", "file2.csv"},
		method:           csvcheckcli.MethodStringSet,
		function:         csvcheckcli.FunctionStringDifferent,
		keepIndex:        true,
		outputDir:        "/path/to/output/dir",
		addTimestamp:     true,
		columnsToUse:     nil,
		columnsToIgnore:  nil,
		autoAlign:        false,
		useCommonColumns: true,
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
-1,-1,-1
1,2,3
4,5,6
7,8,9
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
a,b,c
7,8,9
1,2,3
1,2,3
4,5,6
7,8,9
7,8,9
10,10,10
`)
	res1, res2, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	expected1 := Get2DArrayFromCsvString(fmt.Sprintf(`
a,b,c,%s
-1,-1,-1,1
`, csvcheckcli.IndexColumnName))

	expected2 := Get2DArrayFromCsvString(fmt.Sprintf(`
a,b,c,%s
10,10,10,7
`, csvcheckcli.IndexColumnName))

	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)
	assert.Equal(t, expected2, res2)
}

func TestGetResArraysDifferentMatch(t *testing.T) {
	input := userInputSolid{
		inputDir:         "/path/to/input/dir",
		files:            []string{"file2.csv", "file2.csv"},
		method:           csvcheckcli.MethodStringMatch,
		function:         csvcheckcli.FunctionStringDifferent,
		keepIndex:        false,
		outputDir:        "/path/to/output/dir",
		addTimestamp:     true,
		columnsToUse:     nil,
		columnsToIgnore:  nil,
		autoAlign:        false,
		useCommonColumns: true,
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
a,b,c
7,8,9
1,2,3
1,2,3
4,5,6
7,8,9
7,8,9
10,10,10
`)
	res1, res2, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	expected1 := Get2DArrayFromCsvString(`
a,b,c
`)

	expected2 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
7,8,9
10,10,10
`)

	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)
	assert.Equal(t, expected2, res2)
}

func TestGetResArraysDifferentMatchColumnsToUse(t *testing.T) {
	input := userInputSolid{
		inputDir:         "/path/to/input/dir",
		files:            []string{"file2.csv", "file2.csv"},
		method:           csvcheckcli.MethodStringMatch,
		function:         csvcheckcli.FunctionStringDifferent,
		keepIndex:        false,
		outputDir:        "/path/to/output/dir",
		addTimestamp:     true,
		columnsToUse:     []string{"a"},
		columnsToIgnore:  nil,
		autoAlign:        false,
		useCommonColumns: false,
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,1,z
z,7,z
z,7,z
`)
	res1, res2, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	expected1 := Get2DArrayFromCsvString(`
a,b,c
4,5,6
`)

	expected2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,7,z
`)

	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)
	assert.Equal(t, expected2, res2)
}

func TestGetResArraysDifferentMatchColumnsToUseAutoAlign(t *testing.T) {
	input := userInputSolid{
		inputDir:         "/path/to/input/dir",
		files:            []string{"file2.csv", "file2.csv"},
		method:           csvcheckcli.MethodStringMatch,
		function:         csvcheckcli.FunctionStringDifferent,
		keepIndex:        false,
		outputDir:        "/path/to/output/dir",
		addTimestamp:     true,
		columnsToUse:     []string{"a"},
		columnsToIgnore:  nil,
		autoAlign:        true,
		useCommonColumns: false,
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,1,z
z,7,z
z,7,z
`)
	res1, res2, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	expected1 := Get2DArrayFromCsvString(`
a,b,c
4,5,6
`)

	expected2 := Get2DArrayFromCsvString(`
a,d,z
b,a,c
7,z,z
`)

	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)
	assert.Equal(t, expected2, res2)
}

func TestGetResArraysDifferentMatchUseCommonColumnsAutoAlign(t *testing.T) {
	input := userInputSolid{
		inputDir:         "/path/to/input/dir",
		files:            []string{"file2.csv", "file2.csv"},
		method:           csvcheckcli.MethodStringMatch,
		function:         csvcheckcli.FunctionStringDifferent,
		keepIndex:        false,
		outputDir:        "/path/to/output/dir",
		addTimestamp:     true,
		columnsToUse:     []string{"a"},
		columnsToIgnore:  nil,
		autoAlign:        true,
		useCommonColumns: false,
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,1,z
z,7,z
z,7,z
`)
	res1, res2, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	expected1 := Get2DArrayFromCsvString(`
a,b,c
4,5,6
`)

	expected2 := Get2DArrayFromCsvString(`
a,d,z
b,a,c
7,z,z
`)

	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)
	assert.Equal(t, expected2, res2)
}

func TestGetResArraysDifferentMatchUseCommonColumnsAutoAlignKeepCommonAndDifferentColumns(t *testing.T) {
	input := userInputSolid{
		inputDir:         "/path/to/input/dir",
		files:            []string{"file2.csv", "file2.csv"},
		method:           csvcheckcli.MethodStringMatch,
		function:         csvcheckcli.FunctionStringDifferent,
		keepIndex:        false,
		outputDir:        "/path/to/output/dir",
		addTimestamp:     true,
		columnsToUse:     []string{"a"},
		columnsToIgnore:  nil,
		autoAlign:        true,
		useCommonColumns: false,
		ColumnsToKeep:    []string{"a", "z", "c"},
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,1,z
z,7,z
z,7,z
`)
	res1, res2, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	expected1 := Get2DArrayFromCsvString(`
a,c
4,6
`)

	expected2 := Get2DArrayFromCsvString(`
a,z
b,c
7,z
`)

	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)
	assert.Equal(t, expected2, res2)
}

func TestGetResArraysDifferentMatchUseCommonColumnsAutoAlignKeepCommonAndDifferentColumnsAndRearrange(t *testing.T) {
	input := userInputSolid{
		inputDir:            "/path/to/input/dir",
		files:               []string{"file2.csv", "file2.csv"},
		method:              csvcheckcli.MethodStringMatch,
		function:            csvcheckcli.FunctionStringDifferent,
		keepIndex:           false,
		outputDir:           "/path/to/output/dir",
		addTimestamp:        true,
		columnsToUse:        []string{"a"},
		columnsToIgnore:     nil,
		autoAlign:           true,
		useCommonColumns:    false,
		ColumnsToKeep:       []string{"a", "z", "c"},
		ColumnsArrangement1: []string{"c", "a"},
		ColumnsArrangement2: []string{"z", "a"},
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,1,z
z,7,z
z,7,z
`)
	res1, res2, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	expected1 := Get2DArrayFromCsvString(`
c,a
6,4
`)

	expected2 := Get2DArrayFromCsvString(`
z,a
c,b
z,7
`)

	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)
	assert.Equal(t, expected2, res2)
}

func TestGetResArraysDifferentMatchRearrangeWithNonExistentColumnsForColumnsArrangement1(t *testing.T) {
	input := userInputSolid{
		inputDir:            "/path/to/input/dir",
		files:               []string{"file2.csv", "file2.csv"},
		method:              csvcheckcli.MethodStringMatch,
		function:            csvcheckcli.FunctionStringDifferent,
		keepIndex:           false,
		outputDir:           "/path/to/output/dir",
		addTimestamp:        true,
		columnsToUse:        []string{"a"},
		columnsToIgnore:     nil,
		autoAlign:           true,
		useCommonColumns:    false,
		ColumnsToKeep:       []string{"a", "z", "c"},
		ColumnsArrangement1: []string{"c", "a"},
		ColumnsArrangement2: []string{"z", "a", "p"},
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,1,z
z,7,z
z,7,z
`)
	_, _, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	assert.NotNil(t, err)
}

func TestGetResArraysDifferentMatchRearrangeWithNonExistentColumnsForColumnsArrangement2(t *testing.T) {
	input := userInputSolid{
		inputDir:            "/path/to/input/dir",
		files:               []string{"file2.csv", "file2.csv"},
		method:              csvcheckcli.MethodStringMatch,
		function:            csvcheckcli.FunctionStringDifferent,
		keepIndex:           false,
		outputDir:           "/path/to/output/dir",
		addTimestamp:        true,
		columnsToUse:        []string{"a"},
		columnsToIgnore:     nil,
		autoAlign:           true,
		useCommonColumns:    false,
		ColumnsToKeep:       []string{"a", "z", "c"},
		ColumnsArrangement1: []string{"c", "a", "aaa"},
		ColumnsArrangement2: []string{"z", "a"},
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,1,z
z,7,z
z,7,z
`)
	_, _, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	assert.NotNil(t, err)
}

func TestGetResArraysDifferentMatchRearrangeWithMissingColumnsOnColumnsArrangement1(t *testing.T) {
	input := userInputSolid{
		inputDir:            "/path/to/input/dir",
		files:               []string{"file2.csv", "file2.csv"},
		method:              csvcheckcli.MethodStringMatch,
		function:            csvcheckcli.FunctionStringDifferent,
		keepIndex:           false,
		outputDir:           "/path/to/output/dir",
		addTimestamp:        true,
		columnsToUse:        []string{"a"},
		columnsToIgnore:     nil,
		autoAlign:           true,
		useCommonColumns:    false,
		ColumnsToKeep:       []string{"a", "z", "c"},
		ColumnsArrangement1: []string{"c"},
		ColumnsArrangement2: []string{"z", "a"},
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,1,z
z,7,z
z,7,z
`)
	_, _, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	assert.NotNil(t, err)
}

func TestGetResArraysDifferentMatchRearrangeWithMissingColumnsOnColumnsArrangement2(t *testing.T) {
	input := userInputSolid{
		inputDir:            "/path/to/input/dir",
		files:               []string{"file2.csv", "file2.csv"},
		method:              csvcheckcli.MethodStringMatch,
		function:            csvcheckcli.FunctionStringDifferent,
		keepIndex:           false,
		outputDir:           "/path/to/output/dir",
		addTimestamp:        true,
		columnsToUse:        []string{"a"},
		columnsToIgnore:     nil,
		autoAlign:           true,
		useCommonColumns:    false,
		ColumnsToKeep:       []string{"a", "z", "c"},
		ColumnsArrangement1: []string{"c", "a"},
		ColumnsArrangement2: []string{"a"},
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,1,z
z,7,z
z,7,z
`)
	_, _, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	assert.NotNil(t, err)
}

func TestGetResArraysDifferentMatchUseCommonColumnsAutoAlignDeleteCommonAndDifferentColumnsAndRearrange(t *testing.T) {
	input := userInputSolid{
		inputDir:            "/path/to/input/dir",
		files:               []string{"file2.csv", "file2.csv"},
		method:              csvcheckcli.MethodStringMatch,
		function:            csvcheckcli.FunctionStringDifferent,
		keepIndex:           false,
		outputDir:           "/path/to/output/dir",
		addTimestamp:        true,
		columnsToUse:        []string{"a"},
		columnsToIgnore:     nil,
		autoAlign:           true,
		useCommonColumns:    false,
		ColumnsToKeep:       nil,
		ColumnsToDelete:     []string{"a", "c"},
		ColumnsArrangement1: nil,
		ColumnsArrangement2: nil,
	}.getUserInput()

	arr1 := Get2DArrayFromCsvString(`
a,b,c
1,2,3
4,5,6
7,8,9
`)
	arr2 := Get2DArrayFromCsvString(`
d,a,z
a,b,c
z,1,z
z,7,z
z,7,z
`)
	res1, res2, err := csvcheckcli.GetResArrays(arr1, arr2, input)

	expected1 := Get2DArrayFromCsvString(`
b
5
`)

	expected2 := Get2DArrayFromCsvString(`
d,z
a,c
z,z
`)

	assert.Nil(t, err)
	assert.Equal(t, expected1, res1)
	assert.Equal(t, expected2, res2)
}
