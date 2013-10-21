package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ant0ine/go-json-rest"
	"github.com/bitly/go-simplejson"
	"github.com/coocood/jas"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/xorm"
	"github.com/shxsun/pm25/model"
)

var (
	Token  = "" // needed by http://pm25.in
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
	engine, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8", DBUser, DBPass, DBName))
	if err != nil {
		return
	}
	defer engine.Close()
	if err = engine.Sync(new(model.Record)); err != nil {
		return
	}
	err = InitService()
	if err != nil {
		log.Fatal(err)
	}
	go collect(interval)

	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{"GET", "/:loc", GetRecord},
	)
	//http.Handle("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/", http.FileServer(http.Dir("./")))
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//handler.ServeHTTP(w, r)
	//})
	router := jas.NewRouter(new(Pm25))
	router.BasePath = "/api/v2/"
	http.Handle(router.BasePath, router)

	log.Println("Start pm2.5 service .....")
	return http.ListenAndServe(addr, nil)
}

func InitService() (err error) {
	areas, err := Areas()
	if err != nil {
		return
	}
	for _, loc := range areas {
		log.Printf("load city: %s", loc)
		r, err := pm25(loc)
		if err != nil {
			log.Printf("load error: %s", err)
			continue
		}
		records[loc] = r
	}
	return nil
}

// load all areas from db
func Areas() ([]string, error) {
	areaRecords := make([]model.Record, 0)
	err := engine.Cols("area").GroupBy("area").Find(&areaRecords)
	if err != nil {
		return nil, err
	}
	areas := make([]string, 0, len(areaRecords))
	for _, r := range areaRecords {
		areas = append(areas, r.Area)
	}
	return areas, nil
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
		time.Sleep(dur)
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
				//log.Println(err)
			}
		}
		mu.Unlock()
	}
}
