BIN=yetti
MAIN=github.com/MarkyMan4/yetti/main
PARSER=github.com/MarkyMan4/yetti/parser
LEXER=github.com/MarkyMan4/yetti/lexer

build:
	go build -o $(BIN) main.go

clean:
	rm $(BIN)

test_all:
	go test ./...

test_main:
	go test $(MAIN)

test_parser:
	go test $(PARSER)

test_lexer:
	go test $(LEXER)

