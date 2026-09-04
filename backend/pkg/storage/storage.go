package storage

// Storage 文件存储抽象（开发期本地 mock，后续替换为 OSS 实现）
type Storage interface {
	// Save 保存文件，objectKey 为存储路径，返回对外访问 url
	Save(objectKey string, data []byte) (url string, err error)
	// Delete 删除文件，url 为 Save 返回的对外访问地址
	Delete(url string) error
}
