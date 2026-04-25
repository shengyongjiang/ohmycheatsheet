# omcs — Oh My Cheat Sheet

一个带记忆追踪功能的终端速查表工具。数据源自 [cheat.sh](https://cheat.sh)。

## 安装

```bash
go install github.com/shengyongjiang/ohmycheatsheet/cmd/omcs@latest
```

确保 `$GOPATH/bin` 在你的 `PATH` 中：

```bash
export PATH="$HOME/go/bin:$PATH"  # 添加到 ~/.bashrc 或 ~/.zshrc
```

或从源码构建：

```bash
git clone https://github.com/shengyongjiang/ohmycheatsheet.git
cd ohmycheatsheet
go build -o omcs ./cmd/omcs/
```

## 使用方法

### 查看速查表

```bash
omcs git          # 显示 10 条未记住的随机条目
omcs git --all    # 显示所有条目
omcs git --random # 重新随机排列条目
```

默认视图显示 10 条未记住的条目，以每日固定的随机顺序排列。使用 `--random` 强制重新洗牌。

### 交互模式

```bash
omcs git -i
```

打开 TUI 界面显示所有条目。条目顺序与上次非交互输出一致，所以可以先用 `omcs git` 预览，再用 `omcs git -i` 深入学习。

#### 快捷键

| 按键 | 功能 |
|---|---|
| `j` / `k` / 方向键 | 上下导航 |
| `Left` / `Right` | 循环切换记忆状态 |
| `x` / `X` | 标记为已记住 |
| `Enter` | 标记为待复习 |
| `a` | 切换显示全部/过滤模式 |
| `r` | 重置所有状态 |
| `q` / `Esc` | 退出（自动保存） |
| `?` | 帮助 |

### 记忆状态

- **未记住** (`o`，白色) — 默认状态，始终显示
- **已记住** (`x`，灰色) — 默认隐藏，表示你已掌握该命令
- **待复习** (`+`，红色) — 高亮显示，需要再次复习

### 其他命令

```bash
omcs review          # 闪卡式复习"待复习"条目
omcs review git      # 仅复习 git 相关条目
omcs stats           # 显示所有命令的记忆进度
omcs stats git       # 显示特定命令的进度
omcs list            # 列出所有已追踪的命令
omcs reset git       # 重置 git 的记忆状态
omcs reset --all     # 重置所有内容
omcs completion zsh  # 生成 shell 补全脚本（支持 bash/zsh/fish）
```

## 工作原理

1. 首次使用时从 `cheat.sh` 获取速查表内容，本地缓存 7 天
2. 条目按每日确定性种子随机排列（同一天内保持一致）
3. 记忆状态以 JSON 格式持久化存储在本地
4. 当条目被标记为已记住并隐藏后，相关条目（如 `git` 的 `git-log`、`git-stash`）会自动补充到列表中

## 数据存储

| 数据 | 位置（macOS） |
|---|---|
| 速查表缓存 | `~/Library/Caches/omcs/cheatsh/` |
| 记忆状态 | `~/Library/Application Support/omcs/state.json` |
| 配置文件（可选） | `~/Library/Application Support/omcs/config.json` |
| 随机种子 | `~/Library/Caches/omcs/cheatsh/seeds/` |

Linux 上遵循 XDG 默认路径（`~/.cache/omcs/` 和 `~/.config/omcs/`）。

## 项目结构

```
cmd/              CLI 命令（show、review、stats、list、reset、completion）
internal/
  config/         配置加载
  model/          核心类型（Entry、Page、MemoryState、EntryState）
  parser/         解析 cheat.sh 纯文本格式为结构化条目
  render/         非交互式文本渲染器
  resolver/       通过 cheat.sh 数据源将命令名解析为页面
  shuffle/        基于种子的确定性条目随机排列
  source/         cheat.sh HTTP 客户端及本地缓存
  store/          JSON 格式的状态持久化
  tui/            交互式 TUI（基于 Bubble Tea）
```
