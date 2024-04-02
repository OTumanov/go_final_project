Сие надо запускать примерно так --  
`export TODO_PORT=7540 && export TODO_DBFILE="./scheduler.db" && go run cmd/main.go`  
Тесты --  
`go test -run ^TestApp$ ./tests`  
`go test -run ^TestDB$ ./tests`  
`go test -run ^TestNextDate$ ./tests`  
`go test -run ^TestAddTask$ ./tests`  
`go test -run ^TestTasks$ ./tests`
