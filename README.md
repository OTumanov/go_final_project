Сие надо запускать примерно так --  
`export TODO_PORT=7540 && export TODO_DBFILE="./scheduler.db" && go run main.go`  
Тесты этого этапа --
`go test -run ^TestDB$ ./tests`