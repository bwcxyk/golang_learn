-- 将文件名设置为：ora.sql，删除所有多余的内容包括注释
-- 并更改语句合适的语句

SELECT ROWNUM ROWNUM1, name, age
FROM user
WHERE name = :keyword
