package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleRequest(c *gin.Context, f func(c *gin.Context) *Response) {
	ctx := c.Request.Context()
	if _, ok := ctx.Deadline(); !ok {
		// если у запроса нет дедлайна, то его обработка происходит синхронно,
		// без дополнительных накладных ресурсов
		handleRequestReal(c, f(c))
		return
	}

	doneChan := make(chan *Response)
	go func() {
		doneChan <- f(c)
	}()
	select {
	case <-ctx.Done():
		// Nothing to do because err handled from timeout middleware
	case res := <-doneChan:
		handleRequestReal(c, res)
	}

	// ctx.Done() возвращает канал, который закрывается, когда контекст завершается
	// <-ctx.Done() - чтение из канала, которое завершится, когда канал будет закрыт,
	// сигнализирует о завершении контекста
	// case <-ctx.Done() - сработает, если дедлайн истек или если контекст был отменен вручную

	// case res := <-doneChan - сработает, когда горутина завершает выполнение f(c) и отправляет результат в doneChan
	// раньше, чем контекст завершится
}

func handleRequestReal(c *gin.Context, res *Response) {
	if res.Err == nil {
		statusCode := res.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}
		if res.Data != nil {
			c.JSON(res.StatusCode, res.Data)
		} else {
			c.Status(res.StatusCode)
		}
		return
	}

	var err *ErrorResponse
	err, ok := res.Err.(*ErrorResponse)
	if !ok {
		res.StatusCode = http.StatusInternalServerError
		err = &ErrorResponse{Detail: "An error has occurred, please try again later"}
	}
	c.AbortWithStatusJSON(res.StatusCode, err)
}
