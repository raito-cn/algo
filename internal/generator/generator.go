package generator

import (
	"sync"
)

var template string
var once sync.Once

func GetTemplate() string {
	once.Do(func() {
		template = getTemplate()
	})
	return template
}

type Problem struct {
	Title       string
	Difficulty  string
	Tags        []string
	SolutionURL string
	Score       *uint8
	CreatedAt   string
	UpdatedAt   string
	Slug        string
	Description string
	Solution    string
	Code        *Code
}

type Code struct {
	Language string
	Data     string
}

func getTemplate() string {
	template = `
# {{ problem.Title }}

| 属性 | 内容 |
| ---- | ---- |
| **难度** | {{ problem.Difficulty }} |
| **标签** | {% for t in problem.Tags %}{{ t }}{% if not forloop.Last %}, {% endif %}{% endfor %} |
| **链接** | [在线题目]({{ problem.SolutionURL }}) |
{% if problem.Score != none %}| **评分** | {{ problem.Score }} |{% endif %}
| **创建时间** | {{ problem.CreatedAt }} |
| **更新时间** | {{ problem.UpdatedAt }} |
| **Slug** | {{ problem.Slug }} |

---

## 📖 题目描述

{{ problem.Description|default:"暂无题目描述" }}

---

## 💡 解题思路

{{ problem.Solution|safe|default:"暂无解题思路" }}

---

## 🛠 代码实现

{% if problem.Code != None %}
~~~{{ problem.Code.Language }}
{{ problem.Code.Data|safe }}
~~~

{% else %}
暂无代码实现
{% endif %}
`
	return template
}
