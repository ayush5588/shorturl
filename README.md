# shorturl

ShortURL is a handy tool for transforming long and unwieldy web addresses into concise, easy-to-use links. Simplify link sharing and management with ShortURL.

It uses Redis as the database.

## Installation
Using docker-compose
```
docker-compose up
```
OR<br><br>
You can directly do:

```
docker run -p 6379:6379 -d redis:alpine
go run main.go
```
