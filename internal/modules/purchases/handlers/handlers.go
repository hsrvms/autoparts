// internal/modules/purchases/handlers/handlers.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	purchasemodels "github.com/hsrvms/autoparts/internal/modules/purchases/models"
	"github.com/hsrvms/autoparts/internal/modules/purchases/services"
	"github.com/labstack/echo/v4"
)

type PurchaseHandler struct {
    service services.PurchaseService
}

func NewPurchaseHandler(service services.PurchaseService) *PurchaseHandler {
    return &PurchaseHandler{
        service: service,
    }
}

// GetPurchases handles retrieval of all purchases with optional filtering
func (h *PurchaseHandler) GetPurchases(c echo.Context) error {
    filter := &purchasemodels.PurchaseFilter{}

    // Parse query parameters
    if supplierID := c.QueryParam("supplier_id"); supplierID != "" {
        id, err := strconv.Atoi(supplierID)
        if err == nil {
            filter.SupplierID = &id
        }
    }

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

    if invoiceNumber := c.QueryParam("invoice_number"); invoiceNumber != "" {
        filter.InvoiceNumber = &invoiceNumber
    }

    ctx := c.Request().Context()
    purchases, err := h.service.GetAll(ctx, filter)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, purchases)
}

// GetPurchaseByID handles retrieval of a single purchase
func (h *PurchaseHandler) GetPurchaseByID(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid purchase ID")
    }

    ctx := c.Request().Context()
    purchase, err := h.service.GetByID(ctx, id)
    if err != nil {
        switch err {
        case services.ErrPurchaseNotFound:
            return echo.NewHTTPError(http.StatusNotFound, err.Error())
        default:
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
    }

    return c.JSON(http.StatusOK, purchase)
}

// CreatePurchase handles creation of a new purchase
func (h *PurchaseHandler) CreatePurchase(c echo.Context) error {
    purchase := new(purchasemodels.Purchase)
    if err := c.Bind(purchase); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    ctx := c.Request().Context()
    id, err := h.service.Create(ctx, purchase)
    if err != nil {
        switch err {
        case services.ErrInvalidSupplierID, services.ErrInvalidItemID,
             services.ErrInvalidQuantity, services.ErrInvalidCostPerUnit,
             services.ErrInvalidDate:
            return echo.NewHTTPError(http.StatusBadRequest, err.Error())
        case services.ErrDuplicateInvoiceNumber:
            return echo.NewHTTPError(http.StatusConflict, err.Error())
        default:
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
    }

    purchase.PurchaseID = id
    return c.JSON(http.StatusCreated, purchase)
}

// UpdatePurchase handles updating an existing purchase
func (h *PurchaseHandler) UpdatePurchase(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid purchase ID")
    }

    purchase := new(purchasemodels.Purchase)
    if err := c.Bind(purchase); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    purchase.PurchaseID = id

    ctx := c.Request().Context()
    err = h.service.Update(ctx, purchase)
    if err != nil {
        switch err {
        case services.ErrPurchaseNotFound:
            return echo.NewHTTPError(http.StatusNotFound, err.Error())
        case services.ErrInvalidSupplierID, services.ErrInvalidItemID,
             services.ErrInvalidQuantity, services.ErrInvalidCostPerUnit,
             services.ErrInvalidDate:
            return echo.NewHTTPError(http.StatusBadRequest, err.Error())
        case services.ErrDuplicateInvoiceNumber:
            return echo.NewHTTPError(http.StatusConflict, err.Error())
        default:
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
    }

    return c.JSON(http.StatusOK, purchase)
}

// DeletePurchase handles deletion of a purchase
func (h *PurchaseHandler) DeletePurchase(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid purchase ID")
    }

    ctx := c.Request().Context()
    err = h.service.Delete(ctx, id)
    if err != nil {
        switch err {
        case services.ErrPurchaseNotFound:
            return echo.NewHTTPError(http.StatusNotFound, err.Error())
        default:
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
    }

    return c.NoContent(http.StatusNoContent)
}

// GetSupplierPurchases handles retrieval of all purchases for a supplier
func (h *PurchaseHandler) GetSupplierPurchases(c echo.Context) error {
    supplierID, err := strconv.Atoi(c.Param("supplierId"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid supplier ID")
    }

    ctx := c.Request().Context()
    purchases, err := h.service.GetSupplierPurchases(ctx, supplierID)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, purchases)
}

// GetItemPurchases handles retrieval of all purchases for an item
func (h *PurchaseHandler) GetItemPurchases(c echo.Context) error {
    itemID, err := strconv.Atoi(c.Param("itemId"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid item ID")
    }

    ctx := c.Request().Context()
    purchases, err := h.service.GetItemPurchases(ctx, itemID)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, purchases)
}
