# ora_query

一个简单的 Oracle 批量查询工具：
- 从 Excel 的 `数据源` 工作表读取首列关键字
- 使用 `search.sql` 执行查询（并发）
- 将结果写回同一个 Excel 文件的 `查询结果` 工作表

## 目录说明

- `oracle.go`：主程序
- `config.yaml`：数据库连接配置
- `config_example.yaml`：配置示例
- `search.sql`：实际执行的 SQL
- `search_template.sql`：SQL 模板参考

## 环境要求

- Go 1.22+（建议）
- 可访问的 Oracle 数据库

## 配置

1. 复制并修改配置：

```yaml
database:
  host: 192.168.1.70
  port: 1521
  service_name: ORCL
  user: test_user
  password: 123456
```

2. 保存为 `config.yaml`（与可执行文件/源码在同一目录）。

## Excel 格式要求

- 必须存在工作表：`数据源`
- 第 1 行为表头（会跳过）
- 从第 2 行开始读取第一列作为关键字
- 空行或空白内容会自动跳过

## SQL 要求

- SQL 文件名固定为：`search.sql`
- SQL 中请使用命名参数 `:keyword`
- 例如：

```sql
SELECT *
FROM your_table
WHERE your_column = :keyword
```

注意：SQL 命名参数需与代码一致，当前使用的是 `:keyword`。

## 运行方式

在 `tools/ora_query` 目录执行：

```bash
go run oracle.go 1.xlsx
```

或先编译再运行：

```bash
go build -o oracle.exe oracle.go
./oracle.exe 1.xlsx
```

## 输出结果

- 程序会在输入的 Excel 文件中创建/写入工作表 `查询结果`
- 首行为查询返回的列名
- 后续行为每个关键字对应的查询结果

## 常见问题

1. `请提供 Excel 文件名作为命令行参数`
   - 运行时缺少文件名参数。

2. `无法读取 SQL 文件`
   - 确认 `search.sql` 存在且在当前工作目录。

3. 查询报错但程序继续
   - 程序会记录日志并写入空结果行，优先检查 SQL 和参数名是否匹配（`keyword`）。
