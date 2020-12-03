# xlsx2csv
Simple io.Reader-compatible library which converts an XLSX sheet to CSV row by row

## Example of usage
``` go
import (
    ...
)

func main() {
    file, err := os.Open(path)
    if err != nil {
        log.Fatalln(err)
    }
    defer file.Close()

    reader, err := xlsx2csv.NewReader(file, xlsx2csv.WithName("sheet"), ',')
    if err != nil {
        log.Fatalln(err)
    }

    rawCSV, err := ioutil.ReadAll(reader)
    if err != nil {
        log.Fatalln(err)
    }

    log.Println(string(rawCSV))
}
```