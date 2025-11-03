package model

import (
	"algo/internal/util"
	"algo/pkg/config"
	"fmt"
	"github.com/mozillazg/go-pinyin"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// Difficulty 题目难度
type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

func (d *Difficulty) String() string {
	return string(*d)
}

func (d *Difficulty) Valid() bool {
	toLower := strings.ToLower(d.String())
	switch toLower {
	case "easy", "medium", "hard":
		*d = Difficulty(toLower)
		return true
	default:
		return false
	}
}

// Problem 题目
type Problem struct {
	ID          int64      `gorm:"primaryKey;autoIncrement:false;comment:主键"`       // 主键
	Title       string     `gorm:"not null;size:255;comment:题目名"`                   // 题目名
	Slug        string     `gorm:"uniqueIndex;not null;size:50;comment:题目短id、文件名用"` // 题目短id、文件名用
	Difficulty  Difficulty `gorm:"not null;comment:题目难度"`                           // 题目难度
	SolutionURL string     `gorm:"not null;text;comment:在线题目链接"`                    // 在线题目链接
	Note        string     `gorm:"text;comment:题目笔记"`                               // 题目笔记
	CodePath    string     `gorm:"text;comment:本地代码文件路径"`                           // 本地代码文件路径
	Score       *uint8     `gorm:"comment:题目评分"`                                    // 题目分数
	CreatedAt   time.Time  `gorm:"autoCreateTime;comment:创建时间"`                     // 创建时间
	UpdatedAt   time.Time  `gorm:"autoUpdateTime;comment:更新时间"`                     // 更新时间
	Tags        []*Tag     `gorm:"many2many:problem_tags;"`                         // 标签逻辑
	Description string     `gorm:"text;comment:题目描述"`                               // 题目描述
}

func (p *Problem) SetSlug() {
	a := pinyin.NewArgs()
	var slugParts []string

	for _, r := range p.Title {
		switch {
		case unicode.Is(unicode.Han, r): // 中文字符
			py := pinyin.Pinyin(string(r), a)
			if len(py) > 0 && len(py[0]) > 0 {
				slugParts = append(slugParts, strings.ToLower(py[0][0]))
			}
		case unicode.IsLetter(r) || unicode.IsNumber(r): // 英文/数字
			slugParts = append(slugParts, strings.ToLower(string(r)))
		default: // 空格/标点
			slugParts = append(slugParts, "-")
		}
	}

	slug := strings.Join(slugParts, "")
	// 合并连续的 -
	re := regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	p.Slug = fmt.Sprintf("%04d_%s", p.ID, slug)
}

func (p *Problem) SetCodePath() {
	data, err := os.ReadFile(p.CodePath)
	if err != nil {
		util.GetLog().Error("read code file error")
		panic(err)
	}
	// 2. 获取目标目录
	codeDir := config.GetConfig().Dir.CodeDir
	// 确保目录存在
	if err = os.MkdirAll(codeDir, 0755); err != nil {
		util.GetLog().Error("create code dir error", zap.Error(err))
		panic(err)
	}

	// 3. 构建新文件名和路径
	ext := filepath.Ext(p.CodePath) // 保留原始文件扩展名
	newFileName := fmt.Sprintf("%s_code%s", p.Slug, ext)
	dstPath := filepath.Join(codeDir, newFileName)
	if err = os.WriteFile(dstPath, data, 0644); err != nil {
		util.GetLog().Error("write code file error", zap.Error(err))
		panic(err)
	}

	// 5. 更新 Problem 的 CodePath 为新路径
	p.CodePath = dstPath
}

// Tag 题目标签
type Tag struct {
	ID       int64      `gorm:"primaryKey;autoIncrement:false;comment:主键"` // 主键
	Name     string     `gorm:"uniqueIndex;not null;size:50;comment:标签名"`  // 标签名
	Problems []*Problem `gorm:"many2many:problem_tags;"`                   // 题目
}
