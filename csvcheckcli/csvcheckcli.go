package csvcheckcli

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/BrianWeiHaoMa/csvcheck"

	"github.com/spf13/pflag"
)

const IndexColumnName = "_ind"

const MethodStringMatch = "match"
const MethodStringSet = "set"
const MethodStringDirect = "direct"

const FunctionStringCommon = "common"
const FunctionStringDifferent = "different"

var MethodMappings = map[string]int{
	MethodStringMatch:  csvcheck.MethodMatch,
	MethodStringSet:    csvcheck.MethodSet,
	MethodStringDirect: csvcheck.MethodDirect,
}

type UserInput struct {
	InputDir              *string
	Files                 *[]string
	Method                *string
	Function              *string
	KeepIndex             *bool
	OutputDir             *string
	AddTimestamp          *bool
	ColumnsToUse          *[]string
	ColumnsToIgnore       *[]string
	AutoAlign             *bool
	UseCommonColumns      *bool
	ColumnsToKeep         *[]string
	ColumnsToDelete       *[]string
	ColumnsArrangement1   *[]string
	ColumnsArrangement2   *[]string
	PrintInCsvFormat      *bool
	PrettyFormatMaxLength *int
}

func ParseUserInput(input *UserInput) (UserInput, error) {
	var res UserInput
	if input == nil {
		res = UserInput{}
		res.InputDir = pflag.StringP("inputdir", "d", "", "The directory containing the input files. This will be prepended to the input file paths. Must be given.")
		res.Files = pflag.StringSliceP("files", "f", []string{}, "The input files paths to compare. 2 should be provided.")
		res.Method = pflag.StringP("method", "m", "set", "The method to use for comparison. Options: match, set, direct. By default, set is used.")
		res.Function = pflag.StringP("function", "F", "", "The function to use for comparison. Options: common, different. A function must be given.")
		res.KeepIndex = pflag.BoolP("keepindex", "k", false, fmt.Sprintf("Whether to keep the indices from the original csv of the rows in the result (%s column will be added).", IndexColumnName))
		res.OutputDir = pflag.StringP("outputdir", "o", "", "The directory to write the output files to.")
		res.AddTimestamp = pflag.BoolP("addtimestamp", "t", false, "Whether or not to add a timestamp to the output file name.")
		res.ColumnsToUse = pflag.StringSliceP("usecolumns", "c", nil, "The columns to use for comparison.")
		res.ColumnsToIgnore = pflag.StringSliceP("ignorecolumns", "i", nil, "The columns to ignore for comparison.")
		res.AutoAlign = pflag.BoolP("autoalign", "a", false, "Whether or not to auto align the columns of the csv files. Common columns will be aligned on the left side.")
		res.UseCommonColumns = pflag.BoolP("usecommoncolumns", "C", false, "Whether to use all the common columns between the csv files for comparison.")
		res.ColumnsToKeep = pflag.StringSliceP("keepcolumns", "K", nil, "The columns to keep in the output.")
		res.ColumnsToDelete = pflag.StringSliceP("deletecolumns", "D", nil, "The columns to delete in the output.")
		res.ColumnsArrangement1 = pflag.StringSliceP("columnsarrangement1", "r", nil, "An arrangement for the columns in the first output.")
		res.ColumnsArrangement2 = pflag.StringSliceP("columnsarrangement2", "R", nil, "An arrangement for the columns in the second output.")
		res.PrintInCsvFormat = pflag.BoolP("csv", "p", false, "Whether to print the output in csv format. By default, the output is printed in a columns-aligned.")
		res.PrettyFormatMaxLength = pflag.IntP("prettyformatmaxlength", "l", -1, "The maximum length before truncation of a column entry when printing in pretty format. Negative values mean no limit. By default, there is no limit.")

		pflag.Parse()
	} else {
		res = *input
	}

	if *res.InputDir == "" {
		return UserInput{}, fmt.Errorf("inputdir must be given")
	}

	if len(*res.Files) != 2 {
		return UserInput{}, fmt.Errorf("exactly 2 file paths needed")
	}

	if _, exists := MethodMappings[*res.Method]; !exists {
		return UserInput{}, fmt.Errorf("unsupported method %s", *res.Method)
	}

	columnsCompInputCnt := 0
	if *res.ColumnsToUse != nil {
		columnsCompInputCnt++
	}
	if *res.ColumnsToIgnore != nil {
		columnsCompInputCnt++
	}
	if *res.UseCommonColumns {
		columnsCompInputCnt++
	}
	if columnsCompInputCnt > 1 {
		return UserInput{}, fmt.Errorf("usecolumns, ignorecolumns, and usecommoncolumns cannot be used together")
	}

	columnsResInputCnt := 0
	if *res.ColumnsToKeep != nil {
		columnsResInputCnt++
	}
	if *res.ColumnsToDelete != nil {
		columnsResInputCnt++
	}
	if columnsResInputCnt > 1 {
		return UserInput{}, fmt.Errorf("keepcolumns and deletecolumns cannot be used together")
	}

	switch *res.Function {
	case FunctionStringCommon:
	case FunctionStringDifferent:
	case "":
		return UserInput{}, fmt.Errorf("function must be given")
	default:
		return UserInput{}, fmt.Errorf("unsupported function %s", *res.Function)
	}

	return res, nil
}

// Adds the index to the row on the right side.
func addIndexToRow(row []csvcheck.StringHashable, index int) []csvcheck.StringHashable {
	newLength := len(row) + 1
	res := make([]csvcheck.StringHashable, newLength)
	res[newLength-1] = csvcheck.BasicStringHashable(fmt.Sprintf("%d", index))
	copy(res[:newLength-1], row)
	return res
}

// Gets the result arrays based off of user input.
func GetResArrays(csvArray1, csvArray2 [][]csvcheck.StringHashable, input UserInput) ([][]csvcheck.StringHashable, [][]csvcheck.StringHashable, error) {
	columnsToUse := csvcheck.GetRowFromRow(*input.ColumnsToUse)
	columnsToIgnore := csvcheck.GetRowFromRow(*input.ColumnsToIgnore)

	var err error = nil
	if *input.UseCommonColumns {
		columnsToUse, err = csvcheck.GetCommonColumns(csvArray1, csvArray2)
		if err != nil {
			return nil, nil, err
		}
	}

	if *input.AutoAlign {
		csvArray1, csvArray2, err = csvcheck.AutoAlignCsvArrays(csvArray1, csvArray2)
		if err != nil {
			return nil, nil, err
		}
	}

	options := csvcheck.Options{
		Method:        MethodMappings[*input.Method],
		UseColumns:    columnsToUse,
		IgnoreColumns: columnsToIgnore,
	}

	if *input.KeepIndex {
		options.SortIndices = true
	}

	var res1 = [][]csvcheck.StringHashable{}
	var res2 = [][]csvcheck.StringHashable{}
	var indices1 = []int{}
	var indices2 = []int{}
	switch *input.Function {
	case FunctionStringCommon:
		res1, res2, indices1, indices2, err = csvcheck.GetCommonRows(csvArray1, csvArray2, options)
	case FunctionStringDifferent:
		res1, res2, indices1, indices2, err = csvcheck.GetDifferentRows(csvArray1, csvArray2, options)
	default:
		return nil, nil, fmt.Errorf("unsupported function")
	}

	if err != nil {
		return nil, nil, err
	}

	if *input.KeepIndex {
		res1[0] = append(res1[0], csvcheck.BasicStringHashable(IndexColumnName))
		for i := 1; i < len(res1); i++ {
			res1[i] = addIndexToRow(res1[i], indices1[i])
		}
		res2[0] = append(res2[0], csvcheck.BasicStringHashable(IndexColumnName))
		for i := 1; i < len(res2); i++ {
			res2[i] = addIndexToRow(res2[i], indices2[i])
		}
	}

	columnsToKeep := csvcheck.GetRowFromRow(*input.ColumnsToKeep)
	columnsToDelete := csvcheck.GetRowFromRow(*input.ColumnsToDelete)
	if columnsToKeep == nil {
		columnsToKeep = append(res1[0], res2[0]...)
	}
	if columnsToDelete == nil {
		columnsToDelete = []csvcheck.StringHashable{}
	}

	res1, err = csvcheck.KeepColumns(res1, columnsToKeep)
	if err != nil {
		return nil, nil, err
	}
	res2, err = csvcheck.KeepColumns(res2, columnsToKeep)
	if err != nil {
		return nil, nil, err
	}

	res1, err = csvcheck.IgnoreColumns(res1, columnsToDelete)
	if err != nil {
		return nil, nil, err
	}
	res2, err = csvcheck.IgnoreColumns(res2, columnsToDelete)
	if err != nil {
		return nil, nil, err
	}

	columnsArrangement1 := csvcheck.GetRowFromRow(*input.ColumnsArrangement1)
	columnsArrangement2 := csvcheck.GetRowFromRow(*input.ColumnsArrangement2)
	if *input.ColumnsArrangement1 != nil {
		res1, err = csvcheck.RearrangeColumns(res1, columnsArrangement1)
		if err != nil {
			return nil, nil, err
		}
	}
	if *input.ColumnsArrangement2 != nil {
		res2, err = csvcheck.RearrangeColumns(res2, columnsArrangement2)
		if err != nil {
			return nil, nil, err
		}
	}

	return res1, res2, nil
}

func ReadCsvFile(filePath string) [][]csvcheck.StringHashable {
	file, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	res := make([][]csvcheck.StringHashable, len(records))
	for i, record := range records {
		res[i] = csvcheck.GetRowFromRow(record)
	}

	return res
}

func WriteString(filePath string, content string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
}
