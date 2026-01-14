package main

import (
	"strings"
	"testing"
)

func TestCsvToExcelTsv(t *testing.T) {
	// テストケース
	// 1. 通常のフィールド
	// 2. 0埋め数値 ("00123") -> ="00123" に変換されるべき
	// 3. カンマを含むフィールド -> タブ区切りにおいて引用符で囲まれるべき（csv.Writerが処理）
	// 4. 数値の0 ("0") -> 文字列として扱っても良いが、今回は ="0" になる想定
	inputCSV := `id,name,code
1,Alice,00100
2,"Bob, Jr",050
3,Charlie,0
`

	r := strings.NewReader(inputCSV)
	output, err := CsvToExcelTsv(r)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// 行ごとに分割して検証
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 4 { // ヘッダ + 3データ行
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	// ヘッダ行の検証
	expectedHeader := "id\tname\tcode"
	if lines[0] != expectedHeader {
		t.Errorf("Header mismatch.\nExpected: %q\nGot:      %q", expectedHeader, lines[0])
	}

	// データ行1: 1, Alice, 00100 -> codeが変換されているか
	// encoding/csv の仕様により、フィールド内のダブルクォートはエスケープされます。
	// 元の値: ="00100"
	// TSV出力: "=""00100"""
	// Excelはこの形式を正しく解釈します。
	if !strings.Contains(lines[1], `"=""00100"""`) {
		t.Errorf("Zero padding logic failed for line 1.\nGot: %s", lines[1])
	}

	// データ行2: 2, "Bob, Jr", 050
	// 元の値: ="050"
	// TSV出力: "=""050"""
	if !strings.Contains(lines[2], `"=""050"""`) {
		t.Errorf("Zero padding logic failed for line 2.\nGot: %s", lines[2])
	}
	// カンマを含むフィールドの確認
	if !strings.Contains(lines[2], "Bob, Jr") {
		t.Errorf("Comma handling failed for line 2.\nGot: %s", lines[2])
	}

	// データ行3: 3, Charlie, 0 -> 0単体も ="0" に変換されるロジックになっている
	// 元の値: ="0"
	// TSV出力: "=""0"""
	if !strings.Contains(lines[3], `"=""0"""`) {
		t.Errorf("Zero handling failed for line 3.\nGot: %s", lines[3])
	}
}
