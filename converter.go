package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// CsvToExcelTsv はCSVデータを読み込み、Excel貼り付け用に最適化されたTSV形式に変換します。
func CsvToExcelTsv(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)

	// バッファを作成し、TSV書き込み用のWriterを用意
	var buf bytes.Buffer
	tsvWriter := csv.NewWriter(&buf)
	tsvWriter.Comma = '\t' // 区切り文字をタブに設定
	tsvWriter.UseCRLF = true

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// 空行はスキップ
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 自前のロジックでCSV行をパース（読み込んだままの形を維持するため）
		record := parseDirtyCSVLine(line)

		processedRecord := make([]string, len(record))
		for i, field := range record {
			processedRecord[i] = escapeForExcel(field)
		}

		if err := tsvWriter.Write(processedRecord); err != nil {
			return "", fmt.Errorf("TSV書き込みエラー(行%d): %w", lineCount, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	tsvWriter.Flush()
	if err := tsvWriter.Error(); err != nil {
		return "", fmt.Errorf("バッファフラッシュエラー: %w", err)
	}

	return buf.String(), nil
}

// parseDirtyCSVLine は行儀の悪いCSV行を柔軟にパースします。
func parseDirtyCSVLine(line string) []string {
	parts := strings.Split(line, ",")
	var fields []string
	var currentField strings.Builder
	inQuote := false

	for i, part := range parts {
		if inQuote {
			currentField.WriteString(",")
		}
		currentField.WriteString(part)

		quoteCount := strings.Count(part, `"`)
		if quoteCount%2 != 0 {
			inQuote = !inQuote
		}

		if !inQuote || i == len(parts)-1 {
			fields = append(fields, currentField.String())
			currentField.Reset()
			inQuote = false
		}
	}
	return fields
}

// escapeForExcel はフィールドの値をそのままExcelに表示できるよう変換します。
func escapeForExcel(field string) string {
	// 方針: 入力ファイルから読み込んだ内容を「そのまま」Excelに表示する。
	// そのため、クォート除去や条件分岐は一切行わず、
	// 全ての値を無条件で Excelの数式形式(="...") として出力する。

	// Excel数式内でダブルクォートを正しく表示するためにエスケープ (" -> "")
	escapedField := strings.ReplaceAll(field, `"`, `""`)

	// 全てを数式化して返す
	return fmt.Sprintf(`="%s"`, escapedField)
}
