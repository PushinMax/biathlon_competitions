package repository

import (
	"strings"
	"testing"
	"time"
)

func TestRepo_RegisteredCompetitor(t *testing.T) {
	repo := New("test_config.json")

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"Register new competitor", 1, false},
		{"Register duplicate competitor", 1, true},
		{"Register another competitor", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.RegisteredCompetitor(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisteredCompetitor() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if _, exists := repo.table[tt.id]; !exists {
					t.Errorf("Competitor %d not registered", tt.id)
				}
			}
		})
	}
}

func TestRepo_SetStartTime(t *testing.T) {
	repo := New("test_config.json")
	testTime := time.Now()

	
	repo.RegisteredCompetitor(1)
	repo.RegisteredCompetitor(2)
	repo.table[2].StatusFlag = 1 

	tests := []struct {
		name    string
		id      int
		time    time.Time
		wantErr bool
	}{
		{"Set time for new competitor", 1, testTime, false},
		{"Set time for started competitor", 2, testTime, true},
		{"Set time for non-existent competitor", 3, testTime, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.SetStartTime(tt.id, tt.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetStartTime() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if !repo.table[tt.id].ScheduledStart.Equal(tt.time) {
					t.Errorf("Start time not set correctly")
				}
			}
		})
	}
}

func TestRepo_StartComp(t *testing.T) {
	repo := New("test_config.json")
	startTime := time.Now()

	repo.RegisteredCompetitor(1)
	repo.SetStartTime(1, startTime)
	repo.RegisteredCompetitor(2)
	repo.RegisteredCompetitor(3)
	repo.SetStartTime(3, startTime)
	repo.table[3].StatusFlag = 1 

	tests := []struct {
		name    string
		id      int
		time    time.Time
		wantErr bool
	}{
		{"Start valid competitor", 1, startTime, false},
		{"Start already started competitor", 3, startTime, true},
		{"Start non-existent competitor", 4, startTime, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.StartComp(tt.id, tt.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartComp() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				comp := repo.table[tt.id]
				if comp.StatusFlag != 1 {
					t.Errorf("StatusFlag not set to 1")
				}
				if !comp.ActualStart.Equal(tt.time) {
					t.Errorf("Actual start time not set correctly")
				}
				if len(comp.Laps) != 1 {
					t.Errorf("Initial lap not created")
				}
			}
		})
	}
}

func TestRepo_EndLap(t *testing.T) {
	repo := New("test_config.json")
	repo.cfg.Laps = 3 
	startTime := time.Now()
	lap1Time := startTime.Add(5 * time.Minute)

	repo.RegisteredCompetitor(1)
	repo.SetStartTime(1, startTime)
	repo.StartComp(1, startTime)

	repo.RegisteredCompetitor(2) 
	repo.RegisteredCompetitor(3)
	repo.SetStartTime(3, startTime)
	repo.StartComp(3, startTime)
	repo.table[3].StatusFlag = 2 

	tests := []struct {
		name    string
		id      int
		time    time.Time
		wantErr bool
	}{
		{"Competitor not started", 2, lap1Time, true},
		{"Competitor on firing range", 3, lap1Time, true},
		{"Non-existent competitor", 4, lap1Time, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.EndLap(tt.id, tt.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndLap() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				comp := repo.table[tt.id]
				lastLapIdx := len(comp.Laps) - 1

				if strings.Contains(tt.name, "Final lap") {
					if comp.StatusFlag != 4 {
						t.Errorf("Competitor not marked as finished")
					}
					if comp.TotalTime != tt.time.Sub(comp.ScheduledStart) {
						t.Errorf("Total time not calculated correctly")
					}
				} else {
					if comp.Laps[lastLapIdx-1] != tt.time.Sub(comp.LastLaps) {
						t.Errorf("Lap time not recorded correctly")
					}
					if len(comp.Laps) != lastLapIdx+2 { 
						t.Errorf("New lap not created")
					}
				}
			}
		})
	}
}


func TestRepo_PenaltyLapsFlow(t *testing.T) {
	repo := New("test_config.json")
	startTime := time.Now()
	penaltyStart := startTime.Add(10 * time.Minute)
	repo.RegisteredCompetitor(1)
	repo.SetStartTime(1, startTime)
	repo.StartComp(1, startTime)
	repo.StartRange(1, 1)
	repo.Hit(1, 1)
	repo.LeftRange(1) 

	t.Run("Start penalty laps", func(t *testing.T) {
		err := repo.StartPenatlyLaps(1, penaltyStart)
		if err != nil {
			t.Errorf("StartPenatlyLaps() error = %v", err)
		}

		if !repo.table[1].StartPenaltyLaps.Equal(penaltyStart) {
			t.Errorf("StartPenaltyLaps not set correctly")
		}
	})

}

func TestRepo_Termination(t *testing.T) {
	repo := New("test_config.json")
	startTime := time.Now()

	repo.RegisteredCompetitor(1)
	repo.SetStartTime(1, startTime)
	repo.StartComp(1, startTime)

	t.Run("Terminate competitor", func(t *testing.T) {
		err := repo.Termination(1, "Injury")
		if err != nil {
			t.Errorf("Termination() error = %v", err)
		}

		comp := repo.table[1]
		if comp.StatusFlag != 5 {
			t.Errorf("StatusFlag not set to 5 (disqualified)")
		}
		if comp.Status != "Injury" {
			t.Errorf("Status = %s, want 'Injury'", comp.Status)
		}
	})

	t.Run("Terminate already terminated competitor", func(t *testing.T) {
		err := repo.Termination(1, "Injury")
		if err == nil {
			t.Errorf("Expected error when terminating already terminated competitor")
		}
	})
}

