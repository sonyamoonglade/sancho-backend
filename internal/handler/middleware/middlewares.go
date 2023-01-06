package middleware

type Middlewares struct {
	JWTAuth    *JWTAuthMiddleware
	XRequestID *XRequestIDMiddleware
}

func NewMiddlewares(jwtAuth *JWTAuthMiddleware, xRequestID *XRequestIDMiddleware) *Middlewares {
	return &Middlewares{
		JWTAuth:    jwtAuth,
		XRequestID: xRequestID,
	}
}
