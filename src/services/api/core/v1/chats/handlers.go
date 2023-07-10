package chats

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/artbred/aliasflux/src/domain/models"
	"github.com/artbred/aliasflux/src/domain/providers/godaddy"
	"github.com/artbred/aliasflux/src/domain/providers/zion"
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/pkg/config"
	"github.com/artbred/aliasflux/src/pkg/storages/redisdb"
	"github.com/artbred/aliasflux/src/services/api/internal"
	"github.com/labstack/echo/v4"
	"github.com/twinj/uuid"
	"golang.org/x/time/rate"
	"net/http"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

// CreateChatHandler
// @Summary      Create chat
// @Tags         chats
// @Produce      json
// @Accept 		 json
// @Param        data body CreateChatRequest true "Chat"
// @Success      201  {object} CreateChatResponse
// @Router       /chats/create [post]
func (r *Router) CreateChatHandler(c echo.Context) error {
	req := &CreateChatRequest{}

	if err := internal.ValidateRequest(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: err.Error(),
		})
	}

	ctx := c.Request().Context()

	validationErr, err := models.ValidateSettings(ctx, req.ChatSettings)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to check settings are valid")
		return internal.InternalServerErrorResponse(c)
	}

	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: validationErr.Error(),
		})
	}

	user, err := models.GetUserByID(ctx, req.UserID)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to get user by id")
		return internal.InternalServerErrorResponse(c)
	}

	if user == nil || user.DeletedAt != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: "We couldn't find a user with this id",
		})
	}

	freeChatFeatures, err := models.GetFreeChatFeatures(ctx)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to get free chat features")
		return internal.InternalServerErrorResponse(c)
	}

	if freeChatFeatures == nil || freeChatFeatures.UserMessagesLimit == 0 {
		// TODO: create payment link
		return nil
	}

	chat := &models.Chat{
		ID:           uuid.NewV4().String(),
		UserID:       user.ID,
		Settings:     req.ChatSettings,
		Platform:     req.ChatSettings.Platform,
		ChatFeatures: *freeChatFeatures,
	}

	err = chat.Create(ctx)
	if err != nil {
		common.Logger.WithError(err).WithField("user_id", user.ID).Error("failed to create chat for user")
		return internal.InternalServerErrorResponse(c)
	}

	return c.JSON(http.StatusCreated, CreateChatResponse{
		ChatID: chat.ID,
		BaseResponse: internal.BaseResponse{
			Ok:      true,
			Message: "Chat created successfully!",
		},
	})
}

// GetChatHandler
// @Summary      Get chat
// @Tags         chats
// @Produce      json
// @Accept 		 json
// @Param        id path string true "Chat ID"
// @Param        offset query int true "Offset"
// @Success      200  {object} CreateChatResponse
// @Router       /chats/{id} [get]
func (r *Router) GetChatHandler(c echo.Context) error {
	req := &GetChatRequest{}

	if err := internal.ValidateRequest(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: err.Error(),
		})
	}

	ctx := c.Request().Context()

	chat, err := models.GetChat(ctx, req.ChatID)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to get chat")
		return internal.InternalServerErrorResponse(c)
	}

	err = chat.LoadMessages(ctx, true)

	if chat == nil || chat.DeletedAt != nil {
		return c.JSON(http.StatusBadRequest, internal.BaseResponse{
			Ok:      false,
			Message: "No chat found with this id",
		})
	}

	return c.JSON(http.StatusOK, GetChatResponse{
		Chat:   chat,
		Offset: req.Offset + len(chat.Messages),
		BaseResponse: internal.BaseResponse{
			Ok:      true,
			Message: "Chat retrieved successfully!",
		},
	})
}

// ListPlatformsHandler
// @Summary      Get available chat configurations
// @Tags         chats
// @Produce      json
// @Success      200  {object} CreateChatResponse
// @Router       /chats/platforms [get]
func (r *Router) ListPlatformsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	configs, err := models.ListAvailablePlatformSettings(ctx)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to list available chat settings")
		return internal.InternalServerErrorResponse(c)
	}

	if configs == nil {
		common.Logger.Errorf("no available chat settings found")
		return internal.InternalServerErrorResponse(c)
	}

	return c.JSON(http.StatusOK, ListChatConfigurationsResponse{
		Configurations: configs,
		BaseResponse: internal.BaseResponse{
			Ok:      true,
			Message: "Chat configurations retrieved successfully!",
		},
	})
}

func (r *Router) WebsocketHandler(c echo.Context) error {
	req := &WebsocketConnectRequest{}

	err := internal.ValidateRequest(c, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	chat, err := validateChat(ctx, req.ChatID)
	if err != nil {
		return err
	}

	rdb := redisdb.Connection()

	exists, err := rdb.SIsMember(ctx, redisdb.OpenChats.String(), chat.ID).Result()
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to check if chat is locked")
		return echo.NewHTTPError(http.StatusInternalServerError, "Please try again later")
	}

	if exists {
		return echo.NewHTTPError(http.StatusForbidden, "Chat is locked")
	} else {
		err = rdb.SAdd(ctx, redisdb.OpenChats.String(), chat.ID).Err()
		if err != nil {
			common.Logger.WithError(err).Errorf("failed to add chat %s to open chats", chat.ID)
			return echo.NewHTTPError(http.StatusInternalServerError, "Please try again later")
		}
	}

	ws, err := websocket.Accept(c.Response(), c.Request(), &websocket.AcceptOptions{
		InsecureSkipVerify: config.Debug,
	})

	if err != nil {
		common.Logger.WithError(err).Errorf("failed to accept websocket")
		return echo.NewHTTPError(http.StatusInternalServerError, "Please try again later")
	}

	isClosed := false

	defer func() {
		if !isClosed {
			errClose := ws.Close(websocket.StatusNormalClosure, "websocket closing")
			if errClose != nil {
				common.Logger.WithError(errClose).Error("failed to close websocket")
			}
		}

		err := rdb.SRem(context.Background(), redisdb.OpenChats.String(), chat.ID).Err()
		if err != nil {
			common.Logger.WithError(err).Errorf("failed to remove chat %s from open chats", chat.ID)
		}
	}()

	l := rate.NewLimiter(rate.Every(time.Second*1), 1)

	for {
		err := chatHandler(ctx, chat.ID, ws, l)
		status := websocket.CloseStatus(err)

		if status == websocket.StatusNormalClosure || status == websocket.StatusNoStatusRcvd || errors.As(err, &context.DeadlineExceeded) {
			isClosed = true
			break
		}

		if err != nil {
			common.Logger.WithError(err).Error("failed to handle websocket")
			break
		}
	}

	return nil
}

func validateChat(ctx context.Context, chatID string) (*models.Chat, error) {
	chat, err := models.GetChat(ctx, chatID)
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to get user with chat")
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Please try again later")
	}

	if chat == nil || chat.DeletedAt != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "No chat found with this id")
	}

	if chat.UserMessagesCount > chat.ChatFeatures.UserMessagesLimit {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "You have reached your message limit")
	}

	chat.LoadMessages(ctx, true)

	return chat, nil
}

type ChatResponseDomain struct {
	Message    string                `json:"message"`
	NameDomain []godaddy.NameDomains `json:"name_domain"`
}

func chatHandler(ctx context.Context, chatID string, conn *websocket.Conn, limiter *rate.Limiter) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute) // Close connection after 1 minute of inactivity
	defer cancel()

	err = limiter.Wait(ctx)
	if err != nil {
		return err
	}

	msgType, inputByte, err := conn.Read(ctx)
	if err != nil {
		return err
	}

	chat, err := validateChat(ctx, chatID)
	if err != nil {
		err = conn.Write(ctx, msgType, []byte(err.Error()))
		conn.CloseRead(ctx)
		return
	}

	usrMessageInput := string(inputByte)

	res, err := zion.NewCreateBusinessName(chat, usrMessageInput).Send(context.Background())
	if err != nil {
		common.Logger.WithError(err).Errorf("failed to create business name")
		err = conn.Write(ctx, msgType, []byte("Please try again later"))
		conn.CloseRead(ctx)
		return
	}

	msgUser := models.ChatMessage{
		ChatID:  chat.ID,
		Message: usrMessageInput,
		Sender:  models.SenderUser,
	}

	msgAiBytes, _ := json.Marshal(res)

	msgAI := models.ChatMessage{
		ChatID:  chat.ID,
		Message: string(msgAiBytes),
		Sender:  models.SenderAI,
	}

	err = msgUser.Create(ctx)
	err = msgAI.Create(ctx)

	var tlds string

	for i, tld := range chat.Settings.Tld {
		if i > 0 {
			tlds += ","
		}
		tlds += tld.Name
	}

	var byteMessage []byte

	if len(res.Names) > 0 {
		wg := &sync.WaitGroup{}
		nameDomainChan := make(chan *godaddy.NameDomains, len(res.Names))

		for _, name := range res.Names {
			wg.Add(1)
			go func(businessName string) {
				defer wg.Done()

				nameDomain, err := godaddy.NewClient().SuggestDomains(businessName, tlds)
				if err != nil {
					nameDomainChan <- nil
					common.Logger.WithError(err).Errorf("failed to get domain suggestion")
				}

				nameDomainChan <- nameDomain
			}(name)
		}

		wg.Wait()
		close(nameDomainChan)

		resChat := ChatResponseDomain{
			Message: res.Message,
		}

		for name := range nameDomainChan {
			if name != nil {
				resChat.NameDomain = append(resChat.NameDomain, *name)
			}
		}

		byteMessage, err = json.Marshal(resChat)
	} else {
		byteMessage, err = json.Marshal(res)
	}

	if err != nil {
		common.Logger.WithError(err).Errorf("can't create messages for chat")
		err = conn.Write(ctx, msgType, []byte("Please try again later"))
		conn.CloseRead(ctx)
		return
	}

	return conn.Write(ctx, msgType, byteMessage)
}
