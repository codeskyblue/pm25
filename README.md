pm25
====

get city pm2.5. (I hope one day we will not care about the pm2.5)

## How to deploy
```
go get github.com/shxsun/pm25  # first build
./pm25 -daemon  # start server
```

## How to use
get beijing weather

```
./pm25 -server=$IP:$PORT beijing
```

## Thanks
<http://pm25.in>

## Package use
* <https://github.com/ant0ine/go-json-rest>
* <https://github.com/aybabtme/color>
* <https://github.com/bitly/go-simplejson>
