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
1. Configure the `.env` file.
2. When each wallet reaches the specified limit, it will automatically generate a new wallet to continue mining. The wallet private keys are saved in the 'wal.txt' directory.






