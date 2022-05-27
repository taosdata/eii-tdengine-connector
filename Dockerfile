# Dockerfile for TDengineConnector

ARG EII_VERSION
ARG UBUNTU_IMAGE_VERSION
ARG ARTIFACTS="/artifacts"
ARG TDENGINE_VERSION
FROM ia_common:${EII_VERSION} as common
FROM ia_eiibase:${EII_VERSION} as builder
LABEL description="TDengineConnector image"
WORKDIR /eii
ADD . /eii
# Install TDengine
ENV TDENGINE_PACKAGE="TDengine-server-2.4.0.20-Linux-x64.deb"
RUN dpkg -i ${TDENGINE_PACKAGE}
RUN rm ${TDENGINE_PACKAGE}

# Install TDengine go driver 
RUN mkdir -p ${GOPATH}/src/github.com/taosdata/driver-go/
RUN mv ./driver-go-2.0.0 ${GOPATH}/src/github.com/taosdata/driver-go/v2

# Copy Dependencies
WORKDIR ${GOPATH}/src/IEdgeInsights
ARG CMAKE_INSTALL_PREFIX
ENV CMAKE_INSTALL_PREFIX=${CMAKE_INSTALL_PREFIX}
COPY --from=common ${CMAKE_INSTALL_PREFIX}/include ${CMAKE_INSTALL_PREFIX}/include
COPY --from=common ${CMAKE_INSTALL_PREFIX}/lib ${CMAKE_INSTALL_PREFIX}/lib
COPY --from=common /eii/common/util/util.go common/util/util.go
COPY --from=common ${GOPATH}/src ${GOPATH}/src
COPY --from=common /eii/common/libs/EIIMessageBus/go/EIIMessageBus $GOPATH/src/EIIMessageBus
COPY --from=common /eii/common/libs/ConfigMgr/go/ConfigMgr $GOPATH/src/ConfigMgr


# Compile TDengineConnector
ENV PATH="$PATH:/usr/local/go/bin" \
    PKG_CONFIG_PATH="$PKG_CONFIG_PATH:${CMAKE_INSTALL_PREFIX}/lib/pkgconfig" \
    LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:${CMAKE_INSTALL_PREFIX}/lib"

ENV CGO_CFLAGS="$CGO_FLAGS -I ${CMAKE_INSTALL_PREFIX}/include -O2 -D_FORTIFY_SOURCE=2 -Werror=format-security -fstack-protector-strong -fPIC" \
    CGO_LDFLAGS="$CGO_LDFLAGS -L${CMAKE_INSTALL_PREFIX}/lib -z noexecstack -z relro -z now"

WORKDIR /eii
RUN go build TDengineConnector.go 

FROM ubuntu:${UBUNTU_IMAGE_VERSION}
ARG CMAKE_INSTALL_PREFIX
WORKDIR /eii
ENV TDENGINE_PACKAGE="TDengine-server-2.4.0.20-Linux-x64.deb"
ADD ./${TDENGINE_PACKAGE} /eii
RUN dpkg -i /eii/${TDENGINE_PACKAGE}
RUN rm /eii/${TDENGINE_PACKAGE}
COPY --from=builder /eii/startup.sh /eii/startup.sh
COPY --from=builder /eii/TDengineConnector /eii/TDengineConnector
COPY --from=common ${CMAKE_INSTALL_PREFIX}/lib/libeii*.so ${CMAKE_INSTALL_PREFIX}/lib/
COPY --from=common ${CMAKE_INSTALL_PREFIX}/lib/libsafestring.so ${CMAKE_INSTALL_PREFIX}/lib/
COPY --from=common ${CMAKE_INSTALL_PREFIX}/lib/libzmq.so.5 ${CMAKE_INSTALL_PREFIX}/lib/
COPY --from=common ${CMAKE_INSTALL_PREFIX}/lib/lib*.so.1 ${CMAKE_INSTALL_PREFIX}/lib/
EXPOSE 6030-6042/tcp 
EXPOSE 6030-6042/udp 
ENV PATH="$PATH:/usr/local/go/bin" \
	PKG_CONFIG_PATH="$PKG_CONFIG_PATH:${CMAKE_INSTALL_PREFIX}/lib/pkgconfig" \
	LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:${CMAKE_INSTALL_PREFIX}/lib"
RUN chmod a+x ./startup.sh
ENTRYPOINT ["./startup.sh"]
