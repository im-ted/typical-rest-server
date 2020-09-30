package library

import (
	"github.com/typical-go/typical-rest-server/internal/app/domain/library/controller"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

type (
	// Router to server
	Router struct {
		dig.In
		BookCntrl controller.BookCntrl
	}
)

var _ typrest.Router = (*Router)(nil)

// SetRoute to echo server
func (r *Router) SetRoute(e typrest.Server) {
	group := e.Group("/library")
	typrest.SetRoute(group,
		&r.BookCntrl,
	)
}