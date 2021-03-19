package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yaitoo/sparrow/types"
)

type sharding struct {
	token string
	value string
}

func getShardings(cmd string, from, to *time.Time) []sharding {

	if from != nil {

		if strings.Index(cmd, "/*year*/") > -1 {
			return getYearShardings(from, to)
		} else if strings.Index(cmd, "/*month*/") > -1 {
			return getMonthShardings(from, to)

		} else if strings.Index(cmd, "/*week*/") > -1 {
			return getWeekShardings(from, to)
		} else if strings.Index(cmd, "/*day*/") > -1 {
			return getDayShardings(from, to)
		}
	}

	return []sharding{sharding{token: "", value: ""}}
}

func getYearShardings(from, to *time.Time) []sharding {
	var fromYear, toYear int

	if from == nil {
		return []sharding{}
	}

	fromYear = from.Year()
	toYear = from.Year()

	if to != nil && to.Year() > fromYear {
		toYear = to.Year()
	}

	shardings := make([]sharding, 0, 5)
	for i := 0; fromYear+i <= toYear; i++ {
		sharding := sharding{}
		sharding.token = "/*year*/"
		sharding.value = fmt.Sprint(fromYear + i)

		shardings = append(shardings, sharding)
	}

	return shardings
}

func getMonthShardings(from, to *time.Time) []sharding {
	var fromMonth, toMonth time.Time

	if from == nil {
		return []sharding{}
	}

	fromMonth = time.Date(from.Year(), from.Month(), 0, 0, 0, 0, 0, from.Location())
	toMonth = fromMonth

	if to != nil && to.After(*from) {
		toMonth = time.Date(to.Year(), to.Month(), 0, 0, 0, 0, 0, to.Location())
	}

	shardings := make([]sharding, 0, 5)
	for i := 0; fromMonth.AddDate(0, i, 0).After(toMonth) == false; i++ {
		now := fromMonth.AddDate(0, i, 0)
		sharding := sharding{}
		sharding.token = "/*month*/"
		sharding.value = types.FormatTime(&now, "yyyyMM")

		shardings = append(shardings, sharding)
	}

	return shardings
}

func getWeekShardings(from, to *time.Time) []sharding {
	var fromWeek, toWeek int

	if from == nil {
		return []sharding{}
	}

	fromWeek = getWeek(*from)
	toWeek = fromWeek

	if to != nil && to.After(*from) {
		toWeek = getWeek(*to)
	}

	shardings := make([]sharding, 0, 5)
	for i := 0; fromWeek+i <= toWeek; i++ {
		now := from.AddDate(0, 0, i*7)
		sharding := sharding{}
		sharding.token = "/*week*/"
		sharding.value = strconv.Itoa(getWeek(now))

		shardings = append(shardings, sharding)
	}

	return shardings
}

func getDayShardings(from, to *time.Time) []sharding {
	var fromDay, toDay time.Time

	if from == nil {
		return []sharding{}
	}

	fromDay = *from
	toDay = fromDay

	if to != nil && to.After(*from) {
		toDay = *to
	}

	shardings := make([]sharding, 0, 5)
	for i := 0; fromDay.AddDate(0, 0, i).After(toDay) == false; i++ {
		now := fromDay.AddDate(0, 0, i)
		sharding := sharding{}
		sharding.token = "/*day*/"
		sharding.value = types.FormatTime(&now, "yyyyMMdd")

		shardings = append(shardings, sharding)
	}

	return shardings
}

func getWeek(time time.Time) int {
	year, week := time.ISOWeek()

	return year*100 + week
}
