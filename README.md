# EC
Compile Provisional Ghana 2016 Election results from Electoral Commission Ghana Facebook feed.
Results are output to STDOUT as CSV (a simple output redirection will get a CSV file), 
and also persisted to a MongoDB database (which can be read later to feed a dashboard of some sort).

## Parameters
Change the following in main.go.

1. <app_id>: replace with valid Facebook App ID.
1. <app_secret>: replace with valid Facebook App Secret.
