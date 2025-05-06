package service

import (
	"biathlon/internal/repository"
	"biathlon/internal/schemas"
	"fmt"
	"time"
	"strconv"
	"strings"
)

type Service struct {
	repo *repository.Repo
}

func New(repo *repository.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Execute(event *schemas.Event) (string, error) {
	switch event.EventID {
	case 1:
		err := s.repo.RegisteredCompetitor(event.CompetitorID)
		return fmt.Sprintf("[%s] The competitor(%d) registered", 
			event.Time.Format(schemas.TimeFormat), event.CompetitorID), err
	case 2:
		time, err := time.Parse(schemas.TimeFormat, event.Params[0])
		if err != nil {
			return "", err
		}
		err = s.repo.SetStartTime(event.CompetitorID, time)
		return fmt.Sprintf("[%s] The start time for the competitor(%d) was set by a draw to %s", 
			event.Time.Format(schemas.TimeFormat), event.CompetitorID, event.Params[0]), err
	case 3:
		err := s.repo.StartLine(event.CompetitorID)
		return fmt.Sprintf("[%s] The competitor(%d) is on the start line", 
			event.Time.Format(schemas.TimeFormat), event.CompetitorID), err
	case 4:
		err := s.repo.StartComp(event.CompetitorID, event.Time)
		return fmt.Sprintf("[%s] The competitor(%d) has started", 
			event.Time.Format(schemas.TimeFormat), event.CompetitorID), err
	case 5:
		num, err := strconv.Atoi(event.Params[0])
		if err != nil {
			return "", err
		}
		err = s.repo.StartRange(event.CompetitorID, num)
		return fmt.Sprintf("[%s] The competitor(%d) is on the firing range(%d)", 
			event.Time.Format(schemas.TimeFormat), 
			event.CompetitorID,
			num), err
	case 6:
		num, err := strconv.Atoi(event.Params[0])
		if err != nil {
			return "", err
		}
		err = s.repo.Hit(event.CompetitorID, num)
		return fmt.Sprintf("[%s] The target(%d) has been hit by competitor(%d)", 
			event.Time.Format(schemas.TimeFormat), 
			num, event.CompetitorID), err
	case 7:
		err := s.repo.LeftRange(event.CompetitorID)
		return fmt.Sprintf("[%s] The competitor(%d) left the firing range", 
			event.Time.Format(schemas.TimeFormat), event.CompetitorID), err
	case 8:
		err := s.repo.StartPenatlyLaps(event.CompetitorID, event.Time)
		return fmt.Sprintf("[%s] The competitor(%d) entered the penalty laps", 
			event.Time.Format(schemas.TimeFormat), event.CompetitorID), err
	case 9:
		err := s.repo.EndPenaltyLaps(event.CompetitorID, event.Time)
		return fmt.Sprintf("[%s] The competitor(%d) left the penalty laps", 
			event.Time.Format(schemas.TimeFormat), event.CompetitorID), err
	case 10:
		err := s.repo.EndLap(event.CompetitorID, event.Time)
		return fmt.Sprintf("[%s] The competitor(%d) ended the main lap", 
			event.Time.Format(schemas.TimeFormat), event.CompetitorID), err
	case 11:
		comment := strings.Join(event.Params, " ")
		err := s.repo.Termination(event.CompetitorID, comment)
		return fmt.Sprintf("[%s] The competitor(%d) can`t continue: %s", 
			event.Time.Format(schemas.TimeFormat), event.CompetitorID, comment), err
	}
	return "", fmt.Errorf("unknown event")
}

func (s *Service) GetResults() (string, error) {
	return s.repo.GetReport()
}