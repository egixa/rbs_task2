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
	const gigabyte = 1000 * 1000 * 1000
	const megabyte = 1000 * 1000
	const kilobyte = 1000

	if size > gigabyte {
		return fmt.Sprintf("%.2f гб", float64(size)/(gigabyte))
	} else if size > megabyte {
		return fmt.Sprintf("%.2f мб", float64(size)/(megabyte))
	} else if size > kilobyte {
		return fmt.Sprintf("%.2f кб", float64(size)/(kilobyte))
	}
	return fmt.Sprintf("%d б", size)
}

type File struct {
	Type string
	Name string
	Size int64
}

type Files []*File

func dirSize(path string) (int64, error) {

	var size int64
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Ошибка при обходе директории:", err)
			return err
		}
		size += info.Size()
		return nil
	})
	return size, err
}

func main() {
	start := time.Now()

	rootFolder := flag.String("root", "", "Путь до директории для вывода структуры\n")
	sortOptin := flag.String("sort", "asc", "Параметр сортировки:\n по убыванию -\n по возрастанию -")
	flag.Parse()

	if *rootFolder == "" {
		fmt.Println(time.Now().Format("01-02-2006 15:04:05"), "Отсутствуют данные о местоположении директории.")
		fmt.Println("Ожидаемые параметры вызова программы:")
		flag.PrintDefaults()
		return
	}

	if *sortOptin != "asc" && *sortOptin != "desc" {
		*sortOptin = "asc"
		fmt.Println("Введен некорректный параметр сортировки. По умолчанию будет использована сортировка по возрастанию.")
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
			return s[i].Size > s[j].Size
		})
	case "asc":
		sort.Slice(s, func(i, j int) (less bool) {
			return s[i].Size < s[j].Size
		})
	}

	for _, file := range s {
		fmt.Printf("%s %s Размер: %s\n", file.Type, file.Name, formatSize(file.Size))
	}

	duration := time.Since(start)
	fmt.Println(duration)
}
