package utils

// PageRequest 表示分页查询请求
type PageRequest struct {
	PageNumber int `json:"page"` // 页码，从0开始。
	PageSize   int `json:"size"` // 每页记录数。
}

// Page 表示分页查询结果
type Page[T any] struct {
	TotalPages       int   `json:"totalPages"`
	TotalElements    int64 `json:"totalElements"`
	PageNumber       int   `json:"number"`
	PageSize         int   `json:"size"`
	NumberOfElements int   `json:"numberOfElements"`
	Content          []T   `json:"content"`
}

func (pr *PageRequest) GetStartRowIndex() int {
	return pr.PageNumber * pr.PageSize
}

// MakePage 构造分页查询结果
func MakePage[T any](page, size int, total int64, content []T) Page[T] {
	var totalPages int
	if size <= 0 {
		if total > 0 {
			totalPages = 1
		} else {
			totalPages = 0
		}
	} else {
		totalPages = int(total / int64(size))
		if total%int64(size) == 0 {
			totalPages++
		}
	}

	if content == nil {
		content = make([]T, 0)
	}

	return Page[T]{
		TotalPages:       totalPages,
		TotalElements:    total,
		PageNumber:       page,
		PageSize:         size,
		NumberOfElements: len(content),
		Content:          content,
	}
}

// MakePage2 构造分页查询结果
func MakePage2[T any](content []T) Page[T] {
	var totalPages int
	if len(content) > 0 {
		totalPages = 1
	} else {
		totalPages = 0
	}
	return Page[T]{
		TotalPages:       totalPages,
		TotalElements:    int64(len(content)),
		PageNumber:       len(content),
		PageSize:         len(content),
		NumberOfElements: len(content),
		Content:          content,
	}
}
