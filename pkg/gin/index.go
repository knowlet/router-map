package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 *
 * @api {get} / Check System is alive
 * @apiName GetIndex
 * @apiGroup Index
 * @apiVersion  v1.0.0
 *
 *
 * @apiSuccess (200) {String} ok System Online.
 *
 * @apiSuccessExample {String} Success-Response:
 * ok
 *
 *
 */
func GetIndexHandler(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
