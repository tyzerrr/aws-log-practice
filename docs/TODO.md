# TODO

## Next

- RDS for PostgreSQL を private subnet に作成する
- ECS task security group から RDS security group への通信を許可する
- DB 接続情報を Secrets Manager または SSM Parameter Store で管理する
- ECS task definition に DB 接続情報を渡す
- backend アプリケーションを ECS 上で起動する

## Deployment

- ecspresso を導入する
- ECS task definition / service deploy を ecspresso 管理に寄せる
- GitHub Actions の CI を実アプリ image build / ECR push / ecspresso deploy に更新する
- HTTPS 経由で実アプリケーションの動作確認を行う

## Logging

- CloudWatch Logs の log group 設計を整理する
- Kinesis Data Firehose を追加する
- ECS application logs を Firehose 経由で保存・分析できる構成にする
- Production Ready なログ運用に必要な retention、暗号化、権限、監視を整理する

## Documentation

- terraform-docs を導入する
- Terraform module / environment の README を生成する
- 必要であれば CI で terraform-docs の差分チェックを行う
