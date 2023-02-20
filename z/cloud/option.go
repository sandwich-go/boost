package cloud

//go:generate optionGen  --option_return_previous=false
func PutOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@ContentType(comment="上传文件的类型")
		"ContentType": "application/octet-stream",
		// annotation@ContentDisposition(comment="上传文件的内容描述")
		"ContentDisposition": "",
		// annotation@CacheControl(comment="上传文件的缓存控制")
		"CacheControl": "",
	}
}

//go:generate optionGen  --option_return_previous=false
func StorageOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@Region(comment="云存储的Region")
		"Region": "",
	}
}
