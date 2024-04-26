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


