
http_bench -n 1000 -c 100 -t 3000 -m GET http://localhost:8080/s/\?content\=%E8%87%AA%E5%8A%A8%E6%9C%BA\&offset\=0

go-torch -u http://localhost:8080 -t30 -p > profile-no-cache.svg

go test ./views/search/search_test.go -bench=. -run=none -benchtime=1s
