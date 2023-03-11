package main

import (
	"github.com/gin-gonic/gin"
	// "net/http"
	"fmt"
	"os"

	// "io/ioutil"
	"os/exec"
	"strings"
)

const GOFOLDER = "go_scripts"

func main() {

	router := gin.Default()
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("file")

		if err := os.Mkdir(GOFOLDER, 0755); err != nil {
			// 文件夹已存在
			if !os.IsExist(err) {
				fmt.Println("go脚本文件夹创建失败!")
				return
			}
		}

		dst := "./" + GOFOLDER + "/" + file.Filename

		fmt.Println(dst)
		// 上传文件至指定的完整文件路径
		c.SaveUploadedFile(file, dst)

		// 进行判断上传的是否是go文件还是go编译之后的文件，后续改良需要从context中读取文件类型
		if strings.HasSuffix(file.Filename, ".go") {

			goBuildFileName := strings.Replace(file.Filename, ".go", "", -1)

			cmd := exec.Command("go", "build", "-o", fmt.Sprintf("./go_build/%s", goBuildFileName), dst)
			output, err := cmd.Output()
			if err != nil {
				fmt.Println("执行失败", err)
			}
			fmt.Println(string(output))
		}

	})
	router.Run(":8080")
}
