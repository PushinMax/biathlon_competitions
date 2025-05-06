package repository

import (
	"fmt"
	"time"
)



func (repo *Repo) EndLap(id int, t time.Time) error {
    comp, ok := repo.table[id]
    if !ok {
        return fmt.Errorf("Competitor(%d) is not registered", id)
    }
    if comp.StatusFlag != 1 {
        return fmt.Errorf("Competitor(%d) at other stage", id)
    }
    
    comp.Laps[len(comp.Laps) - 1] = t.Sub(comp.LastLaps)
    if len(comp.Laps) == repo.cfg.Laps {
        // финиш для участника
        comp.TotalTime = t.Sub(comp.ScheduledStart)
        return nil
    }
    comp.LastLaps = t
    comp.Laps = append(comp.Laps, 0)
    

    return nil
}

func (repo *Repo) Termination(id int, comment string) error {
    comp, ok := repo.table[id]
    if !ok {
        return fmt.Errorf("Competitor(%d) is not registered", id)
    }
    if comp.StatusFlag == 5 {
        return fmt.Errorf("Competitor(%d) at other stage", id)
    }
    comp.Status = comment
    comp.StatusFlag = 5
    return nil
}