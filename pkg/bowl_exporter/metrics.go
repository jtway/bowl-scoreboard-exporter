package bowl_exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const Namespace = "bowl_exporter"

type metrics struct {
	score         *prometheus.GaugeVec
	quarter       *prometheus.GaugeVec
	timeRemaining *prometheus.GaugeVec
	winner        *prometheus.GaugeVec
	inProgress    *prometheus.GaugeVec
	homeTeam      *prometheus.GaugeVec
}

func NewMetrics() *metrics {
	m := &metrics{
		score: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "team_score",
			Help:      "current score for team",
		},
			[]string{"bowl_name", "venue", "date", "date_time", "team_location", "team_name", "record"},
		),
		quarter: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "quarter",
			Help:      "current quarter or 1 for not started",
		},
			[]string{"bowl_name", "venue", "date", "date_time", "team_location", "team_name", "record"},
		),
		timeRemaining: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "time_remaining",
			Help:      "Time remaining in quarter",
		},
			[]string{"bowl_name", "venue", "date", "date_time", "team_location", "team_name", "record"},
		),
		winner: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "winner",
			Help:      "1 if the team is the winner, 0 if loser or game not finished",
		},
			[]string{"bowl_name", "venue", "date", "date_time", "team_location", "team_name", "record"},
		),
		inProgress: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "in_progress",
			Help:      "1 if the game is in progress, 0 otherwise",
		},
			[]string{"bowl_name", "venue", "date", "date_time", "team_location", "team_name", "record"},
		),
		homeTeam: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "home_team",
			Help:      "1 is the team is the home team, 0 if away",
		},
			[]string{"bowl_name", "venue", "date", "date_time", "team_location", "team_name", "record"},
		),
	}
	return m
}
