package handlers

import (
	"net/http"
	"strconv"

	suppliermodels "github.com/hsrvms/autoparts/internal/modules/suppliers/models"
	"github.com/hsrvms/autoparts/internal/modules/suppliers/services"
	"github.com/labstack/echo/v4"
)

type SupplierHandler struct {
    service services.SupplierService
}

func NewSupplierHandler(service services.SupplierService) *SupplierHandler {
    return &SupplierHandler{
        service: service,
    }
}

// GetSuppliers handles retrieval of all suppliers with optional filtering
func (h *SupplierHandler) GetSuppliers(c echo.Context) error {
    filter := &suppliermodels.SupplierFilter{}

    // Parse query parameters
    if search := c.QueryParam("search"); search != "" {
        filter.SearchTerm = &search
    }

    if hasActiveItems := c.QueryParam("has_active_items"); hasActiveItems == "true" {
        active := true
        filter.HasActiveItems = &active
    }

    ctx := c.Request().Context()
    suppliers, err := h.service.GetAll(ctx, filter)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, suppliers)
}

// GetSupplierByID handles retrieval of a single supplier
func (h *SupplierHandler) GetSupplierByID(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid supplier ID")
    }

    ctx := c.Request().Context()
    supplier, err := h.service.GetByID(ctx, id)
    if err != nil {
        switch err {
        case services.ErrSupplierNotFound:
            return echo.NewHTTPError(http.StatusNotFound, err.Error())
        default:
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
    }

    return c.JSON(http.StatusOK, supplier)
}

// CreateSupplier handles creation of a new supplier
func (h *SupplierHandler) CreateSupplier(c echo.Context) error {
    supplier := new(suppliermodels.Supplier)
    if err := c.Bind(supplier); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    ctx := c.Request().Context()
    id, err := h.service.Create(ctx, supplier)
    if err != nil {
        switch err {
        case services.ErrDuplicateSupplierName:
            return echo.NewHTTPError(http.StatusConflict, err.Error())
        default:
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
    }

    supplier.SupplierID = id
    return c.JSON(http.StatusCreated, supplier)
}

// UpdateSupplier handles updating an existing supplier
func (h *SupplierHandler) UpdateSupplier(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid supplier ID")
    }

    supplier := new(suppliermodels.Supplier)
    if err := c.Bind(supplier); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    supplier.SupplierID = id

    ctx := c.Request().Context()
    err = h.service.Update(ctx, supplier)
    if err != nil {
        switch err {
        case services.ErrSupplierNotFound:
            return echo.NewHTTPError(http.StatusNotFound, err.Error())
        case services.ErrDuplicateSupplierName:
            return echo.NewHTTPError(http.StatusConflict, err.Error())
        default:
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
    }

    return c.JSON(http.StatusOK, supplier)
}

// DeleteSupplier handles deletion of a supplier
func (h *SupplierHandler) DeleteSupplier(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid supplier ID")
    }

    ctx := c.Request().Context()
    err = h.service.Delete(ctx, id)
    if err != nil {
        switch err {
        case services.ErrSupplierNotFound:
            return echo.NewHTTPError(http.StatusNotFound, err.Error())
        case services.ErrSupplierHasItems:
            return echo.NewHTTPError(http.StatusConflict, err.Error())
        default:
            return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
        }
    }

    return c.NoContent(http.StatusNoContent)
}
