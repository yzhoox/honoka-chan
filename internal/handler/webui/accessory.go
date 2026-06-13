package webui

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	webuischema "honoka-chan/internal/schema/webui"
	"honoka-chan/pkg/db"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const accessoryDataPath = "assets/serverdata/accessory_data.json"

type accessoryInventoryItem struct {
	AccessoryID int `json:"accessory_id"`
	Count       int `json:"count"`
	WornCount   int `json:"worn_count"`
}

func accessoryPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/accessory.html", gin.H{
		"menu": 3,
		"url":  strings.Split(ctx.Request.URL.String(), "?")[0],
	})
}

func accessoryListPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/accessory_list.html", gin.H{
		"menu": 3,
		"url":  strings.Split(ctx.Request.URL.String(), "?")[0],
	})
}

func accessoryData(ctx *gin.Context) {
	data, err := os.ReadFile(accessoryDataPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, webuischema.Msg{
			Code:    1,
			Message: "饰品数据读取失败！",
		})
		return
	}

	ctx.Data(http.StatusOK, "application/json; charset=utf-8", data)
}

func addAccessory(ctx *gin.Context) {
	accessoryID, err := strconv.Atoi(strings.TrimSpace(ctx.PostForm("accessory_id")))
	if err != nil || accessoryID <= 0 {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品 ID 格式有误！",
		})
		return
	}

	quantity, err := strconv.Atoi(strings.TrimSpace(ctx.PostForm("quantity")))
	if err != nil || quantity <= 0 || quantity > 99 {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "数量格式有误！",
		})
		return
	}

	exists, err := db.MainEng.Table("accessory_m").
		Where("accessory_id = ? AND is_material = ? AND rarity = ?", accessoryID, 0, 4).
		Exist()
	if err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品校验失败！",
		})
		return
	}
	if !exists {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品不存在！",
		})
		return
	}

	session := db.UserEng.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "数据库事务开启失败！",
		})
		return
	}

	for range quantity {
		if _, err := session.Insert(&usermodel.UserAccessory{
			UserID:      ctx.GetInt("userid"),
			AccessoryID: accessoryID,
			Exp:         0,
		}); err != nil {
			session.Rollback()
			ctx.JSON(http.StatusOK, webuischema.Msg{
				Code:    1,
				Message: "饰品添加失败！",
			})
			return
		}
	}

	if err := session.Commit(); err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品添加失败！",
		})
		return
	}

	ctx.JSON(http.StatusOK, webuischema.Msg{
		Code:    0,
		Message: "添加成功，请重新进入游戏查看！",
	})
}

func accessoryInventoryData(ctx *gin.Context) {
	userID := ctx.GetInt("userid")

	accessories := []usermodel.UserAccessory{}
	if err := db.UserEng.Table(new(usermodel.UserAccessory)).
		Where("user_id = ?", userID).
		OrderBy("accessory_id ASC, accessory_owning_user_id ASC").
		Find(&accessories); err != nil {
		ctx.JSON(http.StatusInternalServerError, webuischema.Msg{
			Code:    1,
			Message: "饰品库存读取失败！",
		})
		return
	}

	wears := []usermodel.UserAccessoryWear{}
	if err := db.UserEng.Table(new(usermodel.UserAccessoryWear)).
		Where("user_id = ?", userID).
		Find(&wears); err != nil {
		ctx.JSON(http.StatusInternalServerError, webuischema.Msg{
			Code:    1,
			Message: "饰品佩戴信息读取失败！",
		})
		return
	}

	owningIDToAccessoryID := make(map[int]int, len(accessories))
	inventoryMap := make(map[int]*accessoryInventoryItem)
	for _, row := range accessories {
		owningIDToAccessoryID[row.AccessoryOwningUserID] = row.AccessoryID
		if inventoryMap[row.AccessoryID] == nil {
			inventoryMap[row.AccessoryID] = &accessoryInventoryItem{
				AccessoryID: row.AccessoryID,
			}
		}
		inventoryMap[row.AccessoryID].Count++
	}

	for _, wear := range wears {
		accessoryID, ok := owningIDToAccessoryID[wear.AccessoryOwningUserID]
		if !ok {
			continue
		}
		if inventoryMap[accessoryID] == nil {
			inventoryMap[accessoryID] = &accessoryInventoryItem{
				AccessoryID: accessoryID,
			}
		}
		inventoryMap[accessoryID].WornCount++
	}

	result := make([]accessoryInventoryItem, 0, len(inventoryMap))
	for _, item := range inventoryMap {
		result = append(result, *item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].AccessoryID < result[j].AccessoryID
	})

	ctx.JSON(http.StatusOK, result)
}

func deleteAccessory(ctx *gin.Context) {
	userID := ctx.GetInt("userid")

	accessoryID, err := strconv.Atoi(strings.TrimSpace(ctx.PostForm("accessory_id")))
	if err != nil || accessoryID <= 0 {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品 ID 格式有误！",
		})
		return
	}

	quantity, err := strconv.Atoi(strings.TrimSpace(ctx.PostForm("quantity")))
	if err != nil || quantity <= 0 || quantity > 99 {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "数量格式有误！",
		})
		return
	}

	session := db.UserEng.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "数据库事务开启失败！",
		})
		return
	}

	accessories := []usermodel.UserAccessory{}
	if err := session.Table(new(usermodel.UserAccessory)).
		Where("user_id = ? AND accessory_id = ?", userID, accessoryID).
		OrderBy("accessory_owning_user_id ASC").
		Find(&accessories); err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品库存读取失败！",
		})
		return
	}
	if len(accessories) < quantity {
		session.Rollback()
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "删除数量超过持有数量！",
		})
		return
	}

	wears := []usermodel.UserAccessoryWear{}
	if err := session.Table(new(usermodel.UserAccessoryWear)).
		Where("user_id = ?", userID).
		Find(&wears); err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品佩戴信息读取失败！",
		})
		return
	}

	wornSet := make(map[int]struct{}, len(wears))
	for _, wear := range wears {
		wornSet[wear.AccessoryOwningUserID] = struct{}{}
	}

	deleteIDs := make([]int, 0, quantity)
	for _, row := range accessories {
		if _, worn := wornSet[row.AccessoryOwningUserID]; worn {
			continue
		}
		deleteIDs = append(deleteIDs, row.AccessoryOwningUserID)
		if len(deleteIDs) == quantity {
			break
		}
	}
	if len(deleteIDs) < quantity {
		for _, row := range accessories {
			if _, worn := wornSet[row.AccessoryOwningUserID]; !worn {
				continue
			}
			deleteIDs = append(deleteIDs, row.AccessoryOwningUserID)
			if len(deleteIDs) == quantity {
				break
			}
		}
	}

	if len(deleteIDs) != quantity {
		session.Rollback()
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品删除目标计算失败！",
		})
		return
	}

	if _, err := session.Table(new(usermodel.UserAccessoryWear)).
		In("accessory_owning_user_id", deleteIDs).
		Where("user_id = ?", userID).
		Delete(new(usermodel.UserAccessoryWear)); err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品卸下失败！",
		})
		return
	}

	if _, err := session.Table(new(usermodel.UserAccessory)).
		In("accessory_owning_user_id", deleteIDs).
		Where("user_id = ? AND accessory_id = ?", userID, accessoryID).
		Delete(new(usermodel.UserAccessory)); err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品删除失败！",
		})
		return
	}

	if err := session.Commit(); err != nil {
		session.Rollback()
		ctx.JSON(http.StatusOK, webuischema.Msg{
			Code:    1,
			Message: "饰品删除失败！",
		})
		return
	}

	ctx.JSON(http.StatusOK, webuischema.Msg{
		Code:    0,
		Message: "删除成功，请重新登录游戏生效！",
	})
}

func init() {
	router.AddHandler("admin", "GET", "/accessory", middleware.WebAuth, accessoryPage)
	router.AddHandler("admin", "GET", "/accessory/list", middleware.WebAuth, accessoryListPage)
	router.AddHandler("admin", "GET", "/accessory/data", middleware.WebAuth, accessoryData)
	router.AddHandler("admin", "GET", "/accessory/list/data", middleware.WebAuth, accessoryInventoryData)
	router.AddHandler("admin", "POST", "/accessory/add", middleware.WebAuth, addAccessory)
	router.AddHandler("admin", "POST", "/accessory/delete", middleware.WebAuth, deleteAccessory)
}
