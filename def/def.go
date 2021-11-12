package def

import (
	"encoding/json"
	"errors"
	"gram/base"
	"io/ioutil"
	"strings"
)

type Tag struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Def struct {
	Tags        []Tag    `json:"tags"`
	Productions []string `json:"productions"`
}

var def Def

func InitDef() error {
	file, err := ioutil.ReadFile("def/def.json")
	if err != nil {
		return errors.New("文件 def.json 读取失败")
	}

	err = json.Unmarshal(file, &def)
	if err != nil {
		return errors.New("def.json 解析失败")
	}

	return err
}

// GetTags 返回的是程序解析时真正使用的 base.Tag
func GetTags() []base.Tag {
	var tags []base.Tag

	for _, val := range def.Tags {
		tag := base.Tag{
			Type: func(t Tag) int {
				if t.Type == "终结符" {
					return base.TERM
				}
				if t.Type == "非终结符" {
					return base.NONTERM
				}
				return -1
			}(val),
			Value: val.Value,
		}
		tags = append(tags, tag)
	}

	return tags
}

// GetProductions 返回的是程序解析时使用的 base.Production
// 形如 "E → E+T | E–T | T" 的产生式会被拆分为三个 production
func GetProductions() []base.Production {
	var productions []base.Production

	for _, val := range def.Productions {
		// 首先去空格，然后切分
		strs := strings.Split(strings.Replace(val, " ", "", -1), "→")

		// strs被分为两个部分，0为左部元素，1为右部元素

		// 根据 GetTags 找到左部对应的 Tag
		var left base.Tag
		for _, tag := range GetTags() {
			if tag.Value == strs[0] {
				left = tag
				break
			}
		}

		// 右部根据 "｜" 再次划分，然后再遍历加入 productions
		for _, right := range strings.Split(strs[1], "|") {
			var tags []base.Tag
			index := 0
			for index < len(right) {
				for _, tag := range GetTags() {
					if index+len(tag.Value) <= len(right) && tag.Value == right[index:index+len(tag.Value)] {
						tags = append(tags, tag)
						index += len(tag.Value)
						break
					}
				}
			}

			production := base.Production{
				Left:  left,
				Right: tags,
			}

			productions = append(productions, production)
		}

	}

	return productions
}
