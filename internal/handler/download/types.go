package download

type PkgInfo struct {
	PkgType int `xorm:"pkg_type"`
	PkgID   int `xorm:"pkg_id"`
	Order   int `xorm:"pkg_order"`
	Size    int `xorm:"pkg_size"`
}
