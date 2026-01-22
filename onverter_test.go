package main

import (
	"strings"
	"testing"
)

func TestCsvToExcelTsv(t *testing.T) {
	// テストケース
	// ユーザー要望: 入力ファイルの内容を「そのまま」Excelに表示する
	// 1. "abc"       -> Excel表示: "abc"
	// 2. ""00001""   -> Excel表示: ""00001""
	// 3. 00001       -> Excel表示: 00001
	// 4. ""12345""   -> Excel表示: ""12345""

	inputCSV := `col1,col2,col3,col4
"abc",""00001"",00001,""12345""`

	r := strings.NewReader(inputCSV)
	output, err := CsvToExcelTsv(r)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
	dataLine := lines[1]

	// 1. "abc"
	// 期待値: Excelで "abc" と表示されること。
	// そのためのTSV形式: =" ""abc"" " -> CSVエスケープ -> "="" ""abc"" """
	// 少なくとも、元の値 "abc" がそのまま含まれている必要がある
	if !strings.Contains(dataLine, `"abc"`) {
		t.Errorf("Value '\"abc\"' lost.\nGot: %s", dataLine)
	}
	// 数式化されていること
	if !strings.Contains(dataLine, `="`) {
		t.Errorf("Not formulated.\nGot: %s", dataLine)
	}

	// 2. ""00001""
	// 期待値: Excelで ""00001"" と表示されること。
	if !strings.Contains(dataLine, `""00001""`) {
		t.Errorf("Value '\"\"00001\"\"' lost.\nGot: %s", dataLine)
	}

	// 3. 00001
	// 期待値: Excelで 00001 と表示されること。(0落ちせず)
	if !strings.Contains(dataLine, `="00001"`) {
		// TSVエスケープにより "=""00001""" になっている可能性があるが、
		// コア部分 ="00001" は共通して含まれるはず
		t.Errorf("Value '00001' logic failed.\nGot: %s", dataLine)
	}

	// 4. ""12345""
	// 期待値: Excelで ""12345"" と表示されること。
	if !strings.Contains(dataLine, `""12345""`) {
		t.Errorf("Value '\"\"12345\"\"' lost.\nGot: %s", dataLine)
	}
}
