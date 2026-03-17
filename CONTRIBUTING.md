# 贡献指南

## 快速开始

### 前置要求
- Go 1.21 或更高版本
- Bash shell

### 运行脚本

#### 1. 运行测试
```bash
bash script/test.sh
```

#### 2. 编译应用
```bash
bash script/build.sh
```

#### 3. 运行应用
```bash
bash script/run.sh
```

## 工作流程

1. **测试代码**
   ```bash
   bash script/test.sh
   ```
   - 运行所有单元测试
   - 检查覆盖率
   - 输出 PASS/FAIL 状态

2. **构建应用**
   ```bash
   bash script/build.sh
   ```
   - 编译 Go 代码为可执行文件
   - 输出：`robot-order-system`

3. **运行应用**
   ```bash
   bash script/run.sh
   ```
   - 执行 CLI 应用
   - 输出结果到 `result.txt`
   - 包含时间戳 (HH:MM:SS)

## 输出格式

应用输出到 `result.txt` 的格式：

```
[HH:MM:SS] ✅ 订单完成: #1 [VIP] (耗时: 10.0s, 机器人: 1)
[HH:MM:SS] ✅ 订单完成: #2 [REGULAR] (耗时: 10.0s, 机器人: 2)
```

### 时间戳格式
- 格式：`[HH:MM:SS]`
- 用途：跟踪订单完成时间
- 精度：秒级 (Seconds)

## GitHub Actions

所有推送和拉取请求都会自动运行 GitHub Actions：

### 工作流程
1. **Test** - 运行单元测试 (Go 1.21, 1.22)
2. **Build** - 编译 CLI 应用
3. **Run** - 执行应用并验证输出
4. **Quality** - 代码质量检查

### 查看构建状态
访问仓库的 "Actions" 标签查看工作流程运行状态。

## 本地运行工作流程

```bash
# 1. 运行测试
bash script/test.sh

# 2. 编译应用
bash script/build.sh

# 3. 运行应用
bash script/run.sh

# 4. 检查输出
cat result.txt
```

## 提交拉取请求

1. 创建特性分支
   ```bash
   git checkout -b feature/your-feature
   ```

2. 进行更改并提交
   ```bash
   git add .
   git commit -m "Add your feature"
   ```

3. 推送到仓库
   ```bash
   git push origin feature/your-feature
   ```

4. 在 GitHub 上创建拉取请求

5. 等待所有 GitHub Actions 检查通过

## 常见问题

### 构建失败怎么办？
1. 检查 Go 版本：`go version`
2. 检查 Go 模块：`go mod verify`
3. 清理构建缓存：`go clean`

### 测试失败怎么办？
1. 检查错误消息
2. 在本地运行：`go test -v`
3. 查看 GitHub Actions 日志

### 脚本权限错误怎么办？
```bash
chmod +x script/*.sh
```

## 联系方式

有问题？请创建 GitHub Issue。

---

**最后更新**: 2026-03-17
