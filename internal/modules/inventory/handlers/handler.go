package handlers

import (
	"net/http"
	"strconv"

	inventorymodels "github.com/hsrvms/autoparts/internal/modules/inventory/models"
	"github.com/hsrvms/autoparts/internal/modules/inventory/services"
	"github.com/labstack/echo/v4"
)

type InventoryHandler struct {
	service        services.InventoryService
	barcodeService services.BarcodeService
}

func NewInventoryHandler(service services.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		service:        service,
		barcodeService: services.NewBarcodeService(),
	}
}

// GetItems handles the retrieval of items with optional filtering
func (h *InventoryHandler) GetItems(c echo.Context) error {
	filter := &inventorymodels.ItemFilter{}

	// Parse query parameters
	if categoryID := c.QueryParam("category_id"); categoryID != "" {
		id, err := strconv.Atoi(categoryID)
		if err == nil {
			filter.CategoryID = &id
		}
	}

	if supplierID := c.QueryParam("supplier_id"); supplierID != "" {
		id, err := strconv.Atoi(supplierID)
		if err == nil {
			filter.SupplierID = &id
		}
	}

	if partNumber := c.QueryParam("part_number"); partNumber != "" {
		filter.PartNumber = &partNumber
	}

	if search := c.QueryParam("search"); search != "" {
		filter.SearchTerm = &search
	}

	if lowStock := c.QueryParam("low_stock"); lowStock == "true" {
		isLowStock := true
		filter.LowStock = &isLowStock
	}

	if isActive := c.QueryParam("is_active"); isActive != "" {
		active := isActive == "true"
		filter.IsActive = &active
	}

	ctx := c.Request().Context()
	items, err := h.service.GetItems(ctx, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, items)
}

// GetLowStockItems handles the retrieval of items with low stock
func (h *InventoryHandler) GetLowStockItems(c echo.Context) error {
	ctx := c.Request().Context()
	items, err := h.service.GetLowStockItems(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, items)
}

// GetItemByID handles the retrieval of a single item by ID
func (h *InventoryHandler) GetItemByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid item ID")
	}

	ctx := c.Request().Context()
	item, err := h.service.GetItemByID(ctx, id)
	if err != nil {
		if err == services.ErrItemNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, item)
}

// GetItemByBarcode handles the retrieval of a single item by barcode
func (h *InventoryHandler) GetItemByBarcode(c echo.Context) error {
	barcode := c.Param("barcode")
	if barcode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "barcode is required")
	}

	ctx := c.Request().Context()
	item, err := h.service.GetItemByBarcode(ctx, barcode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if item == nil {
		return echo.NewHTTPError(http.StatusNotFound, "item not found")
	}

	return c.JSON(http.StatusOK, item)
}

// CreateItem handles the creation of a new item
func (h *InventoryHandler) CreateItem(c echo.Context) error {
	item := new(inventorymodels.Item)
	if err := c.Bind(item); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	id, err := h.service.CreateItem(ctx, item)
	if err != nil {
		switch err {
		case services.ErrDuplicatePartNumber, services.ErrDuplicateBarcode:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	item.ItemID = id
	return c.JSON(http.StatusCreated, item)
}

// UpdateItem handles the update of an existing item
func (h *InventoryHandler) UpdateItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid item ID")
	}

	item := new(inventorymodels.Item)
	if err := c.Bind(item); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	item.ItemID = id

	ctx := c.Request().Context()
	err = h.service.UpdateItem(ctx, item)
	if err != nil {
		switch err {
		case services.ErrItemNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case services.ErrDuplicatePartNumber, services.ErrDuplicateBarcode:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, item)
}

// DeleteItem handles the deletion of an item
func (h *InventoryHandler) DeleteItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid item ID")
	}

	ctx := c.Request().Context()
	err = h.service.DeleteItem(ctx, id)
	if err != nil {
		switch err {
		case services.ErrItemNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// GetCompatibilities handles the retrieval of compatibilities for an item
func (h *InventoryHandler) GetCompatibilities(c echo.Context) error {
	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid item ID")
	}

	ctx := c.Request().Context()
	compatibilities, err := h.service.GetCompatibilities(ctx, itemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, compatibilities)
}

// AddCompatibility handles adding a new compatibility
func (h *InventoryHandler) AddCompatibility(c echo.Context) error {
	compatibility := new(inventorymodels.Compatibility)
	if err := c.Bind(compatibility); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	id, err := h.service.AddCompatibility(ctx, compatibility)
	if err != nil {
		switch err {
		case services.ErrItemNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case services.ErrCompatibilityExists:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	compatibility.CompatID = id
	return c.JSON(http.StatusCreated, compatibility)
}

// RemoveCompatibility handles removing a compatibility
func (h *InventoryHandler) RemoveCompatibility(c echo.Context) error {
	itemID, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid item ID")
	}

	submodelID, err := strconv.Atoi(c.Param("submodelId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid submodel ID")
	}

	ctx := c.Request().Context()
	err = h.service.RemoveCompatibility(ctx, itemID, submodelID)
	if err != nil {
		if err.Error() == "compatibility not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// GetCompatibleItems handles retrieving all compatible items for a vehicle submodel
func (h *InventoryHandler) GetCompatibleItems(c echo.Context) error {
	submodelID, err := strconv.Atoi(c.Param("submodelId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid submodel ID")
	}

	ctx := c.Request().Context()
	items, err := h.service.GetCompatibleItems(ctx, submodelID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, items)
}

func (h *InventoryHandler) GetBarcodeImage(c echo.Context) error {
	barcode := c.Param("barcode")
	if barcode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "barcode is required")
	}

	imgBytes, err := h.barcodeService.GenerateBarcodeImage(barcode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, "image/png", imgBytes)
}
