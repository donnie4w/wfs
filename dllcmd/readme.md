# WFS Go 动态链接库

一个用 Go 编写的高性能文件存储库，暴露为动态链接库（DLL/.so）供 C/C++ 应用程序使用。

## 功能特性

- **数据追加**: 向文件添加数据，支持可选压缩
- **数据读取**: 高效检索文件内容
- **文件操作**: 重命名和删除文件
- **文件检查**: 检查文件是否存在
- **键列表查询**: 分页获取文件列表
- **服务配置**: 灵活控制 HTTP、Thrift、Admin 服务开关
- **线程安全**: 支持并发访问
- **跨平台**: 支持 Windows、Linux、macOS

## 作用说明
1.   wfs的动态链接库（wfs.dll/wfs.so），可以将wfs内嵌到c/c++程序中，wfs的服务（HTTP、Thrift、Admin）通过初始化函数决定是否开启，默认开启。
2.   动态链接库依旧包含wfs所有的功能特性，如果服务没有开启，那么它会作为一个可以增删改查的库提供给c/c++程序调用。
3.   wfs动态链接库依旧可以读取wfs.json配置，配置数据库类型，配置文件参数等。

## 快速开始

###  C/C++ 调用示例

```c
#include <stdio.h>
#include "wfs.h"

int main() {
    // 初始化（使用默认配置，开启所有服务）
    char* result = Init(NULL);  //  也可以传入服务启动JSON结构字符串参数：{"http": false, "thrift": false, "admin": false}
    if (result) {
        printf("初始化失败: %s\n", result);
        FreeMemory(result);
        return 1;
    }
    
    // 检查初始化状态
    if (IsInit() != 1) {
        printf("未初始化成功\n");
        return 1;
    }
    
    // 使用示例
    unsigned char data[] = {1, 2, 3, 4, 5};
    result = Append("test.txt", data, 5, 0);
    if (result) {
        printf("操作失败: %s\n", result);
        FreeMemory(result);
    } else {
        printf("操作成功\n");
    }
    
    // 检查文件是否存在
    int exists = Has("test.txt");
    if (exists == 1) {
        printf("文件存在\n");
    }
    
    // 获取文件列表
    char* keysJSON = GetKeys(1, 10);  // 从Id：1开始，返回10个key
    if (keysJSON) {
        printf("文件列表: %s\n", keysJSON);
        FreeMemory(keysJSON);
    }
    
    // 清理
    Close();
    return 0;
}
```

## API 参考

### 初始化函数

#### `Init(configJSON)`
- **描述**: 初始化 WFS 库
- **参数**: `configJSON` - JSON配置字符串，NULL表示使用默认配置
- **返回**: `char*` - 成功返回 `NULL`，失败返回错误信息
- **配置示例**:
  ```json
  {
    "http": true,    // 开启HTTP服务
    "thrift": false, // 关闭Thrift服务  
    "admin": true    // 开启Admin服务
  }
  ```

#### `IsInit()`
- **描述**: 检查初始化状态
- **返回**: `int` - 1表示已初始化，0表示未初始化

#### `GetInitStatus()`
- **描述**: 获取初始化状态信息
- **返回**: `char*` - 状态信息字符串

#### `Close()`
- **描述**: 关闭库并释放资源
- **返回**: `void`

### 文件操作函数

#### `Append(name, data, dataLen, compress)`
- **描述**: 向文件追加数据
- **参数**:
  - `name`: 文件名
  - `data`: 数据指针
  - `dataLen`: 数据长度
  - `compress`: 压缩标志 (0=不压缩, 1=压缩)
- **返回**: `char*` - 成功返回 `NULL`，失败返回错误信息

#### `Get(path, resultLen)`
- **描述**: 读取文件数据
- **参数**:
  - `path`: 文件路径
  - `resultLen`: 返回数据长度的指针
- **返回**: `unsigned char*` - 数据指针，需要调用 `FreeMemory` 释放

#### `Delete(path)`
- **描述**: 删除文件
- **参数**: `path` - 文件路径
- **返回**: `char*` - 成功返回 `NULL`，失败返回错误信息

#### `Rename(path, newpath)`
- **描述**: 重命名文件
- **参数**:
  - `path`: 原文件路径
  - `newpath`: 新文件路径
- **返回**: `char*` - 成功返回 `NULL`，失败返回错误信息

#### `Has(path)`
- **描述**: 检查文件是否存在
- **参数**: `path` - 文件路径
- **返回**: `int` - 1表示存在，0表示不存在，-1表示未初始化

#### `GetKeys(fromId, limit)`
- **描述**: 分页获取文件列表
- **参数**:
  - `fromId`: 起始ID
  - `limit`: 返回数量限制
- **返回**: `char*` - JSON格式的文件列表，需要调用 `FreeMemory` 释放
- **返回格式**:
  ```json
  {
    "keys": [
      {"name": "file1.txt", "id": 1},
      {"name": "file2.txt", "id": 2}
    ]
  }
  ```

### 工具函数

#### `FreeMemory(ptr)`
- **描述**: 释放由 `Get` 函数分配的内存
- **参数**: `ptr` - 要释放的内存指针
- **返回**: `void`

#### `FreeString(str)`
- **描述**: 释放由库分配的字符串内存
- **参数**: `str` - 要释放的字符串指针
- **返回**: `void`

## 错误处理

所有函数遵循相同的错误处理模式：
- 成功时返回 `NULL` 或相应的成功值
- 失败时返回错误信息字符串，需要调用 `FreeMemory` 或 `FreeString` 释放

```c
char* result = Append("file.txt", data, length, 0);
if (result) {
    printf("错误: %s\n", result);
    FreeString(result);
} else {
    printf("成功\n");
}

unsigned char* data = Get("file.txt", &len);
if (data) {
    // 使用数据
    FreeMemory(data);  // 必须释放！
}
```

## 服务配置

### 配置选项

| 服务 | 配置字段 | 默认值 | 说明           |
|------|----------|--------|--------------|
| HTTP | `http` | `true` | HTTP文件访问服务   |
| Thrift | `thrift` | `true` | Thrift RPC服务 |
| Admin | `admin` | `true` | 管理控制台服务      |

### 配置示例

```c
// 只开启HTTP服务
char* config = "{\"thrift\":false,\"admin\":false}";
Init(config);

// 只开启Thrift和Admin服务
char* config = "{\"http\":false}";
Init(config);

// 使用默认配置（开启所有服务）
Init(NULL);
```

## 平台支持

| 平台 | 输出文件 | 编译命令 |
|------|----------|----------|
| Windows | `wfs.dll` | `go build -buildmode=c-shared -o wfs.dll dll_wfs.go` |
| Linux | `wfs.so` | `go build -buildmode=c-shared -o wfs.so dll_wfs.go` |
| macOS | `wfs.dylib` | `go build -buildmode=c-shared -o wfs.dylib dll_wfs.go` |

## 内存管理

**重要**: 由库返回的字符串和二进制数据必须正确释放：

```c
// 二进制数据使用 FreeMemory  
unsigned char* data = Get("file.txt", &len);
if (data) {
    // 使用数据
    FreeMemory(data);    // 释放
}

// JSON字符串使用 FreeMemory 或 FreeString
char* json = GetKeys(0, 10);
if (json) {
    // 处理JSON
    FreeMemory(json);    // 两种方式都可以
}
```

## 常见问题

### 1. 初始化失败
- 确保先调用 `Init()` 再使用其他函数
- 检查 `IsInit()` 返回值确认初始化状态

### 2. 内存泄漏
- 所有由库返回的指针都必须正确释放
- 字符串使用 `FreeString`，二进制数据使用 `FreeMemory`

### 3. 配置解析错误
- 确保JSON格式正确
- 未设置的字段使用默认值 `true`

### 4. 线程安全
- 库本身是线程安全的
- 但建议对共享资源进行适当的同步控制


## 技术支持

如有问题请提交 Issue 或联系: donnie4w@gmail.com