# Integrate TDengine with Intel EII

## Install Docker and Docker Compose (v1.28+)

Please refer to Docker [document](https://docs.docker.com/get-docker/) for detail instructions.

## Provision

```bash
cd tdengine/provision
sudo -E ./provision.sh ../docker-compose.yml
```



## Get Docker Images

```bash
cd tdengine
docker-compose pull
```



## Run the  Solution

```bash
docker-compose up -d
```



## Import Grafana dashboard

- Visit to http://localhost:3000 in browser. 
- Login as user 'admin' with password 'admin'
- Import dashboard from **tdengine/grafana/dashboard/dashboard.json**

