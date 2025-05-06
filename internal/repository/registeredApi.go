package repository

import (
	"time"
	"fmt"
)
func (repo *Repo) RegisteredCompetitor(id int) error {
    if _, ok := repo.table[id]; !ok {
        repo.table[id] = &Competitor{
            ID: id,
            StatusFlag: 0,
            Laps: make([]time.Duration, 0),

        }
        return nil
    }

    return fmt.Errorf("Competitor(%d) already exists", id)
}

func (repo *Repo) SetStartTime(id int, t time.Time) error {
    comp, ok := repo.table[id]
    if !ok {
        return fmt.Errorf("Competitor(%d) is not registered", id)
    }
    // Изменить
    if comp.StatusFlag > 0 {
        return fmt.Errorf("Competitor(%d) use already started or missed your start", id)
    }
    
    comp.ScheduledStart = t
    return nil
}
