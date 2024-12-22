// handlers/auth.go
package handlers

import (
    "context"
    "crypto/sha256"
    "encoding/json"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt"
    "survey-platform-server/db"
    "survey-platform-server/models"
    "fmt"
)

var JwtKey = []byte("your_secret_key") // Заменить на мой секретный ключ

type Credentials struct {
    Login    string `json:"login"`
    Password string `json:"password"`
}

type Claims struct {
    UserID int `json:"user_id"`
    jwt.StandardClaims
}

type contextKey string

const userIDKey = contextKey("userID")

func SetUserID(ctx context.Context, userID int) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(userIDKey).(int)
    return userID, ok
}

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

    // Вставка нового пользователя в базу данных
    var userID int
    err = db.DB.QueryRow(
        "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING userid",
        creds.Login, hashedPassword).Scan(&userID)
    if err != nil {
        http.Error(w, "Ошибка при регистрации пользователя", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":  "Успешная регистрация",
        "user_id":  userID,
    })
}

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

    // Проверка пользователя в базе данных
    var user models.User
    err = db.DB.QueryRow(
        "SELECT userid, login, password FROM users WHERE login=$1",
        creds.Login).Scan(&user.UserID, &user.Login, &user.Password)
    if err != nil {
        http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
        return
    }

    if user.Password != hashedPassword {
        http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
        return
    }

    // Создание JWT-токена
    expirationTime := time.Now().Add(24 * time.Hour) // Токен действует 24 часа
    claims := &Claims{
        UserID: user.UserID,
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
        "user_id": user.UserID,
    })
}

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
        Secure:   false, // Установите true при использовании HTTPS
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Успешный выход",
    })
}

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
        "user_id":       claims.UserID,
    })
}
