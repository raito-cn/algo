package cmd

import (
	"algo/internal/db"
	"algo/internal/model"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "[ 添加一个新题目 ] Add a new problem",
	Run:   addProblem,
}

func InitAddCmd() *cobra.Command {
	addCmd.Flags().StringP("title", "t", "", "[ 题目标题 ] Problem title")
	addCmd.Flags().StringP("difficulty", "d", "", "[ 题目难度 ] Problem difficulty (easy|medium|hard)")
	addCmd.Flags().StringP("tags", "g", "", "[ 题目标签，英文逗号分割 ] Problem tags, comma separation")
	addCmd.Flags().StringP("solution", "S", "", "[ 题目在线地址 ] Problem solution URL")
	addCmd.Flags().StringP("note", "n", "", "[ 题目笔记 ] Problem note")
	addCmd.Flags().StringP("codePath", "c", "", "[ 题目代码本地地址 ] Problem code path")
	addCmd.Flags().BoolP("debug", "D", false, "Debug mode")
	addCmd.Flags().StringP("score", "s", "", "[ 题目评分 ] Problem score")
	return addCmd
}

func addProblem(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	title, difficulty, tags, solution, note, codePath, score := getAddCmdParams(cmd, debug)

	diff := model.Difficulty(difficulty)
	valid := diff.Valid()
	if !valid {
		panic("invalid difficulty")
	}

	conn := db.GetDB(debug)
	err := conn.Transaction(func(tx *gorm.DB) error {
		tagsArr, err := CheckTags(tx, tags)
		if err != nil {
			return err
		}
		var count int64
		tx.Model(&model.Problem{}).Select("max(id)").Scan(&count)
		count++
		problem := &model.Problem{
			ID:          count,
			Title:       title,
			Difficulty:  diff,
			Tags:        tagsArr,
			SolutionURL: solution,
			Note:        note,
			CodePath:    codePath,
		}
		if scoreInt, err := strconv.Atoi(score); err == nil {
			if scoreInt >= 0 && scoreInt <= 255 { // uint8 范围检查
				v := uint8(scoreInt)
				problem.Score = &v
			}
		}
		problem.SetSlug()
		problem.SetCodePath()
		if err := tx.Create(&problem).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println("Failed to add problem:", err)
	} else {
		fmt.Println("Problem added successfully")
	}
}

func CheckTags(tx *gorm.DB, tags string) ([]*model.Tag, error) {
	var count int64
	split := strings.Split(tags, ",")
	tagsArr := make([]*model.Tag, 0, len(split))
	for _, tagName := range split {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}
		tagName = strings.ToLower(tagName)

		tx.Model(&model.Tag{}).Select("max(id)").Scan(&count)
		var tag model.Tag
		if err := tx.Where("name = ?", tagName).First(&tag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				count++
				tag = model.Tag{
					ID:   count,
					Name: tagName,
				}
				if err = tx.Create(&tag).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		tagsArr = append(tagsArr, &tag)
	}
	return tagsArr, nil
}

func getAddCmdParams(cmd *cobra.Command, debug bool) (string, string, string, string, string, string, string) {
	title := getCmdParam(cmd, "title", "Title: ", debug)
	difficulty := getCmdParam(cmd, "difficulty", "Difficulty (easy|medium|hard):", debug)
	tags := getCmdParam(cmd, "tags", "Tags: ", debug)
	solution := getCmdParam(cmd, "solution", "Solution URL: ", debug)
	note := getCmdParam(cmd, "note", "Note: ", debug)
	codePath := getCmdParam(cmd, "codePath", "Code Path: ", debug)
	score := getCmdParam(cmd, "score", "Score: ", debug)
	return title, difficulty, tags, solution, note, codePath, score
}
