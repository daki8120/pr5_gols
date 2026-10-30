// Файл formatter.go — форматирование вывода файлов
package main

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
)

// ====================================================
// СТРУКТУРА: один файл
// ====================================================

type FileItem struct {
	Name     string
	Path     string
	Info     fs.FileInfo
	IsDir    bool
	IsExec   bool
	IsHidden bool
}

// ====================================================
// ФОРМАТИРОВАНИЕ РАЗМЕРА
// ====================================================

func humanSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%6.1fG", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%6.1fM", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%6.1fK", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%7d", bytes)
	}
}

// ====================================================
// ФОРМАТИРОВАНИЕ ПРАВ ДОСТУПА
// ====================================================

func formatMode(mode os.FileMode) string {
	var b strings.Builder

	if mode.IsDir() {
		b.WriteByte('d')
	} else if mode&os.ModeSymlink != 0 {
		b.WriteByte('l')
	} else {
		b.WriteByte('-')
	}

	b.WriteByte(bitChar(mode, 0400, 'r'))
	b.WriteByte(bitChar(mode, 0200, 'w'))
	b.WriteByte(bitChar(mode, 0100, 'x'))

	b.WriteByte(bitChar(mode, 0040, 'r'))
	b.WriteByte(bitChar(mode, 0020, 'w'))
	b.WriteByte(bitChar(mode, 0010, 'x'))

	b.WriteByte(bitChar(mode, 0004, 'r'))
	b.WriteByte(bitChar(mode, 0002, 'w'))
	b.WriteByte(bitChar(mode, 0001, 'x'))

	return b.String()
}

func bitChar(mode os.FileMode, bit os.FileMode, char byte) byte {
	if mode&bit != 0 {
		return char
	}
	return '-'
}

// ====================================================
// ФОРМАТИРОВАНИЕ ДАТЫ
// ====================================================

func formatTime(t time.Time) string {
	now := time.Now()
	if t.Year() == now.Year() {
		return t.Format("Jan _2 15:04")
	}
	return t.Format("Jan _2  2006")
}

// ====================================================
// ОПРЕДЕЛЕНИЕ ШИРИНЫ ТЕРМИНАЛА (ЗАДАНИЕ 4, доп.)
// ====================================================

// getTerminalWidth пытается определить ширину терминала
// через переменную окружения COLUMNS, иначе возвращает значение по умолчанию
func getTerminalWidth() int {
	const defaultWidth = 80

	if cols := os.Getenv("COLUMNS"); cols != "" {
		if w, err := strconv.Atoi(cols); err == nil && w > 0 {
			return w
		}
	}

	return defaultWidth
}

// ====================================================
// КОНФИГУРАЦИЯ ВЫВОДА
// ====================================================

type Config struct {
	Long      bool
	All       bool
	Human     bool
	SortTime  bool
	SortSize  bool
	Reverse   bool
	Recursive bool
	Color     bool
	Icons     bool
	OneColumn bool   // ЗАДАНИЕ 1: флаг -1 — по одному файлу на строку
	Filter    string // ЗАДАНИЕ 2: флаг --filter
	JSON      bool   // ЗАДАНИЕ 3 (доп.): флаг --json
}

// ====================================================
// ВЫВОД: КОРОТКИЙ ФОРМАТ (МНОГО КОЛОНОК)
// ====================================================

func printShort(items []FileItem, cfg Config) {
	if len(items) == 0 {
		return
	}

	// ЗАДАНИЕ 1: режим -1 — каждый файл на отдельной строке
	if cfg.OneColumn {
		for _, item := range items {
			name := item.Name
			if item.IsDir {
				name += "/"
			}
			if cfg.Color {
				color, icon := FileColor(item.Name, item.IsDir, item.IsExec)
				if cfg.Icons {
					fmt.Printf("%s%s %s%s\n", color, icon, name, Reset)
				} else {
					fmt.Printf("%s%s%s\n", color, name, Reset)
				}
			} else {
				fmt.Println(name)
			}
		}
		return
	}

	maxWidth := 0
	for _, item := range items {
		name := item.Name
		if item.IsDir {
			name += "/"
		}
		if len(name) > maxWidth {
			maxWidth = len(name)
		}
	}

	// ЗАДАНИЕ 4: используем определённую ширину терминала
	termWidth := getTerminalWidth()
	colWidth := maxWidth + 2
	cols := termWidth / colWidth
	if cols < 1 {
		cols = 1
	}

	for i, item := range items {
		name := item.Name
		if item.IsDir {
			name += "/"
		}

		if cfg.Color {
			color, icon := FileColor(item.Name, item.IsDir, item.IsExec)
			if cfg.Icons {
				fmt.Printf("%s%s %-*s%s", color, icon, colWidth-2, name, Reset)
			} else {
				fmt.Printf("%s%-*s%s", color, colWidth, name, Reset)
			}
		} else {
			fmt.Printf("%-*s", colWidth, name)
		}

		if (i+1)%cols == 0 || i == len(items)-1 {
			fmt.Println()
		}
	}
}

// ====================================================
// ВЫВОД: ПОДРОБНЫЙ ФОРМАТ
// ====================================================

func printLong(items []FileItem, cfg Config) {
	if len(items) == 0 {
		return
	}

	var totalBlocks int64
	for _, item := range items {
		totalBlocks += item.Info.Size() / 512
	}
	fmt.Printf("итого %d\n", totalBlocks)

	for _, item := range items {
		mode := formatMode(item.Info.Mode())

		var sizeStr string
		if cfg.Human {
			sizeStr = humanSize(item.Info.Size())
		} else {
			sizeStr = fmt.Sprintf("%7d", item.Info.Size())
		}

		modTime := formatTime(item.Info.ModTime())

		name := item.Name
		if item.IsDir {
			name += "/"
		}

		if cfg.Color {
			color, icon := FileColor(item.Name, item.IsDir, item.IsExec)
			if cfg.Icons {
				fmt.Printf("%s 1 %-8s %s %s %s%s %s%s\n",
					mode, "user", sizeStr, modTime,
					color, icon, name, Reset)
			} else {
				fmt.Printf("%s 1 %-8s %s %s %s%s%s\n",
					mode, "user", sizeStr, modTime,
					color, name, Reset)
			}
		} else {
			fmt.Printf("%s 1 %-8s %s %s %s\n",
				mode, "user", sizeStr, modTime, name)
		}
	}
}
