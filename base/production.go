package base

import "errors"

// Production 产生式
type Production struct {
	Left  Tag   // 当且仅当 Left.Type == NONTERM
	Right []Tag // 产生式的右部是一个标识符切片
}

// GetProductionsByTag 根据记号，返回该非终结符的所有产生式
func GetProductionsByTag(productions []Production, left Tag) ([]Production, error) {
	if left.Type != NONTERM {
		return nil, errors.New("终结符无产生式")
	}

	var aimedProductions []Production

	for _, production := range productions {
		if production.Left == left {
			aimedProductions = append(aimedProductions, production)
		}
	}

	return aimedProductions, nil
}
