package api

import (
	"context"
	"encoding/json"

	usecase "github.com/arsura/gourney/cmd/usecases"
	"github.com/arsura/gourney/pkg/constant"
	model "github.com/arsura/gourney/pkg/models/mongodb"
	util "github.com/arsura/gourney/pkg/utils"
	validator "github.com/arsura/gourney/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CreatePostReqBody struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type UpdatePostReqBody struct {
	Id      string `json:"id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type PostHandlerProvider interface {
	FindPostByIdHandler(c *fiber.Ctx) error
	CreatePostHandler(c *fiber.Ctx) error
	UpdatePostByIdHandler(c *fiber.Ctx) error
	DeletePostByIdHandler(c *fiber.Ctx) error
}

type postHandler struct {
	postUsecase usecase.PostUsecaseProvider
	validator   *validator.Validator
	logger      *zap.SugaredLogger
}

func NewPostHandler(postUsecase usecase.PostUsecaseProvider, validator *validator.Validator, logger *zap.SugaredLogger) *postHandler {
	return &postHandler{postUsecase, validator, logger}
}

func (h *postHandler) FindPostByIdHandler(c *fiber.Ctx) error {
	var (
		requestId = c.Locals("requestid").(string)
		ctx       = context.WithValue(c.UserContext(), constant.RequestIdKey, requestId)
	)

	id, err := util.StringToObjectId(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id must be an object id",
		})
	}

	result, err := h.postUsecase.FindPostById(ctx, *id)
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
		ctx       = context.WithValue(c.UserContext(), constant.RequestIdKey, requestId)
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

	_, err = h.postUsecase.CreatePost(ctx, &model.Post{
		Title:   post.Title,
		Content: post.Content,
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
		ctx       = context.WithValue(c.UserContext(), constant.RequestIdKey, requestId)
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

	_, err = h.postUsecase.UpdatePostById(ctx, *id, &model.Post{
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
		ctx       = context.WithValue(c.UserContext(), constant.RequestIdKey, requestId)
	)

	id, err := util.StringToObjectId(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "id must be an object id",
		})
	}

	_, err = h.postUsecase.DeletePostById(ctx, *id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": "post not found",
		})

	}
	return c.SendStatus(fiber.StatusNoContent)
}
