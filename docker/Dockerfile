FROM golang

MAINTAINER LYC <imyikong@gmail.com>

RUN mkdir /app/ && \
    git clone https://github.com/MrLYC/ZHRSS.git /app/ZHRSS && \
    make -C /app/ZHRSS/

ADD entry.sh /entry.sh

EXPOSE 80

ENV ADDR :80
ENV URL https://www.zhihu.com/people/mr_lyc
ENV LOCATION Asia/Shanghai
ENV CACHE 600
ENV PATH /

ENTRYPOINT ["/entry.sh"]
