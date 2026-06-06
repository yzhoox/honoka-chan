package costume

import honokautils "honoka-chan/internal/utils"

func CostumeApi(action string) (res any, err error) {
	switch action {
	case "costumeList":
		res, err = costumeList()
	default:
		err = honokautils.NewUnimplementedActionError("costume", action)
	}
	return res, err
}
