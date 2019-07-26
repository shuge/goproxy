TIMESTAMP=`date +%Y%m%d.%H%M%S`
COMMIT=`git log --format="%H" -n 1`

all:
	go build -o goproxy.bin -ldflags "-X main.buildTimestamp=build=$(TIMESTAMP);commit=$(COMMIT)" main.go

clean:
	rm goproxy.bin
