package db

import (
	"fmt"
	"time"

	"github.com/solozyx/object-storage/db/db_mysql"
)

// UserFile : 用户文件表结构体
type UserFile struct {
	UserName    string
	FileHash    string // 文件唯一标识
	FileName    string
	FileSize    int64
	UploadAt    string // 文件上传时间
	LastUpdated string // 文件最后修改时间戳
}

// 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, err := db_mysql.DBConn().Prepare(`insert ignore into tbl_user_file 
 				(user_name,file_sha1,file_name,file_size,upload_at) values (?,?,?,?,?)`)
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		return false
	}
	return true
}

// 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := db_mysql.DBConn().Prepare(`select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name=? limit ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}
	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			// Scan失败直接跳出循环 认为数据不可用
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}
