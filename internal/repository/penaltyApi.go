package repository

import (
	"fmt"
	"time"
)


func (repo *Repo) StartPenatlyLaps(id int, t time.Time) error {
    comp, ok := repo.table[id]
    if !ok {
        return fmt.Errorf("Competitor(%d) is not registered", id)
    }
    if comp.StatusFlag != 3 {
        return fmt.Errorf("Competitor(%d) use already started or missed your start", id)
    }
    comp.StartPenaltyLaps = t
    return nil
}

func (repo *Repo) EndPenaltyLaps(id int, t time.Time) error {
    comp, ok := repo.table[id]
    if !ok {
        return fmt.Errorf("Competitor(%d) is not registered", id)
    }
    if comp.StatusFlag != 3 {
        return fmt.Errorf("Competitor(%d) use already started or missed your start", id)
    }
    comp.PenaltyTime += t.Sub(comp.StartPenaltyLaps)
    comp.TotalPenaltyCount += comp.CurrentPenaltyCount
    comp.StatusFlag = 1 
    return nil
}