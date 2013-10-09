// get latest PM2.5 from http://pm25.in
// Author skyblue.
// -- I hope the sky is blue and the air is clean.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest"
	"github.com/aybabtme/color"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const token = "R9yH3BjLG5g58Z5dUbvn"

var (
	daemon    = flag.Bool("daemon", false, "start as daemon")
	server    = flag.String("server", "115.28.15.5:8077", "server address")
	recyle    = flag.Duration("recyle", time.Minute*10, "data collect recyle duration")
	addr      = flag.String("addr", ":8077", "listen address")
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

var colorLevel = []color.Paint{
	color.GreenPaint,
	color.YellowPaint,
	color.RedPaint,
	color.PurplePaint,
	color.PurplePaint,
}

var faceLevel = []string{
	"^O^",
	"-_-",
	"-_!",
	"-_-!",
	"-_-!!",
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
	r = Record{
		AQI:       aqi,
		PM25:      station.Get("pm2_5").MustInt(),
		PM10:      station.Get("pm10").MustInt(),
		SO2:       station.Get("so2").MustInt(),
		Area:      station.Get("area").MustString(),
		TimePoint: station.Get("time_point").MustString(),
	}
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

func progress(tot int, cur int, paint color.Paint) string {
	brush := color.NewBrush("", paint)
	return "[" + brush(strings.Repeat("#", cur)) + strings.Repeat("-", tot-cur) + "]"
}

func cli(loc string) (err error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/%s", *server, flag.Arg(0)))
	if err != nil {
		return
	}
	record := &Record{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, record)
	if err != nil {
		return
	}
	l := record.AQI / 100
	if l > 5 {
		l = 5
	}
	brush := color.NewBrush("", colorLevel[l])
	stars := (record.AQI + 9) / 10
	//grayBrush := color.NewBrush("", color.LightGrayPaint)
	//bar := "[" + brush(strings.Repeat("#", stars)) + grayBrush(strings.Repeat("-", 50-stars)) + "]"
	bar := progress(50, stars, colorLevel[l])
	fmt.Printf("%-5s %s\n", brush(faceLevel[l]), bar)

	fmt.Printf("%#v\n", *record)
	return
}

func main() {
	flag.Parse()
	if *daemon {
		log.Println("Start pm2.5 service ...")
		go collect(recyle)
		handler := rest.ResourceHandler{}
		handler.SetRoutes(
			rest.Route{"GET", "/:loc", GetRecord},
		)
		http.ListenAndServe(*addr, &handler)
	} else {
		if flag.NArg() != 1 {
			flag.Usage()
			fmt.Printf("[EXAMPLE]\n%s beijing   # will get beijing pm2.5\n", os.Args[0])
			return
		}
		if err := cli(flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
	}
}
