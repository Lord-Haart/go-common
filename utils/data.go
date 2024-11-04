package utils

func GetIdList[T any, ID any](items []T, getId func(T) ID) []ID {
	result := make([]ID, 0, len(items))

	for _, item := range items {
		result = append(result, getId(item))
	}

	return result
}

func LeftJoin[T1 any, T2 any, ID comparable](items1 []T1, items2 []T2, f1 func(T1) ID, f2 func(T2) ID, f3 func(T1, T2)) int {
	c := 0
	for _, item1 := range items1 {
		for _, item2 := range items2 {
			if f1(item1) == f2(item2) {
				f3(item1, item2)
				c++
			}
		}
	}
	return c
}

func LeftJoin2[T1 any, T2 any](items1 []T1, items2 []T2, f1 func(T1, T2) bool, f2 func(T1, T2)) int {
	c := 0
	for _, item1 := range items1 {
		for _, item2 := range items2 {
			if f1(item1, item2) {
				f2(item1, item2)
				c++
			}
		}
	}
	return c
}
