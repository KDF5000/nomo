FROM golang:1.16

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' > /etc/timezone

WORKDIR /opt/nomo
COPY . ./


# build a release version
ENV WORKDIR=/opt/nomo/output
RUN ./build.sh

CMD '/opt/nomo/output/run_wx.sh'; 'bash'
