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

	var buf bytes.Buffer
	tsvWriter := csv.NewWriter(&buf)
	tsvWriter.Comma = '\t'
	tsvWriter.UseCRLF = true

	// 状態管理用変数
	var currentSection string
	splitTargetIndex := -1 // 分割対象の列インデックス (-1は対象なし)

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if strings.TrimSpace(line) == "" {
			continue
		}

		// 1. 行をパース
		record := parseDirtyCSVLine(line)
		if len(record) == 0 {
			continue
		}

		// 2. セクション開始の判定
		firstCol := strings.ToUpper(strings.TrimSpace(record[0]))

		if firstCol == "NET" || firstCol == "JOB" {
			currentSection = firstCol
			splitTargetIndex = -1 // セクションが変わったらリセット
		}

		isHeaderLine := false

		// 3. JOBセクションにおいて、分割対象列(jobname_jes)の位置を特定する
		if currentSection == "JOB" && splitTargetIndex == -1 {
			for i, colName := range record {
				if strings.EqualFold(colName, "jobname_jes") {
					splitTargetIndex = i
					isHeaderLine = true // 見つかったこの行こそがヘッダ行である
					break
				}
			}
		}

		// 4. 列分割・置換処理
		if currentSection == "JOB" && splitTargetIndex != -1 && len(record) > splitTargetIndex {
			if isHeaderLine {
				// ヘッダ行の場合: 列名を指定のもの(jobname_jes1, jobname_jes2)に強制置換
				record = expandColumn(record, splitTargetIndex, "jobname_jes1", "jobname_jes2")
			} else {
				// データ行の場合: 値を "_" で分割して展開
				targetValue := record[splitTargetIndex]
				parts := strings.SplitN(targetValue, "_", 2)
				val1 := parts[0]
				val2 := ""

				// 分割が行われた場合（_が含まれていた場合）、1列目の末尾に "_" を残す
				// 例: "01_zzzz" -> parts=["01", "zzzz"] -> val1="01_", val2="zzzz"
				if len(parts) > 1 {
					val1 += "_"
					val2 = parts[1]
				}

				record = expandColumn(record, splitTargetIndex, val1, val2)
			}
		}

		// 5. Excel形式への変換と書き出し
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

// expandColumn はレコード内の指定インデックスの列を削除し、代わりに2つの値を挿入して列を拡張します。
func expandColumn(record []string, index int, val1, val2 string) []string {
	newRecord := make([]string, 0, len(record)+1)
	newRecord = append(newRecord, record[:index]...)
	newRecord = append(newRecord, val1, val2)
	if index+1 < len(record) {
		newRecord = append(newRecord, record[index+1:]...)
	}
	return newRecord
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
	escapedField := strings.ReplaceAll(field, `"`, `""`)
	return fmt.Sprintf(`="%s"`, escapedField)
}
