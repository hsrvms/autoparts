package handlers

import (
	"net/http"
	"strconv"
	"time"

	salesmodels "github.com/hsrvms/autoparts/internal/modules/sales/models"
	"github.com/hsrvms/autoparts/internal/modules/sales/services"
	"github.com/labstack/echo/v4"
)

type SaleHandler struct {
	service services.SaleService
}

func NewSaleHandler(service services.SaleService) *SaleHandler {
	return &SaleHandler{
		service: service,
	}
}

// GetSales handles retrieval of all sales with optional filtering
func (h *SaleHandler) GetSales(c echo.Context) error {
	filter := &salesmodels.SaleFilter{}

	// Parse query parameters
	if itemID := c.QueryParam("item_id"); itemID != "" {
		id, err := strconv.Atoi(itemID)
		if err == nil {
			filter.ItemID = &id
		}
	}

	if startDate := c.QueryParam("start_date"); startDate != "" {
		if date, err := time.Parse(time.RFC3339, startDate); err == nil {
			filter.StartDate = &date
		}
	}

	if endDate := c.QueryParam("end_date"); endDate != "" {
		if date, err := time.Parse(time.RFC3339, endDate); err == nil {
			filter.EndDate = &date
		}
	}

	if customerName := c.QueryParam("customer_name"); customerName != "" {
		filter.CustomerName = &customerName
	}

	if customerPhone := c.QueryParam("customer_phone"); customerPhone != "" {
		filter.CustomerPhone = &customerPhone
	}

	if customerEmail := c.QueryParam("customer_email"); customerEmail != "" {
		filter.CustomerEmail = &customerEmail
	}

	if transactionNumber := c.QueryParam("transaction_number"); transactionNumber != "" {
		filter.TransactionNumber = &transactionNumber
	}

	if soldBy := c.QueryParam("sold_by"); soldBy != "" {
		filter.SoldBy = &soldBy
	}

	ctx := c.Request().Context()
	sales, err := h.service.GetAll(ctx, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, sales)
}

// GetSaleByID handles retrieval of a single sale
func (h *SaleHandler) GetSaleByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid sale ID")
	}

	ctx := c.Request().Context()
	sale, err := h.service.GetByID(ctx, id)
	if err != nil {
		switch err {
		case services.ErrSaleNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, sale)
}

// CreateSale handles creation of a new sale
func (h *SaleHandler) CreateSale(c echo.Context) error {
	sale := new(salesmodels.Sale)
	if err := c.Bind(sale); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	id, err := h.service.Create(ctx, sale)
	if err != nil {
		switch err {
		case services.ErrInvalidItemID, services.ErrInvalidQuantity,
			services.ErrInvalidPricePerUnit, services.ErrInvalidDate,
			services.ErrInvalidCustomerEmail:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case services.ErrDuplicateTransactionNumber:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case services.ErrInsufficientStock:
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	sale.SaleID = id
	return c.JSON(http.StatusCreated, sale)
}

// UpdateSale handles updating an existing sale
func (h *SaleHandler) UpdateSale(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid sale ID")
	}

	sale := new(salesmodels.Sale)
	if err := c.Bind(sale); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	sale.SaleID = id

	ctx := c.Request().Context()
	err = h.service.Update(ctx, sale)
	if err != nil {
		switch err {
		case services.ErrSaleNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case services.ErrInvalidItemID, services.ErrInvalidQuantity,
			services.ErrInvalidPricePerUnit, services.ErrInvalidDate,
			services.ErrInvalidCustomerEmail:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case services.ErrDuplicateTransactionNumber:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case services.ErrInsufficientStock:
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, sale)
}

// DeleteSale handles deletion of a sale
func (h *SaleHandler) DeleteSale(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid sale ID")
	}

	ctx := c.Request().Context()
	err = h.service.Delete(ctx, id)
	if err != nil {
		switch err {
		case services.ErrSaleNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// GetByTransactionNumber handles retrieval of a sale by transaction number
func (h *SaleHandler) GetByTransactionNumber(c echo.Context) error {
	transactionNumber := c.Param("transactionNumber")
	if transactionNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "transaction number is required")
	}

	ctx := c.Request().Context()
	sale, err := h.service.GetByTransactionNumber(ctx, transactionNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if sale == nil {
		return echo.NewHTTPError(http.StatusNotFound, "sale not found")
	}

	return c.JSON(http.StatusOK, sale)
}

// GetCustomerSales handles retrieval of all sales for a customer
func (h *SaleHandler) GetCustomerSales(c echo.Context) error {
	customerEmail := c.Param("customerEmail")
	if customerEmail == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "customer email is required")
	}

	ctx := c.Request().Context()
	sales, err := h.service.GetCustomerSales(ctx, customerEmail)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, sales)
}
