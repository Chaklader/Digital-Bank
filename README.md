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

```textmate
$ aws configure 
$ cat ~/.aws/credentials 
[default]
aws_access_key_id = XXXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXXX

[arefe]
aws_access_key_id = XXXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXXX

[github]
aws_access_key_id = XXXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXXX
```

<br>

```textmate
$ ls -l ~/.aws/
total 16
-rw-------  1 chaklader  staff   43 Apr 26 15:02 config
-rw-------  1 chaklader  staff  351 Apr 27 08:37 credentials
```

```textmate
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

```textmate
$ aws secretsmanager get-secret-value --secret-id digital_bank

"An error occurred (AccessDeniedException) when calling the GetSecretValue operation: User: arn:aws:iam::095420225548:user/github-ci is 
not authorized to perform: secretsmanager:GetSecretValue on resource: digital_bank
```

After we add the permission for the AWS Secret Manager for the GitHUb-CI user using the `deployment` group, we can read the 
secrets as below:


```textmate
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

```textmate
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


## Elastic Kubernetes Service (EKS)

<br>

AWS EKS (Elastic Kubernetes Service) is a managed service that simplifies deploying, managing, and operating Kubernetes
clusters on AWS. An EKS cluster consists of a control plane managed by AWS and worker nodes managed by the user. The
control plane is responsible for managing the Kubernetes cluster and scheduling applications to run on the worker nodes.
Worker nodes are Amazon EC2 instances that run containerized applications and are registered to the EKS cluster. The worker
nodes can be configured with the desired compute, storage, and networking resources to suit the application's needs. AWS
EKS integrates with other AWS services like VPC, IAM, and Elastic Load Balancing to provide a secure and highly available
environment for running containerized workloads. For the deployment, we created a cluster named `digital-bank` as shown
below.

<br>

![alt text](images/EKS_Cluster.png)

<br>

To create the EKS cluster, we need to create an EKS cluster service role which is an IAM role that needs to be created and
associated with the EKS cluster during its creation. This role allows the EKS control plane to manage and interact with
other AWS services on your behalf, such as creating and managing worker nodes, managing load balancers, and accessing other
resources required for the proper functioning of the EKS cluster. It is recommended to create a dedicated role with the
necessary permissions for the EKS cluster to ensure secure and controlled access to AWS resources.

<br>

![alt text](images/EKS_Cluster_Service_Role.png)

<br>

To create the EKS Open the IAM console and navigate to the Roles section, then select Create role. Under Trusted entity
type, choose AWS service. From the Use cases for other AWS services dropdown list, select EKS, and for your use case,
choose EKS - Cluster, then click Next. On the Add permissions tab, proceed to the next step by clicking Next again. For
Role name, enter a unique name for your role, such as eksClusterRole. Provide a descriptive text for the Description field,
for example, `Amazon EKS - Cluster role`. Finally, review the details and choose Create role to create the Amazon EKS cluster
role in the IAM console.

<br>

Now, we need to add the EKS worker nodes to the cluster - EKS worker nodes are the compute resources within an Amazon EKS
cluster that run the containerized applications and workloads. These worker nodes are Amazon EC2 instances that are registered
to the EKS cluster and managed by the control plane. The worker nodes are responsible for running the Kubernetes node
components, such as the kubelet and kube-proxy, which facilitate communication between the control plane and the containers
running on the nodes. Workers nodes can be configured with the desired compute, storage, and networking resources based
on the application's requirements. They can be launched as part of an Auto Scaling group to automatically scale the cluster
capacity based on the workload demands. We can use the `Add Node Group` button below to append the worker nodes.


<br>

![alt text](images/EKS_Add_Worker_Node.png)

<br>

For the EKS worker node group, we need to create a new IAM role named `AWSEKSNodeRole` with the following policies:

1. AmazonEKS_CNI_Policy
2. AmazonEKSWorkerNodePolicy
3. AmazonEC2ContainerRegistryReadOnly

<br>

![alt text](images/EKS_Worker_Node_IAM.png)

<br>

![alt text](images/EKS_Worker_Node_IAM_APPEND.png)

<br>

When configuring an Amazon EKS worker node group, the EC2 instance type selected plays a crucial role in determining the
number of pods that can be run on each worker node. The instance type defines the maximum network interfaces and the number
of private IPv4 and IPv6 addresses per interface, as shown in the provided table. Additionally, the "eni-max-pods.txt" file
from the AWS labs repository contains a mapping that calculates the maximum number of pods per instance type based on the
available network interfaces and IP addresses. This information helps in selecting the appropriate instance type and configuring
the desired, minimum, and maximum number of worker nodes in the node group to meet the application's resource requirements
and ensure efficient pod scheduling within the EKS cluster.

It informs about the maximum number of pods that can be supported by different EC2 instance types when using Amazon EKS (Elastic
Kubernetes Service). It explains the formula used to calculate the maximum number of pods for each EC2 instance type. The formula is:

```textmate
# of ENI * (# of IPv4 per ENI - 1) + 2
```

The file provides a link to the AWS documentation for more information on using the formula. Finally, the file lists various EC2
instance types and their corresponding maximum number of pods that can be supported. For example, `a1.2xlarge` can support `58` pods,
while `c3.4xlarge` can support `234` pods.

<br>

![alt text](images/EKS_Pods_Per_EC2_Instance.png)

<br>

![alt text](images/EKS_Pods_Per_EC2_Instance_APPEND.png)

<br>


### Testing via GitHub Actions 

<br>

We have GitHub Action workflow created in the file named `.github/workflows/test.yaml` which is designed to run unit tests for 
a Go project. It is triggered whenever code is pushed to the main branch or a pull request is created against the main branch. 
The workflow runs on an Ubuntu runner and includes a PostgreSQL service for testing purposes. The steps involve setting up the 
Go environment, checking out the code, installing the golang-migrate tool, running database migrations using make migrateup, 
and finally executing the unit tests using make test. This workflow helps ensure that the codebase passes all tests before 
merging changes into the main branch.

<br>

## Deployment via GitHub Actions


We have GitHub Action workflow created in the file named `.github/workflows/deploy.yaml` which is designed to run unit tests for
a Go project. The initial version is to create the docker image in the ECR and provided below:


```yaml

name: Deploy to PROD

on:
  push:
    branches: [ main ]

jobs:

  deploy:
    name: Build image
    runs-on: ubuntu-latest

    steps:
      - name: Check out code
        uses: actions/checkout@v2

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v1
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v1

      - name: Load secrets and save to app.env
        run: aws secretsmanager get-secret-value --secret-id Digital_Bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env

      - name: Build, tag, and push image to Amazon ECR
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          ECR_REPOSITORY: digitalbank
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG -t $ECR_REGISTRY/$ECR_REPOSITORY:latest .
          docker push -a $ECR_REGISTRY/$ECR_REPOSITORY

```

After we push the code to the `main` branch, the CI/CD pipeline will run and docker image will be pushed to the ECR as below.

<br>

![alt text](images/ECR_Docker_Image.png)

<br>

Now the image is ready, and we can pull after login to the docker. We need to log in to the docker before we can pull the image:

```textmate
$ aws ecr get-login-password \
    --region us-east-1 \
| docker login \
    --username AWS \
    --password-stdin 366655867831.dkr.ecr.us-east-1.amazonaws.com
```

Run the image locally and can test it using the Postman:  

```textmate
$ docker pull 366655867831.dkr.ecr.us-east-1.amazonaws.com/digitalbank:latest
$ docker run -p 8080:8080 366655867831.dkr.ecr.us-east-1.amazonaws.com/digitalbank:latest
```

Now, we can check the AWS user with the command and see the current username is `GitHUB-CI`

```shell
$ aws sts get-caller-identity

{
    "UserId": "AIDAVKXTERO32BYHSFV6L",
    "Account": "366655867831",
    "Arn": "arn:aws:iam::366655867831:user/GitHUB-CI"
}
```

As we didn't use this user to create the EKS resources, we can't retrieve the cluster info and pods info using the commands
provided below.  

<br>

```textmate
$ kubectl cluster-info

$ kubectl get pods
```

<br>

So, we need to revert to the original user `chaklader` that created these resources and then, we can find these info. 
We need to update the configuration and use the correct context before we proceed. 

<br>

```textmate
$ aws eks update-kubeconfig --name digital-bank --region us-east-1
Updated context arn:aws:eks:us-east-1:366655867831:cluster/digital-bank in /Users/chaklader/.kube/config
```

<br>

```textmate
$ kubectl config use-context arn:aws:eks:us-east-1:366655867831:cluster/digital-bank
Switched to context "arn:aws:eks:us-east-1:366655867831:cluster/digital-bank".
```

Now, we can find the updated configuration, cluster info and pods with the commands provided below.

<br>

```textmate
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

<br>

```textmate
$ kubectl cluster-info
Kubernetes control plane is running at https://E44AED5442512EC56EA2BFBD88920895.gr7.us-east-1.eks.amazonaws.com
CoreDNS is running at https://E44AED5442512EC56EA2BFBD88920895.gr7.us-east-1.eks.amazonaws.com/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

<br>

```textmate
$ kubectl get pods
NAME                                           READY   STATUS    RESTARTS   AGE
cm-acme-http-solver-gxzpw                      1/1     Running   0          2d3h
digital-bank-api-deployment-57c7b975cf-g29j7   1/1     Running   0          2d6h
digital-bank-api-deployment-57c7b975cf-p69qr   1/1     Running   0          2d6h
```


Let's create a file `eks/aws-auth.yaml` to provide the EKS full permission for the GitHUB-CI user as this user needs to 
access the Kubernetes Cluster and manage the deployment. 

<br>

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

To do that, we need to use the default user `chaklader` and then apply the `eks/aws-auth.yaml` to the cluster as below:



<br>

```textmate

$ export AWS_PROFILE=chaklader

$ kubectl apply -f eks/aws-auth.yaml
```

Now, the user `GitHUB-CI` is ready to manage the deployment process and the cluster info and pods can be seen in the CLI 
as before with the default user. We run the `eks/deployment.yaml` for deployment as provided. We can see the deployment 
the `K9s` console as provided. 

```textmate
$ kubectl apply -f eks/deployment.yaml

deployment.apps/digital-bank-api-deployment created
```

<br>

```yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  name: digital-bank-api-deployment
  labels:
    app: digital-bank-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: digital-bank-api
  template:
    metadata:
      labels:
        app: digital-bank-api
    spec:
      containers:
        - name: digital-bank-api
          image: 366655867831.dkr.ecr.us-east-1.amazonaws.com/digitalbank:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
```


<br>

![alt text](images/Deployed_EKS_Cluster.png)

<br>

![alt text](images/Deployed_EKS_Cluster_APPEND.png)

<br>

In Kubernetes, a Service is a method for exposing a network application that is running as one or more Pods in your
cluster. A key aim of Services in Kubernetes is that you don't need to modify your existing application to use an unfamiliar
service discovery mechanism. We can run code in Pods, whether this is a code designed for a cloud-native world, or an older app
you've containerized. We use a Service to make that set of Pods available on the network so that clients can interact with it.

A Kubernetes service is an abstraction that defines a logical set of pods and a policy for accessing them. It acts as a load 
balancer and provides a stable IP address and DNS name for the pods that are part of the service. The service ensures that 
traffic is evenly distributed across the pods, and it also handles scenarios where pods are added or removed, ensuring that 
the service remains available. The deployment, on the other hand, is responsible for creating and managing the pods that run 
your application containers. It ensures that the desired number of replicas are running and automatically handles rolling updates 
and rollbacks. Here's why the service needs to be deployed after the deployment:

- Pod Selection: The service uses label selectors to identify the pods it should proxy traffic to. These labels are defined in 
the deployment's pod template. If the service is created before the deployment, there won't be any pods matching the label 
selectors, and the service won't have any endpoints to forward traffic to.

- Service Discovery: Services provide a stable DNS name and IP address for the pods they represent. If the service is created 
before the pods, there won't be any pods to provide these details for, and other components won't be able to discover and 
connect to the application.

- Load Balancing: The service acts as a load balancer, distributing traffic across the pods. If the service is created before 
the pods, there won't be any pods to load balance the traffic to.

By deploying the service after the deployment, you ensure that the pods are running and ready to receive traffic before 
the service is created. This way, the service can correctly identify the pods, provide service discovery, and load balance 
the traffic across the available pods. In order to access the Kubernetes resources from the outside, we need to deploy the 
service `eks/service.yaml` as below:


<br>

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

<br>

```textmate
$ kubectl apply -f eks/service.yaml

service/digital-bank-api-service configured
```

<br>

![alt text](images/EKS_Service.png)

<br>

In the image above, we can see the service is deployed but the `TYPE` is `ClusterIP` and there is `EXTERNAL-IP` so the Kubernetes
resources can't be access from the outside by us. We can change the `TYPE` to `LoadBalancer`, re-deploy and acquire an `EXTERNAL-IP`
identified as `acefda17bec1049e48640f98a99c2653-885218934.us-east-1.elb.amazonaws.com` below and this host name can be used 
for regular requests as we did with the `localhost`. 

<br>

![alt text](images/EKS_Service_APPEND.png)

<br>

![alt text](images/EKS_Service_URL_NS_LOOKUP_Before_A_RECORD.png)

<br>

![alt text](images/EKS_Service_URL.png)

<br>

### AWS Route 53

<br>

We would like to use a dedicated domain for the deployment process and will use AWS Route 53 for the purpose. AWS Route 53 
is a highly available and scalable cloud Domain Name System (DNS) web service provided by Amazon Web Services (AWS). It allows 
to manage domain names, route end-user requests to internet applications by translating domain names into IP addresses, and 
perform health checks to monitor the availability of your application resources. Route 53 integrates with other AWS services, 
providing reliable and cost-effective domain name management and traffic routing capabilities for your applications running on 
AWS or on-premises infrastructure.

Route 53 has several key components that allow to manage and route traffic to the resources:

- Record Name: This is the domain or subdomain name for which you want to define routing rules.

- Record Type: This specifies the type of DNS record, such as A (IPv4 address), AAAA (IPv6 address), CNAME (alias for 
another domain name), or MX (mail exchange record).

- Route Traffic To: This defines the resource (e.g., EC2 instance, load balancer, or S3 bucket) to which Route 53 should 
route traffic for the specified record name and type.

- Routing Policy**: This determines how Route 53 responds to queries for the record and how it routes traffic to the resources. 
Options include simple, failover, geolocation, latency, multivalue answer, and weighted routing policies.


There are some concepts needs for Route 53 and are essential for managing domain names and DNS settings.

- DNS Settings: Route 53 allows you to configure DNS settings, such as specifying the authoritative nameservers for domain 
and managing resource record sets.

- Domain Nameserver Registrations: When you register a domain with Route 53, it automatically assigns four AWS nameservers to 
your domain, which are used to route traffic to your resources.

- DNSSEC: Route 53 supports DNSSEC, which adds a layer of security to DNS by digitally signing DNS data to protect against 
data spoofing and man-in-the-middle attacks. This ensures the authenticity and integrity of your DNS data.

- A Record: An A record (short for Address record) is a type of DNS record that maps a domain name to an IPv4 address. It 
is one of the most common and essential record types in the Domain Name System (DNS). IP Address Resolution: When a user or 
client attempts to access a website or web application by entering a domain name (e.g., example.com) in a web browser, the 
DNS system needs to translate that domain name into an IP address that the computer can understand. The A record provides 
this mapping between the domain name and the corresponding IPv4 address.

<br>


After purchasing a domain name, we need to set up an A record in Route 53's hosted zone to map the domain name to the IP address 
of the web server or load balancer. We need to Create a hosted zone in Route 53 for the domain name if we haven't already. 
By default, Route 53 creates NS (Name Server) and SOA (Start of Authority) records for the hosted zone. Next, we need to create 
an A record within the hosted zone, specifying your domain name (e.g., example.com or www.example.com) as the Record Name, and the 
public IP address of your web server or load balancer (for our case) as the Value. This A record will resolve your domain name to the 
specified IP address, allowing visitors to access the website or application hosted on that server. 


<br>

![alt text](images/Route_53_A_Record_CREATE.png)

<br>

![alt text](images/Route_53_DASHBOARD.png)

<br>

[//]: # (TODO: The image above needs update with the A Record)

![alt text](images/EKS_Service_URL_NS_LOOKUP_After_A_RECORD.png)

<br>


### Kubernetes Ingress


Kubernetes Ingress is a collection of rules that allow inbound connections to reach the cluster services. It acts as a 
single entry point for incoming traffic, routing it to the appropriate services based on configured rules. Ingress is 
needed in Kubernetes to expose services externally, enabling external clients to access applications running inside the 
cluster. It provides features like load balancing, SSL/TLS termination, and name-based virtual hosting, making it easier 
to manage and secure incoming traffic to services in a Kubernetes cluster. Ingress exposes HTTP and HTTPS routes from outside 
the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource.

Before proceeding to the Ingress, change the type of the `eks/service.yaml` to `ClusterIP` from `LoadBalancer`:

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

We need to deploy the `eks/service.yaml` and then, need to deploy the `eks/ingress.yaml` as provided below:

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

<br>

```textmate
$ kubectl apply -f eks/service.yaml

$ kubectl apply -f eks/ingress.yaml
```

<br>

We have ClusterIP in the deployed service by now, and we can look for service and ingress in the K9s console:

<br>

![alt text](images/EKS_Service_ClusterIP.png)

<br>

![EKS_Ingress_Deploy_INITIAL](images/EKS_Ingress_Deploy_INITIAL.png)

<br>

This Ingress will send traffic to the `digital-bank-api-service` defined as the service: 

<br>

![EKS_Ingress_Sending_Traffic_A_RECORD_UPDATED](images/EKS_Ingress_Sending_Traffic_A_RECORD_UPDATED.png)

<br>

The Nginx Ingress controller is a key component in Kubernetes that provides load balancing, reverse proxy, and advanced 
routing capabilities for incoming traffic to services within the cluster. It acts as a single entry point, distributing 
requests across multiple service replicas, handling TLS/SSL termination, supporting name-based virtual hosting, URL rewriting, 
and integrating seamlessly with Kubernetes Services. By deploying a Nginx Ingress controller, you can simplify traffic management, 
enhance security, and leverage advanced features for efficiently routing and exposing your applications to external clients.
In the previous image, we had no ADDRESS for the Ingress, and we need to deploy the Nginx Ingress Controller for acquiring 
ADDRESS that the host can direct traffic to. 


<br>

![Nginx_ingress_CONTROLLER](images/Nginx_ingress_CONTROLLER.png)

<br>

![Nginx_ingress_CONTROLLER_RUN](images/Nginx_ingress_CONTROLLER_RUN.png)

<br>

```textmate
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.1/deploy/static/provider/aws/deploy.yaml
```

<br>

After the redeployment the Ingress, now we can see the `ADDRESS` and we can use that to update the A-record in the AWS Route 53. 

![Nginx_ingress_CONTROLLER_AFTER_RUN](images/Nginx_ingress_CONTROLLER_AFTER_RUN.png)

<br>

![Route_53_A_Record_UPDATE_INGRESS](images/Route_53_A_Record_UPDATE_INGRESS.png)

<br>

![INGRESS_Update_CHECK](images/INGRESS_Update_CHECK.png)

<br>

`IngressClass` is a Kubernetes resource that allows you to define and specify different types of Ingress controllers within 
a cluster, enabling the simultaneous operation of multiple Ingress controllers. This is particularly useful in multi-tenant 
environments, where different teams or applications require different ingress configurations or features, as well as scenarios 
involving vendor-specific features, migration strategies, or testing and experimentation with various Ingress controllers. 
By leveraging IngressClass, you can better manage and organize Ingress resources, isolate traffic and configurations, and 
take advantage of the capabilities offered by different Ingress controller vendors.

<br>

![Ingress_Class_UPDATED](images/Ingress_Class_UPDATED.png)

<br>

This is the updated `eks/ingress.yaml` with the info about the `IngressClass` provided, and we can see the class of the Ingress
is changed to `Nginx`:

<br>

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

<br>

```textmate
$ kubectl apply -f eks/ingress.yaml
```
<br>



--- Enable TLS and Certificate 



Make the Client/Server communication secure using TLS
Should only be use if the DNS provider has an API to update the records


DNS-01 Challenge
![alt text](images/dns-01_challenge.png)

HTTP 01 Challenge

![alt text](images/http-01_challenge.png)


URL YT: <https://youtu.be/-f4Gbk-U758>

SITE: <https://letsencrypt.org/>




Install Kubernetes cert manager

<br>

```textmate

$ kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.5/cert-manager.yaml
```

<br>

![EKS_Cert_Manager_PODS](images/EKS_Cert_Manager_PODS.png)

<br>


```textmate
$ kubectl get pods --namespace cert-manager

NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-7ddd8cdb9f-bxlsn              1/1     Running   0          22h
cert-manager-cainjector-57cd76c845-2lq2b   1/1     Running   0          22h
cert-manager-webhook-cf8f9f895-8c7bd       1/1     Running   0          22h

```

<br>

![alt text](images/Create_Basic_ACME_ISSUER.png)

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


Update the ingress to find the certificate:


<br>

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

<br>

Now we can check that the TLS is enabled

<br>

![TLS](images/TLS.png)

<br>

![alt text](images/Certificates.png)

<br>

![alt text](images/Certificate_ISSUED.png)

<br>

![alt text](images/TLS_HTTPS_Enabled.png)

<br>

![alt text](images/Postman_Requests_WORKING.png)

<br>



### K9s 

```textmate
$ brew install k9s 
```

In the K9s console, check for the `>Clusterissuer`
In the K9s console, check for the `>secrets` for its private keys

The certificates are still empty:

In the K9s console, check for the `>certificate` for its private keys
In the K9s console, check for the `>certificaterequest` for its private keys


In the K9s console, use:

```textmate
$ configmap
```


Delete the existing deployments in the k9s and then <d>


```textmate
$ deployments 
$ services
```









