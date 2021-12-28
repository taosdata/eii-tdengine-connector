# Dockerfile for TDengineConnector

ARG EII_VERSION
ARG UBUNTU_IMAGE_VERSION
ARG ARTIFACTS="/artifacts"
FROM ia_common:${EII_VERSION} as common
FROM ia_eiibase:${EII_VERSION} as builder
LABEL description="TDengineConnector image"
WORKDIR /root
ADD . /root
RUN dpkg -i TDengine-server-2.3.5.0-beta-Linux-x64.deb
EXPOSE 6030-6042/tcp 
EXPOSE 6030-6042/udp 
#CMD ["taosd"]
#CMD ["./startup.sh"]
ENTRYPOINT ["./startup.sh"]

