package repository

import (
	"biathlon/internal/schemas"
	"fmt"
	"time"
)

type Repo struct {
    table map[int]*Competitor
    cfg schemas.Config
}

type Competitor struct {
	ID int
	ScheduledStart time.Time
    ActualStart    time.Time
    Status         string
    CurrentPenaltyCount int
    TotalPenaltyCount  int
    PenaltyTime    time.Duration
    TotalTime      time.Duration
    LastLaps time.Time
    Laps []time.Duration
    Hits           int
    Shots          int
    Goals []bool
    StatusFlag uint32
    StartPenaltyLaps time.Time
    // 0 - зарегистрирован
    // 1 - на дистанции
    // 2 - на огневом рубеже
    // 3 - на штрафном круге 
    // 4 - финишировал
    // 5 - дисквалификация
}

func New(path string) *Repo {
    return &Repo{
        table: make(map[int]*Competitor),
        cfg: schemas.LoadConfig(path),
    }
}

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

func (repo *Repo) StartLine(id int) error {
    // придумать логику
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

func (repo * Repo) StartPenatlyLaps(id int, t time.Time) error {
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



