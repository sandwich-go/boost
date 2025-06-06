package cloud

import "net/http"

//go:generate optionGen  --option_return_previous=false
func PutOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@ContentType(comment="上传文件的类型")
		"ContentType": "application/octet-stream",
		// annotation@ContentDisposition(comment="上传文件的内容描述")
		"ContentDisposition": "",
		// annotation@CacheControl(comment="上传文件的缓存控制")
		"CacheControl": "",
		// annotation@DisableContentSha256(comment="禁止发送 Content-Sha256, 在非s3场景下，文件较小时，content-sha256 也会出现在文件内容中")
		"DisableContentSha256": false,
		// annotation@CustomHeader(comment="自定义上传时附加的http header")
		"CustomHeader": http.Header(nil),
		// annotation@CustomMeta(comment="自定义上传时的 meta 信息")
		"CustomMeta": map[string]string(nil),
		// annotation@SendContentMd5(comment="gcs需要在上传时，minio 参数中指定 md5-base64")
		"SendContentMd5": false,
		// annotation@FileMD5(comment="文件MD5")
		"FileMD5": "",
	}
}

//go:generate optionGen  --option_return_previous=false
func StorageOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@Region(comment="云存储的Region")
		"Region": "",
		// annotation@StorageType(comment="云存储类型")
		"StorageType": StorageType(""),
	}
}
