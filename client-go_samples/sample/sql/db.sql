DROP TABLE IF EXISTS t_job;
CREATE TABLE t_job(  
    `id` int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `job_name` VARCHAR(100) NOT NULL COMMENT 'job name',
    `pod_name` VARCHAR(100) NOT NULL COMMENT 'pod name',
    `node_name` VARCHAR(100) NOT NULL COMMENT 'node name',
    `status` INT(4) NOT NULL COMMENT 'job status',
    `create_time` DATETIME COMMENT 'Create Time'
) COMMENT '';

DROP TABLE IF EXISTS t_node;
CREATE TABLE t_node(  
    `id` int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `node_name` VARCHAR(100) NOT NULL COMMENT 'node name',
    `host` VARCHAR(100) NOT NULL COMMENT 'host',
    `create_time` DATETIME COMMENT 'Create Time'
) COMMENT '';

insert into 
  t_node (
    node_name, 
    host, 
    create_time
  )
values
  (
    'k8s-worker1.sdns.dev.cloud', 
    '10.69.82.12', 
    now()
  ),  (
    'k8s-worker2.sdns.dev.cloud', 
    '10.69.64.134', 
    now()
  ),  (
    'k8s-worker3.sdns.dev.cloud', 
    '10.69.65.110', 
    now()
  ),  (
    'k8s-worker4.sdns.dev.cloud', 
    '36.103.180.175', 
    now()
  ),  (
    'k8s-worker-4090-1.sdns.dev.cloud', 
    '36.103.180.180', 
    now()
  ),  (
    'k8s-worker-4090-2.sdns.dev.cloud', 
    '36.103.180.181', 
    now()
  ),  (
    'k8s-worker-4090-3.sdns.dev.cloud', 
    '36.103.180.182', 
    now()
  ),  (
    'k8s-worker-4090-4.sdns.dev.cloud', 
    '36.103.180.183', 
    now()
  );

  insert into 
    t_job (
      job_name, 
      pod_name, 
      node_name, 
      status, 
      create_time
    )
  values
    (
      'pytorch-simple', 
      'pytorch-simple-master-0', 
      'k8s-worker3.sdns.dev.cloud', 
      1, 
      now()
    ),(
      'pytorch-simple', 
      'pytorch-simple-worker-0', 
      'k8s-worker2.sdns.dev.cloud', 
      1, 
      now()
    ),(
      'pytorch-simple', 
      'pytorch-simple-worker-1', 
      'k8s-worker3.sdns.dev.cloud', 
      1, 
      now()
    ),(
      'pytorch-simple', 
      'pytorch-simple-worker-2', 
      'k8s-worker4.sdns.dev.cloud', 
      1, 
      now()
    ),(
      'pytorch-hard', 
      'pytorch-hard-master-0', 
      'k8s-worker-4090-1.sdns.dev.cloud', 
      1, 
      now()
    ),(
      'pytorch-hard', 
      'pytorch-hard-worker-0', 
      'k8s-worker-4090-4.sdns.dev.cloud', 
      1, 
      now()
    );