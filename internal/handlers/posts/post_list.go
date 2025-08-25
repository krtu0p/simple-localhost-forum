package posts

import (

    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

func (h *Handler) GetAllPost(c *gin.Context) {
    ctx := c.Request.Context()

    pageIndexStr := c.Query("pageIndex")
    pageSizeStr := c.Query("pageSize")

    pageIndex, err := strconv.Atoi(pageIndexStr)
    if err != nil || pageIndex < 1 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid page index",
        })
        return
    }

    pageSize, err := strconv.Atoi(pageSizeStr)
    if err != nil || pageSize <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid page size",
        })
        return
    }

    response, err := h.postSvc.GetAllPost(ctx, pageSize, pageIndex)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch posts",
        })
        return
    }

    c.JSON(http.StatusOK, response)
}