docker kill ia_tdengineconnector 
docker rm ia_tdengineconnector
docker run -p 6030-6042:6030-6042/tcp --name ia_tdengineconnector --network eii -d taos/ia_tdengineconnector:0.0.1 

