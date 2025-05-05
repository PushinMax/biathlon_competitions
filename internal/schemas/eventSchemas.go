package schemas


import "time"

type Event struct {
    Time     time.Time
    EventID  int
    CompetitorID int
    Params   []string
}

type Config struct {
    Laps        int    `json:"laps"`
    LapLen      int    `json:"lapLen"`
    PenaltyLen  int    `json:"penaltyLen"`
    FiringLines int    `json:"firingLines"`
    Start       string `json:"start"`
    StartDelta  string `json:"startDelta"`
}