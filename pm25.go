// get latest PM2.5 from http://pm25.in
// Author skyblue.
// -- I hope the sky is blue and the air is clean.
package main

import (
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest"
	//"github.com/aybabtme/color"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const token = "5j1znBVAsnSf5xQyNQyq"

var (
	recyle    = flag.Duration("recyle", time.Minute*30, "data collect recyle duration")
	addr      = flag.String("addr", ":8080", "listen address")
	apiFormat = "http://www.pm25.in/api/querys/aqi_details.json?stations=no&city=%s&token=%s"

	records = make(map[string]Record)
	mu      = sync.RWMutex{}
)

type Record struct {
	AQI       int
	PM25      int
	PM10      int
	SO2       int
	Area      string
	TimePoint string
}

func pm25(loc string) (r Record, err error) {
	resp, err := http.Get(fmt.Sprintf(apiFormat, loc, token))
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	j, err := simplejson.NewJson(data)
	if err != nil {
		return
	}
	station := j.GetIndex(0)
	aqi, err := station.Get("aqi").Int()
	if err != nil {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic:%s", e)
		}
	}()
	r.AQI = aqi
	r.PM25 = station.Get("pm2_5").MustInt()
	r.SO2 = station.Get("so2").MustInt()
	r.PM10 = station.Get("pm10").MustInt()
	r.Area = station.Get("area").MustString()
	r.TimePoint = station.Get("time_point").MustString()
	return
}

func collect(dur *time.Duration) {
	for {
		mu.Lock()
		for loc, _ := range records {
			r, err := pm25(loc)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("Update", r)
			records[loc] = r
		}
		mu.Unlock()
		time.Sleep(*dur)
	}
}

func GetRecord(w *rest.ResponseWriter, req *rest.Request) {
	loc := req.PathParam("loc")
	mu.RLock()
	r, exists := records[loc]
	mu.RUnlock()
	if !exists {
		log.Printf("First request '%s'", loc)
		mu.Lock()
		r, _ = pm25(loc)
		records[loc] = r
		mu.Unlock()
	}
	w.WriteJson(r)
}

func main() {
	fmt.Println("Start pm2.5 service ...")
	go collect(recyle)

	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{"GET", "/:loc", GetRecord},
	)
	http.ListenAndServe(*addr, &handler)
}
