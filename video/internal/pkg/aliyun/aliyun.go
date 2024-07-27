package aliyun

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"video/internal/conf"
	"video/internal/pkg/model"

	"github.com/go-kratos/kratos/v2/log"

	ffmpeg "github.com/u2takey/ffmpeg-go"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func NewBucket(aliyun *conf.AliyunOSS) (*oss.Bucket, error) {
	// 初始化 OSS 客户端
	endpoint := aliyun.Endpoint               // OSS 的 endpoint
	accessKeyID := aliyun.AccessKeyId         // 阿里云的 Access Key ID
	accessKeySecret := aliyun.AccessKeySecret // 阿里云的 Access Key Secret
	bucketName := aliyun.BucketName           // OSS 的 Bucket 名称

	// 创建 OSS 客户端对象
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		log.Fatalf("Error creating OSS client: %s", err)
	}

	// 获取存储空间（Bucket）对象
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Fatalf("Error obtaining bucket: %s", err)
	}
	return bucket, nil

}

func UploadFile(bucket *oss.Bucket, videoKafkaMessage *model.VideoKafkaMessage) (string, string, error) {

	currentDir, _ := os.Getwd()
	log.Debug(" UploadFile工作目录:", currentDir, "  videoKafkaMessage.VideoPath:", videoKafkaMessage.VideoPath)

	// 上传文件
	//log.Debug("begin upload file to oss")
	err := bucket.PutObjectFromFile("videos/"+videoKafkaMessage.VideoFileName, videoKafkaMessage.VideoPath, oss.ObjectACL(oss.ACLPublicRead))
	if err != nil {
		log.Debugf("upload file failed: %v", err)
		return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}

	// 获取可播放的视频 URL
	playURL := fmt.Sprintf("https://%s.%s/%s", bucket.BucketName, bucket.GetConfig().Endpoint, "videos/"+videoKafkaMessage.VideoFileName)

	// 获取视频封面 URL	buf := bytes.NewBuffer(nil)
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(playURL).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 3)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Debugf("get video cover failed: %v", err)
		return "", "", err
	}
	coverFilename := strings.TrimSuffix(videoKafkaMessage.VideoPath, ".mp4") + "_cover.jpeg"
	// 上传封面
	err = bucket.PutObject("covers/"+coverFilename, buf)
	if err != nil {
		log.Debugf("upload cover failed: %v", err)
		return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}
	coverURL := fmt.Sprintf("https://%s.%s/%s", bucket.BucketName, bucket.GetConfig().Endpoint, "covers/"+coverFilename)

	return playURL, coverURL, nil
}
