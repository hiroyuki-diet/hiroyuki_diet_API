# AWS デプロイ手順

このドキュメントでは、GitHub Actions を使用して AWS ECS (Fargate) に自動デプロイする方法を説明します。

## 前提条件

- AWS アカウント
- AWS CLI のインストール
- 適切な IAM 権限

## AWS リソースのセットアップ

### 1. ECR リポジトリの作成

```bash
aws ecr create-repository \
  --repository-name hiroyuki-diet-api \
  --region ap-northeast-1
```

作成後、リポジトリの URI をメモしておきます。

### 2. RDS (PostgreSQL) の作成

```bash
aws rds create-db-instance \
  --db-instance-identifier hiroyuki-diet-db \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --master-username admin \
  --master-user-password YOUR_PASSWORD \
  --allocated-storage 20 \
  --vpc-security-group-ids YOUR_SECURITY_GROUP_ID \
  --db-subnet-group-name YOUR_SUBNET_GROUP \
  --publicly-accessible false
```

または、AWS コンソールから作成することもできます。

### 3. Secrets Manager にパスワードを保存

```bash
aws secretsmanager create-secret \
  --name hiroyuki-diet/db-password \
  --description "Database password for Hiroyuki Diet API" \
  --secret-string YOUR_DB_PASSWORD \
  --region ap-northeast-1
```

### 4. CloudWatch Logs グループの作成

```bash
aws logs create-log-group \
  --log-group-name /ecs/hiroyuki-diet-api \
  --region ap-northeast-1
```

### 5. ECS クラスターの作成

```bash
aws ecs create-cluster \
  --cluster-name hiroyuki-diet-cluster \
  --region ap-northeast-1
```

### 6. IAM ロールの作成

#### ecsTaskExecutionRole

このロールは ECS がコンテナを起動するために必要です。

```bash
aws iam create-role \
  --role-name ecsTaskExecutionRole \
  --assume-role-policy-document file://trust-policy.json

aws iam attach-role-policy \
  --role-name ecsTaskExecutionRole \
  --policy-arn arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
```

trust-policy.json:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

Secrets Manager へのアクセス権限も追加:

```bash
aws iam put-role-policy \
  --role-name ecsTaskExecutionRole \
  --policy-name SecretsManagerAccess \
  --policy-document file://secrets-policy.json
```

secrets-policy.json:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Resource": "arn:aws:secretsmanager:ap-northeast-1:YOUR_ACCOUNT_ID:secret:hiroyuki-diet/*"
    }
  ]
}
```

#### ecsTaskRole

このロールはコンテナが AWS サービスにアクセスするために必要です（必要に応じて）。

```bash
aws iam create-role \
  --role-name ecsTaskRole \
  --assume-role-policy-document file://trust-policy.json
```

### 7. タスク定義ファイルの編集

`.aws/task-definition.json` を編集して、以下の値を実際の値に置き換えます:

- `YOUR_ACCOUNT_ID`: AWS アカウント ID
- `YOUR_RDS_ENDPOINT`: RDS のエンドポイント
- `YOUR_DB_USER`: データベースのユーザー名

### 8. ECS タスク定義の登録

```bash
aws ecs register-task-definition \
  --cli-input-json file://.aws/task-definition.json \
  --region ap-northeast-1
```

### 9. Application Load Balancer (ALB) の作成 (オプション)

外部からアクセスする場合は ALB が必要です。

```bash
# ALB の作成
aws elbv2 create-load-balancer \
  --name hiroyuki-diet-alb \
  --subnets YOUR_SUBNET_1 YOUR_SUBNET_2 \
  --security-groups YOUR_SECURITY_GROUP \
  --region ap-northeast-1

# ターゲットグループの作成
aws elbv2 create-target-group \
  --name hiroyuki-diet-tg \
  --protocol HTTP \
  --port 8080 \
  --vpc-id YOUR_VPC_ID \
  --target-type ip \
  --health-check-path / \
  --region ap-northeast-1

# リスナーの作成
aws elbv2 create-listener \
  --load-balancer-arn YOUR_ALB_ARN \
  --protocol HTTP \
  --port 80 \
  --default-actions Type=forward,TargetGroupArn=YOUR_TARGET_GROUP_ARN \
  --region ap-northeast-1
```

### 10. ECS サービスの作成

```bash
aws ecs create-service \
  --cluster hiroyuki-diet-cluster \
  --service-name hiroyuki-diet-service \
  --task-definition hiroyuki-diet-api \
  --desired-count 1 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[YOUR_SUBNET_1,YOUR_SUBNET_2],securityGroups=[YOUR_SECURITY_GROUP],assignPublicIp=ENABLED}" \
  --load-balancers "targetGroupArn=YOUR_TARGET_GROUP_ARN,containerName=hiroyuki-diet-api,containerPort=8080" \
  --region ap-northeast-1
```

ALB を使用しない場合（パブリック IP で直接アクセス）:

```bash
aws ecs create-service \
  --cluster hiroyuki-diet-cluster \
  --service-name hiroyuki-diet-service \
  --task-definition hiroyuki-diet-api \
  --desired-count 1 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[YOUR_SUBNET_1,YOUR_SUBNET_2],securityGroups=[YOUR_SECURITY_GROUP],assignPublicIp=ENABLED}" \
  --region ap-northeast-1
```

## GitHub Secrets の設定

GitHub リポジトリの Settings > Secrets and variables > Actions で以下のシークレットを追加します:

- `AWS_ACCESS_KEY_ID`: AWS アクセスキー ID
- `AWS_SECRET_ACCESS_KEY`: AWS シークレットアクセスキー

### IAM ユーザーの作成と権限設定

GitHub Actions 用の IAM ユーザーを作成します:

```bash
aws iam create-user --user-name github-actions-deployer
```

必要な権限ポリシーを作成（deploy-policy.json）:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "ecr:PutImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ecs:UpdateService",
        "ecs:DescribeServices",
        "ecs:DescribeTaskDefinition",
        "ecs:RegisterTaskDefinition"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "iam:PassRole"
      ],
      "Resource": [
        "arn:aws:iam::YOUR_ACCOUNT_ID:role/ecsTaskExecutionRole",
        "arn:aws:iam::YOUR_ACCOUNT_ID:role/ecsTaskRole"
      ]
    }
  ]
}
```

ポリシーをアタッチ:

```bash
aws iam put-user-policy \
  --user-name github-actions-deployer \
  --policy-name GitHubActionsDeployPolicy \
  --policy-document file://deploy-policy.json
```

アクセスキーを作成:

```bash
aws iam create-access-key --user-name github-actions-deployer
```

このコマンドで出力されたアクセスキー ID とシークレットアクセスキーを GitHub Secrets に設定します。

## デプロイの実行

main ブランチに push すると、自動的にデプロイが実行されます:

```bash
git add .
git commit -m "Setup AWS deployment"
git push origin main
```

GitHub Actions のタブでデプロイの進行状況を確認できます。

## トラブルシューティング

### デプロイが失敗する場合

1. CloudWatch Logs でコンテナのログを確認
2. ECS タスクの状態を確認
3. セキュリティグループの設定を確認（ポート 8080 が開いているか）
4. RDS への接続が可能か確認

### ログの確認方法

```bash
aws logs tail /ecs/hiroyuki-diet-api --follow --region ap-northeast-1
```

### ECS サービスの状態確認

```bash
aws ecs describe-services \
  --cluster hiroyuki-diet-cluster \
  --services hiroyuki-diet-service \
  --region ap-northeast-1
```

## コスト削減のヒント

- 開発環境では Fargate Spot を使用
- 使用していない時間帯は ECS サービスの desired count を 0 に設定
- RDS は t3.micro または t4g.micro を使用
- 不要になったリソースは削除

## その他のデプロイオプション

### App Runner を使用する場合

より簡単にデプロイしたい場合は、App Runner を使用することもできます。
`.github/workflows/deploy.yml` を以下のように変更します:

```yaml
name: Deploy to App Runner

on:
  push:
    branches:
      - main

env:
  AWS_REGION: ap-northeast-1
  ECR_REPOSITORY: hiroyuki-diet-api

jobs:
  deploy:
    name: Deploy to App Runner
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v4
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: ${{ env.AWS_REGION }}

    - name: Login to Amazon ECR
      id: login-ecr
      uses: aws-actions/amazon-ecr-login@v2

    - name: Build and push image to Amazon ECR
      env:
        ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
        IMAGE_TAG: ${{ github.sha }}
      run: |
        docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG ./backend
        docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG

    - name: Deploy to App Runner
      run: |
        aws apprunner update-service \
          --service-arn YOUR_APP_RUNNER_SERVICE_ARN \
          --source-configuration "ImageRepository={ImageIdentifier=$ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG,ImageConfiguration={Port=8080}}" \
          --region $AWS_REGION
```

App Runner の方がセットアップが簡単で、自動スケーリングも組み込まれています。
