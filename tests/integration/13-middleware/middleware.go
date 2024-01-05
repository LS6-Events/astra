package petstore

import "github.com/gin-gonic/gin"

func petsMiddleware(c *gin.Context) {
	apiKey, ok := c.GetQuery("api_key")
	if !ok {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "api_key is missing",
		})
		return
	}

	if apiKey != "1234567890" {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "invalid api_key",
		})
		return
	}

	c.Next()
}

func handlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, ok := c.GetQuery("inline_api_key")
		if !ok {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "api_key is missing",
			})
			return
		}

		if apiKey != "1234567890" {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "invalid api_key",
			})
			return
		}

		c.Next()
	}
}

func headerMiddleware(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Authorization header is missing",
		})
		return
	}

	if authorization != "Bearer 1234567890" {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "invalid Authorization header",
		})
		return
	}

	c.Next()
}
