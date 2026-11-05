package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// readLine читает одну строку с клавиатуры и убирает пробелы и перевод строки по краям.
// Если ввод закончился (Ctrl+D) или чтение сломалось, программа завершается,
// иначе меню крутилось бы бесконечно, получая пустую строку
func readLine(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && line != "" {
			// последняя строка без Enter в конце: отдаём её, следующий вызов завершит программу
			return strings.TrimSpace(line)
		}
		if errors.Is(err, io.EOF) {
			fmt.Println("\nВвод завершён (Ctrl+D). Несохранённые изменения потеряны. Выход.")
		} else {
			fmt.Println("\nОшибка чтения ввода:", err)
		}
		os.Exit(0)
	}
	return strings.TrimSpace(line)
}

// checkText проверяет текстовое поле перед записью в базу:
// оно не должно быть пустым и не должно содержать разделитель "/",
// иначе строка в Base.txt разъедется на лишние поля и студент пропадёт при загрузке
func checkText(s string) error {
	if s == "" {
		return errors.New("поле не может быть пустым")
	}
	if strings.Contains(s, "/") {
		return errors.New("нельзя использовать символ «/», он разделяет поля в файле")
	}
	return nil
}
