package handler

import (
	"database/sql"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Plus 托管资源允许管理员直接在后台维护目录，不必修改 Git 仓库。
func ensurePlusManagedResourceStorage(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS plus_managed_plugins (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		plugin_id VARCHAR(60) NOT NULL UNIQUE, category VARCHAR(20) NOT NULL DEFAULT 'other',
		name VARCHAR(100) NOT NULL, description VARCHAR(500) NOT NULL DEFAULT '', homepage VARCHAR(500) NOT NULL DEFAULT '',
		icon VARCHAR(100) NOT NULL DEFAULT 'ri:puzzle-line', version VARCHAR(40) NOT NULL DEFAULT '1.0.0',
		author_name VARCHAR(100) NOT NULL DEFAULT '', download_url VARCHAR(500) NOT NULL DEFAULT '',
		enabled TINYINT(1) NOT NULL DEFAULT 1, published TINYINT(1) NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		KEY idx_published (published)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Plus后台托管插件目录'`); err != nil {
		return err
	}
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS plus_managed_templates (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		template_key VARCHAR(60) NOT NULL UNIQUE, name VARCHAR(100) NOT NULL, description VARCHAR(500) NOT NULL DEFAULT '',
		version VARCHAR(40) NOT NULL DEFAULT '1.0.0', preview_url VARCHAR(500) NOT NULL DEFAULT '', content_url VARCHAR(500) NOT NULL,
		sha256 CHAR(64) NOT NULL DEFAULT '', schema_version INT NOT NULL DEFAULT 1, author_name VARCHAR(100) NOT NULL DEFAULT '',
		published TINYINT(1) NOT NULL DEFAULT 1, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		KEY idx_published (published)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Plus后台托管首页模板目录'`)
	return err
}

func managedResourceDB() (*sql.DB, error) {
	db, err := openSystemConfigDB()
	if err != nil {
		return nil, err
	}
	if err = ensurePlusManagedResourceStorage(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
func validResourceURL(value string) bool {
	u, err := url.Parse(strings.TrimSpace(value))
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func AdminManagedPluginList(c *gin.Context) {
	db, err := managedResourceDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	defer db.Close()
	rows, err := db.Query(`SELECT id, plugin_id, category, name, description, homepage, icon, version, author_name, download_url, enabled, published FROM plus_managed_plugins ORDER BY id DESC`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "读取托管插件失败"})
		return
	}
	defer rows.Close()
	list := []gin.H{}
	for rows.Next() {
		var id int64
		var pid, cat, name, desc, home, icon, version, author, download string
		var enabled, published bool
		if rows.Scan(&id, &pid, &cat, &name, &desc, &home, &icon, &version, &author, &download, &enabled, &published) == nil {
			list = append(list, gin.H{"id": id, "pluginId": pid, "category": cat, "name": name, "description": desc, "homepage": home, "icon": icon, "version": version, "author": author, "downloadUrl": download, "enabled": enabled, "published": published})
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": gin.H{"list": list}})
}
func AdminManagedPluginSave(c *gin.Context) {
	var req struct {
		ID                                                                                  int64 `json:"id"`
		PluginID, Category, Name, Description, Homepage, Icon, Version, Author, DownloadURL string
		Published                                                                           bool `json:"published"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.PluginID) == "" || strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "插件标识和名称不能为空"})
		return
	}
	if req.DownloadURL != "" && !validResourceURL(req.DownloadURL) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "下载地址必须是有效的 HTTP/HTTPS 地址"})
		return
	}
	if req.Category == "" {
		req.Category = "other"
	}
	if req.Icon == "" {
		req.Icon = "ri:puzzle-line"
	}
	if req.Version == "" {
		req.Version = "1.0.0"
	}
	db, err := managedResourceDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	defer db.Close()
	var result sql.Result
	if req.ID > 0 {
		result, err = db.Exec(`UPDATE plus_managed_plugins SET plugin_id=?,category=?,name=?,description=?,homepage=?,icon=?,version=?,author_name=?,download_url=?,published=? WHERE id=?`, req.PluginID, req.Category, req.Name, req.Description, req.Homepage, req.Icon, req.Version, req.Author, req.DownloadURL, req.Published, req.ID)
	} else {
		result, err = db.Exec(`INSERT INTO plus_managed_plugins (plugin_id,category,name,description,homepage,icon,version,author_name,download_url,published) VALUES (?,?,?,?,?,?,?,?,?,?)`, req.PluginID, req.Category, req.Name, req.Description, req.Homepage, req.Icon, req.Version, req.Author, req.DownloadURL, req.Published)
	}
	_ = result
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "保存插件失败，插件标识可能已存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "插件已保存"})
}
func AdminManagedPluginDelete(c *gin.Context) {
	db, err := managedResourceDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	defer db.Close()
	if _, err = db.Exec("DELETE FROM plus_managed_plugins WHERE id=?", c.Param("id")); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "删除插件失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "插件已删除"})
}

func AdminManagedTemplateList(c *gin.Context) {
	db, err := managedResourceDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	defer db.Close()
	rows, err := db.Query(`SELECT id,template_key,name,description,version,preview_url,content_url,sha256,schema_version,author_name,published FROM plus_managed_templates ORDER BY id DESC`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "读取托管模板失败"})
		return
	}
	defer rows.Close()
	list := []gin.H{}
	for rows.Next() {
		var id int64
		var key, name, desc, ver, preview, content, sha, author string
		var schema int
		var published bool
		if rows.Scan(&id, &key, &name, &desc, &ver, &preview, &content, &sha, &schema, &author, &published) == nil {
			list = append(list, gin.H{"id": id, "templateId": key, "name": name, "description": desc, "version": ver, "previewUrl": preview, "contentUrl": content, "sha256": sha, "schemaVersion": schema, "author": author, "published": published})
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "", "data": gin.H{"list": list}})
}
func AdminManagedTemplateSave(c *gin.Context) {
	var req struct {
		ID                                                                             int64 `json:"id"`
		TemplateID, Name, Description, Version, PreviewURL, ContentURL, SHA256, Author string
		SchemaVersion                                                                  int  `json:"schemaVersion"`
		Published                                                                      bool `json:"published"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.TemplateID) == "" || strings.TrimSpace(req.Name) == "" || !validResourceURL(req.ContentURL) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "模板标识、名称和有效内容地址不能为空"})
		return
	}
	if req.Version == "" {
		req.Version = "1.0.0"
	}
	if req.SchemaVersion == 0 {
		req.SchemaVersion = 1
	}
	db, err := managedResourceDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	defer db.Close()
	if req.ID > 0 {
		_, err = db.Exec(`UPDATE plus_managed_templates SET template_key=?,name=?,description=?,version=?,preview_url=?,content_url=?,sha256=?,schema_version=?,author_name=?,published=? WHERE id=?`, req.TemplateID, req.Name, req.Description, req.Version, req.PreviewURL, req.ContentURL, req.SHA256, req.SchemaVersion, req.Author, req.Published, req.ID)
		_, _ = db.Exec(`UPDATE home_templates SET template_key=?,name=?,description=?,version=?,preview_url=?,template_url=?,sha256=?,schema_version=?,available=? WHERE id=(SELECT id FROM (SELECT id FROM home_templates WHERE template_key=? AND source_type='plus-managed' LIMIT 1) x)`, req.TemplateID, req.Name, req.Description, req.Version, req.PreviewURL, req.ContentURL, req.SHA256, req.SchemaVersion, req.Published, req.TemplateID)
	} else {
		_, err = db.Exec(`INSERT INTO plus_managed_templates (template_key,name,description,version,preview_url,content_url,sha256,schema_version,author_name,published) VALUES (?,?,?,?,?,?,?,?,?,?)`, req.TemplateID, req.Name, req.Description, req.Version, req.PreviewURL, req.ContentURL, req.SHA256, req.SchemaVersion, req.Author, req.Published)
		if err == nil {
			_, err = db.Exec(`INSERT INTO home_templates (catalog_id,template_key,source_id,name,description,version,source_url,source_type,preview_url,author_name,template_url,sha256,schema_version,available) VALUES (NULL,?,?,?, ?,?,'','plus-managed',?,?,?, ?,?,?)`, req.TemplateID, 0, req.TemplateID, req.Name, req.Description, req.Version, req.PreviewURL, req.Author, req.ContentURL, req.SHA256, req.SchemaVersion, req.Published)
		}
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "保存模板失败，模板标识可能已存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "模板已保存"})
}
func AdminManagedTemplateDelete(c *gin.Context) {
	db, err := managedResourceDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "数据库连接失败"})
		return
	}
	defer db.Close()
	if _, err = db.Exec("DELETE FROM plus_managed_templates WHERE id=?", c.Param("id")); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "删除模板失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "模板已删除"})
}
func AdminManagedResourceMenuInit(db *sql.DB) error { return ensurePlusManagedResourceStorage(db) }
