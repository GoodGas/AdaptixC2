package connector

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type DownloadFileId struct {
	File string `json:"file_id"`
}

func (tc *TsConnector) TcGuiDownloadSync(ctx *gin.Context) {
	var (
		downloadFid DownloadFileId
		answer      gin.H
		err         error
	)

	err = ctx.ShouldBindJSON(&downloadFid)
	if err != nil {
		_ = ctx.Error(errors.New("invalid action"))
		return
	}

	filename, content, err := tc.teamserver.TsDownloadSync(downloadFid.File)
	if err != nil {
		answer = gin.H{"message": err.Error(), "ok": false}
	} else {
		encodedContent := base64.StdEncoding.EncodeToString(content)
		filename = strings.Split(filename, "_")[1]

		answer = gin.H{"ok": true, "filename": filename, "content": encodedContent}
	}

	ctx.JSON(http.StatusOK, answer)
}

func (tc *TsConnector) TcGuiDownloadDelete(ctx *gin.Context) {
	var (
		downloadFid DownloadFileId
		answer      gin.H
		err         error
	)

	err = ctx.ShouldBindJSON(&downloadFid)
	if err != nil {
		_ = ctx.Error(errors.New("invalid action"))
		return
	}

	err = tc.teamserver.TsDownloadDelete(downloadFid.File)
	if err != nil {
		answer = gin.H{"message": err.Error(), "ok": false}
	} else {
		answer = gin.H{"message": "file delete", "ok": true}
	}

	ctx.JSON(http.StatusOK, answer)
}
