package gmaws

type UploadToS3Response struct {
	ObjectType string
	ObjectSize int
	BaseUri    string
	RemotePath string
	FullPath   string
}
