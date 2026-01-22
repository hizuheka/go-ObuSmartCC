package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
)

// zeroPaddingRegex は「0から始まる数字列」にマッチします。
// 例: "01", "000123", "0"
var zeroPaddingRegex = regexp.MustCompile(`^0\d*$`)

// CsvToExcelTsv はCSVデータを読み込み、Excel貼り付け用に最適化されたTSV形式に変換します。
// 0落ちを防ぐため、0から始まる数字列はExcelの数式形式(="001")に変換します。
func CsvToExcelTsv(r io.Reader) (string, error) {
	csvReader := csv.NewReader(r)

	// 【修正点】: -1を設定することで、行ごとに列数が異なっていてもエラーにしません。
	csvReader.FieldsPerRecord = -1

	// バッファを作成し、TSV書き込み用のWriterを用意
	var buf bytes.Buffer
	tsvWriter := csv.NewWriter(&buf)
	tsvWriter.Comma = '\t' // 区切り文字をタブに設定

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("CSV読み込みエラー: %w", err)
		}

		// レコード内の各フィールドを加工
		processedRecord := make([]string, len(record))
		for i, field := range record {
			processedRecord[i] = escapeForExcel(field)
		}

		if err := tsvWriter.Write(processedRecord); err != nil {
			return "", fmt.Errorf("TSV書き込みエラー: %w", err)
		}
	}

	tsvWriter.Flush()
	if err := tsvWriter.Error(); err != nil {
		return "", fmt.Errorf("バッファフラッシュエラー: %w", err)
	}

	return buf.String(), nil
}

// escapeForExcel はフィールドの値を確認し、Excelで0落ちや意図しない変換を防ぐ形式に変換します。
func escapeForExcel(field string) string {
	// 0から始まる数字列（"0123", "007"など）の場合
	// Excelの数式記法 ="VALUE" に変換することで、文字列として強制的に認識させる
	if zeroPaddingRegex.MatchString(field) {
		return fmt.Sprintf(`="%s"`, field)
	}

	return field
}
