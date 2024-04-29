# DIGITAL BANK


## Deployment 

<br>

### AWS: Users, Roles and Groups

<br>

An IAM user represents an individual or application that interacts with AWS services, and each user has unique security 
credentials used to authenticate and access AWS resources. Users can be assigned permissions directly or through group 
memberships, and they are typically used for individuals or applications that require long-term access to AWS resources.
We create a new user named `GitHUB-CI` along with `chaklader` and will use that to manage the deployment process. 


<br>

![AWS Users - Normal Ops and Deployment](images/AWS_Users.png)

<br>

![AWS_GitHUB_CI_User](images/AWS_GitHUB_CI_User.png)


<br>

After creating AWS users, it is necessary to install the AWS CLI (Command Line Interface) for managing operations through 
the command line. Additionally, you will need to configure the access keys as outlined in the following steps:

<br>

```shell
$ aws configure 
$ cat ~/.aws/credentials 
[default]
aws_access_key_id = XXXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXXX

[github]
aws_access_key_id = XXXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXXX
```

<br>

```shell
$ ls -l ~/.aws/
total 16
-rw-------  1 chaklader  staff   43 Apr 26 15:02 config
-rw-------  1 chaklader  staff  351 Apr 27 08:37 credentials
```

```shell
$ cat ~/.aws/config 
[default]
region = us-east-1
output = json
```


In my case, I use the IAM user `chaklader` to log in to the AWS console and create the resources like `RDS`, `ECR`, `EKS` etc
and hence, I use the same credentials as the `default` section of the above `config` file. 


An IAM group is a collection of IAM users, and groups are used to simplify permissions management by assigning permissions 
to a group rather than individual users. Users inherit the permissions assigned to the groups they belong to, and groups help 
organize users based on their roles or responsibilities within an organization. The `Deployment` user group permission is 
provided below and the user `GitHUB-CI` needs to be in the `Deployment` user group. I also put IAM user `chaklader` in the same
group for testing purpose. 

<br>

![Deployment_User_Group_Permissions](images/Deployment_User_Group_Permissions.png)

<br>

The `DeploymentGroupEKSPolicy` is described in the `Deployment` user group in AWS dashboard:

<br>

![AWS_DeploymentGroupEKSPolicy_Permission](images/AWS_DeploymentGroupEKSPolicy_Permission.png)

<br>


### RDS

<br>

To create an AWS Postgres RDS instance named `digital-bank`, log in to the AWS Management Console, navigate to the RDS 
service, and click on `Create database`. Choose `Standard Create`, select `PostgreSQL` as the engine, and specify the 
desired version (e.g., `PostgreSQL 16.x-R1`). Configure the DB instance size (e.g., `db.t3.micro`), set the DB instance 
identifier to `digital-bank`, and provide a master username and password. Choose the desired VPC (e.g., `default-vpc-0f6cf7d178eb0c8d8`) 
and subnet group (e.g., `default-vpc-0f6cf7d178eb0c8d8`), and create a new security group or select an existing one. Set 
the public accessibility option to `Yes` if needed. Configure additional settings such as backup retention period, 
maintenance window, and encryption options. Click on `Create database` to initiate the RDS instance creation process. 
Once created, retrieve the connection details (endpoint, port, username, and password) from the AWS Management Console. 

<br>

![alt text](images/RDS.png)

<br>

Configure the security group inbound rules to allow traffic on the PostgreSQL port (default: `5432`) from the desired IP 
range or security group. Finally, update the application's database configuration to use the provided RDS connection details.
Create the AWS Postgres DB and test it with Table plus if the remote connection is working. We need SSL enabled for the 
connection

<br>

![alt text](images/RDS_SG_Inbound_Rules.png)

<br>

We have application secrets provided in the `app.env` file that we need to run server locally, and we need to save these 
secrets in the AWS secrets manager with the production values as below. 

<br>

![alt text](images/Secret_Manager_Credentials.png)

<br>


Initially, the `SecretManagerReadWrite` policy was not included in the `deployment` user group and hence, we received the 
error below:

```shell
$ aws secretsmanager get-secret-value --secret-id digital_bank

"An error occurred (AccessDeniedException) when calling the GetSecretValue operation: User: arn:aws:iam::095420225548:user/github-ci is 
not authorized to perform: secretsmanager:GetSecretValue on resource: digital_bank
```

After we add the permission for the AWS Secret Manager for the GitHUb-CI user using the `deployment` group, we can read the 
secrets as below:


```shell
$ aws secretsmanager get-secret-value --secret-id Digital_Bank
{
    "ARN": "arn:aws:secretsmanager:us-east-1:366655867831:secret:Digital_Bank-UZysxN",
    "Name": "Digital_Bank",
    "VersionId": "eb67f52e-541f-4b30-8dd1-f521432411ea",
    "SecretString": "{\"DB_SOURCE\":\"postgresql://root:OIJIWTiG508B54n88kA7@digital-bank.czzl3uwtdaas.us-east-1.rds.amazonaws.com:5432/digital_bank\",\"DB_DRIVER\":\"postgres\",\"HTTP_SERVER_ADDRESS\":\"0.0.0.0:8080\",\"GRPC_SERVER_ADDRESS\":\"0.0.0.0:9090\",\"TOKEN_SYMMETRIC_KEY\":\"48924940a30b055c3e01a873d05fcec7\",\"MIGRATION_URL\":\"file://db/migration\",\"REDIS_ADDRESS\":\"0.0.0.0:6379\",\"EMAIL_SENDER_NAME\":\"Digital_Bank\",\"EMAIL_SENDER_ADDRESS\":\"digitalbanktest@gmail.com\",\"EMAIL_SENDER_PASSWORD\":\"jekfcygyenvzekke\"}",
    "VersionStages": [
        "AWSCURRENT"
    ],
    "CreatedDate": "2024-04-26T17:02:13.506000+06:00"
}
```

We need to provide these info in the `app.env` as part of the deployment procedure as the same format with the command 
below that will be included in the `deployment.yaml` pipeline. 

<br>

```shell
$ aws secretsmanager get-secret-value --secret-id Digital_Bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]'

DB_SOURCE=postgresql://root:XXXXXXXXX@digital-bank.czzl3uwtdaas.us-east-1.rds.amazonaws.com:5432/digital_bank
DB_DRIVER=postgres
HTTP_SERVER_ADDRESS=0.0.0.0:8080
GRPC_SERVER_ADDRESS=0.0.0.0:9090
TOKEN_SYMMETRIC_KEY=XXXXXXXXX
MIGRATION_URL=file://db/migration
REDIS_ADDRESS=0.0.0.0:6379
EMAIL_SENDER_NAME=Digital_Bank
EMAIL_SENDER_ADDRESS=digitalbanktest@gmail.com
EMAIL_SENDER_PASSWORD=jekfcygyenvzekke
```

<br>


### GitHub Repository Secrets 

<br>

GitHub repo secrets are encrypted environment variables that you can store in a repository on GitHub. These secrets can contain
sensitive information, such as access tokens, API keys, or other credentials, and can be used in GitHub Actions workflows
to securely access external services or resources. Secrets are securely encrypted and can only be accessed by the repository
they are stored in, ensuring that sensitive data is kept confidential and cannot be accessed by unauthorized parties. Using
secrets in your workflows allows you to avoid hard-coding sensitive information in your code, which can improve the security
and maintainability of your projects. We need to set the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` GitHUB repository  
secrets from the settings page of the repository.

<br>

![alt text](images/Github_Repo_Secrets_AWS.png)

<br>

### Dockerfile

<br>

As part of the deployment procedure, we create a `Dockerfile` and put that inside the root of the project. This Dockerfile 
has `2` stages: the `build` stage and the `run` stage. The `build` stage starts from the `golang:1.20-alpine3.19` base image, 
sets the working directory to `/app`, copies the source code, and runs `go build` to compile the Go application into an 
executable binary named `main`. The run stage starts from the `alpine:3.19` base image, copies the compiled `main` 
executable and other files (`app.env`, `start.sh`, `wait-for.sh`, and the `db/migration` directory) to the `/app` directory, 
exposes ports 8080 and 9090, sets the default command to run the `main` executable, and sets the entrypoint to execute the `start.sh` 
script when the container starts. This Dockerfile is likely used for building and running a Go application in a containerized 
environment.

<br>


### Elastic Container Registry (ECR)

<br>

AWS Elastic Container Registry (ECR) is a fully-managed Docker container registry service provided by Amazon Web Services (AWS). 
It allows you to store, manage, and deploy Docker container images securely and efficiently. With ECR, you can push, pull, 
and manage your Docker images from anywhere using the Docker command-line interface or your preferred continuous integration 
and continuous deployment (CI/CD) tools. ECR integrates seamlessly with other AWS services, such as Amazon Elastic Kubernetes 
Service (EKS), Amazon Elastic Container Service (ECS), and AWS Lambda, making it a convenient choice for containerized 
application deployments on the AWS platform. For the deployment purpose, we created a ECR repository named `digitalbank` 
in the ECR as shown below. 

<br>

![alt text](images/ECR.png)

<br>


### GitHub Actions 

<br>

GitHub Actions is a continuous integration and continuous deployment (CI/CD) platform provided by GitHub. It allows you to 
automate your software development workflows, such as building, testing, and deploying your code directly from your GitHub 
repository. You can define custom workflows using YAML syntax, which specify the events that trigger the workflow, the 
jobs to be executed, and the steps within each job. We use GitHub Actions to manage the testing and deployment procedure 
in the AWS cloud infrastructure. 

<br>

### Testing 

<br>

We have GitHub Action workflow created in the file named `.github/workflows/test.yaml` which is designed to run unit tests for 
a Go project. It is triggered whenever code is pushed to the main branch or a pull request is created against the main branch. 
The workflow runs on an Ubuntu runner and includes a PostgreSQL service for testing purposes. The steps involve setting up the 
Go environment, checking out the code, installing the golang-migrate tool, running database migrations using make migrateup, 
and finally executing the unit tests using make test. This workflow helps ensure that the codebase passes all tests before 
merging changes into the main branch.

<br>

## Deployment 

Create a repository in the AWS ECR and run the  `.github/workflows/deploy.yaml` to push the image to the ECR repo. Now the image is ready 
and we can pull after login to the docker. We need to Login to the docker before we can pull the image:

```shell
aws ecr get-login-password \
    --region us-east-1 \
| docker login \
    --username AWS \
    --password-stdin 366655867831.dkr.ecr.us-east-1.amazonaws.com
```

Run the image locally to test it:


```shell
docker run -p 8080:8080 366655867831.dkr.ecr.us-east-1.amazonaws.com/digitalbank:latest
```



## Create EKS cluster


![alt text](images/eks_cluster.png)


Creating the Amazon EKS cluster role 

![alt text](images/eks_cluster_service_ROLE.png)

To create your Amazon EKS cluster role in the IAM console


1. Open the IAM console 
2. Choose Roles, then Create role.
3. Under Trusted entity type, select AWS service.
4. From the Use cases for other AWS services dropdown list, choose EKS.
5. Choose EKS - Cluster for your use case, and then choose Next.
6. On the Add permissions tab, choose Next.
7. For Role name, enter a unique name for your role, such as eksClusterRole.
8. For Description, enter descriptive text such as Amazon EKS - Cluster role.
9. Choose Create role.


- Add worker node to the cluster 

![alt text](images/eks_add_worker_node.png)


- Add EKS Node groups and we need to create a new IAM role for this - The EC2 instance used decide how many pods we can run there

![alt text](images/worker_node_iam_group.png)

We need 3 plocies for the IAM role:

```shell
AmazonEKS_CNI_Policy
AmazonEKSWorkerNodePolicy
AmazonEC2ContainerRegistryReadOnly
```

![alt text](images/AWSEKSNodeRole.png)

<br>


```shell
$ kubectl cluster-info
Kubernetes control plane is running at https://E44AED5442512EC56EA2BFBD88920895.gr7.us-east-1.eks.amazonaws.com
CoreDNS is running at https://E44AED5442512EC56EA2BFBD88920895.gr7.us-east-1.eks.amazonaws.com/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```


```shell
$ aws eks update-kubeconfig --name digital-bank --region us-east-1
Updated context arn:aws:eks:us-east-1:366655867831:cluster/digital-bank in /Users/chaklader/.kube/config
```

```shell
$ kubectl config use-context arn:aws:eks:us-east-1:366655867831:cluster/digital-bank
Switched to context "arn:aws:eks:us-east-1:366655867831:cluster/digital-bank".
```

```shell
$ cat ~/.kube/config 
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: XXXXXXXX
    server: https://E44AED5442512EC56EA2BFBD88920895.gr7.us-east-1.eks.amazonaws.com
  name: arn:aws:eks:us-east-1:366655867831:cluster/digital-bank
contexts:
- context:
    cluster: arn:aws:eks:us-east-1:366655867831:cluster/digital-bank
    user: arn:aws:eks:us-east-1:366655867831:cluster/digital-bank
  name: arn:aws:eks:us-east-1:366655867831:cluster/digital-bank
current-context: arn:aws:eks:us-east-1:366655867831:cluster/digital-bank
kind: Config
preferences: {}
users:
- name: arn:aws:eks:us-east-1:366655867831:cluster/digital-bank
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      args:
      - --region
      - us-east-1
      - eks
      - get-token
      - --cluster-name
      - digital-bank
      - --output
      - json
      command: aws
```

Provide the EKS full permission for the GitHUB-CI user as this user needs to access the Kubernetes Cluster and manage the deployment 


```shell
$ aws sts get-caller-identity

{
    "UserId": "AIDAVKXTERO32BYHSFV6L",
    "Account": "366655867831",
    "Arn": "arn:aws:iam::366655867831:user/GitHUB-CI"
}
```



$ kubectl get pods

$ export AWS_PROFILE=github
Chakladers-MacBook-Pro:Desktop chaklader$ kubectl get pods

Chakladers-MacBook-Pro:Desktop chaklader$ export AWS_PROFILE=default
Chakladers-MacBook-Pro:Desktop chaklader$ kubectl get pods



Provide the EKS full permission for the GitHUB-CI user as this user needs to access the Kubernetes Cluster and manage the deployment.
First, use the default user and then apply the YAML to the cluster as below:


```yaml
apiVersion: v1 
kind: ConfigMap 
metadata: 
  name: aws-auth 
  namespace: kube-system 
data: 
  mapUsers: | 
    - userarn: arn:aws:iam::366655867831:user/GitHUB-CI
      username: GitHUB-CI
      groups:
        - system:masters
```



```shell
$ kubectl apply -f eks/aws-auth.yaml
configmap/aws-auth unchanged

$ kubectl get pods
```

```shell
$ brew install k9s 
```

In the K9s console, use:

```shell
$ configmap
```

In Kubernetes, a Service is a method for exposing a network application that is running as one or more Pods in your
cluster. A key aim of Services in Kubernetes is that you don't need to modify your existing application to use an unfamiliar
service discovery mechanism. You can run code in Pods, whether this is a code designed for a cloud-native world, or an older app
you've containerized. You use a Service to make that set of Pods available on the network so that clients can interact with it.



```shell
$ kubectl apply -f eks/deployment.yaml
deployment.apps/digital-bank-api-deployment created
```


It informs about the maximum number of pods that can be supported by different EC2 instance types when using Amazon EKS (Elastic 
Kubernetes Service). It explains the formula used to calculate the maximum number of pods for each EC2 instance type. The formula is:

```shell
# of ENI * (# of IPv4 per ENI - 1) + 2
```

The file provides a link to the AWS documentation for more information on using the formula.

Finally, the file lists various EC2 instance types and their corresponding maximum number of pods that can be supported. For example, a1.2xlarge can support 58 pods, while c3.4xlarge can support 234 pods.

Delete the existing deployments in the k9s and then <d>


```shell
$ deployments 
$ services
```


![alt text](images/ip_addresses.png)
![alt text](images/ip_addresses_1.png)

In order to access the Kubernetes resources from the outside, we need to deploy the service as below:

```shell
$ kubectl apply -f eks/service.yaml
service/digital-bank-api-service configured
```


![alt text](images/service.png)

Purchase domain name from the AWS Route 53

![alt text](images/default.png)

Route 53 -> Hosted Zones -> Records -> Domain Name and we will find the list of DNS records (NS and SOA)

We will already have 2 records - NS and SOA and need to create an A record

![alt text](images/records.png)

<br>

![alt text](images/before_A_record.png)

<br>

![alt text](images/ns_lookup_before_A-record.png)

Use the LOad balancer IP in the A record

![alt text](images/create_A_Record.png)

![alt text](images/create_A_Record_2.png)

![alt text](images/ns_lookup_after_A-record.png)

Create Record - A record and use the EXTERNAL-IP for the service for the route to traffic field

Ingress

An API object that manages external access to the services in a cluster, typically HTTP.

Ingress may provide load balancing, SSL termination and name-based virtual hosting.

Ingress exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is
controlled
by rules defined on the Ingress resource.

Before proceeding to the Ingress, change the type of the `service.yaml` to ClusterIP from LoadBalancer

We had this before `service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: digital-bank-api-service
spec:
  selector:
    app: digital-bank-api
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

We need to change it to:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: digital-bank-api-service
spec:
  selector:
    app: digital-bank-api
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
```

The `ingress.yaml` is provided below:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: digital-bank-ingress
spec:
  rules:
    #    This is from the A record 
    - host: "api.digital-bank.org"
      http:
        paths:
          - pathType: Prefix
            path: "/"
            backend:
              service:
                #  This comes from the metadata of the service
                name: digital-bank-api-service
                port:
                  #                  the port is from the service.yaml file 
                  number: 80
```

```shell

$ kubectl apply -f eks/service.yaml
 
$ kubectl apply -f eks/ingress.yaml
```

We have ClusterIP in the deployed service NOW - look for service and ingress in the K9s console

![alt text](images/service_ClusterIP.png)

<br>

![alt text](images/ingress_initial_deploy.png)

Ingress is sending the traffic to the `digital-bank-api-service`:

![alt text](images/ingress_send_traffic_service.png)

The ingress doesn't have external traffic to the domain and need to update the A record for that:

Just ingress is not enough and we need to set up an ingress controller - we use Nginx ingress controller

Nginx Ingress Controller

![alt text](images/nginx_ingress_controller.png)

![alt text](images/nginx_ingress_controller_RUN.png)

```shell
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.1/deploy/static/provider/aws/deploy.yaml
```

![alt text](images/ingress_class.png)

Copy the address above and paste to the A-record in the AWS route 53

In the K9s console check `ingressclass`

Provide the address of the Ingress to the A-record

![alt text](images/A-record_update_forIngress.png)

![alt text](images/ingress_update_check.png)

IngressClass resource in the Ingress YAML

![alt text](images/ingress_class_updated.png)

Update the ingress.yaml file:

```yaml
apiVersion: networking.k8s.io/v1
kind: IngressClass
metadata:
  name: nginx
spec:
  controller: k8s.io/ingress-nginx
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: digital-bank-ingress
spec:
  ingressClassName: nginx
  rules:
    - host: "api.digital-bank.org"
      http:
        paths:
          - pathType: Prefix
            path: "/"
            backend:
              service:
                name: digital-bank-api-service
                port:
                  number: 80
```

$ kubectl apply -f eks/ingress.yaml

The class of the ingress is changed to Nginx

Make the Client/Server communication secure using TLS

URL YT: <https://youtu.be/-f4Gbk-U758>

SITE: <https://letsencrypt.org/>

Should only be use if the DNS provider has an API to update the records

![alt text](images/dns-01_challenge.png)

HTTP 01 Challenge

![alt text](images/http-01_challenge.png)

Install Kubernetes cert manager

```shell
$ kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.5/cert-manager.yaml
```

![alt text](images/cert_manager_pods.png)

```shell
$ kubectl get pods --namespace cert-manager

NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-7ddd8cdb9f-bxlsn              1/1     Running   0          22h
cert-manager-cainjector-57cd76c845-2lq2b   1/1     Running   0          22h
cert-manager-webhook-cf8f9f895-8c7bd       1/1     Running   0          22h

```

![alt text](images/create_Basic_ACME_Issuer.png)

Now, deploy the `eks/issuer.yaml` to the Kubernetes:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt
spec:
  acme:
    email: omi.chaklader@gmail.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      # Secret resource that will be used to store the account's private key.
      name: letsencrypt-account-private-key
    # Add a single challenge solver, HTTP01 using nginx
    solvers:
      - http01:
          ingress:
            ingressClassName: nginx
```

In the K9s console, check for the `>Clusterissuer`
In the K9s console, check for the `>secrets` for its private keys

The certificates are still empty:

In the K9s console, check for the `>certificate` for its private keys
In the K9s console, check for the `>certificaterequest` for its private keys

Update the ingress to find the certificate:

```yaml
spec:
  controller: k8s.io/ingress-nginx
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: digital-bank-ingress
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt
spec:
  ingressClassName: nginx
  rules:
    - host: "api.digital-bank.org"
      http:
        paths:
          - pathType: Prefix
            path: "/"
            backend:
              service:
                name: digital-bank-api-service
                port:
                  number: 80
  tls:
    - hosts:
        - api.digital-bank.org
      secretName: digital-bank-api-cert
```

Now we can check that the TLS is enabled

![alt text](images/tls.png)

<br>

![alt text](images/certificates.png)

<br>

![alt text](images/all.png)

![alt text](images/443.png)

<br>

![alt text](images/request.png)







