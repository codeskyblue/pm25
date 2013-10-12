// get latest PM2.5 from http://pm25.in
// Author skyblue.
// -- I hope the sky is blue and the air is clean.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aybabtme/color"
	"github.com/shxsun/pm25/api"
)

var (
	daemon = flag.Bool("daemon", false, "start as daemon")
	server = flag.String("server", "115.28.15.5:8077", "server address")
	addr   = flag.String("addr", ":8077", "listen address")
)

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

func progress(tot int, cur int, paint color.Paint) string {
	brush := color.NewBrush("", paint)
	return "[" + brush(strings.Repeat("#", cur)) + strings.Repeat("-", tot-cur) + "]"
}

func cli(loc string) (err error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/%s", *server, flag.Arg(0)))
	if err != nil {
		return
	}
	record := &api.Record{}
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
	bar := progress(50, stars, colorLevel[l])
	fmt.Printf("%-5s %s\n", brush(faceLevel[l]), bar)

	fmt.Printf("%#v\n", *record)
	return
}

func main() {
	flag.Parse()
	if *daemon {
		http.ListenAndServe(*addr, nil)
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
