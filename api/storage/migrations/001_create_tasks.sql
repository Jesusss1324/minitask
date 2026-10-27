CREATE TABLE tasks (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,

    title VARCHAR(100) NOT NULL,

    description TEXT NOT NULL,

    status ENUM(
        'pending',
        'in_progress',
        'completed'
    ) NOT NULL DEFAULT 'pending',

    priority ENUM(
        'low',
        'medium',
        'high'
    ) NOT NULL DEFAULT 'medium',

    created_at DATETIME(3) NOT NULL,

    updated_at DATETIME(3) NOT NULL,

    PRIMARY KEY (id),

    CONSTRAINT chk_tasks_title_length CHECK (CHAR_LENGTH(TRIM(title)) BETWEEN 3 AND 100),

    CONSTRAINT chk_tasks_description_length CHECK (CHAR_LENGTH(description) <= 1000)
)
DEFAULT CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;