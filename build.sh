# go build -ldflags -H=windowsgui cmd/smartcut/main.go

private_key="../keys/id_rsa_smartcut"
public_key="keys/id_rsa.pub"
binary_name=smartcut_windows_amd64.exe
zip_name=binaries_windows.zip
folder_name=binaries_windows


rm -f $zip_name
rm -rf $folder_name

fyne package -app-id com.mouuff.smartcut -os windows -icon images/icon.png -release

mkdir $folder_name
mv SmartCut.exe $folder_name/$binary_name


rocket-update sign -key $private_key -path $folder_name
rocket-update verify -publicKey $public_key -path $folder_name

cd $folder_name
zip -r ../$zip_name *
cd -

