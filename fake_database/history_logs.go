package fakedatabase

import "time"

type HistoryLog struct {
	ClientId     string
	Source       string
	Target       string
	SourceAmount float64
	TargetAmount float64
	Rate         float64
	CreatedAt    time.Time
}

var HistoryLogs = []HistoryLog{}

func UpdateHistoryLog(clientId string, source string, sourceAmount float64, target string, targetAmount float64, currentRate float64) {
	HistoryLogs =
		append(
			HistoryLogs,
			HistoryLog{
				ClientId:     clientId,
				Source:       source,
				Target:       target,
				SourceAmount: sourceAmount,
				TargetAmount: targetAmount,
				Rate:         currentRate,
				CreatedAt:    time.Now().UTC(),
			},
		)
}
