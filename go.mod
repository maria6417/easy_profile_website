module example.com/go-mod-test

go 1.16

replace local.packages/models => ./09_mongo/models

require (
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)
