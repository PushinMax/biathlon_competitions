package handler

import (
	"biathlon/internal/schemas"
	"os"
	"strings"
	"testing"
	"time"
	"reflect"
)

func TestParseFile(t *testing.T) {
	validContent := `[09:05:59.867] 1 1
[09:15:00.841] 2 1 09:30:00.000
[09:29:45.734] 3 1
[09:30:01.005] 4 1`

	invalidTimeFormat := `[09:05:59.867] 1 1
[invalid_time] 2 1`

	missingComponents := `[09:05:59.867] 1
[09:15:00.841] 2`

	outOfOrder := `[09:15:00.841] 1 1
[09:05:59.867] 2 1`

	tests := []struct {
		name        string
		content     string
		wantEvents  int
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid file",
			content:    validContent,
			wantEvents: 4,
			wantErr:    false,
		},
		{
			name:        "invalid time format",
			content:     invalidTimeFormat,
			wantEvents:  0,
			wantErr:     true,
			errContains: "invalid time format",
		},
		{
			name:        "missing components",
			content:     missingComponents,
			wantEvents:  0,
			wantErr:     true,
			errContains: "missing event components",
		},
		{
			name:        "events out of order",
			content:     outOfOrder,
			wantEvents:  0,
			wantErr:     true,
			errContains: "events out of order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "testfile")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.WriteString(tt.content); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			h := &FileHandler{
				fileName: tmpfile.Name(),
			}

			events, err := h.parseFile()

			if (err != nil) != tt.wantErr {
				t.Errorf("parseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("parseFile() error = %v, should contain %v", err.Error(), tt.errContains)
				}
				return
			}

			if len(events) != tt.wantEvents {
				t.Errorf("parseFile() got %d events, want %d", len(events), tt.wantEvents)
			}
			
			if tt.name == "valid file" {
				wantTime, _ := time.Parse(schemas.TimeFormat, "09:05:59.867")
				wantEvent := schemas.Event{
					Time:         wantTime,
					EventID:      1,
					CompetitorID: 1,
					Params:       []string{},
				}

				if !events[0].Time.Equal(wantEvent.Time) ||
					events[0].EventID != wantEvent.EventID ||
					events[0].CompetitorID != wantEvent.CompetitorID ||
					len(events[0].Params) != len(wantEvent.Params) {
					t.Errorf("parseFile() first event = %v, want %v", events[0], wantEvent)
				}

				wantTime2, _ := time.Parse(schemas.TimeFormat, "09:15:00.841")
				wantEvent2 := schemas.Event{
					Time:         wantTime2,
					EventID:      2,
					CompetitorID: 1,
					Params:       []string{"09:30:00.000"},
				}

				if !events[1].Time.Equal(wantEvent2.Time) ||
					events[1].EventID != wantEvent2.EventID ||
					events[1].CompetitorID != wantEvent2.CompetitorID ||
					!reflect.DeepEqual(events[1].Params, wantEvent2.Params) {
					t.Errorf("parseFile() second event = %v, want %v", events[1], wantEvent2)
				}
			}
		})
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	h := &FileHandler{
		fileName: "nonexistent_file.txt",
	}

	_, err := h.parseFile()
	if err == nil {
		t.Error("parseFile() expected error for nonexistent file, got nil")
	}
}