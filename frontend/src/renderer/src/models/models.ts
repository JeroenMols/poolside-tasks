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
    created_by: string
    description: string
    status: string
    updated_at: string
}