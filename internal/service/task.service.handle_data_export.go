package service

import (
	"base_go_be/internal/constants"
	"base_go_be/internal/model"
	"base_go_be/internal/until"
	"math"
	"time"
)

func (ts *TaskService) summaryTask(tasks []model.Task) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(tasks))

	for _, task := range tasks {
		var username, fullName string
		if task.UserInformation != nil {
			username = task.UserInformation.Username
			fullName = task.UserInformation.FullName
		}
		record := map[string]interface{}{
			"month":      task.CustomCreatedAt.Format("01"),
			"date":       task.CustomCreatedAt.Format("2006/01/02"),
			"username":   username,
			"full_name":  fullName,
			"client":     task.Client,
			"job":        task.Job,
			"item":       task.Item,
			"role":       task.Role,
			"started_at": until.SafeTimeUTC7(task.StartedAt),
			"ended_at":   until.SafeTimeUTC7(task.EndedAt),
			"hour":       task.Hours,
			"minute":     task.Minute,
			"break_time": task.BreakTime,
			"total_time": task.WorkTime,
			"ot":         task.OT,
			"volume":     task.Volume,
			"note":       task.Note,
		}
		results = append(results, record)
	}

	return results
}

// Calculate the amount of JOB and time
func (ts *TaskService) analysisTasksByType(tasks []model.Task, start, end *time.Time, exportType string) []map[string]interface{} {
	dates := until.GetDateRange(start, end)
	dayCount := len(dates)

	result := make([]map[string]interface{}, 0)

	switch exportType {
	case constants.TaskExportProjectTotals:
		grouped := ts.groupByJobAndDay(tasks)
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

	case constants.TaskExportTimeTotals:
		grouped := ts.groupByUserAndDay(tasks)
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

	return result
}

// groups tasks by job name and creation date
func (ts *TaskService) groupByJobAndDay(tasks []model.Task) map[string]map[string]int {
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

func (ts *TaskService) groupByUserAndDay(tasks []model.Task) map[string]map[string]int {
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
