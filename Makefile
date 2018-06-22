all:
	timestamp=`date +%Y%m%d.%H%M%S` commit=`git log --format="%H" -n 1` go build -o goproxy.bin -ldflags "-X main.buildTimestamp=build=$timestamp;commit=$commit"
