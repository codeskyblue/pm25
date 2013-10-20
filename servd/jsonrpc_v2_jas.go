package servd

import (
	"log"

	"github.com/coocood/jas"
)

type Pm25 struct{}

func (*Pm25) Details(ctx *jas.Context) {
	loc, _ := ctx.FindString("loc")

	var err error
	mu.RLock()
	r, exists := records[loc]
	mu.RUnlock()
	if !exists {
		log.Printf("First request '%s'", loc)
		r, err = pm25(loc)
		if err != nil {
			ctx.Error = jas.NewRequestError("city not monitord")
			return
		}
		mu.Lock()
		records[loc] = r
		mu.Unlock()
	}
	ctx.Data = r
}
