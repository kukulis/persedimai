# programa, kuri ieško skrydžių su persėdimais


## dabartinis testinis variantas:

    cd application

    go mod tidy

    go run main.go


## dokeris

Dokeris kol kas panaudotas tik dviejų db sukūrimui: 'persedimai' ir 'test'.

Vėliau gal ir pačią aplikaciją reiks leisti dockerio konteineryje.
 

## tests

To run a separate test:

    cd drafttests
    go test -run TestLoadDbConfig

### special tags for Draft tests:

//go:build draft


Usage:
- Normal test run (skips drafttests): go test ./...
- Run with drafttests: go test -tags=draft ./...

Generating clustering data

    go test -v -timeout 0  -run TestClustersCreator

Dumping data

    mysqldump -P 23314 -u root -h 127.0.0.1 -p test > clusters_32.sql