

Requirements:
you need to install postgresql, goose, go, 

go install will install gator

you will need a config file:  ~/.gatorconfig.json
with a postgres sql server url 
`"db_url":"postgres://hv@localhost:5432/gator?sslmode=disable"`

Running db migrations on local host
run from sql/schema
goose postgres postgres://hv@localhost:5432/gator up

down does what it should

sql query gen
run sqlc generate from roort
