pm25
====

get city pm2.5. (I hope one day we will not care about the pm2.5)

## How to deploy a server
```
go get github.com/shxsun/pm25  # first build
./pm25 -daemon  # start server
```

**Demo** server: <http://goo.gl/92KKWx>

## How to use
get beijing weather

`pm25 beijing # get beijing pm2.5`

or if you want to to use you own server with $IP and $PORT.
```
./pm25 -server=$IP:$PORT beijing
```

## Thanks
<http://pm25.in>

## Package use
* <https://github.com/ant0ine/go-json-rest>
* <https://github.com/aybabtme/color>
* <https://github.com/bitly/go-simplejson>
