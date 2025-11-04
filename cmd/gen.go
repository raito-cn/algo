package cmd

import (
	"algo/internal/db"
	"algo/internal/generator"
	"algo/internal/model"
	"algo/pkg/config"
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var genCmd = &cobra.Command{
	Use:   "gen [slug] [type]",
	Args:  cobra.ExactArgs(2),
	Short: "[ 生成题目文档 ] Generate a problem markdown file",
	Run:   getMarkdown,
}

func InitGenCmd() *cobra.Command {
	genCmd.Flags().BoolP("debug", "D", false, "Debug mode")
	genCmd.Long = `Generator a problem by its slug.
Example:
  algo gen 0001_two-sum pro --debug
  algo gen 1024 codeforces --debug 
`
	return genCmd
}

func getMarkdown(cmd *cobra.Command, args []string) {
	t := "pro"
	if len(args) > 1 {
		t = args[1]
	}
	conn := db.GetDB(false)
	if t == "pro" {
		var problem model.Problem
		if err := conn.Where("slug = ?", args[0]).Preload("Tags").First(&problem).Error; err != nil {
			fmt.Println("Failed to generate problem:", err)
			return
		}
		markdown, err := renderProblemMarkdown(&problem)
		if err != nil {
			fmt.Println("Failed to render problem markdown:", err)
			return
		}
		dir := config.GetConfig().Dir.MarkdownDir + "/" + problem.Difficulty.String()
		if err = os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("Failed to create markdown dir:", err)
			return
		}
		fileName := fmt.Sprintf("%s.md", problem.Slug)
		filePath := filepath.Join(dir, fileName)
		if err = os.WriteFile(filePath, []byte(markdown), 0644); err != nil {
			fmt.Println("Failed to write markdown file:", err)
			return
		}
	} else {
		contest := &model.Contest{}
		if err := conn.Where("title = ? and `type` = ?", args[0], t).Preload("Problems.Tags").First(contest).Error; err != nil {
			fmt.Println("Failed to generate contest:", err)
			return
		}
		var sb strings.Builder
		for _, problem := range contest.Problems {
			markdown, err := renderProblemMarkdown(problem)
			if err != nil {
				fmt.Println("Failed to render problem markdown:", err)
				return
			}
			sb.WriteString(markdown + "\n---\n")
		}
		dir := config.GetConfig().Dir.MarkdownDir + "/" + contest.Type.String()
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("Failed to create markdown dir:", err)
			return
		}
		fileName := fmt.Sprintf("%s.md", contest.Title)
		filePath := filepath.Join(dir, fileName)
		if err := os.WriteFile(filePath, []byte(sb.String()), 0644); err != nil {
			fmt.Println("Failed to write markdown file:", err)
			return
		}
	}
	fmt.Println("Markdown file generated successfully")
}

func renderProblemMarkdown(p *model.Problem) (string, error) {
	tags := make([]string, 0)
	for _, t := range p.Tags {
		tags = append(tags, t.Name)
	}
	created := p.CreatedAt.Format("2006-01-02 15:04:05")
	updated := p.UpdatedAt.Format("2006-01-02 15:04:05")

	language := strings.TrimPrefix(filepath.Ext(p.CodePath), ".")

	data, err := os.ReadFile(p.CodePath)
	if err != nil {
		return "", err
	}
	code := string(data)

	problem := &generator.Problem{
		Title:       p.Title,
		Difficulty:  p.Difficulty.String(),
		Tags:        tags,
		SolutionURL: p.SolutionURL,
		Score:       p.Score,
		CreatedAt:   created,
		UpdatedAt:   updated,
		Slug:        p.Slug,
		Description: p.Description,
		Solution:    p.Note,
		Code: &generator.Code{
			Language: language,
			Data:     code,
		},
	}
	tpl, err := pongo2.FromString(generator.GetTemplate())
	if err != nil {
		return "", err
	}
	out, err := tpl.Execute(pongo2.Context{"problem": problem})
	if err != nil {
		return "", err
	}
	return out, err
}
