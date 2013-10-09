pm25
====

![leaf](images/leaf.png)

Monitor pm2.5. (I hope one day sky is blue)

## How to use
get beijing weather

`pm25 beijing # get beijing pm2.5`

or if you want to to use you own server(for example: 115.123.321.1:8080/weather)
```
./pm25 -server=115.123.321.1:8080/weather beijing
```

## How to deploy a server
```
go get github.com/shxsun/pm25  # build
./pm25 -daemon  # start server
```

**Demo** server: <http://goo.gl/92KKWx>

## Thanks
<http://pm25.in>

## Package use
* <https://github.com/ant0ine/go-json-rest>
* <https://github.com/aybabtme/color>
* <https://github.com/bitly/go-simplejson>
