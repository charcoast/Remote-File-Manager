go build -C ./rfm.com/discovery/ -o $(pwd)/discovery
go build -C ./rfm.com/executors/create/ -o $(pwd)/create
go build -C ./rfm.com/executors/list/ -o $(pwd)/list
go build -C ./rfm.com/executors/read/ -o $(pwd)/read