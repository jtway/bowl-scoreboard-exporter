package bowl_exporter_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jtway/bowl-scoreboard-exporter/pkg/bowl_exporter"
	"github.com/stretchr/testify/require"
)

func TestResponse(t *testing.T) {
	data, err := os.ReadFile("testdata/bowl-scoreboard.json")
	require.NoError(t, err)
	require.NotEmpty(t, data)

	config := &bowl_exporter.Config{
		Endpoint:      "http://site.api.espn.com/apis/site/v2/sports/football/college-football/scoreboard",
		FetchInterval: 60 * time.Second,
	}

	client := bowl_exporter.NewBowlScoreboardClient(config)

	scoreboardResponse, err := client.ParseResponse(data)
	require.NoError(t, err)
	require.NotNil(t, scoreboardResponse)

	require.NotEmpty(t, scoreboardResponse.Events)
	for _, event := range scoreboardResponse.Events {
		fmt.Printf("%s, %s: %s vs %s\n", event.Date, event.GetBowlName(),
			event.Competitions[0].Competitors[0].Team.DisplayName,
			event.Competitions[0].Competitors[1].Team.DisplayName)

	}
}

func TestStatusResponse(t *testing.T) {
	data, err := os.ReadFile("testdata/bowl-scoreboard.json")
	require.NoError(t, err)
	require.NotEmpty(t, data)

	config := &bowl_exporter.Config{
		Endpoint:      "http://site.api.espn.com/apis/site/v2/sports/football/college-football/scoreboard",
		FetchInterval: 60 * time.Second,
	}

	client := bowl_exporter.NewBowlScoreboardClient(config)

	scoreboardResponse, err := client.ParseResponse(data)
	require.NoError(t, err)
	require.NotNil(t, scoreboardResponse)

	require.NotEmpty(t, scoreboardResponse.Events)

	for _, event := range scoreboardResponse.Events {
		competition := event.Competitions[0]
		if competition.Notes[0].Headline == "Myrtle Beach Bowl" {
			require.Equal(t, competition.Status.Type.Name, bowl_exporter.Final)
		}
	}
}
