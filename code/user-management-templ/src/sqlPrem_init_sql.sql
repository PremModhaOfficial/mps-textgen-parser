-- Auto-generated SQL for sqlPrem
-- Schema: iam

CREATE SCHEMA IF NOT EXISTS iam;

-- Table: iam.user
CREATE TABLE IF NOT EXISTS iam.user (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 password VARCHAR(255) NOT NULL,
 user_name VARCHAR(255) NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 user_mail VARCHAR(255) NOT NULL UNIQUE
);


-- Table: iam.roles
CREATE TABLE IF NOT EXISTS iam.roles (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 desc TEXT,
 name VARCHAR(255) NOT NULL
);


-- Junction table: iam.user_roless
CREATE TABLE IF NOT EXISTS iam.user_roless (
 user_id UUID NOT NULL REFERENCES iam.user(id) ON DELETE CASCADE,
 roles_id UUID NOT NULL REFERENCES iam.roles(id) ON DELETE CASCADE,
 assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY (user_id, roles_id)
);

