-- XXL-JOB Admin 表结构请使用官方/hei-boot 的 tables_xxl_job.sql 先初始化。
-- 本文件为 hei-gin 执行器与 Job Handler 种子（AppName = hei-gin-api）。

INSERT INTO xxl_job_group (id, app_name, title, address_type, address_list, update_time)
VALUES (3, 'hei-gin-api', 'HEI Gin API', 0, NULL, now())
ON CONFLICT (id) DO UPDATE SET
    app_name = EXCLUDED.app_name,
    title = EXCLUDED.title,
    update_time = EXCLUDED.update_time;

INSERT INTO xxl_job_info (
    id, job_group, job_desc, add_time, update_time, author, alarm_email,
    schedule_type, schedule_conf, misfire_strategy, executor_route_strategy,
    executor_handler, executor_param, executor_block_strategy, executor_timeout,
    executor_fail_retry_count, glue_type, glue_source, glue_remark, glue_updatetime,
    child_jobid, trigger_status, trigger_last_time, trigger_next_time
) VALUES
(200, 3, '清理超期已注销账号', now(), now(), 'hei', '',
 'CRON', '0 0 3 * * ?', 'DO_NOTHING', 'FIRST',
 'accountPurgeCancelledJob', '15', 'SERIAL_EXECUTION', 0,
 0, 'BEAN', '', 'GLUE代码初始化', now(),
 '', 0, 0, 0),
(201, 3, '按 start_at/end_at 同步 Banner 状态', now(), now(), 'hei', '',
 'CRON', '0 */5 * * * ?', 'DO_NOTHING', 'FIRST',
 'bannerStatusJob', '', 'SERIAL_EXECUTION', 0,
 0, 'BEAN', '', 'GLUE代码初始化', now(),
 '', 0, 0, 0)
ON CONFLICT (id) DO UPDATE SET
    job_desc = EXCLUDED.job_desc,
    executor_handler = EXCLUDED.executor_handler,
    schedule_conf = EXCLUDED.schedule_conf,
    update_time = EXCLUDED.update_time;
