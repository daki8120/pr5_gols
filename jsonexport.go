// Файл jsonexport.go — вывод информации о файлах в формате JSON (ЗАДАНИЕ 3, доп.)
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONFileEntry — структура одной записи для JSON-вывода
type JSONFileEntry struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	IsExec  bool   `json:"is_executable"`
	Mode    string `json:"mode"`
	ModTime string `json:"mod_time"`
}

// printJSON выводит список файлов в формате JSON
func printJSON(items []FileItem) {
	var entries []JSONFileEntry

	for _, item := range items {
		entries = append(entries, JSONFileEntry{
			Name:    item.Name,
			Path:    item.Path,
			Size:    item.Info.Size(),
			IsDir:   item.IsDir,
			IsExec:  item.IsExec,
			Mode:    formatMode(item.Info.Mode()),
			ModTime: item.Info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка сериализации JSON: %v\n", err)
		return
	}

	fmt.Println(string(data))
}
