package notice

import "fmt"

func NoticeApi(action string) (res any, err error) {
	switch action {
	case "noticeMarquee":
		res, err = noticeMarquee()
	default:
		err = fmt.Errorf("unimplemented action: notice: %s", action)
	}
	return res, err
}
