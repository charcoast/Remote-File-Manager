module rfm.com/discovery

go 1.22.7

replace rfm.com/common => ../common

require (
	github.com/rs/cors v1.11.1
	github.com/thoas/go-funk v0.9.3
	golang.org/x/crypto v0.29.0
	gorm.io/driver/sqlite v1.5.6
	gorm.io/gorm v1.25.12
	rfm.com/common v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)
