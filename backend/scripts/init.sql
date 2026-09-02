-- =============================================================
-- Evenstar Writings 数据库初始化脚本
-- 设计依据：项目设计快照文档 第三节（数据库设计）
-- 说明：
--   1. 本脚本是表结构/外键/索引的唯一权威来源，GORM AutoMigrate 仅作辅助
--   2. 本脚本不含管理员记录、不含任何密码或敏感凭证
--      （管理员由后端首次启动时通过环境变量 EVENSTAR_ADMIN_* 引导创建）
--   3. 软删除表：users / articles / comments（deleted_at）
--      硬删除表：categories / tags / 中间表 / uploads
-- =============================================================

CREATE DATABASE IF NOT EXISTS evenstar
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE evenstar;

-- -------------------------------------------------------------
-- 用户表
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  username      VARCHAR(50)  NOT NULL,
  password      VARCHAR(100) NOT NULL COMMENT 'bcrypt',
  nickname      VARCHAR(50)  NULL COMMENT '可空，空则后端生成默认昵称',
  email         VARCHAR(100) NOT NULL,
  avatar        VARCHAR(255) NULL COMMENT '可空，空则前端默认头像',
  role          TINYINT      NOT NULL DEFAULT 2 COMMENT '1=管理员 2=普通用户',
  status        TINYINT      NOT NULL DEFAULT 1 COMMENT '1=正常 2=禁用',
  token_version INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '改密/重置/禁用/删除时+1，JWT 失效机制',
  created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at    DATETIME     NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_username (username),
  UNIQUE KEY uk_email (email),
  KEY idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- -------------------------------------------------------------
-- 文章表
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS articles (
  id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  author_id    BIGINT UNSIGNED NOT NULL,
  title        VARCHAR(200) NOT NULL,
  summary      VARCHAR(500) NULL,
  content      LONGTEXT     NULL COMMENT 'Markdown 原文',
  cover        VARCHAR(255) NULL COMMENT '可空，空则前端默认封面',
  status       TINYINT      NOT NULL DEFAULT 1 COMMENT '1=草稿 2=发布',
  views        INT UNSIGNED NOT NULL DEFAULT 0,
  published_at DATETIME     NULL COMMENT 'status 首次变 2 时写入，前台排序依据',
  created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at   DATETIME     NULL,
  PRIMARY KEY (id),
  KEY idx_status_published (status, published_at),
  KEY idx_author_id (author_id),
  KEY idx_deleted_at (deleted_at),
  CONSTRAINT fk_articles_author FOREIGN KEY (author_id)
    REFERENCES users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- -------------------------------------------------------------
-- 分类表 / 标签表（硬删除）
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS categories (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name       VARCHAR(50) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_name (name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS tags (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name       VARCHAR(50) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_name (name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- -------------------------------------------------------------
-- 文章-分类 / 文章-标签 中间表（联合主键，硬删除）
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS article_categories (
  article_id  BIGINT UNSIGNED NOT NULL,
  category_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (article_id, category_id),
  KEY idx_category_id (category_id),
  CONSTRAINT fk_ac_article FOREIGN KEY (article_id)
    REFERENCES articles (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_ac_category FOREIGN KEY (category_id)
    REFERENCES categories (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS article_tags (
  article_id BIGINT UNSIGNED NOT NULL,
  tag_id     BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (article_id, tag_id),
  KEY idx_tag_id (tag_id),
  CONSTRAINT fk_at_article FOREIGN KEY (article_id)
    REFERENCES articles (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_at_tag FOREIGN KEY (tag_id)
    REFERENCES tags (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- -------------------------------------------------------------
-- 评论表
--   最多两级：parent_id 恒指一级评论（一级为 NULL）
--   reply_to_id 仅展示用（张三 → 李四），不加外键
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS comments (
  id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  article_id  BIGINT UNSIGNED NOT NULL,
  user_id     BIGINT UNSIGNED NULL COMMENT '可空，用户注销后为 NULL，前台显示「已注销用户」',
  parent_id   BIGINT UNSIGNED NULL COMMENT '一级=NULL；二级=所属一级评论 id（恒指一级，不产生三级）',
  reply_to_id BIGINT UNSIGNED NULL COMMENT '实际回复对象，仅展示用，无外键',
  content     VARCHAR(400) NOT NULL COMMENT '评论内容，最多 400 字',
  is_top      TINYINT       NOT NULL DEFAULT 0 COMMENT '1=置顶（仅一级评论可置顶）',
  top_time    DATETIME      NULL COMMENT '置顶时间，多条置顶按它倒序',
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at  DATETIME      NULL,
  PRIMARY KEY (id),
  KEY idx_article_id (article_id),
  KEY idx_parent_id (parent_id),
  KEY idx_deleted_at (deleted_at),
  CONSTRAINT fk_comments_article FOREIGN KEY (article_id)
    REFERENCES articles (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_comments_user FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT fk_comments_parent FOREIGN KEY (parent_id)
    REFERENCES comments (id) ON DELETE NO ACTION ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- -------------------------------------------------------------
-- 上传文件表（硬删除；删除流程：先删 OSS 成功后才删本表记录）
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS uploads (
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id    BIGINT UNSIGNED NOT NULL,
  scene      VARCHAR(16)  NOT NULL COMMENT 'article / avatar',
  filename   VARCHAR(255) NOT NULL,
  url        VARCHAR(512) NOT NULL,
  size       BIGINT       NOT NULL DEFAULT 0 COMMENT '字节',
  mime       VARCHAR(100) NULL,
  created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_user_id (user_id),
  KEY idx_scene (scene),
  CONSTRAINT fk_uploads_user FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE NO ACTION ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
