# EC
Compile Provisional Ghana 2016 Election results from Electoral Commission Ghana Facebook feed.
Results are output to STDOUT as CSV (a simple output redirection will get a CSV file), 
and also persisted to a MongoDB database (which can be read later to feed a dashboard of some sort).

## Run
1. $ go build
1. $ APPID=<appid> APPSECRET=<appsecret> ./ec

or

1. $ APPID=<appid> APPSECRET=<appsecret> go run main.go

