// Файл main.go — точка входа утилиты gols
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ====================================================
// ЧТЕНИЕ ДИРЕКТОРИИ
// ====================================================

func readDir(path string, cfg Config) ([]FileItem, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать '%s': %w", path, err)
	}

	var items []FileItem
	for _, entry := range entries {
		name := entry.Name()

		isHidden := strings.HasPrefix(name, ".")
		if isHidden && !cfg.All {
			continue
		}

		// ЗАДАНИЕ 2: применяем фильтр по паттерну, если задан
		if cfg.Filter != "" {
			matched, err := filepath.Match(cfg.Filter, name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "gols: некорректный паттерн '%s': %v\n", cfg.Filter, err)
			} else if !matched {
				continue
			}
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		isExec := !entry.IsDir() && (info.Mode()&0111 != 0)

		items = append(items, FileItem{
			Name:     name,
			Path:     filepath.Join(path, name),
			Info:     info,
			IsDir:    entry.IsDir(),
			IsExec:   isExec,
			IsHidden: isHidden,
		})
	}

	return items, nil
}

// ====================================================
// СОРТИРОВКА
// ====================================================

func sortItems(items []FileItem, cfg Config) {
	sort.Slice(items, func(i, j int) bool {
		a, b := items[i], items[j]

		// Директории всегда перед файлами
		if a.IsDir != b.IsDir {
			return a.IsDir
		}

		var less bool

		switch {
		case cfg.SortTime:
			less = a.Info.ModTime().After(b.Info.ModTime())
		case cfg.SortSize:
			less = a.Info.Size() > b.Info.Size()
		default:
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		}

		if cfg.Reverse {
			return !less
		}
		return less
	})
}

// ====================================================
// РЕКУРСИВНЫЙ ВЫВОД
// ====================================================

func listRecursive(path string, cfg Config, depth int) {
	if depth > 0 {
		fmt.Printf("\n%s:\n", path)
	}

	items, err := readDir(path, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gols: %v\n", err)
		return
	}

	sortItems(items, cfg)

	// ЗАДАНИЕ 3: если указан флаг --json, выводим JSON и выходим
	if cfg.JSON {
		printJSON(items)
	} else if cfg.Long {
		printLong(items, cfg)
	} else {
		printShort(items, cfg)
	}

	if cfg.Recursive {
		for _, item := range items {
			if item.IsDir {
				listRecursive(item.Path, cfg, depth+1)
			}
		}
	}
}

// ====================================================
// ГЛАВНАЯ ФУНКЦИЯ
// ====================================================

func main() {
	longFlag := flag.Bool("l", false, "подробный вывод (права, размер, дата)")
	allFlag := flag.Bool("a", false, "показывать скрытые файлы (начинаются с .)")
	humanFlag := flag.Bool("h", false, "размеры в удобочитаемом формате (K, M, G)")
	sortTimeFlag := flag.Bool("t", false, "сортировать по времени изменения")
	sortSizeFlag := flag.Bool("S", false, "сортировать по размеру (большие первые)")
	reverseFlag := flag.Bool("r", false, "обратный порядок сортировки")
	recursiveFlag := flag.Bool("R", false, "рекурсивный вывод поддиректорий")
	colorFlag := flag.Bool("color", true, "цветной вывод по типу файла")
	iconsFlag := flag.Bool("icons", true, "показывать иконки рядом с именами")

	// ЗАДАНИЕ 1: флаг -1 — одна колонка
	oneColumnFlag := flag.Bool("1", false, "выводить по одному файлу на строку")

	// ЗАДАНИЕ 2: флаг --filter
	filterFlag := flag.String("filter", "", "показывать только файлы по шаблону (например: *.go)")

	// ЗАДАНИЕ 3 (доп.): флаг --json
	jsonFlag := flag.Bool("json", false, "вывести информацию о файлах в формате JSON")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование: gols [ФЛАГИ] [ПУТЬ...]\n\n")
		fmt.Fprintf(os.Stderr, "Флаги:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nПримеры:\n")
		fmt.Fprintf(os.Stderr, "  gols                              # текущая директория\n")
		fmt.Fprintf(os.Stderr, "  gols -l -h C:\\Users               # подробно, с размерами\n")
		fmt.Fprintf(os.Stderr, "  gols -l -a .                      # все файлы включая скрытые\n")
		fmt.Fprintf(os.Stderr, "  gols -R --color C:\\go_projects    # рекурсивно с цветом\n")
		fmt.Fprintf(os.Stderr, "  gols -1 .                         # по одному файлу на строку\n")
		fmt.Fprintf(os.Stderr, "  gols --filter \"*.go\" .            # только .go файлы\n")
		fmt.Fprintf(os.Stderr, "  gols --json .                     # вывод в формате JSON\n")
	}

	flag.Parse()

	cfg := Config{
		Long:      *longFlag,
		All:       *allFlag,
		Human:     *humanFlag,
		SortTime:  *sortTimeFlag,
		SortSize:  *sortSizeFlag,
		Reverse:   *reverseFlag,
		Recursive: *recursiveFlag,
		Color:     *colorFlag,
		Icons:     *iconsFlag,
		OneColumn: *oneColumnFlag,
		Filter:    *filterFlag,
		JSON:      *jsonFlag,
	}

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for i, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gols: '%s': %v\n", path, err)
			continue
		}

		if len(paths) > 1 {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", path)
		}

		if !info.IsDir() {
			item := FileItem{
				Name:  filepath.Base(path),
				Path:  path,
				Info:  info,
				IsDir: false,
			}
			if cfg.JSON {
				printJSON([]FileItem{item})
			} else if cfg.Long {
				printLong([]FileItem{item}, cfg)
			} else {
				printShort([]FileItem{item}, cfg)
			}
		} else {
			listRecursive(path, cfg, 0)
		}
	}
}
