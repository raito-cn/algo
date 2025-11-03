package cmd

import (
	"algo/internal/db"
	"algo/internal/model"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "[ 查询题目列表 ] List problems by conditions",
	Run:   listProblems,
}

func InitListCmd() *cobra.Command {
	getDB := db.GetDB(false)
	tags := make([]model.Tag, 0)
	getDB.Model(&model.Tag{}).Find(&tags)
	str := make([]string, 0, len(tags))
	for _, tag := range tags {
		str = append(str, tag.Name)
	}

	listCmd.Flags().StringP("title", "t", "", "[ 题目标题 ] Problem title")
	listCmd.Flags().StringP("difficulty", "d", "", "[ 题目难度 ] Problem difficulty (easy|medium|hard)")
	listCmd.Flags().StringP("tags", "g", "", "[ 题目标签，英文逗号分割 ] Problem tags: "+strings.Join(str, "/")+
		"\n"+"Example: -g tag1,tag2")
	listCmd.Flags().BoolP("debug", "D", false, "Debug mode")
	listCmd.Flags().StringP("score", "s", "", "[ 题目评分 ] Problem score")
	listCmd.Flags().IntP("limit", "l", 100, "[ 查询条数 ] Limit number of problems")
	listCmd.Flags().IntP("offset", "o", 0, "[ 查询页码 ] Offset number of problems")
	return listCmd
}

func listProblems(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	titleKeyword := cmd.Flag("title").Value.String()
	difficulty := cmd.Flag("difficulty").Value.String()
	tagsStr := cmd.Flag("tags").Value.String()
	score := cmd.Flag("score").Value.String()
	limit, _ := cmd.Flags().GetInt("limit")
	offset, _ := cmd.Flags().GetInt("offset")

	if limit < 0 || offset < 0 || limit > 100 {
		fmt.Println("Invalid limit or offset")
		return
	}

	conn := db.GetDB(debug)
	query := conn.Preload("Tags").Limit(limit).Offset(offset)
	var problems []model.Problem
	if difficulty != "" {
		diff := model.Difficulty(difficulty)
		if !diff.Valid() {
			fmt.Println("Invalid difficulty, must be easy|medium|hard")
			return
		}
		query = query.Where("difficulty = ?", diff)
	}
	if titleKeyword != "" {
		query = query.Where("title LIKE ?", "%"+titleKeyword+"%")
	}
	if tagsStr != "" {
		tags := strings.Split(tagsStr, ",")
		query = query.Joins("JOIN problem_tags pt ON pt.problem_id = problems.id").
			Joins("JOIN tags t ON t.id = pt.tag_id").
			Where("t.name IN ?", tags)
	}
	if score != "" {
		query = query.Where("score = ?", score)
	}
	if err := query.Order("created_at DESC").Find(&problems).Error; err != nil {
		fmt.Println("Failed to list problems:", err)
		return
	}

	fmt.Println("Total problems:", len(problems))
	for i, p := range problems {
		fmt.Printf("[%d] [%s] %s | Difficulty: %s | Tags: %s\nSolution URL: %s\n",
			i,
			p.Slug,
			p.Title,
			p.Difficulty,
			getTagNames(p.Tags),
			p.SolutionURL,
		)
	}
}

func getTagNames(tags []*model.Tag) string {
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.Name
	}
	return strings.Join(names, ", ")
}
