package helpers

import "Todo/models"

func CreateUserResponse(user models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"created_at": user.CreatedAt,
	}
}

func CreateTodoResponse(todo models.Todo) map[string]interface{} {
	return map[string]interface{}{
		"id":         todo.ID,
		"title":      todo.Title,
		"decription": todo.Description,
		"completed":  todo.Completed,
		"created_at": todo.CreatedAt,
	}
}
