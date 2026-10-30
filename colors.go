// Файл colors.go — ANSI-цвета и иконки для файлов
package main

import (
	"path/filepath"
	"strings"
)

const (
	Reset = "\033[0m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BrightRed   = "\033[91m"
	BrightGreen = "\033[92m"
	BrightBlue  = "\033[94m"
	BrightCyan  = "\033[96m"
	BrightWhite = "\033[97m"

	Bold = "\033[1m"

	BgBlue = "\033[44m"
)

// FileColor возвращает ANSI-цвет и иконку для файла
func FileColor(name string, isDir bool, isExec bool) (color, icon string) {
	if isDir {
		return Bold + BrightBlue, "📁"
	}

	ext := strings.ToLower(filepath.Ext(name))

	switch ext {
	case ".go":
		return BrightCyan, "🐹"
	case ".py":
		return Yellow, "🐍"
	case ".js", ".ts", ".jsx", ".tsx":
		return Yellow, "📜"
	case ".java", ".class":
		return BrightRed, "☕"
	case ".c", ".cpp", ".h", ".hpp":
		return Blue, "⚙️"
	case ".rs":
		return BrightRed, "🦀"
	case ".dart":
		return BrightBlue, "🎯"

	case ".json":
		return Yellow, "📋"
	case ".yaml", ".yml":
		return Green, "📋"
	case ".toml":
		return Green, "📋"
	case ".xml":
		return Cyan, "📋"
	case ".csv":
		return Green, "📊"
	case ".sql":
		return Blue, "🗄️"
	case ".env":
		return Yellow, "🔑"

	case ".md", ".markdown":
		return White, "📝"
	case ".txt":
		return White, "📄"
	case ".pdf":
		return BrightRed, "📕"
	case ".doc", ".docx":
		return Blue, "📘"

	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp":
		return Magenta, "🖼️"

	case ".mp4", ".avi", ".mkv", ".mov":
		return BrightRed, "🎬"
	case ".mp3", ".wav", ".flac", ".ogg":
		return Magenta, "🎵"

	case ".zip", ".tar", ".gz", ".rar", ".7z", ".bz2":
		return BrightRed, "📦"

	case ".exe", ".msi", ".bat", ".cmd", ".ps1", ".sh":
		return BrightGreen, "⚡"
	}

	if isExec {
		return BrightGreen, "⚡"
	}

	return White, "📄"
}
