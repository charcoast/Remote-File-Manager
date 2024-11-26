go build -C ./rfm.com/discovery/ -o $(pwd)/rfm-discovery
go build -C ./rfm.com/executors/create/ -o $(pwd)/rfm-create
go build -C ./rfm.com/executors/list/ -o $(pwd)/rfm-list
go build -C ./rfm.com/executors/read/ -o $(pwd)/rfm-read