# GitHub Actions 工作流程

本仓库配置了自动化的 CI/CD 工作流程，确保每个提交和拉取请求都通过测试和质量检查。

## 工作流程概览

### 文件位置
`.github/workflows/ci.yml`

### 触发条件

工作流程在以下情况自动运行：

1. **推送到主分支或开发分支**
   ```
   push:
     branches: [ main, develop ]
   ```

2. **拉取请求到主分支或开发分支**
   ```
   pull_request:
     branches: [ main, develop ]
   ```

## 工作流程详情

### 1. 测试任务（Test）

**运行环境**: Ubuntu Latest  
**Go 版本**: 1.21, 1.22（矩阵策略）

**步骤**:
- 检出代码
- 设置 Go 环境
- 运行测试：`bash script/test.sh`
- 上传测试结果

**输出**:
- 测试结果保存为工件

### 2. 构建任务（Build）

**运行环境**: Ubuntu Latest  
**Go 版本**: 1.21  
**依赖**: 测试任务通过

**步骤**:
- 检出代码
- 设置 Go 环境
- 构建应用：`bash script/build.sh`
- 验证可执行文件
- 上传可执行文件

**输出**:
- `robot-order-system` 可执行文件

### 3. 运行任务（Run）

**运行环境**: Ubuntu Latest  
**Go 版本**: 1.21  
**依赖**: 构建任务通过

**步骤**:
- 检出代码
- 设置 Go 环境
- 运行应用：`bash script/run.sh`
- 显示结果
- 验证时间戳格式
- 上传 result.txt

**输出**:
- `result.txt` 包含应用输出和时间戳

### 4. 质量检查任务（Quality）

**运行环境**: Ubuntu Latest  
**Go 版本**: 1.21

**步骤**:
- 检出代码
- 设置 Go 环境
- 运行 golangci-lint
- 检查代码格式
- 运行 go vet

**输出**:
- 代码质量报告

## 查看构建状态

### 在 GitHub 上

1. 访问仓库的 **Actions** 标签
2. 查看最新的工作流程运行
3. 点击查看详细信息

### 在 README 中

添加构建状态徽章到 README：

```markdown
[![CI/CD](https://github.com/[owner]/[repo]/actions/workflows/ci.yml/badge.svg)](https://github.com/[owner]/[repo]/actions)
```

## 故障排除

### 测试失败

如果测试任务失败：

1. 检查 "Test" 任务的日志
2. 查看错误消息
3. 在本地运行：`bash script/test.sh`
4. 修复问题并重新推送

### 构建失败

如果构建任务失败：

1. 检查 "Build" 任务的日志
2. 查看编译错误
3. 在本地运行：`bash script/build.sh`
4. 修复问题并重新推送

### 运行失败

如果运行任务失败：

1. 检查 "Run" 任务的日志
2. 查看执行时的错误
3. 下载 result.txt 工件查看输出
4. 修复问题并重新推送

## 工件

工作流程生成的工件可在以下位置下载：

1. **test-results-go-{version}**
   - 单元测试结果

2. **robot-order-system-executable**
   - 编译的 CLI 应用程序

3. **application-results**
   - result.txt 输出

### 下载工件

1. 打开工作流程运行页面
2. 向下滚动到 "Artifacts" 部分
3. 点击要下载的工件

## 时间戳验证

工作流程在 "Run" 任务中验证时间戳格式：

```bash
if grep -q '\[[0-9][0-9]:[0-9][0-9]:[0-9][0-9]\]' result.txt; then
  echo "✅ Timestamp format verified (HH:MM:SS)"
fi
```

格式：`[HH:MM:SS]`

## 环境变量

当前未设置特定的环境变量。如需添加，编辑 `ci.yml`：

```yaml
env:
  GO111MODULE: on
```

## 性能优化

### 缓存

工作流程可以缓存 Go 模块：

```yaml
- uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

## 安全考虑

1. 不在工作流程中提交敏感信息
2. 使用 GitHub Secrets 存储敏感数据
3. 定期审查工作流程配置
4. 使用 CODEOWNERS 文件管理审核者

## 扩展工作流程

要添加新的工作流程步骤：

1. 编辑 `.github/workflows/ci.yml`
2. 添加新的任务或步骤
3. 提交更改
4. 工作流程自动运行

## 常用命令

### 本地测试整个工作流程

使用 `act` 工具在本地运行：

```bash
# 安装 act
brew install act

# 运行整个工作流程
act

# 运行特定工作
act -j test
act -j build
act -j run
```

## 相关文档

- [GitHub Actions 官方文档](https://docs.github.com/en/actions)
- [Go 工作流程示例](https://github.com/actions/setup-go)
- [上传/下载工件](https://github.com/actions/upload-artifact)

---

**文档版本**: 1.0  
**最后更新**: 2026-03-17
