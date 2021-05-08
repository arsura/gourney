package histories

import (
	"sort"

	fakedatabase "github.com/arsura/moonbase-service/fake_database"
	"github.com/gofiber/fiber/v2"
)

type GetHistoryLogsQuery struct {
	Source string `query:"source"`
	Target string `query:"target"`
	Skip   int32  `query:"skip"`
	Limit  int32  `query:"limit"`
	Sort   string `query:"sort"`
}

func GetHistoryLogsHandler(c *fiber.Ctx) error {
	query := new(GetHistoryLogsQuery)
	if err := c.QueryParser(query); err != nil {
		return err
	}
	sortedField := query.Sort
	sortedHistoryLogs := fakedatabase.HistoryLogs
	if sortedField == "createdAt.dsc" {
		sort.SliceStable(sortedHistoryLogs, func(i, j int) bool {
			return sortedHistoryLogs[i].CreatedAt.After(sortedHistoryLogs[j].CreatedAt)
		})
	}
	return c.JSON(&fiber.Map{"data": []fakedatabase.HistoryLog(sortedHistoryLogs)})
}
