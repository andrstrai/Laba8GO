package main

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type students struct {
	ID           uint
	name         string
	databirthday string
	institute    string
	stipend      float32
	GPA          float32
}

// метод для вывода красивой строки стуктуры
func (st students) String() string {
	return fmt.Sprintf("ID студента: %d\nИмя студента: %s\nДата рождения: %s\n"+
		"Институт: %s\nСтипендия: %.2f\nСредний балл: %.2f",
		st.ID, st.name, st.databirthday, st.institute, st.stipend, st.GPA)
}

// функция для генерации id
func max_id(students_list []students) uint {
	if len(students_list) == 0 {
		return 0
	}
	max := slices.MaxFunc(students_list, func(a, b students) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return max.ID
}

// Проверка строки на соответствие шаблону даты
func is_valid_date(dateStr string) bool {
	// Шаблон: 02 - день, 01 - месяц, 2006 - год
	_, err := time.Parse("2.1.2006", dateStr)
	return err == nil
}

// функция добавления студента в срез
func record(student *[]students, reader *bufio.Reader) {
	id := max_id(*student) + 1

	fmt.Print("Укажите имя студента: ")
	name := readLine(reader)
	for err := checkText(name); err != nil; err = checkText(name) {
		fmt.Printf("Ошибка: %v. Введите имя еще раз: ", err)
		name = readLine(reader)
	}

	fmt.Print("Укажите дату рождения студента в формате ДД.ММ.ГГГГ: ")
	data := readLine(reader)
	for !is_valid_date(data) {
		fmt.Print("Ошибка ввода даты. Введите еще раз (ДД.ММ.ГГГГ): ")
		data = readLine(reader)
	}

	fmt.Print("Укажите институт студента: ")
	inst := readLine(reader)
	for err := checkText(inst); err != nil; err = checkText(inst) {
		fmt.Printf("Ошибка: %v. Введите институт еще раз: ", err)
		inst = readLine(reader)
	}

	fmt.Print("Укажите стипендию студента: ")
	stiStr := readLine(reader)
	sti, err := strconv.ParseFloat(stiStr, 32)
	for err != nil || sti < 0 {
		fmt.Print("Ошибка ввода стипендии. Введите еще раз неотрицательное число: ")
		stiStr = readLine(reader)
		sti, err = strconv.ParseFloat(stiStr, 32)
	}
	STI := float32(sti)
	fmt.Print("Укажите средний балл студента: ")
	gpaStr := readLine(reader)
	gpa, err := strconv.ParseFloat(gpaStr, 32)
	for err != nil || gpa < 2 || gpa > 5 {
		fmt.Print("Ошибка ввода среднего балла. Введите число от 2 до 5: ")
		gpaStr = readLine(reader)
		gpa, err = strconv.ParseFloat(gpaStr, 32)
	}
	GPA := float32(gpa)
	*student = append(*student, students{id, name, data, inst, STI, GPA})
	fmt.Printf("Студент с ID %d добавлен. Не забудьте сохранить изменения (пункт 5).\n", id)
}

// функция, которая сортирует слайс студентов по среднему баллу по убыванию
func sort_by_grades(students_list []students) {
	slices.SortFunc(students_list, func(a, b students) int {
		switch {
		case a.GPA < b.GPA:
			return 1
		case a.GPA > b.GPA:
			return -1
		default:
			return 0
		}
	})
	fmt.Println("Данные отсортированы!")
}

// функция, которая сортирует слайс студентов по размеру стипендии по убыванию
func sort_by_stipend(students_list []students) {
	slices.SortFunc(students_list, func(a, b students) int {
		switch {
		case a.stipend < b.stipend:
			return 1
		case a.stipend > b.stipend:
			return -1
		default:
			return 0
		}
	})
	fmt.Println("Данные отсортированы!")
}

// функция, которая возвращает студентов указанного института
func filter_by_institute(students_list []students, inst string) []students {
	result := []students{}
	for _, st := range students_list {
		if strings.EqualFold(st.institute, inst) {
			result = append(result, st)
		}
	}
	return result
}

// функция для подсчёта средней стипендии
func average_stipend(students_list []students) float32 {
	if len(students_list) == 0 {
		return 0
	}
	sum := float32(0.0)
	for _, st := range students_list {
		sum += st.stipend
	}
	return sum / float32(len(students_list))
}

// функция, которая возвращает студентов со стипендией выше средней
func filter_stipend_above_average(students_list []students) []students {
	avg := average_stipend(students_list)
	result := []students{}
	for _, st := range students_list {
		if st.stipend > avg {
			result = append(result, st)
		}
	}
	return result
}

// функция для вывода студентов таблицей
func print_all(student []students) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 6, ' ', 0)
	fmt.Fprintln(w, "ID\tФИО\tДата рождения\tИнститут\tРазмер стипендий\tСредний балл")
	fmt.Fprintln(w, "--\t------\t-----\t----\t----\t")
	if len(student) != 0 {
		for i := 0; i < len(student); i++ {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%f\t%f\n", student[i].ID, student[i].name, student[i].databirthday, student[i].institute,
				student[i].stipend, student[i].GPA)
		}
		w.Flush()
	}
	fmt.Println("Всего студентов:", len(student))
	fmt.Println("---")
}

// функция для сохранения измененнного слайса в файл
func save_to_file(students_list []students) error {
	f, err := os.Create("Base.tmp")
	if err != nil {
		fmt.Println("Ошибка записи!")
		return err
	}
	w := bufio.NewWriter(f)
	for _, student := range students_list {
		fmt.Fprintf(w, "%d/%s/%s/%s/%.2f/%.2f\n",
			student.ID, student.name, student.databirthday, student.institute, student.stipend, student.GPA)
	}
	if err := w.Flush(); err != nil {
		f.Close()
		os.Remove("Base.tmp")
		fmt.Println("Ошибка записи!")
		return err
	}
	f.Close()
	if err := os.Rename("Base.tmp", "Base.txt"); err != nil {
		os.Remove("Base.tmp")
		fmt.Println("Ошибка записи!")
		return err
	}
	fmt.Println("Изменения сохранены!")
	return nil
}

// функция для очистки экрана (работает и на Windows, и на macOS/Linux)
func clear_screen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// основная функция
func main() {
	all_students := []students{}

	file, err := os.Open("Base.txt")
	if err != nil {
		f, err := os.Create("Base.txt")
		if err != nil {
			fmt.Println("Не удалось создать файл:", err)
			return
		}
		f.Close()
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fields := strings.Split(scanner.Text(), "/")
			if len(fields) < 6 {
				continue
			}
			id, err1 := strconv.Atoi(fields[0])
			sti, err2 := strconv.ParseFloat(fields[4], 32)
			gpa, err3 := strconv.ParseFloat(fields[5], 32)
			if err1 != nil || err2 != nil || err3 != nil {
				fmt.Println("Пропущена некорректная строка:", scanner.Text())
				continue
			}
			STI := float32(sti)
			GPA := float32(gpa)
			all_students = append(all_students, students{uint(id), fields[1], fields[2], fields[3], STI, GPA})
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Ошибка чтения:", err)
		}
		file.Close()
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Для начала работы пожалуйста, нажмите Enter...")
		readLine(reader)
		clear_screen()

		fmt.Println("Меню команд:")
		fmt.Println("1 - Показать список всех студентов")
		fmt.Println("2 - Запись студента в базу данных")
		fmt.Println("3 - Сортировать студентов по среднему баллу (по убыванию)")
		fmt.Println("4 - Сортировать студентов по размеру стипендии (по убыванию)")
		fmt.Println("5 - Сохранить изменения в файл")
		fmt.Println("6 - Показать студентов выбранного института")
		fmt.Println("7 - Показать студентов со стипендией выше средней")
		fmt.Println("8 - Изменить данные студента")
		fmt.Println("9 - Удалить студента")
		fmt.Println("11 - Отчёт в файл: топ-10 по среднему баллу")
		fmt.Println("12 - Отчёт в файл: стипендия выше средней")
		fmt.Println("13 - Отчёт в файл: студенты выбранного института")
		fmt.Println("14 - Отчёт в файл: сводка по институтам")
		fmt.Println("0 - Завершить работу")
		fmt.Println("Выберите пункт из меню управления: ")

		input := readLine(reader)

		a, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Некорректный ввод! Пожалуйста, введите число.")
			continue
		}

		switch a {
		case 0:
			return
		case 1:
			print_all(all_students)
		case 2:
			record(&all_students, reader)
		case 3:
			sort_by_grades(all_students)
		case 4:
			sort_by_stipend(all_students)
		case 5:
			save_to_file(all_students)
		case 6:
			fmt.Print("Введите институт: ")
			inst := readLine(reader)
			print_all(filter_by_institute(all_students, inst))
		case 7:
			fmt.Printf("Средняя стипендия: %.2f\n", average_stipend(all_students))
			print_all(filter_stipend_above_average(all_students))
		case 8:
			fmt.Println("Список всех студентов:")
			print_all(all_students)
			fmt.Println("\nВведите ID студента")
			input := readLine(reader)

			check_id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Некорректный ввод! Пожалуйста, введите число.")
				continue
			}
			check_st, err := StudentByID(all_students, check_id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			ChangeStudent(check_st, reader)
		case 9:
			fmt.Println("Список всех студентов:")
			print_all(all_students)
			fmt.Print("\nВведите ID студента для удаления: ")
			input := readLine(reader)

			del_id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Некорректный ввод! Пожалуйста, введите число.")
				continue
			}
			if err := DeleteStudentByID(&all_students, del_id); err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Студент с ID %d удален. Не забудьте сохранить изменения (пункт 5).\n", del_id)
		case 11:
			err := ReportTop(all_students, 10, reportTopFile)
			printReportResult(reportTopFile, err)
		case 12:
			err := ReportAboveAverage(all_students, reportAboveAvgFile)
			printReportResult(reportAboveAvgFile, err)
		case 13:
			fmt.Print("Введите институт: ")
			inst := readLine(reader)
			err := ReportByInstitute(all_students, inst, reportInstituteFile)
			printReportResult(reportInstituteFile, err)
		case 14:
			err := ReportInstitutesSummary(all_students, reportSummaryFile)
			printReportResult(reportSummaryFile, err)
		default:
			fmt.Println("Нет такого пункта меню!")
		}
	}
}
