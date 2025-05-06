package repository

import (
	"fmt"
	"time"
)


func (repo *Repo) StartLine(id int) error {
    return nil
}

func (repo *Repo) StartComp(id int, t time.Time) error {
    comp, ok := repo.table[id]
    if !ok {
        return fmt.Errorf("Competitor(%d) is not registered", id)
    }
    if comp.StatusFlag > 0 {
        return fmt.Errorf("Competitor(%d) use already started or missed your start", id)
    }
    comp.StatusFlag = 1

    comp.ActualStart = t
    comp.LastLaps = comp.ScheduledStart
    comp.Laps = append(comp.Laps, 0)
    // реализовать доп логику проверки вермени старта

    return nil
}