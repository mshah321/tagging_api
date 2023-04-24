package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type NearByTagsController struct {
	google_tags []string
	key         string
}

func (handle *NearByTagsController) get_nearby_tags(ctx *gin.Context) {
	lat, err := strconv.ParseFloat(ctx.Query("lat"), 64)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "Latitude must be a number",
		})
		return
	}

	lng, err := strconv.ParseFloat(ctx.Query("lng"), 64)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "Longitude must be a number",
		})
		return
	}

	radius, err := strconv.ParseInt(ctx.Query("radius"), 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "Radius must be a number",
		})
		return
	}

	ok := validate_coordinates(lat, lng)
	if !ok {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "Invalid coordinates",
		})
		return
	}

	response, err := get_nearby_tags_count(handle.google_tags, lat, lng, int(radius), handle.key)

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	ctx.JSON(200, response)
}
