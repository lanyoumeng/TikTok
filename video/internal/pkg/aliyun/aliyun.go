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

	// 直接上传文件
	err := bucket.PutObjectFromFile("videos/"+videoKafkaMessage.VideoFileName, videoKafkaMessage.VideoPath, oss.ObjectACL(oss.ACLPublicRead))
	if err != nil {
		log.Debugf("upload file failed: %v", err)
		return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}

	//分片上传
	//// 将本地文件分片，且分片数量指定为3。
	//locaFilename := videoKafkaMessage.VideoPath
	////objectName := fmt.Sprintf("videos/%s", videoKafkaMessage.VideoFileName)
	//objectName := videoKafkaMessage.VideoFileName
	//chunks, err := oss.SplitFileByPartNum(locaFilename, 3)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}
	//fd, err := os.Open(locaFilename)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}
	//defer fd.Close()
	//
	//// 指定过期时间。
	//expires := time.Date(2049, time.January, 10, 23, 0, 0, 0, time.UTC)
	//// 如果需要在初始化分片时设置请求头，请参考以下示例代码。
	//options := []oss.Option{
	//	oss.MetadataDirective(oss.MetaReplace),
	//	oss.Expires(expires),
	//	// 指定该Object被下载时的网页缓存行为。
	//	// oss.CacheControl("no-cache"),
	//	// 指定该Object被下载时的名称。
	//	// oss.ContentDisposition("attachment;filename=FileName.txt"),        ,
	//	// 指定对返回的Key进行编码，目前支持URL编码。
	//	// oss.EncodingType("url"),
	//	// 指定Object的存储类型。
	//	// oss.ObjectStorageClass(oss.StorageStandard),
	//}
	//
	//// 步骤1：初始化一个分片上传事件。
	//imur, err := bucket.InitiateMultipartUpload(objectName, options...)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}
	//// 步骤2：上传分片。
	//var parts []oss.UploadPart
	//for _, chunk := range chunks {
	//	fd.Seek(chunk.Offset, os.SEEK_SET)
	//	// 调用UploadPart方法上传每个分片。
	//	part, err := bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
	//	if err != nil {
	//		fmt.Println("Error:", err)
	//		os.Exit(-1)
	//	}
	//	parts = append(parts, part)
	//}
	//
	//// 步骤3：完成分片上传。
	//cmur, err := bucket.CompleteMultipartUpload(imur, parts)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}
	//fmt.Println("cmur:", cmur)

	// 获取可播放的视频 URL
	playURL := fmt.Sprintf("https://%s.%s/%s", bucket.BucketName, bucket.GetConfig().Endpoint, "videos/"+videoKafkaMessage.VideoFileName)
	//playURL := fmt.Sprintf("https://%s.%s/%s", bucket.BucketName, bucket.GetConfig().Endpoint, videoKafkaMessage.VideoFileName)
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
