package api

import (
	"context"
	"encoding/json"

	"github.com/arsura/gourney/pkg/constant"
	model "github.com/arsura/gourney/pkg/models"
	usecase "github.com/arsura/gourney/pkg/usecases"
	util "github.com/arsura/gourney/pkg/utils"
	validator "github.com/arsura/gourney/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CreatePostReqBody struct {
	Title             string                      `json:"title" validate:"required"`
	Content           string                      `json:"content" validate:"required"`
	SocialNetworkType model.PostSocialNetworkType `json:"social_network_type" validate:"oneof=facebook twitter"`
}

type UpdatePostReqBody struct {
	Id                    string                      `json:"id" validate:"required"`
	Title                 string                      `json:"title" validate:"required"`
	Content               string                      `json:"content" validate:"required"`
	PostSocialNetworkType model.PostSocialNetworkType `json:"social_network_type" validate:"required"`
}

type PostHandlerProvider interface {
	FindPostByIdHandler(c *fiber.Ctx) error
	CreatePostHandler(c *fiber.Ctx) error
	UpdatePostByIdHandler(c *fiber.Ctx) error
	DeletePostByIdHandler(c *fiber.Ctx) error
	CountPostBySocialNetworkTypeHandler(c *fiber.Ctx) error
}

type postHandler struct {
	postUseCase usecase.PostUseCaseProvider
	validator   *validator.Validator
	logger      *zap.SugaredLogger
}

func NewPostHandler(postUseCase usecase.PostUseCaseProvider, validator *validator.Validator, logger *zap.SugaredLogger) *postHandler {
	return &postHandler{postUseCase, validator, logger}
}

func (h *postHandler) FindPostByIdHandler(c *fiber.Ctx) error {
	var (
		requestId = c.Locals("requestid").(string)
		ctx       = context.WithValue(c.UserContext(), constant.REQUEST_ID_KEY, requestId)
	)

	id, err := util.StringToObjectId(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id must be an object id",
		})
	}

	result, err := h.postUseCase.FindPostById(ctx, *id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": "post not found",
		})

	}
	return c.JSON(result)
}

func (h *postHandler) CreatePostHandler(c *fiber.Ctx) error {
	var (
		post      = &CreatePostReqBody{}
		requestId = c.Locals("requestid").(string)
		ctx       = context.WithValue(c.UserContext(), constant.REQUEST_ID_KEY, requestId)
	)

	err := json.Unmarshal(c.Body(), post)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": util.UnmarshalErrorParser(err),
		})
	}

	if err := h.validator.Validate.Struct(*post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": h.validator.TransError(err),
		})
	}

	_, err = h.postUseCase.CreatePost(ctx, &model.Post{
		Title:             post.Title,
		Content:           post.Content,
		SocialNetworkType: post.SocialNetworkType,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "failed to create post",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *postHandler) UpdatePostByIdHandler(c *fiber.Ctx) error {
	var (
		post      = &UpdatePostReqBody{}
		requestId = c.Locals("requestid").(string)
		ctx       = context.WithValue(c.UserContext(), constant.REQUEST_ID_KEY, requestId)
	)

	err := json.Unmarshal(c.Body(), post)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": util.UnmarshalErrorParser(err),
		})
	}

	if err := h.validator.Validate.Struct(*post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": h.validator.TransError(err),
		})
	}

	id, err := util.StringToObjectId(post.Id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id must be an object id",
		})
	}

	_, err = h.postUseCase.UpdatePostById(ctx, *id, &model.Post{
		Title:   post.Title,
		Content: post.Content,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "failed to update post",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *postHandler) DeletePostByIdHandler(c *fiber.Ctx) error {
	var (
		requestId = c.Locals("requestid").(string)
		ctx       = context.WithValue(c.UserContext(), constant.REQUEST_ID_KEY, requestId)
	)

	id, err := util.StringToObjectId(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id must be an object id",
		})
	}

	_, err = h.postUseCase.DeletePostById(ctx, *id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": "post not found",
		})

	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *postHandler) CountPostBySocialNetworkTypeHandler(c *fiber.Ctx) error {
	var (
		requestId = c.Locals("requestid").(string)
		ctx       = context.WithValue(c.UserContext(), constant.REQUEST_ID_KEY, requestId)
	)

	result, err := h.postUseCase.CountPostBySocialNetworkType(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})

	}
	return c.JSON(result)
}
