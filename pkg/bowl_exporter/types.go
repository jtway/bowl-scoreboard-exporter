package bowl_exporter

import (
	"strings"
	"time"
)

const (
	RFC3339Short    = "2006-01-02T10:00Z07:00"
	RFC1123DateOnly = "Mon, 02 Jan 2006"
)

// Only going to deal with the bowl specifics right now
type ScoreboardResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	Id        string `json:"id"`
	Uid       string `json:"uid"`
	Date      string `json:"date"` // Really need to just figure out date parsing
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Season    struct {
		Year int    `json:"year"`
		Type int    `json:"type"`
		Slug string `json:"slug"`
	} `json:"season"`
	Week struct {
		Number int `json:"number"`
	} `json:"week"`
	Competitions []Competition `json:"competitions"`
	// Links
	Weather struct {
		DisplayName     string `json:"displayName"`
		Temperature     uint32 `json:"temperature"`
		HighTemperature uint32 `json:"highTemperature"`
	}
}

func (e *Event) GetBowlDate() (string, error) {
	parsedDate, err := e.GetBowlDatetime()
	if err != nil {
		return "", err
	}
	return parsedDate.Format(RFC1123DateOnly), nil
}

func (e *Event) GetBowlDatetime() (time.Time, error) {
	pieces := strings.Split(e.Date, "Z")
	date := pieces[0] + ":00Z" + pieces[1]
	return time.Parse(time.RFC3339, date)
}

func (e *Event) GetBowlName() string {
	if len(e.Competitions) == 0 {
		return ""
	}
	return e.Competitions[0].GetBowlName()
}

func (e *Event) GetGameStatus() GameResult {
	if len(e.Competitions) == 0 {
		return ""
	}
	return e.Competitions[0].GetGameStatus()
}

func (e *Event) GetVenue() Venue {
	return e.Competitions[0].GetVenue()
}

type Competition struct {
	// Will want stadium, location, notes for the bowl name
	Competitors []Competitor `json:"competitors"`
	Date        string       `json:"date"`
	Id          string       `json:"id"`
	NeutralSite bool         `json:"neutralSite"`
	Venue       Venue        `json:"venue"`
	Notes       []Note       `json:"notes"`
	Status      Status       `json:"status"`
	StartDate   string       `json:"startDate"`
}

func (c *Competition) GetBowlName() string {
	return c.Notes[0].Headline
}

func (c *Competition) GetGameStatus() GameResult {
	return c.Status.Type.Name
}

func (c *Competition) GetVenue() Venue {
	return c.Venue
}

type Note struct {
	Type     string `json:"type"`
	Headline string `json:"headline"`
}

type Venue struct {
	Id       string `json:"id"`
	FullName string `json:"fullName"`
	Address  struct {
		City  string `json:"city"`
		State string `json:"state"`
	}
	Capacity uint32 `json:"capacity"`
	Indoor   bool   `json:"indoor"`
}

type Competitor struct {
	Id       string   `json:"id"`
	Uid      string   `json:"uid"`
	Type     string   `json:"type"`     // May not need this
	HomeAway string   `json:"homeAway"` // Would be nice to translate this to an enum
	Team     Team     `json:"team"`
	Score    string   `json:"score"`
	Records  []Record `json:"records"`
	Winner   *bool    `json:"winner"`
}

func (c *Competitor) GetOverallRecord() string {
	for _, record := range c.Records {
		if record.Name == "overvall" {
			return record.Summary
		}
	}
	return "0-0"
}

type Record struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Type         string `json:"type"`
	Summary      string `json:"summary"`
}

type Team struct {
	Id               string `json:"id"`
	Uid              string `json:"uid"`
	Location         string `json:"location"`
	Name             string `json:"name"`
	Abbreviation     string `json:"abbreviation"`
	DisplayName      string `json:"displayName"`
	ShortDisplayName string `json:"shortDisplayName"`
	Color            string `json:"color"`
	AlternateColor   string `json:"alternateColor"`
	Logo             string `json:"logo"` // Really a URL
}

type Status struct {
	Clock   float32 `json:"clock"`
	Display string  `json:"displayClock"`
	Period  int     `json:"period"`
	Type    struct {
		Id          string     `json:"id"`
		Name        GameResult `json:"name"`
		State       string     `json:"state"`
		Completed   bool       `json:"completed"`
		Description string     `json:"description"`
		Detail      string     `json:"detail"`
		ShortDetail string     `json:"shortDetail"`
	}
}

type GameResult string

const (
	Unknown    GameResult = ""
	Final      GameResult = "STATUS_FINAL"
	Scheduled  GameResult = "STATUS_SCHEDULED"
	InProgress GameResult = "STATUS_IN_PROGRESS"
)
