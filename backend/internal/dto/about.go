package dto

// AboutContent 前台「关于」页可编辑内容
// intro 为 Markdown 文本（空行分段）；stack 为技术栈标签；footnote 为文末一行小注。
type AboutContent struct {
	Intro    string   `json:"intro" binding:"max=8000"`
	Stack    []string `json:"stack"`
	Footnote string   `json:"footnote" binding:"max=800"`
}
