#!/bin/bash

# Базовый URL API
BASE_URL="http://localhost:8080/api"
JWT_TOKEN=""
USER_ID=""

echo "=== Testing User Rewards API ==="

# Функция для вывода цветного результата
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "✅ $2"
    else
        echo -e "❌ $2"
        echo "Response: $3"
    fi
}

# 1. Health check
echo ""
echo "1. Testing health check..."
HEALTH_RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/health_response.json "$BASE_URL/health")
HTTP_CODE=${HEALTH_RESPONSE: -3}
RESPONSE_CONTENT=$(cat /tmp/health_response.json)

if [ $HTTP_CODE -eq 200 ]; then
    print_result 0 "Health check passed"
    echo "   Response: $RESPONSE_CONTENT"
else
    print_result 1 "Health check failed" "HTTP $HTTP_CODE: $RESPONSE_CONTENT"
    exit 1
fi

# 2. Получение списка задач
echo ""
echo "2. Testing get tasks list..."
TASKS_RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/tasks_response.json "$BASE_URL/tasks")
HTTP_CODE=${TASKS_RESPONSE: -3}
RESPONSE_CONTENT=$(cat /tmp/tasks_response.json)

if [ $HTTP_CODE -eq 200 ]; then
    print_result 0 "Get tasks passed"
    echo "   Available tasks retrieved"
else
    print_result 1 "Get tasks failed" "HTTP $HTTP_CODE: $RESPONSE_CONTENT"
fi

# 3. Регистрация пользователя
echo ""
echo "3. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/register_response.json \
    -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "testuser",
        "email": "test@example.com",
        "password": "password123"
    }')
HTTP_CODE=${REGISTER_RESPONSE: -3}
RESPONSE_CONTENT=$(cat /tmp/register_response.json)

if [ $HTTP_CODE -eq 201 ]; then
    print_result 0 "User registration passed"
    # Извлекаем токен и ID пользователя
    JWT_TOKEN=$(echo $RESPONSE_CONTENT | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    USER_ID=$(echo $RESPONSE_CONTENT | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "   User ID: $USER_ID"
    echo "   Token received"
else
    print_result 1 "User registration failed" "HTTP $HTTP_CODE: $RESPONSE_CONTENT"
    echo "   Trying login instead..."
    LOGIN_RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/login_response.json \
        -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "testuser",
            "password": "password123"
        }')
    HTTP_CODE=${LOGIN_RESPONSE: -3}
    RESPONSE_CONTENT=$(cat /tmp/login_response.json)
    
    if [ $HTTP_CODE -eq 200 ]; then
        print_result 0 "User login passed"
        JWT_TOKEN=$(echo $RESPONSE_CONTENT | grep -o '"token":"[^"]*' | cut -d'"' -f4)
        USER_ID=$(echo $RESPONSE_CONTENT | grep -o '"id":"[^"]*' | cut -d'"' -f4)
        echo "   User ID: $USER_ID"
        echo "   Token received"
    else
        print_result 1 "User login also failed" "HTTP $HTTP_CODE: $RESPONSE_CONTENT"
        exit 1
    fi
fi

# 4. Получение статуса пользователя
echo ""
echo "4. Testing get user status..."
STATUS_RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/status_response.json \
    -X GET "$BASE_URL/users/$USER_ID/status" \
    -H "Authorization: Bearer $JWT_TOKEN")
HTTP_CODE=${STATUS_RESPONSE: -3}
RESPONSE_CONTENT=$(cat /tmp/status_response.json)

if [ $HTTP_CODE -eq 200 ]; then
    print_result 0 "Get user status passed"
    BALANCE=$(echo $RESPONSE_CONTENT | grep -o '"balance":[0-9]*' | cut -d':' -f2)
    echo "   User balance: $BALANCE points"
else
    print_result 1 "Get user status failed" "HTTP $HTTP_CODE: $RESPONSE_CONTENT"
fi

# 5. Выполнение задачи
echo ""
echo "5. Testing complete task..."
COMPLETE_RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/complete_response.json \
    -X POST "$BASE_URL/users/$USER_ID/task/complete" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"task_id": "1"}')
HTTP_CODE=${COMPLETE_RESPONSE: -3}
RESPONSE_CONTENT=$(cat /tmp/complete_response.json)

if [ $HTTP_CODE -eq 200 ]; then
    print_result 0 "Complete task passed"
    echo "   Task completed successfully"
else
    print_result 1 "Complete task failed" "HTTP $HTTP_CODE: $RESPONSE_CONTENT"
fi

# 6. Получение обновленного статуса пользователя
echo ""
echo "6. Testing get updated user status..."
STATUS_RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/status2_response.json \
    -X GET "$BASE_URL/users/$USER_ID/status" \
    -H "Authorization: Bearer $JWT_TOKEN")
HTTP_CODE=${STATUS_RESPONSE: -3}
RESPONSE_CONTENT=$(cat /tmp/status2_response.json)

if [ $HTTP_CODE -eq 200 ]; then
    print_result 0 "Get updated user status passed"
    BALANCE=$(echo $RESPONSE_CONTENT | grep -o '"balance":[0-9]*' | cut -d':' -f2)
    echo "   Updated balance: $BALANCE points"
else
    print_result 1 "Get updated user status failed" "HTTP $HTTP_CODE: $RESPONSE_CONTENT"
fi

# 7. Получение leaderboard
echo ""
echo "7. Testing get leaderboard..."
LEADERBOARD_RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/leaderboard_response.json \
    -X GET "$BASE_URL/users/leaderboard?limit=5" \
    -H "Authorization: Bearer $JWT_TOKEN")
HTTP_CODE=${LEADERBOARD_RESPONSE: -3}
RESPONSE_CONTENT=$(cat /tmp/leaderboard_response.json)

if [ $HTTP_CODE -eq 200 ]; then
    print_result 0 "Get leaderboard passed"
    echo "   Leaderboard retrieved successfully"
else
    print_result 1 "Get leaderboard failed" "HTTP $HTTP_CODE: $RESPONSE_CONTENT"
fi

# 8. Тестирование установки реферера (нужен второй пользователь)
echo ""
echo "8. Testing set referrer (creating second user)..."
REGISTER_RESPONSE2=$(curl -s -w "%{http_code}" -o /tmp/register2_response.json \
    -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "testuser2",
        "email": "test2@example.com",
        "password": "password123"
    }')
HTTP_CODE=${REGISTER_RESPONSE2: -3}
RESPONSE_CONTENT=$(cat /tmp/register2_response.json)

if [ $HTTP_CODE -eq 201 ]; then
    USER_ID2=$(echo $RESPONSE_CONTENT | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    JWT_TOKEN2=$(echo $RESPONSE_CONTENT | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "   Second user created: $USER_ID2"
    
    # Устанавливаем реферера
    REFERRER_RESPONSE=$(curl -s -w "%{http_code}" -o /tmp/referrer_response.json \
        -X POST "$BASE_URL/users/$USER_ID2/referrer" \
        -H "Authorization: Bearer $JWT_TOKEN2" \
        -H "Content-Type: application/json" \
        -d "{\"referrer_id\": \"$USER_ID\"}")
    HTTP_CODE=${REFERRER_RESPONSE: -3}
    RESPONSE_CONTENT=$(cat /tmp/referrer_response.json)
    
    if [ $HTTP_CODE -eq 200 ]; then
        print_result 0 "Set referrer passed"
        echo "   Referrer set successfully"
    else
        print_result 1 "Set referrer failed" "HTTP $HTTP_CODE: $RESPONSE_CONTENT"
    fi
else
    print_result 1 "Second user creation failed" "HTTP $HTTP_CODE: $RESPONSE_CONTENT"
fi

echo ""
echo "=== API Testing Completed ==="