package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func getOSPath() string {
	switch runtime.GOOS {
	case "windows":
		root := os.Getenv("SystemRoot")
		if root == "" {
			root = `C:\Windows` // используем raw-строку, чтобы не экранировать слеши
		}
		return filepath.Join(root, "System32", "ntoskrnl.exe")

	case "linux":
		matches, _ := filepath.Glob("/boot/vmlinuz*")
		if len(matches) > 0 {
			return matches[0]
		}
		if _, err := os.Stat("/proc/version"); err == nil {
			return "/proc/version"
		}
		return "unknown (no kernel found in /boot)"

	case "darwin":
		return "/System/Library/Kernels/kernel"

	default:
		return "unknown OS: " + runtime.GOOS
	}
}

// printTree печатает дерево каталогов с найденными файлами
func printTree(path string, depth int) {
	indent := strings.Repeat("  ", depth)
	info, err := os.Stat(path)
	if err != nil {
		return
		
	}
	if info.IsDir() {
		fmt.Println(indent + "📁 " + filepath.Base(path) + "/")
	} else {
		fmt.Println(indent + "📄 " + filepath.Base(path))
	}
}

func main() {
	fmt.Println("=== Finder ===")
	fmt.Println("OS:", runtime.GOOS)
	fmt.Println("OS kernel path:", getOSPath())
	fmt.Println()

	if len(os.Args) < 3 {
		fmt.Println("Использование программы:")
		fmt.Println("  finder <папка> <расширение>")
		fmt.Println("Пример:")
		fmt.Println("  finder ./images .png")
		fmt.Println("  finder C:\\Users .txt")
		return
	}

	folder := os.Args[1]
	wantExt := os.Args[2]

	// Проверяем, что расширение начинается с точки
	if !strings.HasPrefix(wantExt, ".") {
		wantExt = "." + wantExt
	}

	// Проверяем существование папки
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		fmt.Printf(" Папка не найдена: %s\n", folder)
		return
	}

	fmt.Printf("🔎 Поиск файлов с расширением '%s' в '%s'\n\n", wantExt, folder)

	found := 0
	var foundPaths []string

	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // пропускаем недоступные папки
		}
		if !info.IsDir() && strings.EqualFold(filepath.Ext(path), wantExt) {
			foundPaths = append(foundPaths, path)
			found++
		}
		return nil
	})

	if found == 0 {
		fmt.Println("Ничего не найдено.")
		return
	}

	// Выводим список найденных файлов
	fmt.Println("Найденные файлы:")
	for _, p := range foundPaths {
		fmt.Println("  →", p)
	}

	fmt.Printf("\n Найдено файлов: %d\n", found)
}