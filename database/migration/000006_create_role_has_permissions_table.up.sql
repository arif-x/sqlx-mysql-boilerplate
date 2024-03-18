CREATE TABLE IF NOT EXISTS role_has_permissions (
    role_uuid CHAR(36) NOT NULL,
    permission_uuid CHAR(36) NOT NULL,
    PRIMARY KEY (role_uuid, permission_uuid)
);
