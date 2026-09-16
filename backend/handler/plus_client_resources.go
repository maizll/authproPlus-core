package handler

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// PlusClientAppID resolves a public app key to an enabled application. The
// app key scopes the read-only catalog; future activation tokens should provide
// the stronger user/device entitlement check for paid downloads.
func PlusClientAppID(c *gin.Context) (int64, bool) {
	appKey := strings.TrimSpace(c.Param("appKey"))
	if appKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "应用标识不能为空"})
		return 0, false
	}
	db, err := openSystemConfigDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return 0, false
	}
	var appID int64
	if err := db.QueryRow("SELECT id FROM apps WHERE app_key = ? AND enabled = 1", appKey).Scan(&appID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "应用不存在或已禁用"})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "读取应用失败"})
		}
		return 0, false
	}
	return appID, true
}

func PlusClientPlugins(c *gin.Context) {
	appID, ok := PlusClientAppID(c)
	if !ok {
		return
	}
	db, err := openSystemConfigDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	rows, err := db.Query(`SELECT plugin_key, name, description, category, version, package_url, sha256,
		pricing_type, price, currency, trial_days, requires_license
		FROM app_plugins WHERE app_id = ? AND enabled = 1 AND published = 1 ORDER BY category, name`, appID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "读取插件目录失败"})
		return
	}
	defer rows.Close()
	items := make([]gin.H, 0)
	for rows.Next() {
		var key, name, description, category, version, packageURL, sha256, pricingType, currency string
		var price float64
		var trialDays int
		var requiresLicense bool
		if err := rows.Scan(&key, &name, &description, &category, &version, &packageURL, &sha256, &pricingType, &price, &currency, &trialDays, &requiresLicense); err != nil {
			continue
		}
		items = append(items, gin.H{"id": key, "name": name, "description": description, "category": category, "version": version, "downloadUrl": packageURL, "sha256": sha256, "pricingType": pricingType, "price": price, "currency": currency, "trialDays": trialDays, "requiresLicense": requiresLicense})
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": gin.H{"appId": appID, "list": items}})
}

func PlusClientTemplates(c *gin.Context) {
	appID, ok := PlusClientAppID(c)
	if !ok {
		return
	}
	db, err := openSystemConfigDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	rows, err := db.Query(`SELECT template_key, name, description, version, schema_version, content_path, preview_path, sha256
		FROM app_home_templates WHERE app_id = ? AND enabled = 1 AND published = 1 ORDER BY name`, appID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "读取模板目录失败"})
		return
	}
	defer rows.Close()
	items := make([]gin.H, 0)
	for rows.Next() {
		var key, name, description, version, contentPath, previewPath, sha256 string
		var schemaVersion int
		if err := rows.Scan(&key, &name, &description, &version, &schemaVersion, &contentPath, &previewPath, &sha256); err != nil {
			continue
		}
		items = append(items, gin.H{"id": key, "name": name, "description": description, "version": version, "schemaVersion": schemaVersion, "contentUrl": contentPath, "previewUrl": previewPath, "sha256": sha256})
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": gin.H{"appId": appID, "list": items}})
}

func PlusClientAdvertisements(c *gin.Context) {
	appID, ok := PlusClientAppID(c)
	if !ok {
		return
	}
	position := strings.TrimSpace(c.Query("position"))
	if position == "" {
		position = "home-banner"
	}
	now := time.Now()
	db, err := openSystemConfigDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	rows, err := db.Query(`SELECT id, title, image_url, destination_url, description, weight, start_at, end_at
		FROM app_advertisements WHERE app_id = ? AND position = ? AND enabled = 1
		AND (start_at IS NULL OR start_at <= ?) AND (end_at IS NULL OR end_at >= ?)
		ORDER BY weight DESC, id ASC`, appID, position, now, now)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "读取广告失败"})
		return
	}
	defer rows.Close()
	records := make([]gin.H, 0)
	for rows.Next() {
		var id int64
		var title, imageURL, destinationURL, description string
		var weight int
		var startAt, endAt sql.NullTime
		if err := rows.Scan(&id, &title, &imageURL, &destinationURL, &description, &weight, &startAt, &endAt); err != nil {
			continue
		}
		records = append(records, gin.H{"id": id, "title": title, "imageUrl": imageURL, "destinationUrl": destinationURL, "position": position, "weight": weight, "description": description, "startAt": startAt.Time.Format(time.RFC3339), "endAt": endAt.Time.Format(time.RFC3339)})
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": gin.H{"appId": appID, "records": records}})
}
