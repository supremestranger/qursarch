// middleware/auth.go
package middleware

import (
    "net/http"

    "github.com/golang-jwt/jwt"
    "survey-platform-server/handlers"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Получение токена из куки
        c, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                http.Error(w, "Неавторизованный доступ: нет токена", http.StatusUnauthorized)
                return
            }
            http.Error(w, "Ошибка при получении куки", http.StatusBadRequest)
            return
        }

        tokenStr := c.Value

        claims := &handlers.Claims{}

        // Парсинг токена
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return handlers.JwtKey, nil
        })
        if err != nil {
            if err == jwt.ErrSignatureInvalid {
                http.Error(w, "Неверная подпись токена", http.StatusUnauthorized)
                return
            }
            http.Error(w, "Невалидный токен", http.StatusBadRequest)
            return
        }
        if !token.Valid {
            http.Error(w, "Невалидный токен", http.StatusUnauthorized)
            return
        }

        // Добавление UserID в контекст запроса
        ctx := r.Context()
        ctx = handlers.SetUserID(ctx, claims.UserID)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}
