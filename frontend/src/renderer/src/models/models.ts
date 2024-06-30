type RegisterResponse = {
    account_number: string
}

type LogInResponse = {
    access_token: string
}

type CreateTodoListResponse = {
    todo_list_id: string
}

type TodoListReponses = {
    todos: TodoItem[]
}

type TodoItem = {
    id: string
    created_by: string
    description: string
    status: TodoStatus
    updated_at: string
}

type TodoStatus = 'todo' | 'ongoing' | 'done'