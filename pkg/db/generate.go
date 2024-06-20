package db

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i Transactor -o ./mocks/ -s "_minimock.go"
//go:generate ../../../bin/minimock -i Tx -o ./mocks/ -s "_minimock.go"
