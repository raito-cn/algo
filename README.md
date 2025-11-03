# ALGO - 高效、结构化地纪录算法学习过程的 CLI 系统

ALGO 是一个轻量级命令行系统，帮助你高效管理、整理与导出算法题笔记。
适合独立学习者、ACM选手、LeetCode爱好者使用。
---

## 功能概览
通过命令行管理算法题数据，例如：
```bash
$ algo add -t "Two Sum" -d easy -g array -S "https://leetcode.com/problems/two-sum" -n "用 map 存索引即可，O(n)" -c golang.go -s 5 --debug true
$ algo list --tag array
$ algo stat
```
系统自动生成结构化的 Markdown 知识库：
```bash
/notes
  /easy/001_two_sum.md
  /medium/098_validate_bst.md
  /hard/212_word_search.md

```
## 核心功能模块设计
### 1. 题目管理

- **新增题目**：标题、难度、语言、标签、笔记、代码路径等
- **修改/删除题目**：支持按ID或标题操作
- **查询功能**：按难度、标签、关键字筛选

### 2. 笔记与导出

- 自动生成 Markdown 文件，每题一个
- 生成汇总索引 `README.md`
- 支持统计：难度分布、标签分布、总题量等

### 3. 可选增强功能（计划中）

- 从 LeetCode / Codeforces API 自动拉题信息
- 自动检测本地代码是否存在、是否可编译
- 与 GitHub/Gist 同步上传，打造云端知识库

## CLI命令设计
```bash
algo
├── add           # 新增题目
├── edit          # 修改题目
├── remove        # 删除题目
├── list          # 按条件列出题目
├── stat          # 统计信息
├── gen           # 生成 Markdown 笔记
└── sync          # (可选) 同步至 GitHub
```

## 技术栈设计

| 模块         | 技术选型                                             |
|------------|--------------------------------------------------|
| CLI框架      | [spf13/cobra](https://github.com/spf13/cobra)    |
| 数据存储       | SQLite + [sqlx](https://github.com/jmoiron/sqlx) |
| Markdown生成 | Go `text/template`                               |
| 未来扩展       | LeetCode API / GitHub API                        |

## 项目结构示例

```bash
algo-cli/
├── cmd/                 # 命令定义 (add, list, stat 等)
├── internal/
│   ├── db/              # SQLite 数据库逻辑
│   ├── model/           # Problem 数据结构
│   └── generator/       # Markdown 导出逻辑
├── notes/               # 自动生成的 Markdown 笔记
├── go.mod
└── main.go
```

## 示例输出
示例生成的 Markdown 文件：

```markdown
# Two Sum

- 难度: easy
- 标签: array, hashmap
- 链接: https://leetcode.com/problems/two-sum

---

## 解题思路
用 map 存索引即可，O(n)

---

*创建时间: 2025-10-29*
```

## 未来计划
- 支持 algo sync 上传笔记到 GitHub
- algo open 一键打开题目网页
- Web UI 可视化版本
- 导出为 Notion/Obsidian 格式

## 作者
一个热爱算法与工具开发的独立开发者。<br/>
目标是：**让算法笔记更结构化、更轻量、更私有。**