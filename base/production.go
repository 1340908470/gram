package base

import "errors"

// Production 产生式
type Production struct {
	Left  Tag   // 当且仅当 Left.Type == NONTERM
	Right []Tag // 产生式的右部是一个标识符切片
}

// GetProductionsByTag 根据记号，返回该非终结符的所有产生式【返回的是已经消除了左递归的式子】
func GetProductionsByTag(productions []Production, left Tag) ([]Production, error) {
	if left.Type != NONTERM {
		return nil, errors.New("终结符无产生式")
	}

	var aimedProductions []Production

	// 解析生成式
	for _, production := range productions {
		if production.Left == left {
			aimedProductions = append(aimedProductions, production)
		}
	}

	var newProductions []Production
	var tmpTag Tag
	// 消除左递归: E -> Ef | g
	for _, production := range aimedProductions {
		if production.Left == production.Right[0] {
			// 如果换tmpTag前，先添加上一个tmpTag的空产生式
			if tmpTag.Value != "" && tmpTag != production.Left {
				pr := Production{
					Left: tmpTag,
					Right: []Tag{{
						Type:  TERM,
						Value: "ε",
					}},
				}
				newProductions = append(newProductions, pr)
				// 将新添加的 E' 加入tags中
				AddTag(tmpTag)
			}

			// E -> Ef
			tmpTag = production.Left // 记下有左递归的符号
			// 加 ‘
			newTag := Tag{
				Type:  NONTERM,
				Value: tmpTag.Value + "'",
			}

			// 删除第一个Tag，再在最后追加一个newTag
			production.Left = newTag
			production.Right = production.Right[1:]
			production.Right = append(production.Right, newTag)
		} else if production.Left == tmpTag {
			// E -> g
			production.Right = append(production.Right, production.Left)
		}

		newProductions = append(newProductions, production)
	}

	return newProductions, nil
}
