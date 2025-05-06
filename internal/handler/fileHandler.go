package handler

import (
	"biathlon/internal/schemas"
	"biathlon/internal/service"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type FileHandler struct {
	service *service.Service
	fileName string
}


func New(s *service.Service, fileName string) *FileHandler {
	return &FileHandler{
		service: s,
		fileName: fileName,
	}
}

func (h *FileHandler) Start() error {
	events, err := h.parseFile()
	if err != nil {
		return err
	}
	for _, event := range events {
		msg, err := h.service.Execute(&event)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println(msg)
	}
	report, err := h.service.GetResults()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(report)
	return nil
}

func (h * FileHandler) parseFile() ([]schemas.Event, error) {
	file, err := os.Open(h.fileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var events []schemas.Event
    var lastTime time.Time
    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }

        parts := strings.Split(line, "] ")
        if len(parts) < 2 {
            return nil, fmt.Errorf("invalid format at line %d", lineNum)
        }

        timeStr := strings.TrimPrefix(parts[0], "[")
        eventTime, err := time.Parse(schemas.TimeFormat, timeStr)
        if err != nil {
            return nil, fmt.Errorf("invalid time format at line %d: %v", lineNum, err)
        }

        rest := strings.Fields(parts[1])
        if len(rest) < 2 {
            return nil, fmt.Errorf("missing event components at line %d", lineNum)
        }

        eventID, err := strconv.Atoi(rest[0])
        if err != nil {
            return nil, fmt.Errorf("invalid event ID at line %d", lineNum)
        }

        competitorID, err := strconv.Atoi(rest[1])
        if err != nil {
            return nil, fmt.Errorf("invalid competitor ID at line %d", lineNum)
        }

        if !lastTime.IsZero() && eventTime.Before(lastTime) {
            return nil, fmt.Errorf("events out of order at line %d", lineNum)
        }
        lastTime = eventTime

        var params []string
        if len(rest) > 2 {
            params = rest[2:]
        }

        events = append(events, schemas.Event{
            Time:         eventTime,
            EventID:      eventID,
            CompetitorID: competitorID,
            Params:       params,
        })
    }

    return events, scanner.Err()
}