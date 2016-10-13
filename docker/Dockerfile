FROM golang

MAINTAINER LYC <imyikong@gmail.com>

WORKDIR /

RUN git clone https://github.com/MrLYC/ZHRSS.git && \
    make -C ZHRSS bootstrap && \
    make -C ZHRSS

EXPOSE 80

ENV URL https://www.zhihu.com/people/mr_lyc
ENV ADDR :80
ENV LOCATION Asia/Shanghai
ENV CACHE 600
ENV PATH /zhihu.xml

ENTRYPOINT ["/ZHRSS/bin/zhrss", "-url=${URL}", "-addr=${ADDR}", "-location=${LOCATION}", "-cache=${CACHE}", "-path=${PATH}"]

