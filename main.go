package main

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/xuri/excelize/v2"
)

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

type Task1 struct {
	UserInformation UserInfo  `json:"user_information"`
	Minute          int       `json:"minute"`
	Job             string    `json:"job"`
	CreatedAt       time.Time `json:"created_at"`
}

type TaskStatisticRequestDto struct {
	StartDate *time.Time `form:"start_date" time_format:"2006-01-02" time_utc:"true"`
	EndDate   *time.Time `form:"end_date"   time_format:"2006-01-02" time_utc:"true"`
}

// -------------------- PHẦN XỬ LÝ DỮ LIỆU --------------------

func getDateRange(start, end *time.Time) []string {
	const layout = "2006-01-02"
	var dates []string
	for d := *start; !d.After(*end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format(layout))
	}
	return dates
}

func groupByJobAndDay(tasks []Task1) map[string]map[string]int {
	const layout = "2006-01-02"
	m := make(map[string]map[string]int)
	for _, t := range tasks {
		job := t.Job
		if job == "" {
			job = "(blank)"
		}
		day := t.CreatedAt.Format(layout)
		if _, ok := m[job]; !ok {
			m[job] = make(map[string]int)
		}
		m[job][day]++
	}
	return m
}

func groupByUserAndDay(tasks []Task1) map[string]map[string]int {
	const layout = "2006-01-02"
	m := make(map[string]map[string]int)
	for _, t := range tasks {
		name := t.UserInformation.FullName
		if name == "" {
			name = "(blank)"
		}
		day := t.CreatedAt.Format(layout)
		if _, ok := m[name]; !ok {
			m[name] = make(map[string]int)
		}
		m[name][day] += t.Minute
	}
	return m
}

func summarizeTasksByType(tasks []Task1, start, end *time.Time, exportType string) ([]map[string]interface{}, []string) {
	dates := getDateRange(start, end)
	dayCount := len(dates)

	result := make([]map[string]interface{}, 0)

	switch exportType {
	case "PROJECT_TOTAL":
		grouped := groupByJobAndDay(tasks)
		for job, days := range grouped {
			record := map[string]interface{}{"job": job}
			total := 0
			for _, d := range dates {
				val := days[d]
				record[d] = val
				total += val
			}
			record["grand_total"] = total
			result = append(result, record)
		}

	case "TIME_TOTAL":
		grouped := groupByUserAndDay(tasks)
		for user, days := range grouped {
			record := map[string]interface{}{"full_name": user}
			total := 0
			for _, d := range dates {
				val := days[d]
				record[d] = val
				total += val
			}
			record["average"] = math.Round((float64(total)/float64(dayCount))*10) / 10
			result = append(result, record)
		}
	}

	return result, dates
}

// -------------------- PHẦN EXPORT EXCEL --------------------

func exportExcelTimeTotal(data []map[string]interface{}, dates []string) error {
	f := excelize.NewFile()
	sheet := "Report"
	_ = f.SetSheetName(f.GetSheetName(0), sheet)
	_ = f.SetDefaultFont("Arial")

	// Style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "#000000", Family: "Arial"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DAF2D0"}, Pattern: 1},
	})

	cellStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	boldStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Family: "Arial"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	//========== BẢNG 1: HOURS ==========
	_ = f.SetCellValue(sheet, "A1", "Total Time (Hours)")
	_ = f.SetCellValue(sheet, "A2", "Full name")
	for i, d := range dates {
		col, _ := excelize.ColumnNumberToName(i + 2)
		parsed, _ := time.Parse("2006-01-02", d)
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s2", col), parsed.Format("1/2"))
	}
	lastCol, _ := excelize.ColumnNumberToName(len(dates) + 2)
	_ = f.SetCellValue(sheet, lastCol+"2", "Trung bình")
	_ = f.SetCellStyle(sheet, "A2", lastCol+"2", headerStyle)

	row := 3
	for _, rec := range data {
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), rec["full_name"])
		total := 0.0
		for i, d := range dates {
			col, _ := excelize.ColumnNumberToName(i + 2)
			val := float64(rec[d].(int)) / 60.0
			val = math.Round(val*10) / 10
			total += val
			_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), val)
			_ = f.SetCellStyle(sheet, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("%s%d", col, row), cellStyle)
		}
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", lastCol, row), math.Round(total/float64(len(dates))*10)/10)
		_ = f.SetCellStyle(sheet, fmt.Sprintf("%s%d", lastCol, row), fmt.Sprintf("%s%d", lastCol, row), boldStyle)
		row++
	}

	//========== BẢNG 2: MINUTES ==========
	row += 1
	_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total Time (Minutes)")
	_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row+1), "Full name")
	for i, d := range dates {
		col, _ := excelize.ColumnNumberToName(i + 2)
		parsed, _ := time.Parse("2006-01-02", d)
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row+1), parsed.Format("1/2"))
	}
	_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", lastCol, row+1), "Trung bình")
	_ = f.SetCellStyle(sheet, fmt.Sprintf("A%d", row+1), lastCol+fmt.Sprint(row+1), headerStyle)

	row += 2
	for _, rec := range data {
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), rec["full_name"])
		for i, d := range dates {
			col, _ := excelize.ColumnNumberToName(i + 2)
			val := rec[d]
			_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), val)
			_ = f.SetCellStyle(sheet, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("%s%d", col, row), cellStyle)
		}
		// Ghi cột "Trung bình" (đã có sẵn trong rec["average"])
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", lastCol, row), rec["average"])
		_ = f.SetCellStyle(sheet, fmt.Sprintf("%s%d", lastCol, row), fmt.Sprintf("%s%d", lastCol, row), boldStyle)
		row++
	}

	_ = f.SetColWidth(sheet, "A", "A", 30)

	_ = f.SetColWidth(sheet, "B", lastCol, 10)
	_ = f.SetColWidth(sheet, lastCol, lastCol, 15)
	return f.SaveAs("report.xlsx")
}

func exportExcelJob(data []map[string]interface{}, dates []string) error {
	f := excelize.NewFile()
	oldSheet := f.GetSheetName(0)
	sheet := "Project Totals"
	_ = f.SetSheetName(oldSheet, sheet)
	_ = f.SetDefaultFont("Arial")

	// Độ rộng cột
	_ = f.SetColWidth(sheet, "A", "A", 13)

	// Style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "#000000", Family: "Arial"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DAF2D0"}, Pattern: 1},
	})
	cellStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Family: "Arial"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	boldStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Family: "Arial"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	// Header (bắt đầu từ A1)
	_ = f.SetCellValue(sheet, "A1", "Job")
	for i, d := range dates {
		col, _ := excelize.ColumnNumberToName(i + 2)
		parsed, _ := time.Parse("2006-01-02", d)
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s1", col), parsed.Format("1/2"))
	}
	lastCol, _ := excelize.ColumnNumberToName(len(dates) + 2)
	_ = f.SetCellValue(sheet, lastCol+"1", "Grand Total")
	_ = f.SetCellStyle(sheet, "A1", lastCol+"1", headerStyle)

	// Data (bắt đầu từ dòng 2)
	row := 2
	for _, rec := range data {
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), rec["job"])
		for i, d := range dates {
			col, _ := excelize.ColumnNumberToName(i + 2)
			val := rec[d]
			_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), val)
			_ = f.SetCellStyle(sheet, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("%s%d", col, row), cellStyle)
		}
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", lastCol, row), rec["grand_total"])
		_ = f.SetCellStyle(sheet, fmt.Sprintf("%s%d", lastCol, row), fmt.Sprintf("%s%d", lastCol, row), boldStyle)
		row++
	}
	_ = f.SetColWidth(sheet, "B", lastCol, 10)
	_ = f.SetColWidth(sheet, lastCol, lastCol, 15)

	return f.SaveAs("project_report.xlsx")
}

// -------------------- MAIN --------------------

func main() {
	data := []Task1{
		{UserInformation: UserInfo{FullName: "Nguyễn Thị Tố Hoa"}, Job: "Backend", Minute: 453, CreatedAt: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)},
		{UserInformation: UserInfo{FullName: "Phạm Thị An"}, Job: "Backend", Minute: 221, CreatedAt: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)},
		{UserInformation: UserInfo{FullName: "Đỗ Thị Hồng Phong"}, Job: "Frontend", Minute: 452, CreatedAt: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)},
		{UserInformation: UserInfo{FullName: "Lê Thị Liễu"}, Job: "Frontend", Minute: 460, CreatedAt: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)},
		{UserInformation: UserInfo{FullName: "Lê Thị Khánh Linh"}, Job: "DevOps", Minute: 325, CreatedAt: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)},
		{UserInformation: UserInfo{FullName: "Lê Hồ Uyên Phương"}, Job: "QA", Minute: 370, CreatedAt: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)},

		// thêm 1 vài ngày khác để test multi-date
		{UserInformation: UserInfo{FullName: "Nguyễn Thị Tố Hoa"}, Job: "Backend", Minute: 440, CreatedAt: time.Date(2025, 8, 2, 0, 0, 0, 0, time.UTC)},
		{UserInformation: UserInfo{FullName: "Phạm Thị An"}, Job: "Backend", Minute: 230, CreatedAt: time.Date(2025, 8, 2, 0, 0, 0, 0, time.UTC)},
		{UserInformation: UserInfo{FullName: "Đỗ Thị Hồng Phong"}, Job: "Frontend", Minute: 450, CreatedAt: time.Date(2025, 8, 2, 0, 0, 0, 0, time.UTC)},
		{UserInformation: UserInfo{FullName: "Lê Thị Liễu"}, Job: "Frontend", Minute: 470, CreatedAt: time.Date(2025, 8, 2, 0, 0, 0, 0, time.UTC)},
	}
	fmt.Println(data)
	start := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC)
	req := TaskStatisticRequestDto{StartDate: &start, EndDate: &end}

	result, dates := summarizeTasksByType(data, req.StartDate, req.EndDate, "PROJECT_TOTAL")
	fmt.Println(result)

	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
	fmt.Println(dates)

	//if err := exportExcelTimeTotal(result, dates); err != nil {
	//	fmt.Println("❌ Lỗi export:", err)
	//	return
	//}

	if err := exportExcelJob(result, dates); err != nil {
		fmt.Println("❌ Lỗi export:", err)
		return
	}
	fmt.Println("✅ Export thành công -> report.xlsx")
}
