package service

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/xuri/excelize/v2"
)

func exportExcelSummary(taskSummaries []map[string]interface{}, headers []string) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Tasks"
	_ = f.SetSheetName("Sheet1", sheet)
	_ = f.SetDefaultFont("Arial")

	// Style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "#000000", Family: "Arial"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DAF2D0"}, Pattern: 1},
	})

	// --- Header ---
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}
	lastCol, _ := excelize.ColumnNumberToName(len(headers))
	_ = f.SetCellStyle(sheet, "A1", lastCol+"1", headerStyle)
	_ = f.SetRowHeight(sheet, 1, 25)
	_ = f.SetColWidth(sheet, "B", "B", 15)
	_ = f.SetColWidth(sheet, "C", "C", 15)
	_ = f.SetColWidth(sheet, "D", "D", 25)
	_ = f.SetColWidth(sheet, "G", "G", 25)
	_ = f.SetColWidth(sheet, "M", "M", 15)
	_ = f.SetColWidth(sheet, "N", "N", 15)
	_ = f.SetColWidth(sheet, "P", "P", 15)
	_ = f.SetColWidth(sheet, lastCol, lastCol, 40)

	// --- Data rows ---
	for i, t := range taskSummaries {
		row := i + 2
		values := []interface{}{
			t["month"], t["date"], t["username"], t["full_name"], t["client"], t["job"], t["item"], t["role"],
			t["started_at"], t["ended_at"], t["hour"], t["minute"], t["break_time"], t["total_time"],
			t["ot"], t["volume"], t["note"],
		}
		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			_ = f.SetCellValue(sheet, cell, v)
		}
	}

	// --- Export to buffer ---
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return &buf, nil
}

func exportExcelJob(data []map[string]interface{}, dates []string) (*bytes.Buffer, error) {
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

	// --- Export to buffer ---
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return &buf, nil
}

func exportExcelTimeTotal(data []map[string]interface{}, dates []string) (*bytes.Buffer, error) {
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

	// --- Export to buffer ---
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return &buf, nil
}
