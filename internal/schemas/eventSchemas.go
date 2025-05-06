package schemas


import (
	"time"
	"io/ioutil"
	"encoding/json"
)

type Event struct {
    Time     time.Time
    EventID  int
    CompetitorID int
    Params   []string
}

var TimeFormat = "15:04:05.000"

type Config struct {
    Laps        int    `json:"laps"`
    LapLen      int    `json:"lapLen"`
    PenaltyLen  int    `json:"penaltyLen"`
    FiringLines int    `json:"firingLines"`
    Start       string `json:"start"`
    StartDelta  string `json:"startDelta"`
}

func LoadConfig(path string) Config {
    file, _ := ioutil.ReadFile(path)
    var config Config
    json.Unmarshal(file, &config)
    return config
}

