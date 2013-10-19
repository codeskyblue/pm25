package model

import (
	"time"
)

type Record struct {
	Area      string    `xorm:"unique(k)"`
	TimePoint time.Time `xorm:"unique(k)"`
	Aqi       int
	Pm25      int
	Pm10      int
	So2       int
	No2       int
	Co        int
	O3        int
}
