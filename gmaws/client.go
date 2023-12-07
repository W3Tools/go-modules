package gmaws

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	gm "github.com/W3Tools/go-modules"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type AwsClient struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Region          string `json:"region"`

	S3 S3Client
}

type S3Client struct {
	Client   *s3.S3
	Bucket   string `json:"bucket"`
	BaseUri  string `json:"base_uri"`
	Subpath  string `json:"subpath"`
	Endpoint string `json:"endpoint"`
}

func (c *AwsClient) NewS3Client() error {
	credentials := credentials.NewStaticCredentials(c.AccessKeyId, c.AccessKeySecret, "")

	awsCfg := &aws.Config{
		Region:      aws.String(c.Region),
		Credentials: credentials,
	}

	if !strings.EqualFold(c.S3.Endpoint, "") {
		awsCfg.Endpoint = aws.String(c.S3.Endpoint)
	}

	s, err := session.NewSession(awsCfg)
	if err != nil {
		return err
	}

	if c.S3.Client == nil {
		c.S3.Client = s3.New(s)
	}
	return nil
}

func (c *AwsClient) UploadObjectToS3(fileBytes []byte, remotePath string) (ret *UploadToS3Response, err error) {
	keyName := remotePath
	if !strings.EqualFold(c.S3.Subpath, "") {
		keyName = fmt.Sprintf("%v/%v", c.S3.Subpath, remotePath)
	}
	_, err = c.S3.Client.PutObjectWithContext(context.Background(), &s3.PutObjectInput{
		Bucket:       aws.String(c.S3.Bucket),
		Key:          aws.String(keyName),
		Body:         bytes.NewReader(fileBytes),
		ContentType:  aws.String(http.DetectContentType(fileBytes)),
		StorageClass: aws.String("STANDARD"),
	})
	if err != nil {
		return nil, err
	}

	ret = &UploadToS3Response{
		ObjectType: http.DetectContentType(fileBytes),
		ObjectSize: len(fileBytes),
		BaseUri:    c.S3.BaseUri,
		RemotePath: remotePath,
		FullPath:   fmt.Sprintf("%s/%s", c.S3.BaseUri, remotePath),
	}
	return
}

func (c *AwsClient) UploadObjectToS3WithRandomKey(fileBytes []byte) (ret *UploadToS3Response, err error) {
	return c.UploadObjectToS3(fileBytes, uuid.New().String())
}

func (c *AwsClient) UploadSingleFileToS3WithRandomKey(filePath string) (ret *UploadToS3Response, err error) {
	data, err := gm.ReadFileBytes(filePath)
	if err != nil {
		return nil, err
	}

	return c.UploadObjectToS3WithRandomKey(data)
}

func (c *AwsClient) UploadSingleFileToS3WithSpecifyKey(filePath string, keyName string) (ret *UploadToS3Response, err error) {
	data, err := gm.ReadFileBytes(filePath)
	if err != nil {
		return nil, err
	}

	return c.UploadObjectToS3(data, keyName)
}
