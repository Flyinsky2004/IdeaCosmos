package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIndexCoverList(c *gin.Context) {
	pageIndex := c.Query("pageIndex")
	pageIndexInt, _ := strconv.Atoi(pageIndex)
	var results []pojo.Project
	err := config.MysqlDataBase.Preload("Team").Limit(10).Offset(pageIndexInt).Find(&results).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "查询数据库时发生错误"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(results))
}
