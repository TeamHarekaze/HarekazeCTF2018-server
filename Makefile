ifeq ($(OS),Windows_NT)
	LIBINSTALL = .\src\lib_install.sh
	BUILD = go build -o .\dist\harekaze.exe .\src
	RUN = go run .\src\harekaze.go
	CLEAN = go clean .\src
else
	LIBINSTALL = sh ./src/lib_install.sh
	BUILD = go build -o ./dist/harekaze ./src
	RUN = go run .\src\harekaze.go
	CLEAN = go clean ./src
endif

libinstall:
	${LIBINSTALL}

build:
	${BUILD}

clean:
	${CLEAN}