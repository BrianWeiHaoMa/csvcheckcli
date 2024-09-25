# csvcheckcli
A command-line tool for comparing the rows of different csv files.

## Installation
Check the releases and install the executable directly or clone the repository
and the compile manually with go using
```
git clone https://github.com/BrianWeiHaoMa/csvcheckcli.git
cd ./csvcheckcli
go build
```

## Options
Use ./csvcheckcli -h (or ./csvcheckcli.exe -h depending on what OS you are using) to view the options
```
  -t, --addtimestamp                      Whether or not to add a timestamp to the output file name.
  -a, --autoalign                         Whether or not to auto align the columns of the csv files. Common columns will be aligned on the left side.
  -r, --columnsarrangement1 stringArray   An arrangement for the columns in the first output.
  -R, --columnsarrangement2 stringArray   An arrangement for the columns in the second output.
  -p, --csv                               Whether to print the output in csv format. By default, the output is printed in a columns-aligned.
  -D, --deletecolumns stringArray         The columns to delete in the output.
  -f, --files stringArray                 The input files paths to compare. 2 should be provided.
  -F, --function string                   The function to use for comparison. Options: common, different. A function must be given.
  -i, --ignorecolumns stringArray         The columns to ignore for comparison.
  -d, --inputdir string                   The directory containing the input files. This will be prepended to the input file paths. Must be given.
  -K, --keepcolumns stringArray           The columns to keep in the output.
  -k, --keepindex                         Whether to keep the indices from the original csv of the rows in the result (_ind column will be added).
  -m, --method string                     The method to use for comparison. Options: match, set, direct. By default, set is used. (default "set")
  -o, --outputdir string                  The directory to write the output files to.
  -l, --prettyformatmaxlength int         The maximum length before truncation of a column entry when printing in pretty format. Negative values mean no limit. By default, there is no limit. (default -1)
  -c, --usecolumns stringArray            The columns to use for comparison.
  -C, --usecommoncolumns                  Whether to use all the common columns between the csv files for comparison.
```

## Examples
We will use the input files csv1.csv and csv2.csv for these examples.

csv1.csv
| a  | b  | c  |
|----|----|----|
| 1  | 2  | 3  |
| 4  | 5  | 6  |
| 7  | 8  | 9  |
| 10 | 11 | 12 |

csv2.csv
| a  | b  | c  |
|----|----|----|
| 0  | 0  | 0  |
| 1  | 2  | 3  |
| 1  | 2  | 3  |
| 5  | 5  | 5  |
| 11 | 11 | 11 |

### Example 1:
#### Input:
```
.\csvcheckcli.exe -d .\input_files\ -k -f csv1.csv,csv2.csv -F common -o output_files
```

#### Output:
```
Start time: 2024-09-25 13:51:22

Results for file csv1.csv:
a  b  c  _ind
1  2  3  1

Results for file csv2.csv:
a  b  c  _ind
1  2  3  2
1  2  3  3

Results written to output_files\csvcheck_csv1.csv and output_files\csvcheck_csv2.csv.
```

### Example 2:
#### Input:
```
.\csvcheckcli.exe -d .\input_files\ -k -f csv1.csv,csv2.csv -F different -r c,b,a,_ind -R _ind,c,b,a -p -o output_files
```

#### Output:
```
Start time: 2024-09-25 13:54:42

Results for file csv1.csv:
c,b,a,_ind
6,5,4,2
9,8,7,3
12,11,10,4

Results for file csv2.csv:
_ind,c,b,a
1,0,0,0
4,5,5,5
5,11,11,11

Results written to output_files\csvcheck_csv1.csv and output_files\csvcheck_csv2.csv.
```