# 部署指南

本文档说明如何将项目推送到 GitHub 并通过 GitHub Actions 进行 CI/CD 验证。

## 前置条件

1. GitHub 账户
2. Git 已安装
3. GitHub 仓库已创建
4. Go 1.21+ 已安装（本地测试）

## 步骤 1：创建 GitHub 仓库

### 在 GitHub 上创建仓库

1. 访问 [GitHub](https://github.com/new)
2. 填写仓库名称：`robot-order-system`
3. 选择 "Public" 或 "Private"
4. 点击 "Create repository"

## 步骤 2：配置 Git 远程

```bash
# 进入项目目录
cd /Users/wangbingfan/IdeaProjects/go/robotOrderOpts

# 添加远程仓库
git remote add origin https://github.com/YOUR_USERNAME/robot-order-system.git

# 验证远程配置
git remote -v
```

## 步骤 3：推送代码到 GitHub

### 推送主分支

```bash
# 推送主分支
git push -u origin main

# 推送开发分支
git push -u origin develop

# 推送特性分支
git push -u origin feature/interactive-cli
```

### 验证推送

1. 访问 GitHub 仓库
2. 查看 "Code" 标签
3. 确认所有文件都已上传

## 步骤 4：监看 GitHub Actions

### 查看工作流程运行

1. 进入仓库
2. 点击 "Actions" 标签
3. 查看最近的工作流程运行

### 工作流程状态

工作流程运行时会显示以下任务：

- ✅ **Test** - 单元测试（Go 1.21, 1.22）
- ✅ **Build** - 编译 CLI 应用
- ✅ **Run** - 执行应用
- ✅ **Quality** - 代码质量检查

### 查看失败原因

如果工作流程失败：

1. 点击失败的任务
2. 展开 "Run" 步骤查看错误
3. 根据错误消息进行修复
4. 重新推送代码

## 步骤 5：创建拉取请求

### 创建 PR（从 feature 分支到 main）

```bash
# 确保所有本地提交都已推送
git push origin feature/interactive-cli
```

### 在 GitHub 上创建 PR

1. 进入仓库
2. 点击 "Pull requests" 标签
3. 点击 "New pull request"
4. 设置：
   - 基础分支：`main`
   - 比较分支：`feature/interactive-cli`
5. 点击 "Create pull request"

### PR 审查

创建 PR 后：

1. GitHub Actions 会自动运行所有检查
2. 检查通过后会显示 ✅ 绿色标记
3. 检查失败会显示 ❌ 红色标记

## 步骤 6：合并 PR

### 手动合并（需要审查）

1. 所有 GitHub Actions 检查通过 ✅
2. 获得必要的审查批准
3. 点击 "Merge pull request"
4. 确认合并

### 自动合并（可选）

启用自动合并：

1. 进入 PR 页面
2. 点击 "Enable auto-merge"
3. 选择合并方法（Squash, Merge, Rebase）
4. 确认

## 步骤 7：发布版本

### 创建版本标签

```bash
# 创建标签
git tag -a v1.0.0 -m "Release version 1.0.0"

# 推送标签
git push origin v1.0.0
```

### 在 GitHub 上创建 Release

1. 进入仓库的 "Releases" 页面
2. 点击 "Draft a new release"
3. 选择版本标签（v1.0.0）
4. 填写发布说明
5. 点击 "Publish release"

## 故障排除

### 推送失败

```bash
# 错误：remote: Repository not found
# 解决：检查 URL 和权限

# 错误：Permission denied
# 解决：配置 SSH 密钥或使用 HTTPS token
```

### GitHub Actions 失败

### 常见错误

#### 1. Go 版本不兼容

```
error: 'go' version or module issue
```

**解决**:
- 检查 `go.mod` 中的 Go 版本
- 确保 `.github/workflows/ci.yml` 中的版本一致

#### 2. 脚本权限问题

```
bash: script/test.sh: Permission denied
```

**解决**:
```bash
chmod +x script/*.sh
git add script/*.sh
git commit -m "Fix script permissions"
git push origin main
```

#### 3. 时间戳验证失败

```
Timestamp format not found
```

**解决**:
- 检查输出格式是否为 `[HH:MM:SS]`
- 查看 `result.txt` 文件
- 验证 `output.go` 中的时间戳函数

### PR 检查失败

1. 检查 GitHub Actions 日志
2. 修复代码问题
3. 推送修复
4. PR 会自动重新运行检查

## 最佳实践

### 分支策略

```
main (生产)
  ↑
develop (开发)
  ↑
feature/* (特性分支)
```

### 提交信息

```
<类型>: <描述>

<详细说明>

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>
```

类型:
- feat: 新功能
- fix: 错误修复
- docs: 文档更新
- test: 测试添加/修改
- chore: 维护

### 代码审查

1. 至少一个审查批准
2. 所有 GitHub Actions 检查通过
3. 无冲突合并
4. 清晰的提交历史

## 环境变量

### 本地设置（可选）

```bash
# 创建 .env 文件
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org
```

### GitHub Secrets（可选）

如需添加敏感信息：

1. 进入仓库设置
2. "Secrets and variables" > "Actions"
3. 点击 "New repository secret"
4. 在工作流程中使用：`${{ secrets.SECRET_NAME }}`

## 自动化清单

- ✅ 脚本文件可执行权限
- ✅ .gitignore 配置正确
- ✅ GitHub Actions 工作流程已定义
- ✅ PR 模板已创建
- ✅ 时间戳功能已实现
- ✅ 输出到 result.txt

## 监控和维护

### 定期检查

1. 查看 Actions 运行历史
2. 监控测试覆盖率
3. 检查构建时间趋势
4. 审查代码质量报告

### 更新依赖

```bash
# 检查更新
go list -u -m all

# 更新所有依赖
go get -u ./...

# 更新特定包
go get -u github.com/package/name

# 提交更新
git add go.mod go.sum
git commit -m "chore: update dependencies"
git push origin main
```

## 下一步

- 🔄 设置分支保护规则
- 📊 配置代码覆盖率报告
- 🚀 设置自动发布
- 📧 配置通知

---

**文档版本**: 1.0  
**最后更新**: 2026-03-17
