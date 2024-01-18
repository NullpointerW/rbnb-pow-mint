rbnb-pow-mint-tool
-----------------------
## install

```shell
git clone 
go build 
```

## Installation-free: run with docker

cd to project folder with dockerfile
```bash
docker build -t rbnb-pow-mint .
```

Optional: Run the container with limited computational resource.
```bash
docker run -it -m 50g --cpus=20 rbnb-pow-mint
```

## usage
输入钱包地址并设置难度,当每个钱包打到4500时，会自动生产新的钱包继续挖掘，钱包私钥保存在目录下wal.txt下
