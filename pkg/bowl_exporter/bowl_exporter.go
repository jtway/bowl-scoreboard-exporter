package bowl_exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type BowlScoreboardClient struct {
	client  *http.Client
	config  *Config
	metrics *metrics

	scoreboard *ScoreboardResponse
}

func NewBowlScoreboardClient(config *Config) *BowlScoreboardClient {
	bowlScoreboardClient := BowlScoreboardClient{
		config: config,
	}
	bowlScoreboardClient.client = &http.Client{
		Timeout: time.Second * 10,
	}
	bowlScoreboardClient.metrics = NewMetrics()

	return &bowlScoreboardClient
}

func (bsc *BowlScoreboardClient) ParseResponse(responseBody []byte) (*ScoreboardResponse, error) {

	var scoreboardResponse ScoreboardResponse
	if err := json.Unmarshal(responseBody, &scoreboardResponse); err != nil { // Parse []byte to the go struct pointer
		return nil, fmt.Errorf("can not unmarshal JSON. Err: %w", err)
	}

	sort.Slice(scoreboardResponse.Events, func(i, j int) bool {
		// Can't think of a better solution than if we can't parse a date
		// just maintain order.
		firstDate, err := scoreboardResponse.Events[i].GetBowlDatetime()
		if err != nil {
			fmt.Printf("Error parsing first Bowl Date: %s\n", err.Error())
			return false
		}
		secondDate, err := scoreboardResponse.Events[j].GetBowlDatetime()
		if err != nil {
			fmt.Printf("Error parsing second Bowl Date: %s\n", err.Error())
			return true
		}
		return secondDate.Before(firstDate)
	})
	return &scoreboardResponse, nil
}

func (bsc *BowlScoreboardClient) GetScoreboard() (*ScoreboardResponse, error) {

	request, err := http.NewRequest("GET", bsc.config.Endpoint, nil)
	if err != nil {
		return nil, err
	}

	response, err := bsc.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body) // response body is []byte
	if err != nil {
		return nil, err
	}

	scoreboardResponse, err := bsc.ParseResponse(body)
	if err != nil {
		return nil, err
	}

	return scoreboardResponse, nil
}

func (bsc *BowlScoreboardClient) recordMetrics(scoreboad *ScoreboardResponse) {
	for _, event := range scoreboad.Events {

		competition := event.Competitions[0]
		bowlDate, _ := event.GetBowlDate()
		bowlDateTime := event.Date
		bowlName := event.GetBowlName()
		venueName := event.GetVenue().FullName
		for _, competitor := range competition.Competitors {
			teamLocation := competitor.Team.Location
			teamName := competitor.Team.Name
			record := competitor.GetOverallRecord()

			score, _ := strconv.Atoi(competitor.Score)

			bsc.metrics.score.WithLabelValues(bowlName, venueName, bowlDate,
				bowlDateTime, competitor.HomeAway, teamName, record).Set(float64(score))

			bsc.metrics.quarter.WithLabelValues(bowlName, venueName, bowlDate,
				bowlDateTime, competitor.HomeAway, teamLocation, teamName, record).Set(float64(competition.Status.Period))

			bsc.metrics.timeRemaining.WithLabelValues(bowlName, venueName, bowlDate,
				bowlDateTime, competitor.HomeAway, teamLocation, teamName, record).Set(float64(competition.Status.Clock))

			winner := 0
			if competitor.Winner != nil && *competitor.Winner {
				winner = 1
			}
			bsc.metrics.winner.WithLabelValues(bowlName, venueName, bowlDate,
				bowlDateTime, competitor.HomeAway, teamLocation, teamName, record).Set(float64(winner))

			inProgress := 0
			if competition.GetGameStatus() == InProgress {
				inProgress = 1
			}
			bsc.metrics.inProgress.WithLabelValues(bowlName, venueName, bowlDate,
				bowlDateTime, competitor.HomeAway, teamLocation, teamName, record).Set(float64(inProgress))
		}
	}

}

func (bsc *BowlScoreboardClient) GameInProgress() bool {
	for _, event := range bsc.scoreboard.Events {
		gameStatus := event.GetGameStatus()
		if gameStatus == InProgress {
			return true
		}
		bowlDate, err := event.GetBowlDatetime()
		if err != nil {
			fmt.Printf("Unable to parse bowl date: %s\n", err)
			continue
		}
		// If the game start time is before now and the games not over, this means
		// we haven't made a request since kick off, and need to
		if bowlDate.Before(time.Now()) && gameStatus != Final {
			return true
		}
	}
	return false
}

func (bsc *BowlScoreboardClient) Run() {

	go func() {
		for {
			if bsc.scoreboard == nil || bsc.GameInProgress() {
				var err error
				bsc.scoreboard, err = bsc.GetScoreboard()
				if err != nil {
					panic(fmt.Errorf("unable to retrieve bowl scoreboard, %w", err))
				}
			}
			bsc.recordMetrics(bsc.scoreboard)
			time.Sleep(bsc.config.FetchInterval)
		}
	}()
}
