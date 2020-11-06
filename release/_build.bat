cd ../cli
set KECVERSION=v1

set GOARCH=386

set GOOS=linux
go build -o ../release/kec-%KECVERSION%-linux-x86
set GOOS=darwin
go build -o ../release/kec-%KECVERSION%-darwin-x86
set GOOS=windows
go build -o ../release/kec-%KECVERSION%-windows-x86.exe

set GOARCH=amd64

set GOOS=linux
go build -o ../release/kec-%KECVERSION%-linux-x64
set GOOS=darwin
go build -o ../release/kec-%KECVERSION%-darwin-x64
set GOOS=windows
go build -o ../release/kec-%KECVERSION%-windows-x64.exe

cd ../release