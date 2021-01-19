package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func readCSV(fn string) (map[string]float64, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	ret := map[string]float64{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		log.Println(record)
		if record[0] == "Datum" {
			continue
		}
		fl, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		ret[record[0]] = fl

	}
	return ret, nil
}

func getCurHandler(cur string) func(*gin.Context) {
	di, err := readCSV(fmt.Sprintf("2020%s.csv", cur))
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		dat := c.Param("dat")
		v, ok := di[dat]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "date not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{cur: fmt.Sprintf("%.3f", v)})
	}
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	r.GET("/usd/:dat", getCurHandler("usd"))
	r.GET("/eur/:dat", getCurHandler("eur"))
	r.Run("localhost:1444")
}
