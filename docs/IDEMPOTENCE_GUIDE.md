# Kafka-King 生产者幂等性禁用功能使用指南

## 📖 功能说明

本功能用于解决 **Kafka 3.0+ 客户端连接低版本 Kafka 服务端**时的兼容性问题。

### 问题场景

当使用 Kafka 3.0+ 客户端连接到 Kafka 2.x 或更早版本的服务端时，可能会遇到以下错误：

```
Cluster authorization failed
```

**原因：**
- Kafka 3.0+ 客户端默认启用**幂等性生产者**（Idempotent Producer）
- 低版本 Kafka 服务端上的用户默认**没有幂等写入权限**（IdempotentWrite）
- 导致连接失败

### 解决方案

通过 Kafka-King 提供的**"禁用生产者幂等性"**开关，可以在客户端侧关闭幂等性，绕过权限问题。

---

## 🎯 操作步骤

### 步骤 1: 打开连接管理页面

1. 启动 Kafka-King 应用
2. 点击左侧菜单的 **"集群"** 或 **"Cluster"** 选项
3. 如果是新建连接，点击 **"添加集群"** 按钮
4. 如果是编辑现有连接，点击对应连接卡片上的 **"编辑"** 按钮

### 步骤 2: 找到幂等性配置开关

在连接配置抽屉中，向下滚动找到以下配置项：

```
┌─────────────────────────────────────────────────┐
│ 基本配置                                         │
│ ├─ 昵称: [___________________]                 │
│ ├─ 连接地址: [___________________]              │
│ ├─ 使用 TLS: [○]                               │
│ ├─ 使用 SASL: [○]                              │
│                                                  │
│ 兼容性配置                                       │
│ ├─ 禁用生产者幂等性: [○]  ← 这里！              │
│ │   用于解决高版本客户端连接低版本服务端的       │
│ │   兼容性问题。当Kafka 3.0+客户端连接低版本    │
│ │   服务端时，如遇到"Cluster authorization      │
│ │   failed"错误，请开启此选项。                 │
│                                                  │
│ ├─ kerberos: [○]                                │
└─────────────────────────────────────────────────┘
```

**位置说明：**
- 该开关位于 **SASL 配置之后**，**Kerberos 配置之前**
- 标签：**"禁用生产者幂等性"** (中文) / **"Disable Producer Idempotence"** (英文)
- 提示文字会说明使用场景

### 步骤 3: 启用禁用幂等性开关

1. 点击 **"禁用生产者幂等性"** 右侧的开关按钮
2. 开关变为 **蓝色/激活状态** 表示已启用
3. 此时将**禁用生产者幂等性**，允许连接到低版本 Kafka 服务端

### 步骤 4: 测试连接

1. 点击抽屉底部的 **"连接测试"** 按钮
2. 等待测试结果：
   - ✅ **成功**: 显示 "连接成功" 消息
   - ❌ **失败**: 显示具体的错误信息

### 步骤 5: 保存配置

1. 确认测试成功后，点击 **"保存"** 按钮
2. 配置将持久化保存到 `~/.kafka-king/config.yaml`
3. 下次使用该连接时，会自动应用幂等性禁用设置

### 步骤 6: 使用连接

1. 返回集群列表页面
2. 点击刚才配置的集群卡片
3. 连接时会自动应用**禁用幂等性**配置
4. 现在可以正常使用生产者功能了

---

## 📋 配置示例

### YAML 配置文件示例

保存后的配置文件 (`~/.kafka-king/config.yaml`) 格式：

```yaml
width: 1600
height: 870
language: zh
theme: dark
connects:
  - id: 1
    name: "生产环境 Kafka 2.8"
    bootstrap_servers: "broker1:9092,broker2:9092,broker3:9092"
    tls: "disable"
    sasl: "enable"
    sasl_mechanism: "PLAIN"
    sasl_user: "producer_user"
    sasl_pwd: "********"
    disable_idempotence: "enable"  # ← 关键配置：禁用幂等性
    use_ssh: "disable"
```

### JSON 请求示例（内部数据流）

前端发送给后端的配置：

```json
{
  "id": 1,
  "name": "生产环境 Kafka 2.8",
  "bootstrap_servers": "broker1:9092,broker2:9092",
  "tls": "disable",
  "sasl": "enable",
  "sasl_mechanism": "PLAIN",
  "sasl_user": "producer_user",
  "sasl_pwd": "password123",
  "disable_idempotence": "enable",
  "use_ssh": "disable"
}
```

---

## 🔍 验证功能是否生效

### 方法 1: 检查配置文件

查看配置文件是否包含 `disable_idempotence: "enable"`：

```bash
cat ~/.kafka-king/config.yaml | grep disable_idempotence
```

**预期输出：**
```
disable_idempotence: enable
```

### 方法 2: 查看应用日志

启动 Kafka-King 后，查看日志（如果有的话）是否包含幂等性禁用的相关信息。

### 方法 3: 测试连接

1. 使用配置连接到低版本 Kafka 服务端
2. 如果之前失败，现在应该能成功连接
3. 尝试发送消息，验证生产者功能正常

---

## ⚠️ 注意事项

### 何时应该启用此选项？

**建议启用的场景：**
- ✅ Kafka 服务端版本 < 3.0（如 2.8.x, 2.7.x）
- ✅ 遇到 "Cluster authorization failed" 错误
- ✅ 用户没有 IdempotentWrite 权限
- ✅ 无法在服务端授予幂等写入权限

**建议不启用的场景：**
- ❌ Kafka 服务端版本 >= 3.0
- ❌ 用户已有幂等写入权限
- ❌ 需要确保消息精确一次投递（Exactly-Once）

### 禁用幂等性的影响

**优点：**
- ✅ 解决版本兼容性问题
- ✅ 可以连接到低版本服务端
- ✅ 不需要修改服务端权限配置

**缺点：**
- ⚠️ 失去幂等性保证
- ⚠️ 可能产生重复消息（网络重试时）
- ⚠️ 无法使用事务性生产者

**建议：**
如果您的业务场景对消息唯一性有严格要求，建议：
1. 升级 Kafka 服务端到 3.0+
2. 或在服务端为用户授予 IdempotentWrite 权限

---

## 🧪 测试验证

### 单元测试

项目包含完整的单元测试，验证功能正确性：

```bash
cd app/backend/service
go test -v -run TestDisableIdempotenceConfiguration
```

**测试覆盖：**
- ✅ 配置未设置时的默认行为
- ✅ disable_idempotence="disable" 的行为
- ✅ disable_idempotence="enable" 的行为
- ✅ 与其他配置（SASL, TLS）的兼容性
- ✅ 向后兼容性（旧配置）
- ✅ UI 数据流验证
- ✅ Kafka 版本兼容性场景

### 手动测试步骤

1. **准备测试环境**
   - Kafka 2.8.x 服务端（或其他低版本）
   - 创建测试用户（无 IdempotentWrite 权限）

2. **测试场景 1: 不启用开关（预期失败）**
   - 创建连接，不启用"禁用生产者幂等性"
   - 尝试连接
   - **预期结果**: 报错 "Cluster authorization failed"

3. **测试场景 2: 启用开关（预期成功）**
   - 编辑连接，启用"禁用生产者幂等性"
   - 尝试连接
   - **预期结果**: 连接成功
   - 尝试发送消息
   - **预期结果**: 消息发送成功

---

## 🐛 故障排除

### 问题 1: 找不到幂等性开关

**可能原因：**
- 使用的版本太旧，不包含此功能
- 前端构建未包含最新代码

**解决方案：**
1. 确认使用的是包含此功能的版本（检查 git commit）
2. 重新构建应用：
   ```bash
   cd app
   wails build -clean
   ```
3. 或从 GitHub Actions 下载最新构建的版本

### 问题 2: 开关设置后不生效

**可能原因：**
- 配置未保存
- 使用了旧的连接配置

**解决方案：**
1. 确认点击了"保存"按钮
2. 检查配置文件：`cat ~/.kafka-king/config.yaml`
3. 确认 `disable_idempotence: "enable"` 存在
4. 重新选择连接（点击集群卡片）

### 问题 3: 启用后仍然报错

**可能原因：**
- 不是幂等性权限问题，而是其他原因
- 网络连接问题
- SASL 认证配置错误

**解决方案：**
1. 查看完整的错误信息
2. 检查网络连接：`telnet broker-host 9092`
3. 验证 SASL 用户名密码正确
4. 检查 TLS 配置（如果使用 TLS）

---

## 📚 参考资料

### 官方文档

- [Apache Kafka Producer Configuration](https://kafka.apache.org/documentation/#producerconfigs)
- [Kafka Idempotent Producer](https://kafka.apache.org/documentation/#semantics)

### 相关 Issue

- [Cloudera CDP Known Issues](https://docs.cloudera.com/cdp-private-cloud-base/7.1.8/runtime-release-notes/topics/rt-pvc-known-issues-kafka.html)
- [Apache JIRA RANGER-3809](https://issues.apache.org/jira/browse/RANGER-3809)

### Spring Boot 等价配置

在 Spring Boot 中，等价的配置为：

```yaml
spring:
  kafka:
    producer:
      properties:
        enable.idempotence: false
```

---

## 🎓 原理说明

### 幂等性生产者是什么？

幂等性生产者确保即使消息重试，也不会在 Kafka 中产生重复消息。

**工作原理：**
1. 每个生产者会被分配一个唯一的 Producer ID (PID)
2. 每条消息带有序列号 (Sequence Number)
3. Broker 会检测并去重重复的消息

**要求：**
- 用户需要有 `IdempotentWrite` 权限
- Kafka 版本 >= 0.11.0

### 本功能实现原理

**后端实现** (`app/backend/service/kafka.go:346`)：

```go
if connCopy["disable_idempotence"] == "enable" {
    config = append(config, kgo.DisableIdempotentWrite())
}
```

**前端实现** (`app/frontend/src/components/Conn.vue:189-194`)：

```vue
<n-form-item :label="t('conn.disable_idempotence')" path="disable_idempotence">
  <n-flex vertical>
    <n-switch :round="false"
              checked-value="enable"
              unchecked-value="disable"
              v-model:value="currentNode.disable_idempotence"/>
    <n-text depth="3" style="font-size: 12px;">
      {{ t('conn.disable_idempotence_tip') }}
    </n-text>
  </n-flex>
</n-form-item>
```

**数据流：**
```
用户界面 (Vue.js)
    ↓ (v-model binding)
currentNode.disable_idempotence = "enable"
    ↓ (保存配置)
config.yaml
    ↓ (读取配置)
SetConnect(conn map[string]any)
    ↓ (应用配置)
kgo.DisableIdempotentWrite()
    ↓ (创建客户端)
Kafka Producer (幂等性已禁用)
```

---

## ✅ 测试结果

所有单元测试通过：

```
✓ 幂等性未配置 - 默认启用
✓ 幂等性配置为 disable - 启用幂等
✓ 幂等性配置为 enable - 禁用幂等
✓ 完整配置 - 禁用幂等
✓ 配置持久化验证
✓ 向后兼容性验证
✓ UI 数据流验证
✓ Kafka 3.0+ 客户端连接 2.x 服务端
✓ Kafka 3.0+ 客户端连接 3.0+ 服务端
✓ Kafka 2.x 客户端连接 2.x 服务端
```

---

## 💬 常见问题

**Q: 这个功能会影响现有的连接吗？**
A: 不会。现有连接如果没有设置 `disable_idempotence`，将使用默认设置（启用幂等性）。

**Q: 可以动态切换吗？**
A: 可以。编辑连接配置，修改开关后保存，重新连接即可生效。

**Q: 对性能有影响吗？**
A: 配置检查的性能影响可以忽略不计（纳秒级）。禁用幂等性后，可能会略微提升吞吐量（减少了序列号检查）。

**Q: 生产环境建议使用吗？**
A: 仅在遇到版本兼容性问题时使用。如果可以，建议升级服务端或授予权限，保持幂等性。

---

## 📧 反馈与支持

如果您在使用过程中遇到问题，请：

1. 检查本文档的"故障排除"部分
2. 查看 GitHub Issues: https://github.com/Bronya0/Kafka-King/issues
3. 提交新的 Issue（附带错误日志和配置信息）

---

**文档版本**: v1.0
**最后更新**: 2025-11-13
**适用版本**: Kafka-King v1.1.0+
