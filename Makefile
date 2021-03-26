all:
	go build  -ldflags "-s -w"  -o api-go   main.go 
	mv -f api-go /usr/local/nginx/html/api/