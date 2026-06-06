package notice

import honokautils "honoka-chan/internal/utils"

func NoticeApi(action string) (res any, err error) {
	switch action {
	case "noticeMarquee":
		res, err = noticeMarquee()
	default:
		err = honokautils.NewUnimplementedActionError("notice", action)
	}
	return res, err
}
