package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"time"
)

// Имена файлов, в которые записываются отчёты
const (
	reportTopFile       = "report_top10.txt"
	reportAboveAvgFile  = "report_above_average.txt"
	reportInstituteFile = "report_institute.txt"
	reportSummaryFile   = "report_summary.txt"
)

// Формат даты и времени в шапке отчёта: день.месяц.год часы:минуты
const reportDateLayout = "02.01.2006 15:04"

// createReport создаёт файл отчёта, пишет в него шапку,
// затем вызывает body, чтобы дописать содержимое, и сохраняет файл
func createReport(path, title string, body func(w io.Writer)) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("не удалось создать файл %s: %w", path, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("не удалось закрыть файл %s: %w", path, cerr)
		}
	}()

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, title)
	fmt.Fprintln(w, "Дата создания:", time.Now().Format(reportDateLayout))
	fmt.Fprintln(w)
	body(w)

	if err := w.Flush(); err != nil {
		return fmt.Errorf("не удалось записать файл %s: %w", path, err)
	}
	return nil
}

// writeStudentsTable пишет таблицу студентов с выровненными столбцами
func writeStudentsTable(w io.Writer, list []students) {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "№\tID\tФИО\tДата рождения\tИнститут\tСтипендия\tСредний балл")
	for i, st := range list {
		fmt.Fprintf(tw, "%d\t%d\t%s\t%s\t%s\t%.2f\t%.2f\n",
			i+1, st.ID, st.name, st.databirthday, st.institute, st.stipend, st.GPA)
	}
	tw.Flush()
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Всего в отчёте:", len(list))
}

// ReportTop записывает в файл n студентов с самым высоким средним баллом
func ReportTop(list []students, n int, path string) error {
	sorted := slices.Clone(list)
	slices.SortStableFunc(sorted, func(a, b students) int {
		return cmp.Compare(b.GPA, a.GPA) // b и a поменяны местами: по убыванию
	})
	top := sorted[:min(n, len(sorted))]

	title := fmt.Sprintf("Топ-%d студентов по среднему баллу", n)
	return createReport(path, title, func(w io.Writer) {
		writeStudentsTable(w, top)
	})
}

// ReportAboveAverage записывает в файл студентов со стипендией выше средней
func ReportAboveAverage(list []students, path string) error {
	avg := average_stipend(list)
	above := filter_stipend_above_average(list)

	return createReport(path, "Студенты со стипендией выше средней", func(w io.Writer) {
		fmt.Fprintf(w, "Средняя стипендия: %.2f\n\n", avg)
		writeStudentsTable(w, above)
	})
}

// ReportByInstitute записывает в файл студентов указанного института
func ReportByInstitute(list []students, inst string, path string) error {
	found := filter_by_institute(list, inst)

	return createReport(path, "Студенты института "+inst, func(w io.Writer) {
		writeStudentsTable(w, found)
	})
}

// instituteStat — накопленные данные по одному институту
type instituteStat struct {
	count      int
	sumGPA     float32
	sumStipend float32
}

// ReportInstitutesSummary записывает в файл сводку по каждому институту:
// сколько студентов, средний балл и средняя стипендия
func ReportInstitutesSummary(list []students, path string) error {
	stats := map[string]*instituteStat{}
	for _, st := range list {
		key := strings.ToUpper(st.institute) // "икнк" и "ИКНК" считаем одним институтом
		s, ok := stats[key]
		if !ok {
			s = &instituteStat{}
			stats[key] = s
		}
		s.count++
		s.sumGPA += st.GPA
		s.sumStipend += st.stipend
	}

	names := make([]string, 0, len(stats))
	for name := range stats {
		names = append(names, name)
	}
	slices.Sort(names)

	return createReport(path, "Сводка по институтам", func(w io.Writer) {
		tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
		fmt.Fprintln(tw, "Институт\tСтудентов\tСредний балл\tСредняя стипендия")
		for _, name := range names {
			s := stats[name]
			n := float32(s.count)
			fmt.Fprintf(tw, "%s\t%d\t%.2f\t%.2f\n", name, s.count, s.sumGPA/n, s.sumStipend/n)
		}
		tw.Flush()
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Всего институтов:", len(names), "| всего студентов:", len(list))
	})
}

// printReportResult сообщает в консоль, получилось ли сохранить отчёт
func printReportResult(path string, err error) {
	if err != nil {
		fmt.Println("Ошибка при создании отчёта:", err)
		return
	}
	fmt.Println("Отчёт сохранён в файл", path)
}
