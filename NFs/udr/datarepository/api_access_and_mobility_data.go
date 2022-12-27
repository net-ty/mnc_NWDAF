/*
 * Nudr_DataRepository API OpenAPI file
 *
 * Unified Data Repository Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package datarepository

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateAccessAndMobilityData - Creates and updates the access and mobility exposure data for a UE
func CreateAccessAndMobilityData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteAccessAndMobilityData - Deletes the access and mobility exposure data for a UE
func DeleteAccessAndMobilityData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// QueryAccessAndMobilityData - Retrieves the access and mobility exposure data for a UE
func QueryAccessAndMobilityData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
