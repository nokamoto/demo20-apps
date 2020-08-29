CREATE DATABASE cloud;

USE cloud;

CREATE TABLE resourcemanager_project (
    project_key BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id VARCHAR(256) UNIQUE,
    display_name VARCHAR(256)
);

CREATE TABLE compute_instance (
    instance_key BIGINT PRIMARY KEY AUTO_INCREMENT,
    instance_id VARCHAR(256) UNIQUE,
    parent_id VARCHAR(256),
    labels VARCHAR(256),
    FOREIGN KEY (parent_id) REFERENCES resourcemanager_project (project_id)
);

CREATE TABLE iam_permission (
    permission_key BIGINT PRIMARY KEY AUTO_INCREMENT,
    permission_id VARCHAR(256) UNIQUE
);

CREATE TABLE rdb_cluster (
    cluster_key BIGINT PRIMARY KEY AUTO_INCREMENT,
    cluster_id VARCHAR(256) UNIQUE,
    parent_id VARCHAR(256),
    replicas INT,
    FOREIGN KEY (parent_id) REFERENCES resourcemanager_project (project_id)
);

CREATE TABLE rdb_cluster_instance (
    cluster_key BIGINT,
    instance_id VARCHAR(256),
    FOREIGN KEY (cluster_key) REFERENCES rdb_cluster (cluster_key),
    FOREIGN KEY (instance_id) REFERENCES compute_instance (instance_id)
);

INSERT INTO resourcemanager_project (
    project_id,
    display_name
) VALUES
    ('todo', 'todo display name');
