pm25
====

![leaf](images/leaf.png)

Monitor pm2.5. (I hope one day sky is blue)

This is just my interest of Beijing weather.

## How to use
get beijing weather

`pm25 beijing # get beijing pm2.5`

or if you want to to use you own server, for example: 115.28.15.5:8077
```
./pm25 -addr=115.28.15.5:8080 beijing
```

## How to deploy a server
This need a mysql-server. Just offer a dbname. The program will CreateTable itself.

```
go get github.com/shxsun/pm25  # build
./pm25 -daemon  # start server
```

**Demo** server: <http://goo.gl/92KKWx>

## Links
* api resource: <http://pm25.in>
* what is AQI: <http://m.guokr.com/post/431588/>
* history data of beijing: <http://aqicn.org/city/beijing/m/>

## Package use
* <https://github.com/ant0ine/go-json-rest>
* <https://github.com/aybabtme/color>
* <https://github.com/bitly/go-simplejson>
