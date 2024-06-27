package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"time"
)

/*
	func dirWalk(files []fs.FileInfo) {
		for _, file := range files {
			if file.IsDir() == true {
				filepath.Walk("./"+file.Name(), func(wPath string, info os.FileInfo, err error) error {

					return
				})

				continue
			} else {
				fmt.Println("Файл -", file.Name(), "Размер -", file.Size())
				continue
			}

		}
	}
func sortFolder(sortOptin *string, files []) []{

}*/

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

	// Выводим имена файлов и директорий
	for _, file := range files {
		filePath := filepath.Join(*rootFolder, file.Name())

		if file.IsDir() {
			dirSize, err := dirSize(filePath)
			if err != nil {
				return
			}
			fmt.Println("Директория -", file.Name(), "Размер -", dirSize)
			continue
		} else {
			fmt.Println("Файл -", file.Name(), "Размер -", file.Size())
			continue
		}

	}

	duration := time.Since(start)
	fmt.Println(duration, sortOptin)
}
