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

