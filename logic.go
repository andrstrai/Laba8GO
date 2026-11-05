package main

import (
	"bufio"
	"fmt"
	"strconv"
)

func StudentByID(students_list []students, id int) (*students, error) {
	if id <= 0 {
		return nil, fmt.Errorf("Нет студента с ID = %d", id)
	}
	for i := range students_list {
		if students_list[i].ID == uint(id) {
			return &students_list[i], nil
		}
	}
	return nil, fmt.Errorf("Нет студента с ID = %d", id)
}

func ChangeStudent(student *students, reader *bufio.Reader) {
	for {
		clear_screen()
		fmt.Println("Текущий студент:")
		fmt.Println(student)
		fmt.Println("\nВведите номер параметра, который хотите изменить:")
		fmt.Println("1 - имя")
		fmt.Println("2 - дата рождения")
		fmt.Println("3 - институт")
		fmt.Println("4 - стипендия")
		fmt.Println("5 - средний балл")
		fmt.Println("0 - завершить изменения")
		input := readLine(reader)

		num_menu, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Некорректный ввод! Пожалуйста, введите число.")
			continue
		}
		switch num_menu {
		case 0:
			fmt.Println("Изменения внесены. Не забудьте сохранить их в файл (пункт 5).")
			return
		case 1:
			fmt.Print("Укажите имя студента: ")
			name := readLine(reader)
			if err := checkText(name); err != nil {
				fmt.Println("Ошибка:", err)
				fmt.Println("Для продолжения работы пожалуйста, нажмите Enter...")
				readLine(reader)
				clear_screen()
				continue
			}
			student.name = name
		case 2:
			fmt.Print("Укажите дату рождения студента в формате ДД.ММ.ГГГГ: ")
			data := readLine(reader)
			if !is_valid_date(data) {
				fmt.Println("Ошибка ввода даты. Введите еще раз (ДД.ММ.ГГГГ):")
				fmt.Println("Для продолжения работы пожалуйста, нажмите Enter...")
				readLine(reader)
				clear_screen()
				continue
			}
			student.databirthday = data
		case 3:
			fmt.Print("Укажите институт студента: ")
			inst := readLine(reader)
			if err := checkText(inst); err != nil {
				fmt.Println("Ошибка:", err)
				fmt.Println("Для продолжения работы пожалуйста, нажмите Enter...")
				readLine(reader)
				clear_screen()
				continue
			}
			student.institute = inst
		case 4:
			fmt.Print("Укажите стипендию студента: ")
			stiStr := readLine(reader)
			sti, err := strconv.ParseFloat(stiStr, 32)
			if err != nil || sti < 0 {
				fmt.Println("Ошибка ввода стипендии. Введите неотрицательное число.")
				fmt.Println("Для продолжения работы пожалуйста, нажмите Enter...")
				readLine(reader)
				clear_screen()
				continue
			}
			student.stipend = float32(sti)
		case 5:
			fmt.Print("Укажите средний балл студента: ")
			gpaStr := readLine(reader)
			gpa, err := strconv.ParseFloat(gpaStr, 32)
			if err != nil || gpa < 2 || gpa > 5 {
				fmt.Println("Ошибка ввода среднего балла. Введите число от 2 до 5.")
				fmt.Println("Для продолжения работы пожалуйста, нажмите Enter...")
				readLine(reader)
				clear_screen()
				continue
			}
			student.GPA = float32(gpa)
		default:
			fmt.Println("Нет такого пункта меню!")
			fmt.Println("Для продолжения работы пожалуйста, нажмите Enter...")
			readLine(reader)
			clear_screen()
		}
	}
}

// функция, которая удаляет студента с указанным ID и возвращает ошибку, если такого нет
func DeleteStudentByID(students_list *[]students, id int) error {
	if id <= 0 {
		return fmt.Errorf("Нет студента с ID = %d", id)
	}
	for i, st := range *students_list {
		if st.ID == uint(id) {
			*students_list = append((*students_list)[:i], (*students_list)[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Нет студента с ID = %d", id)
}
