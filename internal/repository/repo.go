package repository

import (
	"biathlon/internal/schemas"
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