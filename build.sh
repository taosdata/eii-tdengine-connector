. .env
docker build . -t taos/ia_tdengineconnector:0.0.1  --build-arg EII_VERSION=$EII_VERSION --build-arg CMAKE_INSTALL_PREFIX=$EII_INSTALL_PATH
