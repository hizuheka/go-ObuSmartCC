package main

import (
	"fmt"

	"github.com/atotto/clipboard"
)

// ClipboardService はクリップボードへの書き込み機能を提供するインターフェースです。
// テスト時にモックに差し替えることを可能にします。
type ClipboardService interface {
	WriteAll(text string) error
}

// SystemClipboard は実際のシステムクリップボードを使用する実装です。
type SystemClipboard struct{}

// NewSystemClipboard は SystemClipboard の新しいインスタンスを返します。
func NewSystemClipboard() *SystemClipboard {
	return &SystemClipboard{}
}

// WriteAll は指定されたテキストをクリップボードに書き込みます。
func (c *SystemClipboard) WriteAll(text string) error {
	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("クリップボード書き込み失敗: %w", err)
	}
	return nil
}

// MockClipboard はテスト用のインメモリ実装です。
type MockClipboard struct {
	LastContent string
}

func (m *MockClipboard) WriteAll(text string) error {
	m.LastContent = text
	return nil
}
