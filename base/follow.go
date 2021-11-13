package base

// GetFollow 根据推导式的左部，得到其对应的Follow集
func GetFollow(left Tag, depth int) []Tag {
	// 限制递归深度
	if depth > 5 {
		return []Tag{}
	}
	// Follow集合
	var follow []Tag
	// 遍历查找所有右部中包含left的产生式
	for _, productions := range GetProdMap() {
		for _, production := range productions {
			for i, tag := range production.Right {
				// 不能在最后一个，因为最后一个的话，之后就没有符号了
				if tag == left && i != len(production.Right)-1 {
					// 如果之后是终结符，则直接加入
					if production.Right[i+1].Type == TERM && !HasReTags(production.Right[i+1], follow) {
						follow = append(follow, production.Right[i+1])
					}
					// 如果之后是非终结符，则该非终结符的first集加入follow集
					if production.Right[i+1].Type == NONTERM {
						for _, t := range GetFirst(production.Right[i+1]) {
							if t.Value != "ε" && !HasReTags(t, follow) {
								follow = append(follow, t)
							}
						}
					}
				}

				// 如果在最后一个，或者之后的符号的first集包含空，则继承左部的Follow集
				if i == len(production.Right)-1 || HasReTags(Tag{
					Type:  TERM,
					Value: "ε",
				}, GetFirst(production.Right[i+1])) {
					for _, t := range GetFollow(production.Left, depth+1) {
						if !HasReTags(t, follow) {
							follow = append(follow, t)
						}
					}
				}
			}
		}
	}

	return follow
}
