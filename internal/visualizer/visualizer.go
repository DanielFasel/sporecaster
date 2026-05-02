package visualizer

import (
	"net/http"

	"github.com/DanielFasel/sporecaster/internal/spore"
	vizserver "github.com/DanielFasel/sporecaster/internal/visualizer/server"
)

func Serve(s *spore.Spore, addr string) error {
	return http.ListenAndServe(addr, vizserver.Handler(s))
}
