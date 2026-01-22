package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// LoadFileAsUtf8 は指定されたパスのファイルを読み込み、
// 文字コードを自動判定してUTF-8に変換したReaderを返します。
// 巨大なファイルを扱わない前提のため、一度メモリに読み込んで判定を行います。
func LoadFileAsUtf8(filePath string) (io.Reader, error) {
	// ファイル全体をバイト列として読み込む
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込み失敗: %w", err)
	}

	// 文字コード自動判定
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(content)
	if err != nil {
		// 判定に失敗した場合は、変換せずにそのまま返す（UTF-8と仮定）
		return bytes.NewReader(content), nil
	}

	// 判定結果に基づいて処理を分岐
	// Shift_JISの場合のみ変換処理を行う
	if result.Charset == "Shift_JIS" {
		// Shift-JISからUTF-8への変換デコーダーを作成
		decoder := japanese.ShiftJIS.NewDecoder()

		// バイト列を変換
		utf8Content, _, err := transform.Bytes(decoder, content)
		if err != nil {
			return nil, fmt.Errorf("Shift-JISからの変換に失敗しました: %w", err)
		}

		return bytes.NewReader(utf8Content), nil
	}

	// EUC-JPなども必要であればここに追加可能ですが、
	// それ以外（UTF-8やASCII）はそのまま返す
	return bytes.NewReader(content), nil
}
