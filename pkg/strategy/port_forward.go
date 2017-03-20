package strategy

import (
    "github.com/nebtex/menshend/pkg/pfclient/chisel/server"
    "net/http"
    . "github.com/nebtex/menshend/pkg/utils"
    "net/url"
    "strings"
    "github.com/nebtex/menshend/pkg/resolvers"
    vault "github.com/hashicorp/vault/api"
)

type PortForward struct {
}

//PortForward ..
func (r *PortForward) Execute(rs resolvers.Resolver, tokenInfo *vault.Secret) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Context().Value("IsBrowserRequest").(bool) {
            panic(BadRequest.WithUserMessage("This endpoint cant be used from the browser, download the menshend client"))
        }
        b := rs.Resolve(tokenInfo)
        if !b.Passed() {
            panic(NotAuthorized.WithUserMessage(b.Error().Error()))
        }
        var err error
        URL, err := url.Parse(b.BaseUrl())
        HttpCheckPanic(err, InternalError)
        host := strings.Split(URL.Host, ":")[0]
        remote := host + ":" + r.Header.Get("X-Menshend-Port-Forward")
        chiselServer, err := chserver.NewServer(&chserver.Config{
            Remote:remote,
        })
        HttpCheckPanic(err, InternalError)
        chiselServer.HandleHTTP(w, r)
    })
}


