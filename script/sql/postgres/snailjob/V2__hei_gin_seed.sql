-- HEI Gin SnailJob seed: independent namespace + group + three Go jobs + console password admin/123456
-- Console hashes as SHA256(MD5(plain)); password below is SHA256(MD5('123456')).

UPDATE sj_system_user
SET password = 'cdf4a007e2b02a0c49fc9b7ccfbb8a10c644f635e1765dcf2a7ab794ddc7edac',
    update_dt = now()
WHERE username = 'admin';

INSERT INTO sj_namespace (name, unique_id, description, create_dt, update_dt, deleted)
SELECT
    'HEI Gin',
    'c8f1a2b3d4e5461789abcdef01234567',
    'HEI Gin isolated namespace',
    now(),
    now(),
    0
WHERE NOT EXISTS (
    SELECT 1 FROM sj_namespace WHERE unique_id = 'c8f1a2b3d4e5461789abcdef01234567'
);

INSERT INTO sj_group_config (
    namespace_id, group_name, description, token, group_status, version,
    group_partition, id_generator_mode, init_scene, create_dt, update_dt
)
SELECT
    'c8f1a2b3d4e5461789abcdef01234567',
    'hei_gin_admin',
    'HEI Gin Admin executor group',
    'SJ_heiGinAdminToken1234567890abcd',
    1,
    1,
    0,
    1,
    1,
    now(),
    now()
WHERE NOT EXISTS (
    SELECT 1 FROM sj_group_config
    WHERE namespace_id = 'c8f1a2b3d4e5461789abcdef01234567'
      AND group_name = 'hei_gin_admin'
);

-- executor_type=3 Go; trigger_type=1 CRON; job_status=1 enabled; task_type=1 cluster; route_key=4
INSERT INTO sj_job (
    namespace_id, biz_id, group_name, job_name, args_str, args_type,
    next_trigger_at, job_status, task_type, route_key, executor_type, executor_info,
    trigger_type, trigger_interval, block_strategy, executor_timeout, max_retry_times,
    parallel_num, retry_interval, bucket_index, resident, notify_ids, description, deleted,
    create_dt, update_dt
)
SELECT
    'c8f1a2b3d4e5461789abcdef01234567',
    'hei-gin-accountPurgeCancelledJob',
    'hei_gin_admin',
    '清理超期已注销账号',
    '15',
    1,
    (EXTRACT(EPOCH FROM now()) * 1000)::bigint,
    1, 1, 4, 3, 'accountPurgeCancelledJob',
    1, '0 0 3 * * ?', 1, 0, 0,
    1, 0, 0, 0, '', 'Purge cancelled accounts past retention', 0,
    now(), now()
WHERE NOT EXISTS (
    SELECT 1 FROM sj_job
    WHERE namespace_id = 'c8f1a2b3d4e5461789abcdef01234567'
      AND biz_id = 'hei-gin-accountPurgeCancelledJob'
);

INSERT INTO sj_job (
    namespace_id, biz_id, group_name, job_name, args_str, args_type,
    next_trigger_at, job_status, task_type, route_key, executor_type, executor_info,
    trigger_type, trigger_interval, block_strategy, executor_timeout, max_retry_times,
    parallel_num, retry_interval, bucket_index, resident, notify_ids, description, deleted,
    create_dt, update_dt
)
SELECT
    'c8f1a2b3d4e5461789abcdef01234567',
    'hei-gin-bannerStatusJob',
    'hei_gin_admin',
    '同步 Banner 状态',
    NULL,
    1,
    (EXTRACT(EPOCH FROM now()) * 1000)::bigint,
    1, 1, 4, 3, 'bannerStatusJob',
    1, '0 */5 * * * ?', 1, 0, 0,
    1, 0, 0, 0, '', 'Sync banner ENABLED/DISABLED by start_at/end_at', 0,
    now(), now()
WHERE NOT EXISTS (
    SELECT 1 FROM sj_job
    WHERE namespace_id = 'c8f1a2b3d4e5461789abcdef01234567'
      AND biz_id = 'hei-gin-bannerStatusJob'
);

INSERT INTO sj_job (
    namespace_id, biz_id, group_name, job_name, args_str, args_type,
    next_trigger_at, job_status, task_type, route_key, executor_type, executor_info,
    trigger_type, trigger_interval, block_strategy, executor_timeout, max_retry_times,
    parallel_num, retry_interval, bucket_index, resident, notify_ids, description, deleted,
    create_dt, update_dt
)
SELECT
    'c8f1a2b3d4e5461789abcdef01234567',
    'hei-gin-auditAlertJob',
    'hei_gin_admin',
    '审计量级告警',
    NULL,
    1,
    (EXTRACT(EPOCH FROM now()) * 1000)::bigint,
    1, 1, 4, 3, 'auditAlertJob',
    1, '0 */2 * * * ?', 1, 0, 0,
    1, 0, 0, 0, '', 'Audit volume alert into sys_alert_log', 0,
    now(), now()
WHERE NOT EXISTS (
    SELECT 1 FROM sj_job
    WHERE namespace_id = 'c8f1a2b3d4e5461789abcdef01234567'
      AND biz_id = 'hei-gin-auditAlertJob'
);
