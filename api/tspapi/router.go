package tspapi

import (
	"HeteroAntColonySystem/api/tspapi/dto"
	"HeteroAntColonySystem/api/utils/ginerr"
	"HeteroAntColonySystem/pkg/tsplib"
	"errors"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TSPRouter struct {
	parser *tsplib.TSPLIBParser
	fs     fs.FS
}

func NewRouter(r *gin.RouterGroup, parser *tsplib.TSPLIBParser, fs fs.FS) *TSPRouter {
	tsp := &TSPRouter{
		parser: parser,
		fs:     fs,
	}
	gr := r.Group("/tsp")

	gr.GET("/:file", tsp.GetTSP)
	gr.GET("/files", tsp.ListTSP)
	gr.POST("/parse", tsp.ParseTSP)
	gr.POST("/stats", tsp.GraphStats)

	return tsp
}

func (r *TSPRouter) GetTSP(c *gin.Context) {
	req, err := dto.DeserializeGetTSPRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	f, err := r.fs.Open(req.File)
	if err != nil {
		switch {
		case errors.Is(err, fs.ErrNotExist):
			c.JSON(http.StatusNotFound, ginerr.ErrJSONBody(err))
		case errors.Is(err, fs.ErrPermission):
			c.JSON(http.StatusForbidden, ginerr.ErrJSONBody(err))
		default:
			c.JSON(http.StatusInternalServerError, ginerr.ErrJSONBody(err))
		}
	}
	defer f.Close()

	g, err := r.parser.Parse(f)
	if err != nil {
		switch {
		case errors.Is(err, tsplib.ErrInvalidFormat):
			c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		case errors.Is(err, tsplib.ErrInvalidData):
			c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		case errors.Is(err, tsplib.ErrSectionNotFound):
			c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		case errors.Is(err, tsplib.ErrAdapterNotFound):
			c.JSON(http.StatusNotImplemented, ginerr.ErrJSONBody(err))
		case errors.Is(err, tsplib.ErrUnsupportedType):
			c.JSON(http.StatusNotImplemented, ginerr.ErrJSONBody(err))
		default:
			c.JSON(http.StatusInternalServerError, ginerr.ErrJSONBody(err))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SerializeGetTSPResponse(c, g))
}

func (r *TSPRouter) ParseTSP(c *gin.Context) {
	req, err := dto.DeserializeParseTSPRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	g, err := r.parser.Parse(req.File)
	if err != nil {
		switch {
		case errors.Is(err, tsplib.ErrInvalidFormat):
			c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		case errors.Is(err, tsplib.ErrInvalidData):
			c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		case errors.Is(err, tsplib.ErrSectionNotFound):
			c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		case errors.Is(err, tsplib.ErrAdapterNotFound):
			c.JSON(http.StatusNotImplemented, ginerr.ErrJSONBody(err))
		case errors.Is(err, tsplib.ErrUnsupportedType):
			c.JSON(http.StatusNotImplemented, ginerr.ErrJSONBody(err))
		default:
			c.JSON(http.StatusInternalServerError, ginerr.ErrJSONBody(err))
		}
		return
	}

	c.JSON(http.StatusOK, dto.SerializeParseTSPResponse(c, g))
}

func (r *TSPRouter) ListTSP(c *gin.Context) {
	files, err := fs.ReadDir(r.fs, ".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginerr.ErrJSONBody(err))
		return
	}

	fnames := make([]string, 0, len(files))
	for _, f := range files {
		if f.Type().IsRegular() && len(f.Name()) > 4 && f.Name()[len(f.Name())-4:] == ".tsp" {
			fnames = append(fnames, f.Name())
		}
	}

	c.JSON(http.StatusOK, dto.SerializeListTSPResponse(c, fnames))
}

func (r *TSPRouter) GraphStats(c *gin.Context) {
	g, err := dto.DeserializeGetGraphStatsRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginerr.ErrJSONBody(err))
		return
	}

	stats := CalculateGraphStats(g)
	c.JSON(http.StatusOK, dto.SerializeGetGraphStatsResponse(c, stats))
}
