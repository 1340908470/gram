package base

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

type TmpTag struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Def struct {
	Tags        []TmpTag `json:"tags"`
	Productions []string `json:"productions"`

	Input  string `json:"input"`
	Method string `json:"method"`
}

var def Def
var tags []Tag
var productions []Production
var prodMap = make(map[Tag][]Production)

// InitDef 初始化操作，需在程序的入口处执行，以将json文件的内容读到内存中去
func InitDef() error {
	file, err := ioutil.ReadFile("base/def.json")
	if err != nil {
		return errors.New("文件 def.json 读取失败")
	}

	err = json.Unmarshal(file, &def)
	if err != nil {
		return errors.New("def.json 解析失败")
	}

	// 从文件解析初始Tags
	for _, val := range def.Tags {
		tag := Tag{
			Type: func(t TmpTag) int {
				if t.Type == "终结符" {
					return TERM
				}
				if t.Type == "非终结符" {
					return NONTERM
				}
				return -1
			}(val),
			Value: val.Value,
		}
		tags = append(tags, tag)
	}

	// 从文件解析初始Productions
	for _, val := range def.Productions {
		// 首先去空格，然后切分
		strs := strings.Split(strings.Replace(val, " ", "", -1), "→")

		// strs被分为两个部分，0为左部元素，1为右部元素

		// 根据 GetTags 找到左部对应的 Tag
		var left Tag
		for _, tag := range GetTags() {
			if tag.Value == strs[0] {
				left = tag
				break
			}
		}

		// 右部根据 "｜" 再次划分，然后再遍历加入 productions
		for _, right := range strings.Split(strs[1], "|") {
			var tags []Tag
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

			production := Production{
				Left:  left,
				Right: tags,
			}

			productions = append(productions, production)
		}

	}

	// 填充prodMap
	for _, production := range productions {
		prodMap[production.Left] = append(prodMap[production.Left], production)
	}

	// 在input后面追加 $
	def.Input += "$"

	return err
}

// RemoveLeftRecursion 消除左递归
func RemoveLeftRecursion() {
	// 先删除原有prodMap
	prodMap = make(map[Tag][]Production)

	// 遍历消除左递归
	for _, tag := range GetTags() {
		if tag.Type == NONTERM {
			pros, err := GetProductionsByTag(GetProductions(), tag)
			if err != nil {
				panic(err)
			}
			for _, pro := range pros {
				prodMap[pro.Left] = append(prodMap[pro.Left], pro)
			}
		}
	}

}

// GenerateExtension 构造拓广文法
func GenerateExtension() {
	// 在productions的最前面添加 E' -> E （因为程序默认第一个产生式的左部是起始符）
	oriLeft := GetProductions()[0].Left
	production := Production{
		Left: Tag{
			Type:  NONTERM,
			Value: oriLeft.Value + "'",
		},
		Right: []Tag{oriLeft},
	}
	productions = append([]Production{production}, productions...)

	// 更新prodMap

	// 先删除原有prodMap
	prodMap = make(map[Tag][]Production)

	// 填充prodMap
	for _, p := range productions {
		prodMap[p.Left] = append(prodMap[p.Left], p)
	}
}

func AddTag(tag Tag) {
	tags = append(tags, tag)
}

// GetTags 返回的是程序解析时真正使用的 base.Tag
func GetTags() []Tag {
	return tags
}

// GetProdMap 获得以生成式左部为键的生成式mao
func GetProdMap() map[Tag][]Production {
	return prodMap
}

// GetProductions 返回的是程序解析时使用的 base.Production
// 形如 "E → E+T | E–T | T" 的产生式会被拆分为三个 production
func GetProductions() []Production {
	return productions
}
