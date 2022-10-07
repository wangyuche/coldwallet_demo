# coldwallet_demo

install ganache-cli
```
npm install -g ganache-cli
```
起一個模擬器在本機
```
ganache-cli
```
可用 Store.sol當測試
```
docker pull ethereum/solc:0.4.24
docker run --rm -v $(pwd)/contract:/root ethereum/solc:0.4.24 --abi --bin /root/Store.sol -o /root/build

docker pull ethereum/client-go:alltools-latest
docker run --rm -v $(pwd)/contract:/root ethereum/client-go:alltools-latest abigen --bin=/root/build/Store.bin --abi=/root/build/Store.abi --pkg=store --out=/root/Store.go
```
用來實作USDT
```
docker pull ethereum/solc:0.8.4
docker run --rm -v $(pwd)/contract:/root ethereum/solc:0.8.4 --abi --bin /root/ERC20.sol -o /root/build


docker pull ethereum/client-go:alltools-latest
docker run --rm -v $(pwd)/contract:/root ethereum/client-go:alltools-latest abigen --bin=/root/build/ERC20.bin --abi=/root/build/ERC20.abi --pkg=usdt --out=/root/ERC20.go
```