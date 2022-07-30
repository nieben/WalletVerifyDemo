# WalletVerifyDemo

#### Init
```
git clone git@github.com:nieben/WalletVerifyDemo.git
cd WalletVerifyDemo
go mod tidy
```

#### Start Server: 
```
# The default serve port is 1323, you can change it by adding WALLET_VERIFY_PORT environment variables
go build
WalletVerifyDemo
```

#### Start Client: 
```
go run client/client.go
```