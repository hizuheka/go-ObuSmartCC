package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

// App はアプリケーションの依存関係を保持する構造体です。
type App struct {
	clipboard ClipboardService
	logger    *slog.Logger
}

// NewApp はアプリケーションを初期化します。
func NewApp(cb ClipboardService, logger *slog.Logger) *App {
	return &App{
		clipboard: cb,
		logger:    logger,
	}
}

// Run はアプリケーションのメインロジックを実行します。
func (a *App) Run(filePath string) error {
	a.logger.Info("ファイルの読み込みとエンコード判定を開始します", "path", filePath)

	// 【修正】os.Openではなく、文字コード判定・変換機能付きの関数を使用
	// 内部でファイルを読み込むため、ここでファイルのClose処理は不要（Readerが返る）
	utf8Reader, err := LoadFileAsUtf8(filePath)
	if err != nil {
		return fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	// CSVをExcel互換TSVに変換
	// utf8ReaderはすでにUTF-8に変換されているため、csvパッケージで問題なく扱える
	tsvContent, err := CsvToExcelTsv(utf8Reader)
	if err != nil {
		return fmt.Errorf("変換処理に失敗しました: %w", err)
	}

	a.logger.Info("データの変換が完了しました", "bytes", len(tsvContent))

	// クリップボードにコピー
	if err := a.clipboard.WriteAll(tsvContent); err != nil {
		return err
	}

	a.logger.Info("クリップボードへのコピーが完了しました")
	return nil
}

// waitExit はユーザーがEnterキーを押すまで待機します。
// エラー発生時に、ウィンドウが即座に閉じてメッセージが読めなくなるのを防ぎます。
func waitExit() {
	fmt.Print("\n終了するにはEnterキーを押してください...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}

func main() {
	// ロガーの設定
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// CLI引数の解析
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		logger.Error("エラー: CSVファイルが指定されていません。")
		fmt.Println("使用法: CSVファイルをこの実行ファイル(.exe)にドラッグ＆ドロップしてください。")
		waitExit() // 引数エラー時は待機する
		return
	}

	// 複数のファイルがドロップされた場合、最初の1つだけを処理する
	filePath := args[0]

	// 依存関係の構築
	cbService := NewSystemClipboard()
	app := NewApp(cbService, logger)

	// アプリケーション実行
	if err := app.Run(filePath); err != nil {
		logger.Error("処理中にエラーが発生しました", "error", err)
		waitExit() // 実行時エラー時は待機する
		return
	}

	// 成功時は待機せずに即終了する
}
