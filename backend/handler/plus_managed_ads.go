package handler

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ensurePlusManagedAdStorage(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS plus_managed_ads (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		ad_key VARCHAR(80) NOT NULL UNIQUE, title VARCHAR(160) NOT NULL, image_url VARCHAR(500) NOT NULL DEFAULT '', destination_url VARCHAR(500) NOT NULL DEFAULT '',
		position VARCHAR(30) NOT NULL DEFAULT 'home-banner', weight INT NOT NULL DEFAULT 0, start_at DATETIME NULL, end_at DATETIME NULL,
		description VARCHAR(500) NOT NULL DEFAULT '', published TINYINT(1) NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		KEY idx_ad_position (position), KEY idx_ad_published (published)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Plus后台托管广告'`)
	return err
}

func validAdPosition(value string) bool {
	return value == "home-banner" || value == "sidebar" || value == "popup"
}

type managedAdPayload struct {
	ID             int64  `json:"id"`
	AdKey          string `json:"adKey"`
	Title          string `json:"title"`
	ImageURL       string `json:"imageUrl"`
	DestinationURL string `json:"destinationUrl"`
	Position       string `json:"position"`
	Weight         int    `json:"weight"`
	StartAt        string `json:"startAt"`
	EndAt          string `json:"endAt"`
	Description    string `json:"description"`
	Published      bool   `json:"published"`
}

func managedAdResponse(row scanner) (gin.H, error) {
	var ad managedAdPayload
	var start, end sql.NullTime
	if err := row.Scan(&ad.ID, &ad.AdKey, &ad.Title, &ad.ImageURL, &ad.DestinationURL, &ad.Position, &ad.Weight, &start, &end, &ad.Description, &ad.Published); err != nil {
		return nil, err
	}
	if start.Valid {
		ad.StartAt = start.Time.Format(time.RFC3339)
	}
	if end.Valid {
		ad.EndAt = end.Time.Format(time.RFC3339)
	}
	return gin.H{"id": ad.ID, "adKey": ad.AdKey, "title": ad.Title, "imageUrl": ad.ImageURL, "destinationUrl": ad.DestinationURL, "position": ad.Position, "weight": ad.Weight, "startAt": ad.StartAt, "endAt": ad.EndAt, "description": ad.Description, "published": ad.Published}, nil
}

type scanner interface{ Scan(...any) error }

func AdminManagedAdList(c *gin.Context) {
	db, err := managedResourceDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	defer db.Close()
	rows, err := db.Query(`SELECT id,ad_key,title,image_url,destination_url,position,weight,start_at,end_at,description,published FROM plus_managed_ads ORDER BY weight DESC,id DESC`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "读取广告失败"})
		return
	}
	defer rows.Close()
	list := []gin.H{}
	for rows.Next() {
		item, e := managedAdResponse(rows)
		if e == nil {
			list = append(list, item)
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": gin.H{"list": list}})
}

func AdminManagedAdSave(c *gin.Context) {
	var req managedAdPayload
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.AdKey) == "" || strings.TrimSpace(req.Title) == "" || !validAdPosition(req.Position) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "广告标识、标题或广告位不合法"})
		return
	}
	db, err := managedResourceDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	defer db.Close()
	var start, end any
	if strings.TrimSpace(req.StartAt) != "" {
		start = req.StartAt
	}
	if strings.TrimSpace(req.EndAt) != "" {
		end = req.EndAt
	}
	if req.ID > 0 {
		_, err = db.Exec(`UPDATE plus_managed_ads SET ad_key=?,title=?,image_url=?,destination_url=?,position=?,weight=?,start_at=?,end_at=?,description=?,published=? WHERE id=?`, req.AdKey, req.Title, req.ImageURL, req.DestinationURL, req.Position, req.Weight, start, end, req.Description, req.Published, req.ID)
	} else {
		_, err = db.Exec(`INSERT INTO plus_managed_ads (ad_key,title,image_url,destination_url,position,weight,start_at,end_at,description,published) VALUES (?,?,?,?,?,?,?,?,?,?)`, req.AdKey, req.Title, req.ImageURL, req.DestinationURL, req.Position, req.Weight, start, end, req.Description, req.Published)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "保存广告失败，广告标识可能已存在"})
		return
	}
	invalidateManagedAdCache()
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "广告已保存"})
}
func AdminManagedAdDelete(c *gin.Context) {
	db, err := managedResourceDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	defer db.Close()
	if _, err = db.Exec("DELETE FROM plus_managed_ads WHERE id=?", c.Param("id")); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "删除广告失败"})
		return
	}
	invalidateManagedAdCache()
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "广告已删除"})
}

func invalidateManagedAdCache() {
	advertisementCacheL.Lock()
	advertisementCache = map[string]advertisementCacheEntry{}
	advertisementCacheL.Unlock()
}

func managedAdsForPosition(position string) []advertisementRecord {
	db, err := managedResourceDB()
	if err != nil {
		return nil
	}
	defer db.Close()
	rows, err := db.Query(`SELECT ad_key,title,image_url,destination_url,position,weight,start_at,end_at,description FROM plus_managed_ads WHERE published=1 AND position=? ORDER BY weight DESC,id DESC`, position)
	if err != nil {
		return nil
	}
	defer rows.Close()
	now := time.Now()
	result := []advertisementRecord{}
	for rows.Next() {
		var ad advertisementRecord
		var start, end sql.NullTime
		if rows.Scan(&ad.ID, &ad.Title, &ad.ImageURL, &ad.DestinationURL, &ad.Position, &ad.Weight, &start, &end, &ad.Description) != nil {
			continue
		}
		if start.Valid {
			ad.StartAt = start.Time.Format(time.RFC3339)
		}
		if end.Valid {
			ad.EndAt = end.Time.Format(time.RFC3339)
		}
		if advertisementInWindow(ad, now) {
			result = append(result, ad)
		}
	}
	return result
}
