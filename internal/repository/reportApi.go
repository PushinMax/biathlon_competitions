package repository

import (
	//"biathlon/internal/schemas"
	"fmt"
	"strings"
	"time"
	"sort"
	"log"
)

func (repo *Repo) GetReport() (string, error) {
	competitors := make([]struct{
		id int
		t time.Duration
	}, 0)

	for _, comp := range repo.table {
		competitors = append(competitors, struct{id int; t time.Duration}{
			id: comp.ID,
			t: comp.TotalTime,
		})
	}
	sort.Slice(competitors, func(i, j int) bool {
		if competitors[i].t == 0 {
			return false
		}
		return competitors[i].t < competitors[j].t
	})

	var builder strings.Builder
	for i, c := range competitors {
		comp, ok := repo.table[c.id]
		if !ok {
			log.Fatal("1111")
		}
		var laps strings.Builder
		for i, v := range comp.Laps {
			if v == 0 {
				laps.WriteString(
					"{,}",
				)
				break
			}
			laps.WriteString(fmt.Sprintf(
				"{%s, %s}",
				formatDuration(v),
				fmt.Sprintf("%.3f", float64(repo.cfg.LapLen) / v.Seconds()),
			))
			if i != len(comp.Laps) - 1 {
				laps.WriteString(", ")
			}
		}
		var totalTime string
		if comp.TotalTime == 0 {
			if comp.StatusFlag == 0 {
				totalTime = "NotStarted"
			}
			if comp.StatusFlag == 5 {
				totalTime = "NotFinished"	
			}
		} else {
			totalTime = formatDuration(comp.TotalTime)
		}

		var penaltyTime string
		if comp.PenaltyTime == 0 {
			penaltyTime = fmt.Sprintf(
				"{%s,}",
				formatDuration(comp.PenaltyTime),
			)
		} else {
			penaltyTime = fmt.Sprintf(
				"{%s, %s}",
				formatDuration(comp.PenaltyTime),
				fmt.Sprintf("%.3f", float64(repo.cfg.PenaltyLen) / comp.PenaltyTime.Seconds()),
			)
		}

		builder.WriteString(
			fmt.Sprintf(
				"[%s] %d [%s] %s %d/%d",
				totalTime,
				comp.ID,
				laps.String(),
				penaltyTime,
				comp.Hits,
				comp.Shots,
			),
		)
		if i != len(repo.table) {
			builder.WriteByte('\n')
		}
	}
	return builder.String(), nil
}

func formatDuration(d time.Duration) string {
    totalMs := d.Milliseconds()
    ms := totalMs % 1000
    totalSec := totalMs / 1000
    sec := totalSec % 60
    totalMin := totalSec / 60
    min := totalMin % 60
    hour := totalMin / 60

    return fmt.Sprintf("%02d:%02d:%02d.%03d", hour, min, sec, ms)
}