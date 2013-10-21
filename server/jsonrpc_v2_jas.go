package server

import (
	"github.com/coocood/jas"
	"github.com/shxsun/pm25/model"
	"log"
)

type Pm25 struct{}

func (*Pm25) CityList(ctx *jas.Context) {
	mu.RLock()
	defer mu.RUnlock()
	cities := make([]string, 0, len(records))
	for c, _ := range records {
		cities = append(cities, c)
	}
	ctx.Data = cities
}

func (*Pm25) Details(ctx *jas.Context) {
	loc, _ := ctx.FindString("loc")

	r, err := reqRecord(loc)
	if err != nil {
		ctx.Error = jas.NewRequestError("city not monitord")
		return
	}
	ctx.Data = r
}

func reqRecord(loc string) (r model.Record, err error) {
	mu.RLock()
	r, exists := records[loc]
	mu.RUnlock()
	if !exists {
		log.Printf("First request '%s'", loc)
		r, err = pm25(loc)
		if err != nil {
			return
		}
		mu.Lock()
		records[loc] = r
		mu.Unlock()
	}
	return
}
