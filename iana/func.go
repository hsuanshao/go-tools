package iana

import (
	"sort"

	tzm "github.com/hsuanshao/go-tools/iana/entity/models"
)

func getTimezones() tzm.Timezones {
	timezones := tzm.Timezones{}

	for _, tz := range ianaTimeZones {
		timezones = append(timezones, tzm.IANATimezone(tz))
	}

	sort.Sort(timezones)

	return timezones
}
