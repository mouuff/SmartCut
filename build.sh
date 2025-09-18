
binary_name=smartcut_windows_amd64.exe
zip_name=binaries_windows.zip
go build -ldflags -H=windowsgui cmd/smartcut/main.go
mv main.exe $binary_name
rm $zip_name
zip -r $zip_name $binary_name