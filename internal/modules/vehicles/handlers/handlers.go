package handlers

import (
	"net/http"
	"strconv"

	vehiclemodels "github.com/hsrvms/autoparts/internal/modules/vehicles/models"
	"github.com/hsrvms/autoparts/internal/modules/vehicles/services"
	"github.com/labstack/echo/v4"
)

type VehicleHandler struct {
	service services.VehicleService
}

func NewVehicleHandler(service services.VehicleService) *VehicleHandler {
	return &VehicleHandler{
		service: service,
	}
}

// Make handlers
func (h *VehicleHandler) GetAllMakes(c echo.Context) error {
	ctx := c.Request().Context()
	makes, err := h.service.GetAllMakes(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, makes)
}

func (h *VehicleHandler) GetMakeByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid make ID")
	}

	ctx := c.Request().Context()
	make, err := h.service.GetMakeByID(ctx, id)
	if err != nil {
		if err == services.ErrMakeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, make)
}

func (h *VehicleHandler) CreateMake(c echo.Context) error {
	make := new(vehiclemodels.Make)
	if err := c.Bind(make); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	id, err := h.service.CreateMake(ctx, make)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	make.MakeID = id
	return c.JSON(http.StatusCreated, make)
}

func (h *VehicleHandler) UpdateMake(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid make ID")
	}

	make := new(vehiclemodels.Make)
	if err := c.Bind(make); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	make.MakeID = id

	ctx := c.Request().Context()
	err = h.service.UpdateMake(ctx, make)
	if err != nil {
		switch err {
		case services.ErrMakeNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, make)
}

func (h *VehicleHandler) DeleteMake(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid make ID")
	}

	ctx := c.Request().Context()
	err = h.service.DeleteMake(ctx, id)
	if err != nil {
		switch err {
		case services.ErrMakeNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// Model handlers
func (h *VehicleHandler) GetAllModels(c echo.Context) error {
	ctx := c.Request().Context()
	models, err := h.service.GetAllModels(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, models)
}

func (h *VehicleHandler) GetModelsByMake(c echo.Context) error {
	makeID, err := strconv.Atoi(c.Param("makeId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid make ID")
	}

	ctx := c.Request().Context()
	models, err := h.service.GetModelsByMake(ctx, makeID)
	if err != nil {
		if err == services.ErrMakeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, models)
}

func (h *VehicleHandler) GetModelByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid model ID")
	}

	ctx := c.Request().Context()
	model, err := h.service.GetModelByID(ctx, id)
	if err != nil {
		if err == services.ErrModelNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model)
}

func (h *VehicleHandler) CreateModel(c echo.Context) error {
	model := new(vehiclemodels.Model)
	if err := c.Bind(model); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	id, err := h.service.CreateModel(ctx, model)
	if err != nil {
		if err == services.ErrMakeNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	model.ModelID = id
	return c.JSON(http.StatusCreated, model)
}

func (h *VehicleHandler) UpdateModel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid model ID")
	}

	model := new(vehiclemodels.Model)
	if err := c.Bind(model); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	model.ModelID = id

	ctx := c.Request().Context()
	err = h.service.UpdateModel(ctx, model)
	if err != nil {
		switch err {
		case services.ErrModelNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case services.ErrMakeNotFound:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, model)
}

func (h *VehicleHandler) DeleteModel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid model ID")
	}

	ctx := c.Request().Context()
	err = h.service.DeleteModel(ctx, id)
	if err != nil {
		switch err {
		case services.ErrModelNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// Submodel handlers
func (h *VehicleHandler) GetAllSubmodels(c echo.Context) error {
	ctx := c.Request().Context()
	submodels, err := h.service.GetAllSubmodels(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, submodels)
}

func (h *VehicleHandler) GetSubmodelsByModel(c echo.Context) error {
	modelID, err := strconv.Atoi(c.Param("modelId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid model ID")
	}

	ctx := c.Request().Context()
	submodels, err := h.service.GetSubmodelsByModel(ctx, modelID)
	if err != nil {
		if err == services.ErrModelNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, submodels)
}

func (h *VehicleHandler) GetSubmodelByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid submodel ID")
	}

	ctx := c.Request().Context()
	submodel, err := h.service.GetSubmodelByID(ctx, id)
	if err != nil {
		if err == services.ErrSubmodelNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, submodel)
}

func (h *VehicleHandler) CreateSubmodel(c echo.Context) error {
	submodel := new(vehiclemodels.Submodel)
	if err := c.Bind(submodel); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	id, err := h.service.CreateSubmodel(ctx, submodel)
	if err != nil {
		if err == services.ErrModelNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	submodel.SubmodelID = id
	return c.JSON(http.StatusCreated, submodel)
}

func (h *VehicleHandler) UpdateSubmodel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid submodel ID")
	}

	submodel := new(vehiclemodels.Submodel)
	if err := c.Bind(submodel); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	submodel.SubmodelID = id

	ctx := c.Request().Context()
	err = h.service.UpdateSubmodel(ctx, submodel)
	if err != nil {
		switch err {
		case services.ErrSubmodelNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case services.ErrModelNotFound:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, submodel)
}

func (h *VehicleHandler) DeleteSubmodel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid submodel ID")
	}

	ctx := c.Request().Context()
	err = h.service.DeleteSubmodel(ctx, id)
	if err != nil {
		switch err {
		case services.ErrSubmodelNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusNoContent)
}
