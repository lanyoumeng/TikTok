package errno

var (

	// ErrUserNotFound 表示未找到用户.
	ErrUserNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.UserNotFound", Message: "User was not found."}

	// Errhavefollowed	   已经关注
	Errhavefollowed = &Errno{HTTP: 400, Code: "FailedOperation.havefollowed", Message: "You have followed."}

	// Errnotfollowed	   未关注
	Errnotfollowed = &Errno{HTTP: 400, Code: "FailedOperation.notfollowed", Message: "You have not followed."}
)
