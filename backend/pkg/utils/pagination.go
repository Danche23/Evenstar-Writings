package utils

// ClampPage 归一化分页参数，防止异常/恶意入参导致全表扫描或负偏移。
//   - page < 1       → 1
//   - size < 1       → defSize
//   - size > maxSize → maxSize
//
// 说明：GORM 中 Limit(-1) 表示不限制条数，因此 size 必须钳制为正数上限。
func ClampPage(page, size, defSize, maxSize int) (int, int) {
	if defSize < 1 {
		defSize = 10
	}
	if maxSize < defSize {
		maxSize = defSize
	}
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = defSize
	}
	if size > maxSize {
		size = maxSize
	}
	return page, size
}
