# 生产化验收清单

上线前至少逐项确认：

- 配置文件存在，且 `app / db / redis / token` 关键字段完整
- 服务在配置缺失场景会明确启动失败，不会以空配置继续运行
- `plugin-im` 业务层不再直接访问 `GlobalCrossHub`
- WebSocket `allowed_origins` 已配置，未授权来源无法连接
- WebSocket `trusted_proxies` 已按部署拓扑配置，伪造转发头不会绕过限流
- `/health/live` 返回 200
- `/health/ready` 在依赖正常时返回 200，在 MySQL / Redis / 插件异常时返回 503
- `make test-all` 通过
- 迁移命令在目标环境可执行，且配置缺失时会失败退出
- Swagger 是否启用符合环境要求，生产若启用则已配置 Basic Auth

建议补充记录：

- 最近一次全量测试时间
- 最近一次迁移执行时间
- WS 来源白名单与可信代理配置值
- readiness 接入的网关或编排平台配置
