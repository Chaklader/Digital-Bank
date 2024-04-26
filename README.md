## DIGITAL BANK 

## Database

Enter inside the DB `simple_bank` in the docker container:

```shell
docker exec -it postgres psql -U root -d simple_bank
```

## Deploying the App to Production 

Before we proceed to the deployment in the production, we need to create a Dockerfile and build the docker image with the command:

```shell
docker build -t simplebank:latest . 
```

As the Postgres docker image is build with the `bank-network`, we can run the image with the command provide below:

```shell
docker run --name simplebank -p 8080:8080 --network=bank-network -e DB_SOURCE="postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable"   -d simplebank:latest
```

```shell
aws ecr get-login-password \
    --region us-east-1 \
| docker login \
    --username AWS \
    --password-stdin 366655867831.dkr.ecr.us-east-1.amazonaws.com
```


```shell
docker run -p 8080:8080 366655867831.dkr.ecr.us-east-1.amazonaws.com/digitalbank:fe59fc55e3d8ac8f585e3aa5f2471f1ce6f2b3b6
```

