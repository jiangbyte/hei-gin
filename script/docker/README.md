# 本地调度中心（SnailJob）

API 进程内嵌 **SnailJob** Go 客户端（见根目录 README）。下列 compose 仅用于本地调试 Server；生产请连外部 Server。

前置：本机已有 Postgres，并准备库 `snail_job`、角色（默认脚本期望 `admin` / `123456`，可按环境变量覆盖）。

```bash
# 1) 建库后执行一次 Flyway 迁移（含独立 namespace / hei_gin_admin / 三个 Go Job）
./script/docker/snailjob-flyway.sh

# 2) 启动 Server
docker compose -f script/docker/docker-compose.snailjob.yml up -d
```

控制台：`http://127.0.0.1:9189/snail-job`（默认 `admin` / `123456`）。

`config.yaml`：

```yaml
snail_job:
  enabled: true
  server_host: 127.0.0.1
  server_port: "17888"
  host_ip: 127.0.0.1
  host_port: "17889"
  namespace: c8f1a2b3d4e5461789abcdef01234567
  group_name: hei_gin_admin
  token: SJ_heiGinAdminToken1234567890abcd
```

执行器随 `go run ./app/cmd/api` 启动；Handler 名须与控制台 `executor_info` 一致（`executor_type=3` Go）：

| JobHandler | 模块 |
|------------|------|
| `accountPurgeCancelledJob` | iam/account |
| `bannerStatusJob` | sys/banner |
| `auditAlertJob` | sys/audit |
