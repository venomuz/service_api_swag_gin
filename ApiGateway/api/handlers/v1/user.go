package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/venomuz/service_api_swag_gin/ApiGateway/api/model"
	pb "github.com/venomuz/service_api_swag_gin/ApiGateway/genproto"
	l "github.com/venomuz/service_api_swag_gin/ApiGateway/pkg/logger"
	mail "github.com/venomuz/service_api_swag_gin/ApiGateway/pkg/mail"
	"google.golang.org/protobuf/encoding/protojson"
	"math/rand"
	"net/http"
	"time"
)

// CreateUser creates users
// @Summary      Create an account
// @Description  This api is for creating user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 		 user  body model.Useri true "user body"
// @Success      200  {string}  Ok
// @Router       /v1/users [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	body := pb.Useri{}
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().Create(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser gets user by id
// @Summary      Get an account
// @Description  This api is for getting user
// @Tags         user
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  model.Useri
// @Router       /v1/users/{id} [get]
func (h *handlerV1) GetUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetByID(
		ctx, &pb.GetIdFromUser{Id: guid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser delete user by id
// @Summary      Delete an account
// @Description  This api is for delete user
// @Tags         user
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  model.Id
// @Router       /v1/users/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	fmt.Println(guid)
	_, err := h.serviceManager.UserService().DeleteByID(
		ctx, &pb.GetIdFromUserID{Id: guid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}
	_ = h.redisStorage.Set("123qwe", "123")
	res, _ := h.redisStorage.Get("qwe")
	c.JSON(http.StatusOK, res)
}

// CheckReg CheckUserAnd creates users
// @Summary      Create an account with check
// @Description  This api is for creating user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 		 user  body model.Useri true "user body"
// @Success      200  {string}  Ok
// @Router       /v1/users/check [post]
func (h *handlerV1) CheckReg(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	body := pb.Useri{}
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().CheckLoginMail(ctx, &pb.Check{Key: "login", Value: body.Login})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}
	num := rand.Intn(999999)
	if response.Status == false {
		err := mail.SendMail(int32(num), body.Email[0])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to send user", l.Error(err))
			return
		}
	}
	h.redisStorage.Set(string(num))
	c.JSON(http.StatusCreated, response)
}
