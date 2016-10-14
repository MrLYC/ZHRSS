FROM alpine

MAINTAINER LYC <imyikong@gmail.com>

ADD entry.sh /entry.sh

RUN apk update && \
    apk add go git make && \
    mkdir /app/ && \
    git clone https://github.com/MrLYC/ZHRSS.git /app/ZHRSS && \
    make -C /app/ZHRSS/ && \
    mv /app/ZHRSS/bin/zhrss /app/zhrss && \
    apk del go git make && \
    rm -r /app/ZHRSS/

EXPOSE 80

ENV ADDR :80
ENV URL https://www.zhihu.com/people/mr_lyc
ENV LOCATION Asia/Shanghai
ENV CACHE 600
ENV PATH /

ENTRYPOINT ["/entry.sh"]
