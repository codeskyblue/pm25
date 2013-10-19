package servd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ant0ine/go-json-rest"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/xorm"
	"github.com/shxsun/pm25/model"
)

// needed by http://pm25.in
var (
	Token  = "R9yH3BjLG5g58Z5dUbvn"
	DBUser = "root"
	DBPass = "toor"
	DBName = "pm25"
)

const apiFormat = "http://www.pm25.in/api/querys/aqi_details.json?stations=no&city=%s&token=%s"

var (
	records = make(map[string]model.Record)
	mu      = sync.RWMutex{}
	engine  *xorm.Engine
)

func Run(addr string, interval time.Duration) (err error) {
	log.Println("Start pm2.5 service .....")
	go collect(interval)
	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{"GET", "/:loc", GetRecord},
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	engine, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8", DBUser, DBPass, DBName))
	if err != nil {
		return
	}
	defer engine.Close()
	err = engine.Sync(new(model.Record))
	if err != nil {
		return
	}
	return http.ListenAndServe(addr, nil)
}

func pm25(loc string) (r model.Record, err error) {
	resp, err := http.Get(fmt.Sprintf(apiFormat, loc, Token))
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
	r = model.Record{
		Aqi:  aqi,
		Pm25: station.Get("pm2_5").MustInt(),
		Pm10: station.Get("pm10").MustInt(),
		So2:  station.Get("so2").MustInt(),
		Co:   station.Get("co").MustInt(),
		O3:   station.Get("o3").MustInt(),
		No2:  station.Get("no2").MustInt(),
		Area: loc, //station.Get("area").MustString(),
	}
	timeStr := station.Get("time_point").MustString()
	TimePoint, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		return
	}
	r.TimePoint = TimePoint
	return
}

func collect(dur time.Duration) {
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

			_, err = engine.Insert(r)
			if err != nil {
				log.Println(err)
			}
		}
		mu.Unlock()
		time.Sleep(dur)
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
