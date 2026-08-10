# 本地 XXL-JOB Admin

生产请使用外部调度中心；本地可用：

```bash
# 1) 在业务库执行官方 tables_xxl_job.sql 初始化表结构
# 2) 执行种子
psql "$DATABASE_URL" -f script/sql/xxl_job_hei_gin_seed.sql

# 3) 启动 Admin（镜像仅作本地调试）
docker compose -f script/docker/docker-compose.xxl-job.yml up -d
```

控制台：`http://127.0.0.1:9004/xxl-job-admin`（默认 `admin` / `123456`）

`config.yaml`：

```yaml
xxl_job:
  enabled: true
  access_token: default_token
  admin:
    addresses: http://127.0.0.1:9004/xxl-job-admin
  executor:
    appname: hei-gin-api   # 须与种子 xxl_job_group.app_name 一致
    port: 9999
```

API 进程内嵌执行器（`go run ./app/cmd/api`），Handler：

| JobHandler | 模块 |
|------------|------|
| `bannerStatusJob` | sys/banner |
| `accountPurgeCancelledJob` | iam/account |
