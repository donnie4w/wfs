# 数据迁移工具

本工具用于将WFS元数据批量迁移至新数据库中，支持断点续传、进度记录和性能监控，适用于大规模数据迁移场景。

## 功能特性

- ✅ 支持断点续传：意外中断后可从上次位置继续迁移
- ✅ 自动记录迁移进度和最后处理的 key
- ✅ 实时输出迁移速度和进度
- ✅ 可配置每次迁移的批量大小（batch size）
- ✅ 迁移完成后生成完成标记文件

## 目录结构

假设你的数据目录如下：

```
/wfsdata/wfsdb/
├── MANIFEST-000001
├── CURRENT
├── LOCK
├── LOG
└── wfs.db           # 目标数据库文件
```

其中：
- `wfsdb/` 是 LevelDB 的数据目录
- `wfs.db` 是目标数据库文件（若不存在会自动创建）

## 使用方法

### 基本用法

```bash
./migration <leveldb_directory> <desc_db_file> [batch_size]
```

#### 参数说明：

| 参数 | 说明 |
|------|------|
| `leveldb_directory` | LevelDB 数据目录路径（包含 `.log`, `MANIFEST-*` 等文件） |
| `desc_db_file` | 目标数据库文件路径 |
| `batch_size` | （可选）每次迁移的记录数，默认为 `10000` |

### 示例命令

```bash
# 示例1：使用默认批量大小（10000）
./migration /wfsdata/wfsdb /wfsdata/wfsdb/wfs.db

# 示例2：指定批量大小为 5000
./migration /wfsdata/wfsdb /wfsdata/wfsdb/wfs.db 5000
```

## 断点续传机制

- 程序会自动在当前目录生成 `migration_state.json` 文件，记录：
- 每分钟自动保存一次状态
- 下次运行时会自动检测并从中断处继续

> ⚠️ 注意：请勿手动修改或删除 `migration_state.json`，否则可能导致重复迁移或数据丢失。

## 完成标记

迁移完成后，程序会生成一个 `migration_complete.txt` 文件，内容如下：

```
Migration completed at: 2025-09-28 10:20:30, Total records: 328450
```

可用于自动化脚本判断迁移是否成功完成。


## 注意事项

1. **确保 LevelDB 未被其他进程占用**，否则可能导致读取失败。
2. 批量大小（`batch_size`）可根据内存和性能需求调整：
    - 数值越大，内存占用越高，但吞吐量更高
    - 建议值：5000 ~ 50000