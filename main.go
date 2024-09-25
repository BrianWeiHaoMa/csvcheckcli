package main

import (
	"csvcheckcli/csvcheckcli"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/BrianWeiHaoMa/csvcheck"
)

func main() {
	input, err := csvcheckcli.ParseUserInput(nil)
	if err != nil {
		log.Fatalf("error parsing input:\n%s", err)
	}

	csvPath1 := filepath.Join(*input.InputDir, (*input.Files)[0])
	csvPath2 := filepath.Join(*input.InputDir, (*input.Files)[1])

	csvArray1 := csvcheckcli.ReadCsvFile(csvPath1)
	csvArray2 := csvcheckcli.ReadCsvFile(csvPath2)

	fileName1 := filepath.Base(csvPath1)
	fileName2 := filepath.Base(csvPath2)

	currentTime := time.Now()
	fmt.Printf("Start time: %s\n\n", currentTime.Format("2006-01-02 15:04:05"))

	res1, res2, err := csvcheckcli.GetResArrays(csvArray1, csvArray2, input)
	if err != nil {
		log.Fatal(err)
	}

	resString1, _ := csvcheck.StringFormatCsvArray(res1)
	resString2, _ := csvcheck.StringFormatCsvArray(res2)
	if *input.PrintInCsvFormat {
		fmt.Printf("Results for file %s:\n%s\n", fileName1, resString1)
		fmt.Printf("Results for file %s:\n%s\n", fileName2, resString2)
	} else {
		prettyResString1, _ := csvcheck.PrettyFormatCsvArray(res1, 2, *input.PrettyFormatMaxLength)
		prettyResString2, _ := csvcheck.PrettyFormatCsvArray(res2, 2, *input.PrettyFormatMaxLength)
		fmt.Printf("Results for file %s:\n%s\n", fileName1, prettyResString1)
		fmt.Printf("Results for file %s:\n%s\n", fileName2, prettyResString2)
	}

	if *input.OutputDir != "" {
		fileNameNoExt1 := fileName1[:len(fileName1)-len(filepath.Ext(fileName1))]
		fileNameNoExt2 := fileName2[:len(fileName2)-len(filepath.Ext(fileName2))]

		resFileName1 := fmt.Sprintf("csvcheck_%s.csv", fileNameNoExt1)
		resFileName2 := fmt.Sprintf("csvcheck_%s.csv", fileNameNoExt2)
		if *input.AddTimestamp {
			timeString := currentTime.Format("2006_01_02_15_04_05")
			resFileName1 = fmt.Sprintf("csvcheck_%s_%s.csv", fileNameNoExt1, timeString)
			resFileName2 = fmt.Sprintf("csvcheck_%s_%s.csv", fileNameNoExt2, timeString)
		}

		outputPath1 := filepath.Join(*input.OutputDir, resFileName1)
		outputPath2 := filepath.Join(*input.OutputDir, resFileName2)

		csvcheckcli.WriteString(outputPath1, resString1)
		csvcheckcli.WriteString(outputPath2, resString2)

		fmt.Printf("Results written to %s and %s.\n", outputPath1, outputPath2)
	}
}
