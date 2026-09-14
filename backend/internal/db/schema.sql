-- GROW Point System DDL (SQL Server)

CREATE TABLE categories (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    is_active BIT DEFAULT 1,
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE(),
    deleted_at DATETIME NULL
);

CREATE TABLE activities (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    category_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    default_points INT NOT NULL DEFAULT 0,
    is_custom_input BIT DEFAULT 0,
    is_active BIT DEFAULT 1,
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE(),
    deleted_at DATETIME NULL,
    CONSTRAINT FK_activities_categories FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE admin_users (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    npk VARCHAR(50) NOT NULL UNIQUE,
    is_active BIT DEFAULT 1,
    created_at DATETIME DEFAULT GETDATE(),
    deleted_at DATETIME NULL
);

CREATE TABLE activity_submissions (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    grow_id VARCHAR(50) NOT NULL,
    npk VARCHAR(50) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    department VARCHAR(100) NOT NULL,
    activity_id BIGINT NOT NULL,
    activity_date DATE NOT NULL,
    custom_reference VARCHAR(255) NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    points_awarded INT DEFAULT 0,
    admin_notes TEXT NULL,
    reviewed_by_npk VARCHAR(50) NULL,
    reviewed_at DATETIME NULL,
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE(),
    deleted_at DATETIME NULL,
    CONSTRAINT FK_submissions_activities FOREIGN KEY (activity_id) REFERENCES activities(id),
    CONSTRAINT UQ_submissions_grow_id UNIQUE (grow_id)
);

CREATE INDEX IX_submissions_npk ON activity_submissions(npk);
CREATE INDEX IX_submissions_status ON activity_submissions(status);
CREATE INDEX IX_submissions_activity_date ON activity_submissions(activity_date);

CREATE TABLE submission_evidence (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    submission_id BIGINT NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    original_filename VARCHAR(255) NULL,
    file_size_kb INT NULL,
    sort_order TINYINT DEFAULT 0,
    created_at DATETIME DEFAULT GETDATE(),
    deleted_at DATETIME NULL,
    CONSTRAINT FK_evidence_submission FOREIGN KEY (submission_id) REFERENCES activity_submissions(id)
);

CREATE INDEX IX_evidence_submission ON submission_evidence(submission_id);

CREATE TABLE rewards (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    points_required INT NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    is_active BIT DEFAULT 1,
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE(),
    deleted_at DATETIME NULL
);

CREATE TABLE reward_redemptions (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    npk VARCHAR(50) NOT NULL,
    reward_id BIGINT NOT NULL,
    points_spent INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PROCESSED',
    created_at DATETIME DEFAULT GETDATE(),
    deleted_at DATETIME NULL,
    CONSTRAINT FK_redemptions_rewards FOREIGN KEY (reward_id) REFERENCES rewards(id)
);

CREATE INDEX IX_redemptions_npk ON reward_redemptions(npk);

CREATE TABLE refresh_tokens (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    npk VARCHAR(50) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    expires_at DATETIME NOT NULL,
    revoked_at DATETIME NULL,
    created_at DATETIME DEFAULT GETDATE(),
    user_agent VARCHAR(255) NULL
);

CREATE INDEX IX_refresh_npk ON refresh_tokens(npk);

CREATE TABLE login_attempts (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    npk VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45) NULL,
    success BIT NOT NULL,
    created_at DATETIME DEFAULT GETDATE()
);

CREATE INDEX IX_login_attempts_npk_time ON login_attempts(npk, created_at);
