package core_cookie

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

func WithForwardedMetadata(r *http.Request) context.Context {
	pairs := []string{}
	if v := r.Header.Get("Cookie"); v != "" {
		pairs = append(pairs, "cookie", v)
	}
	if v := r.Header.Get("Authorization"); v != "" {
		pairs = append(pairs, "authorization", v)
	}
	if v := r.Header.Get("User-Agent"); v != "" {
		pairs = append(pairs, "user-agent", v)
	}
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		pairs = append(pairs, "x-client-ip", strings.Split(v, ",")[0])
	} else {
		pairs = append(pairs, "x-client-ip", r.RemoteAddr)
	}
	if len(pairs) == 0 {
		return r.Context()
	}
	md := metadata.Pairs(pairs...)
	return metadata.NewOutgoingContext(r.Context(), md)
}

func WriteSetCookiesFromMetadata(w http.ResponseWriter, md metadata.MD) {
	for _, sc := range md.Get("set-cookie") {
		w.Header().Add("Set-Cookie", sc)
	}
}
