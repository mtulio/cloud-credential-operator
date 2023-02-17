# AWS Security Token Service - Steps to Migrate the OIDC issuer from public S3 Bucket to CloudFront Distribution

- [Prerequisites](#Prerequisites)
- [Test the existing token](#test-token)
- [Setup new OIDC with CloudFront Distribution](#setup)
    - [Create the CloudFront Distribution](#setup-cloudfront)
    - [Create and patch the new OIDC discovery documents and JWKS](#setup-oidc-documents)
    - [Create the new OIDC using CloudFront Distribution](#setup-oidc-idp)
- [Backing up current cluster data](#backup)
- [Patch the cluster to use the new OIDC](#patch-cluster)
- [Revoke public access to S3 Bucket](#revoke-s3-public-access)
- [Rollback to OIDC with S3 Public URL](#rollback)
- [Rollback to OIDC with S3 Public URL](#delete)

## Prerequisites <a name="Prerequisites"></a>

<!-- - Create an AWS cluster with STS (without CloudFront) [tmp section]

> https://mtulio.net/playbooks/openshift/ocp-aws-cco-sts-install-quickly/

```bash
CLUSTER_VERSION="4.12.2" &&\
  CLUSTER_NAME="oidc-pub" &&\
  CLUSTER_BASE_DOMAIN="devcluster.openshift.com" &&\
  create_cluster $CLUSTER_NAME
``` -->

- An OCP Cluster created on AWS with manual authentication mode with STS using S3 Bucket as OIDC URL (issuerURL)


- Environment variables exported:
```bash
export CLUSTER_NAME="<ChangeMe>"
export OIDC_BUCKET_HOST=$(oc get authentication cluster -o jsonpath={'.spec.serviceAccountIssuer'} | awk -F'://' '{print$2}')
export OIDC_BUCKET_NAME=$(oc get authentication cluster -o jsonpath={'.spec.serviceAccountIssuer'} | awk -F'://' '{print$2}' |awk -F'.' '{print$1}')
export CLUSTER_REGION=$(oc get authentication cluster -o jsonpath={'.spec.serviceAccountIssuer'} | awk -F'://' '{print$2}' |awk -F'.' '{print$3}')
```

- Make sure you can reach (read) the S3 Bucket created with the default name

```bash
aws s3 ls --region $CLUSTER_REGION s3://$OIDC_BUCKET_NAME
```

- A clean work directory: a lot of files will be created, make sure you switched to a new work directory to save the files properly (it can be used in the future for rollback)

## Test the existing token <a name="test-token"></a>

This section is to make sure everything is working correctly in your existing environment.

The steps described below will use the credentials provided to the `machine-api` component, trying to assume the role using `aws-cli`.

It's expected that the existing token will be able to authenticate in AWS. Otherwise, you must abort the operations and do not try to run any step described in the next sections

```bash
## test existing token
# Get Token path from AWS credentials mounted to pod
TOKEN_PATH=$(oc get secrets aws-cloud-credentials \
    -n openshift-machine-api \
    -o jsonpath='{.data.credentials}' |\
    base64 -d |\
    grep ^web_identity_token_file |\
    awk '{print$3}')

# Get Controler's pod
CAPI_POD=$(oc get pods -n openshift-machine-api \
    -l api=clusterapi \
    -o jsonpath='{.items[*].metadata.name}')

# Extract tokens from the pod
TOKEN=$(oc exec -n openshift-machine-api ${CAPI_POD} \
    -c machine-controller -- cat ${TOKEN_PATH})

echo $TOKEN | awk -F. '{ print $2 }' | base64 -d 2>/dev/null | jq .iss

IAM_ROLE=$(oc get secrets aws-cloud-credentials \
    -n openshift-machine-api \
    -o jsonpath='{.data.credentials}' |\
    base64 -d |\
    grep ^role_arn |\
    awk '{print$3}')

echo $IAM_ROLE

aws sts assume-role-with-web-identity \
    --role-arn "${IAM_ROLE}" \
    --role-session-name "my-session" \
    --web-identity-token "${TOKEN}"
```

## Setup new OIDC with CloudFront Distribution <a name="setup"></a>

### Create the CloudFront Distribution <a name="setup-cloudfront"></a>

- Create the Origin Access Identity (OAI)

```bash
export DIR_CCO="./"
export OIDC_BUCKET_PATH="/pvt"

export OAI_CLOUDFRONT_ID=$(aws cloudfront create-cloud-front-origin-access-identity \
    --cloud-front-origin-access-identity-config \
    CallerReference="${OIDC_BUCKET_NAME}",Comment="OAI-${OIDC_BUCKET_NAME}" \
    | jq -r .CloudFrontOriginAccessIdentity.Id)
```

- Create the CloudFront Distribution
```bash
# Should be updated to master before merging
#wget https://raw.githubusercontent.com/openshift/cloud-credential-operator/master/docs/sts-oidc-cloudfront.json.tpl
wget https://raw.githubusercontent.com/mtulio/cloud-credential-operator/doc-sts-updates/docs/sts-oidc-cloudfront.json.tpl

cat sts-oidc-cloudfront.json.tpl \
   | envsubst \
   > ${DIR_CCO}/oidc-cloudfront.json


export CLOUDFRONT_HOST=$(aws cloudfront create-distribution-with-tags \
    --distribution-config-with-tags \
    file://${DIR_CCO}/oidc-cloudfront.json \
    | jq -r .Distribution.DomainName)

echo ${CLOUDFRONT_HOST}
echo ${OIDC_BUCKET_HOST}
```

### Create and patch the new OIDC discovery documents and JWKS<a name="setup-oidc-documents"></a>

- Download the current OIDC files (discovery document and JWKS) to the local disk:

```bash
aws s3 sync s3://${OIDC_BUCKET_NAME} ./bucket
```

- Create the new path `/pvt` under the local directory:

```bash
# update the reference
mkdir bucket/pvt/
$ cp -rvf bucket/keys.json bucket/.well-known/ bucket/pvt/
'bucket/keys.json' -> 'bucket/pvt/keys.json'
'bucket/.well-known/' -> 'bucket/pvt/.well-known'
'bucket/.well-known/openid-configuration' -> 'bucket/pvt/.well-known/openid-configuration'

$ ls -a bucket/pvt/
.  ..  keys.json  .well-known

```

- Patch the new documents with the CloudFront Distribution Domain name:
```bash
sed -i "s/$OIDC_BUCKET_HOST/$CLOUDFRONT_HOST/g" bucket/pvt/.well-known/openid-configuration
```

- Upload the patched files to the Bucket with the new object prefix `/pvt`

```bash
aws s3 sync ./bucket/pvt s3://${OIDC_BUCKET_NAME}/pvt
```

The new object path, `/pvt`, must be accessed by CloudFront through OAI. The Bucket Policy will be added to allow that operation from CloudFront Distribution.

- Download the existing template to create the Bucket Policy

```bash
#TODO: replace to master
#wget https://raw.githubusercontent.com/openshift/cloud-credential-operator/master/docs/sts-oidc-bucket-policy.json.tpl
wget https://raw.githubusercontent.com/mtulio/cloud-credential-operator/doc-sts-updates/docs/sts-oidc-bucket-policy.json.tpl
```

- Create the Bucket Policy and apply it

```bash
cat sts-oidc-bucket-policy.json.tpl \
   | envsubst \
   > ${DIR_CCO}/oidc-bucket-policy.json

aws s3api put-bucket-policy \
    --bucket ${OIDC_BUCKET_NAME} \
    --policy file://${DIR_CCO}/oidc-bucket-policy.json
```

Now the CloudFront Distribution must have access to the Bucket object `/pvt/keys.json`, test it:

```
$ curl https://$CLOUDFRONT_HOST/keys.json
```

### Create the new OIDC using CloudFront Distribution<a name="setup-oidc-idp"></a>

- Extract the service account signer public key, to generate the IdP by `ccoctl`:

```bash
oc get configmap bound-sa-token-signing-certs \
    --namespace openshift-kube-apiserver \
    --output json \
    | jq --raw-output '.data["service-account-001.pub"]' \
    > serviceaccount-signer.public
```

- Generate the IdP files into the local directory `new-oidc`:

```bash
./ccoctl aws create-identity-provider \
    --name=${CLUSTER_NAME} \
    --region=${CLUSTER_REGION} \
    --public-key-file=${PWD}/serviceaccount-signer.public \
    --output-dir=new-oidc/ \
    --dry-run
```

- Patch the IdP OIDC to the new Domain name

```bash
sed -i "s/$OIDC_BUCKET_HOST/$CLOUDFRONT_HOST/g" new-oidc/04-iam-identity-provider.json
```

- Discover the thumbprint for the keys from the CloudFront Distribution URL:

> AWS Docs - Getting the Thumbprint: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_create_oidc_verify-thumbprint.html

```bash
openssl s_client -servername $CLOUDFRONT_HOST -showcerts -connect $CLOUDFRONT_HOST:443 </dev/null | openssl x509 -outform pem > certificate.crt

export CERT_THUMBPRINT=$(openssl x509 -in certificate.crt -fingerprint -sha1 -noout | awk -F'=' '{print$2}' | tr -d ':')

jq -r ".ThumbprintList=[\"$CERT_THUMBPRINT\"]" ${PWD}/new-oidc/04-iam-identity-provider.json \
    > ${PWD}/new-oidc/04-iam-identity-provider-new.json
```

- Create the identity provider AWS OIDC:

```bash
aws iam create-open-id-connect-provider \
    --cli-input-json file://${PWD}/new-oidc/04-iam-identity-provider-new.json \
    > ${PWD}/new-oidc//04-iam-identity-provider-object.json 

# Describe or get from file
#export OIDC_ARN=$(aws iam list-open-id-connect-providers |jq -r ".OpenIDConnectProviderList[] | select(.Arn | endswith(\"${CLOUDFRONT_URI}\")).Arn")

export OIDC_ARN=$(jq -r .OpenIDConnectProviderArn ${PWD}/new-oidc//04-iam-identity-provider-object.json)

echo ${OIDC_ARN}
```

## Backup existing state <a name="backup"></a>

- Get Objects and existing tests

```bash
export CURRENT_PATH=$PWD/current-cluster
mkdir $CURRENT_PATH
oc get authentication -o yaml |tee -a $CURRENT_PATH/authentication.yaml

aws sts assume-role-with-web-identity \
    --role-arn "${IAM_ROLE}" \
    --role-session-name "my-session" \
    --web-identity-token "${TOKEN}" \
    | jq -r '.Credentials=""' \
    | tee ${CURRENT_PATH}/identities.json
```

- Save the current IAM Roles

```bash
aws iam list-roles \
    | jq -r  ".Roles[] | select(.RoleName | startswith(\"${CLUSTER_NAME}-openshift\"))" \
    | tee ${CURRENT_PATH}/iam-roles.json
```

## Patch the cluster to use the new OIDC<a name="patch-cluster"></a>

- Patch the trusted policy documents with the new OIDC URL

```bash
sed "s/$OIDC_BUCKET_HOST/$CLOUDFRONT_HOST/g" ${CURRENT_PATH}/iam-roles.json \
    | tee iam-roles-new.json
```

- Patch the IAM Roles Trusted Policy documents

> NOTE 1: from here, the cluster will lose access to the integrated components (machine-api, image registry, CSI, ...)

> NOTE 2: The script below should be run carefully, it was created to show the current and desired policies. If you find anything that does not match the expected changes, abort it immediately.

> Helper [`aws iam get-role`](https://docs.aws.amazon.com/cli/latest/reference/iam/get-role.html)

> Helper [`aws iam update-assume-role-policy`](https://docs.aws.amazon.com/cli/latest/reference/iam/update-assume-role-policy.html)

```bash
for ROLE_NAME in $(jq -r .RoleName iam-roles-new.json);
do
    echo -e ">>>>>\n#> (1) CURRENT IAM Role \"$ROLE_NAME\":";
    aws iam get-role --role-name $ROLE_NAME | jq .Role;

    echo -e "\n#> (2) NEW IAM Role \"$ROLE_NAME\" AssumeRolePolicyDocument:";
    jq -r ". | select(.RoleName == \"$ROLE_NAME\").AssumeRolePolicyDocument" iam-roles-new.json \
        | tee ${PWD}/iam-roles-new-$ROLE_NAME.json
    
    read -p "ATTENTION: The AssumeRolePolicyDocument for IAM Role(1) will be patched to the value of (2). Do you want to continue? [y/n]: " answer
    if [ -z "$answer" ] || [ "$answer" != "y" ]
    then
        echo "answer[$answer]. Canceling the operation.";
        break
    fi
    echo "Patching..."
    aws iam update-assume-role-policy \
        --role-name $ROLE_NAME \
        --policy-document file:///${PWD}/iam-roles-new-$ROLE_NAME.json
    echo "Done! Return code=$?"
done
```

- Patch the issuer URL to the new OIDC URL on the Authentication object:

```bash
oc patch authentication cluster \
    --type=merge \
    -p "{\"spec\":{\"serviceAccountIssuer\":\"https://${CLOUDFRONT_HOST}\"}}"
```

- Wait for the kube-apiserver rollout

> Wait to clean the `PROGRESSING=TRUE`. It could take some minutes to start, and several to finish.

```bash
$ oc get co kube-apiserver -w
$ oc get pods -n openshift-kube-apiserver -l apiserver=true -w
```

- Restart all pods:

```bash
for I in $(oc get ns -o jsonpath='{range .items[*]} {.metadata.name}{"\n"} {end}'); \
      do oc delete pods --all -n $I; \
      sleep 1; \
      done
```

- Test the new token

> Repeat the steps in the section "Test the existing token"

> Make sure the JWT token has the CloudFront Distribution Domain name as the Issuer URL, field `.iss`

> Make sure you can assume the role correctly and the signer will be the CloudFront: `.Provider` in the answer from `assume-role-with-web-identity`

If you have completed these steps successfully, the cluster is using the new identity provider with AWS CloudFront.

## Revoke public access to the S3 Bucket<a name="revoke-s3-public-access"></a>

- Change the default policy blocking public access to the bucket:

```bash
aws s3api put-public-access-block \
    --bucket ${OIDC_BUCKET_NAME} \
    --public-access-block-configuration \
    BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true
```

- Test it (expected to fail [HTTP 403])

```bash
curl -vvv https://$OIDC_BUCKET_HOST/keys.json
```

## Rollback to OIDC with S3 Public URL<a name="rollback"></a>

When the process to migrate to a private bucket have been failed, and you want to roll back to the OIDC issuer URL pointing to the public S3 Bucket, you must follow the steps below.

- Reopen the bucket policy

```bash
aws s3api put-public-access-block \
    --bucket ${OIDC_BUCKET_NAME} \
    --public-access-block-configuration \
    BlockPublicAcls=false,IgnorePublicAcls=false,BlockPublicPolicy=false,RestrictPublicBuckets=false
```

- Replace the `serviceAccountIssuer` with the S3 Bucket's URL

```bash
oc patch authentication cluster \
    --type=merge \
    -p "{\"spec\":{\"serviceAccountIssuer\":\"https://${OIDC_BUCKET_HOST}\"}}"

```

- Patch the Assume Role policy (Trusted Policy)

```bash
# Patch the roles back to S3
for ROLE_NAME in $(jq -r .RoleName ${CURRENT_PATH}/iam-roles.json);
do
    echo -e ">>>>>\n#> (1) CURRENT IAM Role \"$ROLE_NAME\":";
    aws iam get-role --role-name $ROLE_NAME | jq .Role;

    echo -e "\n#> (2) NEW IAM Role \"$ROLE_NAME\" AssumeRolePolicyDocument:";
    jq -r ". | select(.RoleName == \"$ROLE_NAME\").AssumeRolePolicyDocument" ${CURRENT_PATH}/iam-roles.json \
        | tee ${CURRENT_PATH}/iam-roles-rollback-$ROLE_NAME.json
    
    read -p "ATTENTION: The AssumeRolePolicyDocument for IAM Role(1) will be patched to the value of (2). Do you want to continue? [y/n]: " answer
    if [ -z "$answer" ] || [ "$answer" != "y" ]
    then
        echo "answer[$answer]. Canceling the operation.";
        break
    fi
    echo "Patching..."
    aws iam update-assume-role-policy \
        --role-name $ROLE_NAME \
        --policy-document file:///${CURRENT_PATH}/iam-roles-rollback-$ROLE_NAME.json
    echo "Done! Return code=$?"
done
```

- Wait for the kube-apiserver to apply the configuration (PROGRESSING=FALSE)

```bash
$ oc get co kube-apiserver -w
$ oc get pods -n openshift-kube-apiserver -l apiserver=true -w
```

- Restart all pods
```bash
# Restart all pods
for I in $(oc get ns -o jsonpath='{range .items[*]} {.metadata.name}{"\n"} {end}'); \
      do oc delete pods --all -n $I; \
      sleep 1; \
      done
```

## Delete the old OIDC IdP<a name="delete"></a>

If you run successfully the steps and tested them, you can remove the old OIDC pointing to the S3.
