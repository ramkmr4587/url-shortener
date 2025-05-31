package handler

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"url-shortener/internal/service"
	"url-shortener/model"

	restful "github.com/emicklei/go-restful/v3"
)

type Handler struct {
	URLService *service.URLService
}

func NewHandler(svc *service.URLService) *Handler {
	return &Handler{URLService: svc}
}

func (h *Handler) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/shorten").To(h.Shorten))
	ws.Route(ws.GET("/r/{short}").To(h.Redirect))
	ws.Route(ws.GET("/metrics").To(h.Metrics))

	container.Add(ws)
}

func (h *Handler) Shorten(req *restful.Request, resp *restful.Response) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in Shorten: %v\n", r) // Debug log
			resp.WriteError(http.StatusInternalServerError, r.(error))
		}
	}()
	var in model.URLRequest
	if err := req.ReadEntity(&in); err != nil {
		fmt.Printf("ReadEntity error: %v\n", err) // Debug log
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	fmt.Printf("Parsed URLRequest: %+v\n", in) // Debug log
	short := h.URLService.ShortenURL(in.OriginalURL)
	fmt.Printf("Shortened URL: %s\n", short) // Debug log
	resp.WriteEntity(model.URLResponse{ShortURL: short})
}

func (h *Handler) Redirect(req *restful.Request, resp *restful.Response) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in Redirect: %v\n", r) // Debug log
			resp.WriteError(http.StatusInternalServerError, r.(error))
		}
	}()
	short := strings.TrimSpace(req.PathParameter("short"))
	original, ok := h.URLService.GetOriginalURL(short)
	if !ok {
		resp.WriteErrorString(http.StatusNotFound, "short URL not found")
		return
	}
	resp.AddHeader("Location", original)
	resp.WriteHeader(http.StatusMovedPermanently)
}

func (h *Handler) Metrics(req *restful.Request, resp *restful.Response) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in Metrics: %v\n", r) // Debug log
			resp.WriteError(http.StatusInternalServerError, r.(error))
		}
	}()
	domains := h.URLService.GetTopDomains(3)

	type stat struct {
		Domain string `json:"domain"`
		Count  int    `json:"count"`
	}

	var list []stat
	for k, v := range domains {
		list = append(list, stat{k, v})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Count > list[j].Count
	})

	if len(list) > 3 {
		list = list[:3]
	}

	resp.WriteEntity(list)
}
