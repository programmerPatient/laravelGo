/*
 * @Description:字符串操作
 * @Author: mali
 * @Date: 2023-07-19 11:41:43
 * @LastEditTime: 2023-07-31 15:13:30
 * @LastEditors: VSCode
 * @Reference:
 */
package helper

import (
	"strings"
)

//字符串分段
func SplitTextIntoParagraphs(text string, maxLength int) []string {
	paragraphs := []string{}
	words := strings.Fields(text) // 将文本拆分为单词

	currentParagraph := ""
	currentLength := 0

	for _, word := range words {
		// 如果当前段落长度加上当前单词长度超过了最大长度，则开始一个新的段落
		if currentLength+len(word)+1 > maxLength {
			paragraphs = append(paragraphs, currentParagraph)
			currentParagraph = ""
			currentLength = 0
		}

		// 将当前单词添加到当前段落中
		currentParagraph += word + " "
		currentLength += len(word) + 1
	}

	// 将最后一个段落添加到结果中
	paragraphs = append(paragraphs, currentParagraph)

	return paragraphs
}
