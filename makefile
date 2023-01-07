files := $(shell find src/*.go)

battleship : $(files)
	go build -o battleship $(files)

clean:
	rm battleship
