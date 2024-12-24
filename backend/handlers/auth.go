// handlers/auth.go
package handlers

import (
    "context"
    "crypto/sha256"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt"
    "survey-platform-server/db"
    "survey-platform-server/models"
)

var JwtKey = []byte("your_secret_key") // Замените на безопасный ключ

type Credentials struct {
    Login    string `json:"login"`
    Password string `json:"password"`
}

type Claims struct {
    AdminID int `json:"admin_id"`
    jwt.StandardClaims
}

type contextKey string

const adminIDKey = contextKey("adminID")

func SetAdminID(ctx context.Context, adminID int) context.Context {
    return context.WithValue(ctx, adminIDKey, adminID)
}

func GetAdminID(ctx context.Context) (int, bool) {
    adminID, ok := ctx.Value(adminIDKey).(int)
    return adminID, ok
}

// RegisterHandler регистрирует нового администратора
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
        return
    }

    if creds.Login == "" || creds.Password == "" {
        http.Error(w, "Логин и пароль обязательны", http.StatusBadRequest)
        return
    }

    // Хэширование пароля
    hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(creds.Password)))

    // Вставка нового администратора в базу данных
    var adminID int
    err = db.DB.QueryRow(
        "INSERT INTO Admins (Login, Password) VALUES ($1, $2) RETURNING AdminID",
        creds.Login, hashedPassword).Scan(&adminID)
    if err != nil {
        http.Error(w, "Ошибка при регистрации администратора", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":  "Успешная регистрация",
        "admin_id": adminID,
    })
}

// LoginHandler осуществляет вход администратора
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
        return
    }

    if creds.Login == "" || creds.Password == "" {
        http.Error(w, "Логин и пароль обязательны", http.StatusBadRequest)
        return
    }

    // Хэширование пароля
    hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(creds.Password)))

    // Проверка администратора в базе данных
    var admin models.Admin
    err = db.DB.QueryRow(
        "SELECT AdminID, Login, Password FROM Admins WHERE Login=$1",
        creds.Login).Scan(&admin.AdminID, &admin.Login, &admin.Password)
    if err != nil {
        http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
        return
    }

    if admin.Password != hashedPassword {
        http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
        return
    }

    // Создание JWT-токена
    expirationTime := time.Now().Add(10*time.Second) // Токен действует 24 часа
    claims := &Claims{
        AdminID: admin.AdminID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(JwtKey)
    if err != nil {
        http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
        return
    }

    // Установка куки
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    tokenString,
        Expires:  expirationTime,
        HttpOnly: true,
        Secure:   false, // Установите true при использовании HTTPS
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Успешный вход",
        "admin_id": admin.AdminID,
    })
}

// LogoutHandler осуществляет выход администратора
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    // Установка куки с истекшим временем
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    "",
        Expires:  time.Unix(0, 0),
        HttpOnly: true,
        Secure:   false, // RIP https :(
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Успешный выход",
    })
}

// проверяет аутентификацию админа
func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

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

    claims := &Claims{}

    // Парсинг токена
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return JwtKey, nil
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

    // Токен валиден
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "authenticated": true,
        "admin_id":      claims.AdminID,
    })
}

// Middleware для аутентификации администратора
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

        claims := &Claims{}

        // Парсинг токена
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return JwtKey, nil
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

        // Добавление AdminID в контекст запроса
        ctx := context.WithValue(r.Context(), adminIDKey, claims.AdminID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
