-- 设置管理员用户脚本
-- 将第一个注册用户设置为管理员

-- Usage: mysql -h 127.0.0.1 -u root -p"123456" -e "ALTER TABLE user ADD COLUMN is_admin TINYINT(1) DEFAULT 0;"

-- 设置第一个用户为管理员
UPDATE user SET is_admin =1 WHERE id = 1;

-- 查看所有用户
SELECT * FROM user;
