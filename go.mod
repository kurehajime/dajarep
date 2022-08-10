module github.com/kurehajime/dajarep

go 1.19

require (
	github.com/ikawaha/kagome-dict/ipa v1.0.2
	github.com/ikawaha/kagome/v2 v2.4.4
)

require (
	github.com/ikawaha/kagome-dict v1.0.2 // indirect
	golang.org/x/net v0.0.0-20220809184613-07c6da5e1ced
	golang.org/x/text v0.3.7
)

replace (
    github.com/kurehajime/dajarep => ./
)