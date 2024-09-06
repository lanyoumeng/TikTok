package aliyun

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"os"
	"os/exec"
	"strings"
	"sync"
	"video/internal/conf"
	"video/internal/pkg/model"

	"github.com/go-kratos/kratos/v2/log"

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

	//压缩视频
	VideoPath, err := compressVideo(videoKafkaMessage.VideoPath)
	if err != nil {
		log.Debug("compressVideo false:", err)
		return "", "", err
	}

	coverName := strings.TrimSuffix(videoKafkaMessage.VideoFileName, ".mp4") + "_cover.jpg"
	// 截取视频封面
	coverPath, err := generateVideoCover(VideoPath)
	if err != nil {
		log.Debug("generateVideoCover false:", err)
		return "", "", err
	}
	log.Debug("coverPath:", coverPath)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// 直接上传文件
		err := bucket.PutObjectFromFile("videos/"+videoKafkaMessage.VideoFileName, videoKafkaMessage.VideoPath, oss.ObjectACL(oss.ACLPublicRead))
		if err != nil {
			log.Debugf("upload file failed: %v", err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		// 直接上传文件
		err := bucket.PutObjectFromFile("covers/"+coverName, coverPath, oss.ObjectACL(oss.ACLPublicRead))
		if err != nil {
			log.Debugf("upload file failed: %v", err)
			return
		}
	}()

	// 获取可播放的视频 URL
	playURL := fmt.Sprintf("https://%s.%s/%s", bucket.BucketName, bucket.GetConfig().Endpoint, "videos/"+videoKafkaMessage.VideoFileName)
	coverURL := fmt.Sprintf("https://%s.%s/%s", bucket.BucketName, bucket.GetConfig().Endpoint, "covers/"+coverName)

	wg.Wait()

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

	return playURL, coverURL, nil
}

// 压缩视频
func compressVideo(inputVideoPath string) (string, error) {
	outputVideoPath := strings.TrimSuffix(inputVideoPath, ".mp4") + "CMP.mp4"
	command := []string{
		"-i", inputVideoPath,
		"-c:v", "libx264",
		//"-b:v", "1M", // 使用比特率代替 -crf
		"-crf", "43", //较高的CRF值会减小文件大小但可能降低视频质量
		"-y", // This option enables overwriting without asking
		outputVideoPath,
	}
	cmd := exec.Command("ffmpeg", command...)

	//// 打开已存在的日志文件，如果不存在则创建
	//logFile, err := os.OpenFile("logs/video.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	return outputVideoPath, err
	//}
	//defer logFile.Close()
	//
	//cmd.Stderr = logFile // 将stderr重定向到指定的日志文件
	// cmd.Stderr = os.Stderr // 将stderr重定向到控制台以查看错误消息

	err := cmd.Run()
	if err != nil {
		logger.Debug("Error compressing video:", err)
		return outputVideoPath, err
	}

	return outputVideoPath, nil
}

// 截取封面
func generateVideoCover(videoPath string) (string, error) {
	// 使用 ffmpeg 获取视频的第一帧作为封面
	coverFilePath := strings.TrimSuffix(videoPath, ".mp4") + "_cover.jpg"
	command := []string{
		"-i", videoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		coverFilePath,
	}
	cmd := exec.Command("ffmpeg", command...)

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error generating cover:", err)
		return "", err
	}
	return coverFilePath, nil
}
