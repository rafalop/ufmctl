# ufmctl
cli for interacting with UFM API

# building
## install go (ubuntu|mac)
```
snap|brew install go
```
You can also install using upstream tarball etc.

## prep module
```
go mod init ufmctl
go mod tidy
```

## build
```
go build ufmctl
```

# run examples
```
./ufmctl -h
./ufmctl --endpoint https://10.0.0.1 --insecure --username root --password s3cr3t pkeys list
./ufmctl --endpoint https://10.0.0.1 --insecure --username root --password s3cr3t ports list
./ufmctl --endpoint https://10.0.0.1 --insecure --username root --password s3cr3t systems list
```

# notes
This will create a file (by default `ufm-cookies.txt`) in local dir to prevent re-authing, but will not delete it! fairly insecure if you don't delete the cookie after use, but loosely based on what nvidia has documented here https://docs.nvidia.com/networking/display/ufmenterpriserestapiv6152/rest+api+complementary+information ... a work in progress.
