package main

import (
	"strings"
	"testing"
)

func TestCsvToExcelTsv(t *testing.T) {
	// テストケース
	// 1. セクションマーカーとヘッダが別行の場合でも正しく動作するか
	// 2. ヘッダ行の jobname_jes が jobname_jes1, jobname_jes2 に置換されるか
	// 3. データ行の分割が正しく行われ、1列目に "_" が残るか

	inputCSV := `JOB
jobname,jobnumber,jobparam,jobname_jes,val1
"C999.BAT","1","AUTO","01_C299_TITLE","val"`

	r := strings.NewReader(inputCSV)
	output, err := CsvToExcelTsv(r)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// Line 1: JOB (セクションマーカー)
	if !strings.Contains(lines[0], `="JOB"`) {
		t.Errorf("Section marker corrupted.\nGot: %s", lines[0])
	}

	// Line 2: ヘッダ行
	// jobname_jes1, jobname_jes2 に置換されているか
	if !strings.Contains(lines[1], `="jobname_jes1"`) || !strings.Contains(lines[1], `="jobname_jes2"`) {
		t.Errorf("Header replacement failed.\nGot: %s", lines[1])
	}

	// Line 3: データ行
	// "01_C299_TITLE" -> "01_", "C299_TITLE"
	// 要望通り、1列目にアンダースコアが残っていることを確認
	if !strings.Contains(lines[2], `="01_"`) {
		t.Errorf("Data split part 1 failed (underscore check).\nGot: %s", lines[2])
	}
	if !strings.Contains(lines[2], `="C299_TITLE"`) {
		t.Errorf("Data split part 2 failed.\nGot: %s", lines[2])
	}
}
