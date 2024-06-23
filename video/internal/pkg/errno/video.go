package errno

var (

	//redis zset--videoAll中没有视频id数据
	ErrRedisVIdNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.VideoNotFound", Message: "redis zset--videoAll中没有视频id数据"}

	//redis hash--videoInfo中没有视频信息
	ErrRedisVInfoNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.VideoInfoNotFound", Message: "redis hash--videoInfo中没有视频信息"}

	//redis set--publishVids中没有发布视频id数据
	ErrRedisPublishVidsNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.PublishVidsNotFound", Message: "redis set--publishVids中没有发布视频id数据"}
)
