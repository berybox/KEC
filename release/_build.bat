cd ../cli
set KECVERSION=v1

set GOARCH=386

set GOOS=linux
go build -o ../release/kec-%KECVERSION%-%GOOS%-x86
c:\utils\7-Zip\7z.exe a -tzip -mmt8 -sdel ../release/kec-%KECVERSION%-%GOOS%-x86.zip ../release/kec-%KECVERSION%-%GOOS%-x86

set GOOS=darwin
go build -o ../release/kec-%KECVERSION%-%GOOS%-x86
c:\utils\7-Zip\7z.exe a -tzip -mmt8 -sdel ../release/kec-%KECVERSION%-%GOOS%-x86.zip ../release/kec-%KECVERSION%-%GOOS%-x86

set GOOS=windows
go build -o ../release/kec-%KECVERSION%-%GOOS%-x86.exe
c:\utils\7-Zip\7z.exe a -tzip -mmt8 -sdel ../release/kec-%KECVERSION%-%GOOS%-x86.zip ../release/kec-%KECVERSION%-%GOOS%-x86.exe

set GOARCH=amd64

set GOOS=linux
go build -o ../release/kec-%KECVERSION%-%GOOS%-x64
c:\utils\7-Zip\7z.exe a -tzip -mmt8 -sdel ../release/kec-%KECVERSION%-%GOOS%-x64.zip ../release/kec-%KECVERSION%-%GOOS%-x64

set GOOS=darwin
go build -o ../release/kec-%KECVERSION%-%GOOS%-x64
c:\utils\7-Zip\7z.exe a -tzip -mmt8 -sdel ../release/kec-%KECVERSION%-%GOOS%-x64.zip ../release/kec-%KECVERSION%-%GOOS%-x64

set GOOS=windows
go build -o ../release/kec-%KECVERSION%-%GOOS%-x64.exe
c:\utils\7-Zip\7z.exe a -tzip -mmt8 -sdel ../release/kec-%KECVERSION%-%GOOS%-x64.zip ../release/kec-%KECVERSION%-%GOOS%-x64.exe

cd ../release

