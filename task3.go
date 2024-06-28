package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"time"
)

func formatSize(size int64) string {
	if size > 1<<30 {
		return fmt.Sprintf("%.2f гб", float64(size)/(1<<30))
	} else if size > 1<<20 {
		return fmt.Sprintf("%.2f мб", float64(size)/(1<<20))
	} else if size > 1<<10 {
		return fmt.Sprintf("%.2f кб", float64(size)/(1<<10))
	}
	return fmt.Sprintf("%d б", size)
}

type File struct {
	Type string
	Name string
	Size int64
}

type Files []*File

func (f Files) Len() int      { return len(f) }
func (f Files) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

type sortAsk struct{ Files }

func (s sortAsk) Less(i, j int) bool {
	return s.Files[i].Size < s.Files[j].Size
}

type sortDesc struct{ Files }

func (s sortDesc) Less(i, j int) bool {
	return s.Files[i].Size > s.Files[j].Size
}

func dirSize(path string) (int64, error) {

	var size int64
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Ошибка при обходе директории:", err)
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func main() {
	start := time.Now()

	rootFolder := flag.String("root", "", "Путь до директории для вывода структуры\n")
	sortOptin := flag.String("sort", "desc", "Параметр сортировки:\n по убыванию -\n по возрастанию -")
	flag.Parse()

	if *rootFolder == "" {
		fmt.Println(time.Now().Format("01-02-2006 15:04:05"), "Отсутствуют данные о местоположении директории.")
		fmt.Println("Ожидаемые параметры вызова программы:")
		flag.PrintDefaults()
		return
	}

	_, err := os.Stat(*rootFolder)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Ошибка при обнаружении директории:", err)
			return
		}
	}

	// Открываем директорию
	dir, err := os.Open(*rootFolder)
	if err != nil {
		fmt.Println("Ошибка при открытии директории", err)
		return
	}
	defer dir.Close()

	// Получаем список файлов и директорий
	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Ошибка при прочтении директории", err)
		return
	}

	s := []File{}

	// Выводим имена файлов и директорий
	for _, file := range files {
		filename := file.Name()
		fileSize := file.Size()
		filePath := filepath.Join(*rootFolder, filename)

		if file.IsDir() {
			dirSize, err := dirSize(filePath)
			if err != nil {
				return
			}
			s = append(s, File{"Директория", filename, dirSize})
			continue
		} else {
			s = append(s, File{"Файл", filename, fileSize})
			continue
		}
	}

	switch *sortOptin {
	case "desc":
		sort.Slice(s, func(i, j int) (less bool) {
			return s[i].Name < s[j].Name
		})
	case "asc":
		sort.Slice(s, func(i, j int) (less bool) {
			return s[i].Name < s[j].Name
		})
	}
	for _, file := range s {
		fmt.Printf("%s %s Размер: %s\n", file.Type, file.Name, formatSize(file.Size))
	}

	duration := time.Since(start)
	fmt.Println(duration)
}
