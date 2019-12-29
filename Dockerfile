FROM golang:alpine as build

ENV TZ=Europe/Kiev
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk update update &&\
    apk add --update alpine-sdk linux-headers git zlib-dev openssl-dev gperf php php-ctype cmake &&\
    cd /tmp/ &&\
    git clone https://github.com/tdlib/td.git &&\
    cd td &&\
    mkdir build &&\
    cd build &&\
    export CXXFLAGS="" &&\
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=/usr/local .. &&\
    cmake --build . --target install &&\
    cd ../.. &&\
    ls -l /usr/local && rm -rf /tmp/td


WORKDIR ${GOPATH}/src/github.com/mikhno-s/eva_tg_bot

COPY . .
RUN go build  .

# FROM alpine as final
# COPY --from=build  /usr/local/ /usr/local/
# COPY --from=build /go/src/github.com/mikhno-s/eva_tg_bot /usr/local/bin

ENTRYPOINT [ "./eva_tg_bot" ]
