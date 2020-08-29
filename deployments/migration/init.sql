CREATE DATABASE cloud;

USE cloud;

CREATE TABLE resourcemanager_project (
    project_key BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id VARCHAR(256),
    display_name VARCHAR(256)
);

CREATE TABLE compute_instance (
    instance_key BIGINT PRIMARY KEY AUTO_INCREMENT,
    instance_id VARCHAR(256),
    parent_key BIGINT,
    labels VARCHAR(256),
    FOREIGN KEY (parent_key) REFERENCES resourcemanager_project (project_key)
);

CREATE TABLE iam_permission (
    permission_key BIGINT PRIMARY KEY AUTO_INCREMENT,
    permission_id VARCHAR(256)
);

INSERT INTO resourcemanager_project (
    project_id,
    display_name
) VALUES
    ('todo', 'todo display name');
