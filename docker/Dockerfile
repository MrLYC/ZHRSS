FROM alpine

MAINTAINER LYC <imyikong@gmail.com>

ADD entry.sh /entry.sh

RUN apk update && \
    apk add git make go && \
    git clone https://github.com/MrLYC/ZHRSS.git && \
    make -C /ZHRSS/ && \
    mv /ZHRSS/bin/zhrss /bin/zhrss && \
    cp /usr/lib/go/lib/time/zoneinfo.zip / && \
    apk del make git go && \
    mkdir -p /usr/lib/go/lib/time/ && \
    mv /zoneinfo.zip /usr/lib/go/lib/time/zoneinfo.zip && \
    apk add ca-certificates && \
    rm -r /ZHRSS/

EXPOSE 80

ENV ADDR :80
ENV URL https://www.zhihu.com/people/mr_lyc
ENV LOCATION Asia/Shanghai
ENV CACHE 600
ENV PATH /

ENTRYPOINT ["/entry.sh"]
