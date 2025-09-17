
binary_name=$(go run cmd/smartcut/main.go -printBin=true)
go build -ldflags -H=windowsgui cmd/smartcut/main.go
mv main.exe $binary_name.exe
rm binaries.zip
zip -r binaries.zip $binary_name.exe