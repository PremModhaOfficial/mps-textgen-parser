-- ============================================================
-- _UserManagement_ Schema
-- Tables: _users_, _roles_, _user_roles_ (join)
-- ============================================================

CREATE TABLE IF NOT EXISTS _users_ (
    _id_            BIGSERIAL       PRIMARY KEY,
    _username_      TEXT            NOT NULL UNIQUE,
    _email_         TEXT            NOT NULL UNIQUE,
    _password_      TEXT            NOT NULL,
    _first_name_    TEXT            NOT NULL DEFAULT '',
    _last_name_     TEXT            NOT NULL DEFAULT '',
    _is_active_     BOOLEAN         NOT NULL DEFAULT true,
    _created_at_    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    _updated_at_    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS _roles_ (
    _id_            BIGSERIAL       PRIMARY KEY,
    _name_          TEXT            NOT NULL UNIQUE,
    _description_   TEXT            NOT NULL DEFAULT '',
    _created_at_    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS _user_roles_ (
    _user_id_       BIGINT          NOT NULL REFERENCES _users_(_id_) ON DELETE CASCADE,
    _role_id_       BIGINT          NOT NULL REFERENCES _roles_(_id_) ON DELETE CASCADE,
    _assigned_at_   TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (_user_id_, _role_id_)
);

-- Indexes
CREATE INDEX IF NOT EXISTS _idx_users_username_     ON _users_(_username_);
CREATE INDEX IF NOT EXISTS _idx_users_email_        ON _users_(_email_);
CREATE INDEX IF NOT EXISTS _idx_users_is_active_    ON _users_(_is_active_);
CREATE INDEX IF NOT EXISTS _idx_roles_name_         ON _roles_(_name_);
CREATE INDEX IF NOT EXISTS _idx_user_roles_user_id_ ON _user_roles_(_user_id_);
CREATE INDEX IF NOT EXISTS _idx_user_roles_role_id_ ON _user_roles_(_role_id_);
