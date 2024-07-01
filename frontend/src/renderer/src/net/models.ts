export type RegisterRequest = {
  name: string
}

export type RegisterResponse = {
  user_id: string
}

export type LogInRequest = {
  user_id: string
}

export type LogInResponse = {
  access_token: string
}

export type CreateTodoListResponse = {
  todo_list_id: string
}

export type TodoListReponses = {
  todo_list_id: string
  todos: TodoItem[]
}

export type TodoItem = {
  id: string
  created_by: string
  description: string
  status: TodoStatus
  updated_at: string
}

export type ErrorResponse = {
  error: string
}

export type TodoStatus = 'todo' | 'ongoing' | 'done'
