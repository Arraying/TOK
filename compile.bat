REM This is used to compile the binary to a variety of different operating systems and architectures.
REM I am on Windows, which is why this is a batch file. Thus, this is used to easily cross compile everything.
REM Here is a list of compatible GOOS operating systems separated by a space: "android darwin dragonfly freebsd linux nacl netbsd openbsd plan9 solaris windows"
SET GOOS=windows
SET GOARCH=amd64
go build .