package utils

// PageRequest 表示分页查询请求
type PageRequest struct {
	PageNumber int `json:"page"` // 页码，从0开始。
	PageSize   int `json:"size"` // 每页记录数。
}

// Page 表示分页查询结果
type Page[T any] struct {
	TotalPages       int   `json:"totalPages,omitempty"`
	TotalElements    int64 `json:"totalElements,omitempty"`
	Number           int   `json:"number,omitempty"`
	Size             int   `json:"size,omitempty"`
	NumberOfElements int   `json:"numberOfElements,omitempty"`
	Content          []T   `json:"content,omitempty"`
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
		Number:           page,
		Size:             size,
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
		Number:           len(content),
		Size:             len(content),
		NumberOfElements: len(content),
		Content:          content,
	}
}
