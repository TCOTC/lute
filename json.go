// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

import (
	"encoding/json"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

// ParseJSON 用于解析 jsonStr 生成 Markdown 抽象语法树。
func (lute *Lute) ParseJSON(jsonStr string) (ret *parse.Tree) {
	var children []map[string]interface{}
	err := json.Unmarshal(util.StrToBytes(jsonStr), &children)
	if nil != err {
		return
	}

	ret = &parse.Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument}, Context: &parse.Context{Option: lute.Options}}
	ret.Context.Tip = ret.Root
	for _, child := range children {
		lute.genASTByJSON(child, ret)
	}
	return
}

func (lute *Lute) genASTByJSON(jsonNode interface{}, tree *parse.Tree) {
	n := jsonNode.(map[string]interface{})
	typ := n["type"].(string)
	node := &ast.Node{Type: ast.Str2NodeType(typ)}
	switch node.Type {
	case ast.NodeDocument:
	case ast.NodeParagraph:
	case ast.NodeText:
		node.Tokens = util.StrToBytes(n["val"].(string))
	}
	tree.Context.Tip.AppendChild(node)
	tree.Context.Tip = node
	defer tree.Context.ParentTip()

	if nil == n["children"] {
		return
	}
	children := n["children"].([]interface{})
	for _, child := range children {
		lute.genASTByJSON(child, tree)
	}
}

// RenderJSON 用于渲染 JSON 格式数据。
func (lute *Lute) RenderJSON(markdown string) (retJSON string) {
	tree := parse.Parse("", []byte(markdown), lute.Options)
	renderer := render.NewJSONRenderer(tree)
	output := renderer.Render()
	retJSON = string(output)
	return
}
