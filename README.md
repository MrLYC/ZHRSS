# ZHRSS

feed your zhihu time line.

github repo: [https://github.com/MrLYC/ZHRSS](https://github.com/MrLYC/ZHRSS)

docker hub: [https://hub.docker.com/r/mrlyc/zhrss/](https://hub.docker.com/r/mrlyc/zhrss/)



## Quick start

### Source code

```shell
git clone https://github.com/MrLYC/ZHRSS.git
cd ZHRSS
make
./bin/zhrss --help
.bin/zhrss -url https://www.zhihu.com/people/<user_id> -path /zhihu.xml
```

visit your time line feed details: [http://127.0.0.1:8080/zhihu.xml](http://127.0.0.1:8080/zhihu.xml).



### Docker 

```shell
docker run mrlyc/zhrss -p 8080:80 -e URL=https://www.zhihu.com/people/<user_id> -e PATH=/zhihu.xml
```

visit: http://127.0.0.1:8080/zhihu.xml



## DEMO

[刘奕聪 - 知乎](https://zhrss.arukascloud.io/)

