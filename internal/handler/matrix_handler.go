package handler

import (
	"github.com/gofiber/fiber/v2"
	"go-api/internal/client"
)

type MatrixHandler struct {
	nodeClient *client.NodeClient
}

func NewMatrixHandler(nodeClient *client.NodeClient) *MatrixHandler {
	return &MatrixHandler{nodeClient: nodeClient}
}

type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

func (h *MatrixHandler) ProcessMatrix(c *fiber.Ctx) error {
	var req MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido, se esperaba: { \"matrix\": [[...]] }",
		})
	}

	if len(req.Matrix) == 0 || len(req.Matrix[0]) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "la matriz no puede estar vacía",
		})
	}

	cols := len(req.Matrix[0])
	for _, row := range req.Matrix {
		if len(row) != cols {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "la matriz debe ser rectangular",
			})
		}
	}

	// Placeholder para Householder QR
	placeholderQR := fiber.Map{
		"q": req.Matrix,
		"r": req.Matrix,
	}

	stats, err := h.nodeClient.SendStats(placeholderQR)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "falla al comunicar con node-api",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"qr":    placeholderQR,
		"stats": stats,
	})
}