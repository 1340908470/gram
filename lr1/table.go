package lr1

import (
	"errors"
	"fmt"
	"github.com/liushuochen/gotable"
	"gram/base"
	"reflect"
)

const (
	SHIFT = iota
	REDUCE
	ACC
	GOTO
)

// AG Action or Goto
type AG struct {
	Type  int
	Value int
}

// LRTable LR分析表，最外层为切片，索引表示Group的Index，内层为map，键为tag，值为AG
type LRTable []map[base.Tag]AG

// ToString 转为字符串
func (a AG) ToString() string {
	if a.Type == SHIFT {
		return fmt.Sprintf("S%v", a.Value)
	}
	if a.Type == REDUCE {
		return fmt.Sprintf("R%v", a.Value)
	}
	if a.Type == ACC {
		return fmt.Sprintf("ACC")
	}
	if a.Type == GOTO {
		return fmt.Sprintf("%v", a.Value)
	}
	return ""
}

// GenerateLRTable 根据Group、GroupRelation生成表格
func GenerateLRTable() LRTable {
	var lRTable LRTable
	for range Groups {
		lRTable = append(lRTable, make(map[base.Tag]AG))
	}

	// 根据GroupRelation写 SHIFT 和 GOTO
	for index, m := range GroupRelation {
		for tag, i := range m {
			if i.Tag.Type == base.NONTERM {
				lRTable[index][tag] = AG{
					Type:  GOTO,
					Value: i.State,
				}
			}
			if i.Tag.Type == base.TERM {
				lRTable[index][tag] = AG{
					Type:  SHIFT,
					Value: i.State,
				}
			}

		}
	}

	// 遍历Group写ACC 和 R
	for i, group := range Groups {
		for _, d := range group.PDs {
			if reflect.DeepEqual(d.Left, base.GetProductions()[0].Left) &&
				d.DotIndex == len(d.Right) {
				lRTable[i][base.Tag{
					Type:  base.TERM,
					Value: "$",
				}] = AG{
					Type:  ACC,
					Value: 0,
				}
				break
			} else if d.DotIndex == len(d.Right) {
				for _, tag := range base.GetFollow(d.Left) {
					lRTable[i][tag] = AG{
						Type:  REDUCE,
						Value: d.Production.GetIndex(),
					}
				}
			}
		}
	}

	return lRTable
}

func PrintLRTable(table LRTable) error {
	colNames := []string{" "}
	for _, tag := range base.GetTags() {
		if !base.IsEmptyTag(tag) {
			colNames = append(colNames, tag.Value)
		} else if base.IsEmptyTag(tag) {
			colNames = append(colNames, "$")
		}
	}
	gt, err := gotable.Create(colNames...)
	if err != nil {
		return errors.New("创建表格失败")
	}

	for i, m := range table {
		row := make(map[string]string)
		row[" "] = fmt.Sprintf("%v", i)

		for t, production := range m {
			row[t.Value] = production.ToString()
		}

		err = gt.AddRow(row)
		if err != nil {
			return err
		}
	}

	fmt.Printf("-------------------------------------- 打印LR分析表 -------------------------------------\n")
	gt.PrintTable()
	fmt.Printf("-----------------------------------------------------------------------------------------\n\n")

	return err
}
