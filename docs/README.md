# Kafka-King 文档目录

本目录包含 Kafka-King 项目的各类文档。

## 📚 文档列表

### 📖 [幂等性禁用功能使用指南](./IDEMPOTENCE_GUIDE.md)
**用途**: 解决 Kafka 3.0+ 客户端连接低版本服务端的兼容性问题

**内容包括：**
- 功能说明和使用场景
- 详细的操作步骤（带图文说明）
- 配置示例和验证方法
- 故障排除指南
- 测试验证结果
- 原理说明

**适用人群：**
- 遇到 "Cluster authorization failed" 错误的用户
- 需要连接低版本 Kafka 服务端的用户
- 想了解幂等性配置的开发者

---

### 📷 截图目录 (./snap/)

包含项目功能截图，用于展示和文档说明。

---

### 🌍 多语言 README (./readme/)

- `readme-zh.md` - 简体中文
- `readme-ja.md` - 日本语
- `readme-ru.md` - русский
- `readme-ko.md` - 한국어

---

## 🔍 快速查找

### 我遇到了 "Cluster authorization failed" 错误
👉 查看 [幂等性禁用功能使用指南](./IDEMPOTENCE_GUIDE.md)

### 我想了解如何配置连接
👉 查看 [根目录 README](../readme.md) 的"Download"部分

### 我想了解如何构建项目
👉 查看 [构建指南](../BUILD.md)

### 我想查看测试用例
👉 查看 [单元测试文件](../app/backend/service/idempotence_test.go)

---

## 📝 文档贡献

如果您发现文档有误或需要改进，欢迎：

1. 提交 Issue: https://github.com/Bronya0/Kafka-King/issues
2. 提交 Pull Request
3. 发送邮件至: tangssst@163.com

---

## 📄 许可证

所有文档遵循 Apache License 2.0 许可证。
