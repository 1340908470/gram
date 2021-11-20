package lr1

import (
	"fmt"
	"gram/base"
	"reflect"
)

// Group 项目集，如I0
type Group struct {
	Index int                      // 项目集编号
	PDs   []base.ProductionWithDot // 项目集中的产生式合集
}

// Groups 项目集规范族
var Groups []Group

type RelationB struct {
	Tag   base.Tag
	State int
}

// GroupRelation 用于表示项目集之间的关系，是一个map的数组，数组的索引表示Group.Index，
// map中以Tag为键，以对应的Group.Index 为值
var GroupRelation []map[base.Tag]RelationB

// ExpandGroup 根据传入的pd，扩充生成group，并在Groups中查重
// 如果还没有重复的，则将group添加到Groups中
// 无论是否有重复，都返回对应项目集在Groups中的索引
func ExpandGroup(pd base.ProductionWithDot) int {
	group := Group{
		Index: len(Groups),
		PDs:   []base.ProductionWithDot{pd},
	}

	// 如果点已经在最后了，则没有必要扩充了
	if pd.DotIndex < len(pd.Right) {
		// 扩充生成group
		index := 0
		for index < len(group.PDs) {
			if group.PDs[index].Right[group.PDs[index].DotIndex].Type == base.NONTERM {
				ExpandGroupRE(&group, index)
			}
			index++
		}

	}

	// 查重
	for i, g := range Groups {
		if reflect.DeepEqual(g.PDs, group.PDs) {
			return i
		}
	}

	// 如果是新的，则添加至Groups
	Groups = append(Groups, group)
	return group.Index
}

// RemoveRE 去重，包括Group和GroupRelation
func RemoveRE() {
	for i1, g1 := range Groups {
		for i2, g2 := range Groups {
			if reflect.DeepEqual(g1.PDs, g2.PDs) && i1 != i2 {
				// 从Groups中删掉
				Groups = append(Groups[:i2], Groups[i2+1:]...)
				// 把GroupRelation中，等于i2的改为i1，大于i2的减一
				for gri, m := range GroupRelation {
					for tag, i := range m {
						if i.State == i2 {
							GroupRelation[gri][tag] = RelationB{
								Tag:   GroupRelation[gri][tag].Tag,
								State: i1,
							}
						}
						if i.State > i2 {
							GroupRelation[gri][tag] = RelationB{
								Tag:   GroupRelation[gri][tag].Tag,
								State: GroupRelation[gri][tag].State - 1,
							}
						}
					}
				}
				// 把GroupRelation的i2行删掉
				GroupRelation = append(GroupRelation[:i2], GroupRelation[i2+1:]...)
			}
		}
	}
}

// MergeGroup 传入索引切片，合并group到索引小的那个（第一个）
func MergeGroup(is []int) {
	for _, i := range is[1:] {
		// 先把 Groups[i] 的 pds 放到 Groups[is[0]] 里头
		Groups[is[0]].PDs = append(Groups[is[0]].PDs, Groups[i].PDs...)
		// 再把i之后的往前移动一位
		Groups = append(Groups[:i], Groups[i+1:]...)
	}
}

// ExpandGroupRE 递归的扩充group, gi为传入groups中的要分析的group的索引
func ExpandGroupRE(group *Group, gi int) {
	// 此时，查找所有以点后的非终结符为左部，点在最前的式子
	if group.PDs[gi].Right[group.PDs[gi].DotIndex].Type == base.NONTERM {
		// 以点后的非终结符为左部的产生式
		productions := base.GetProdMap()[group.PDs[gi].Right[group.PDs[gi].DotIndex]]

		// 构造tags
		var tags []base.Tag
		// 如果点后面的后面没有符号了，则直接继承group.PDs[gi]的Tags
		if group.PDs[gi].DotIndex+1 >= len(group.PDs[gi].Right) {
			tags = group.PDs[gi].Tags
		}
		// 如果点后面的后面有符号，则使用该符号的FIRST集
		if group.PDs[gi].DotIndex+1 < len(group.PDs[gi].Right) {
			tags = base.GetFirst(group.PDs[gi].Right[group.PDs[gi].DotIndex+1])
		}

		// 根据productions生成productionsWithDot
		for _, production := range productions {
			reFlag := 0 // 用于检测是否有重
			// 首先在group中查重，如果已经有该产生式了，则只需要补充tags，而不用扩充pds
			for j, d := range group.PDs {
				if reflect.DeepEqual(d.Production, production) && d.DotIndex == 0 {
					for _, t := range tags {
						if !base.HasReTags(t, d.Tags) {
							group.PDs[j].Tags = append(group.PDs[j].Tags, t)
						}
					}
					reFlag++
					break
				}
			}

			// 如果依旧为0，则说明没有查到重复
			if reFlag == 0 {
				// 如果原来没有该产生式，则创建并添加至group
				pd := base.ProductionWithDot{
					Production: production,
					DotIndex:   0,
					Tags:       tags,
				}
				group.PDs = append(group.PDs, pd)
			}
		}
	}
}

func PrintFamily() {
	for _, group := range Groups {
		fmt.Printf("I%v:\n", group.Index)
		for _, pD := range group.PDs {
			fmt.Printf("%v\n", pD.ToString())
		}
		fmt.Printf("\n")
	}
}

// GenerateFamily 根据 prodMap 生成项目集规范族
func GenerateFamily() {
	startP := base.GetProductions()[0]
	startPd := base.ProductionWithDot{
		Production: startP,
		DotIndex:   0,
		Tags: []base.Tag{{
			Type:  base.TERM,
			Value: "$",
		}},
	}

	// 生成I0
	ExpandGroup(startPd)

	index := 0
	for index < len(Groups) {
		GroupRelation = append(GroupRelation, make(map[base.Tag]RelationB))
		for _, pd := range Groups[index].PDs {
			// 如果点不在最后，则让点后移
			if pd.DotIndex < len(pd.Right) {
				npd := base.ProductionWithDot{
					Production: pd.Production,
					DotIndex:   pd.DotIndex + 1,
					Tags:       pd.Tags,
				}
				// 获得新生成的group的索引
				nIndex := ExpandGroup(npd)
				// 如果已经有了，则合并
				if GroupRelation[index][pd.Right[pd.DotIndex]].State != 0 {
					MergeGroup([]int{GroupRelation[index][pd.Right[pd.DotIndex]].State, nIndex})
				} else {
					GroupRelation[index][pd.Right[pd.DotIndex]] = RelationB{
						Tag:   pd.Right[pd.DotIndex],
						State: nIndex,
					}
				}
			}
		}
		index++
	}

	// 最后去重
	RemoveRE()
}
