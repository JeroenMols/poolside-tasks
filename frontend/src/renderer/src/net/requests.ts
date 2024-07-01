import {
    type RegisterResponse,
    type RegisterRequest,
    type ErrorResponse,
    type LogInResponse,
    type LogInRequest,
    type CreateTodoListResponse,
    type GetTodoListResponse,
    type TodoItem,
    type CreateTodoRequest,
    type TodoStatus,
    type UpdateTodoRequest
} from "./models"

export async function createAccount(userName: string): Promise<RegisterResponse | ErrorResponse> {
    if (userName.length < 5) {
        return { error: 'Name must be at least 5 characters long' }
    }
    return await doRequest<RegisterRequest, RegisterResponse>(
        'http://localhost:8080/users/register',
        'POST',
        { name: userName }
    )
}

export async function logIn(userId: string): Promise<LogInResponse | ErrorResponse> {
    return await doRequest<LogInRequest, LogInResponse>(
        'http://localhost:8080/users/login',
        'POST',
        { user_id: userId }
    )
}

export async function createTodoList(accessToken: string): Promise<CreateTodoListResponse | ErrorResponse> {
    return await doRequestWithAuth<any, CreateTodoListResponse>(
        'http://localhost:8080/todolists',
        'POST',
        accessToken,
        {},
    )
}

export async function getTodoList(accessToken: string, listId: string): Promise<GetTodoListResponse | ErrorResponse> {
    return await doRequestWithAuth<undefined, GetTodoListResponse>(
        `http://localhost:8080/todolists/${listId}`,
        'GET',
        accessToken,
        undefined,
    )
}

export async function createTodoItem(accessToken: string, listId: string, description: string): Promise<TodoItem | ErrorResponse> {
    if (description.length < 1) {
        return { error: 'Description cannot be empty' }
    }
    return await doRequestWithAuth<CreateTodoRequest, TodoItem>(
        `http://localhost:8080/todos`,
        'POST',
        accessToken,
        { description: description, todo_list_id: listId },
    )
}

export async function updateTodoItem(accessToken: string, todoId: string, status: TodoStatus): Promise<TodoItem | ErrorResponse> {
    return await doRequestWithAuth<UpdateTodoRequest, TodoItem>(
        `http://localhost:8080/todos/${todoId}`,
        'PUT',
        accessToken,
        { status: status },
    )
}

async function doRequest<T, K>(url: string, method: string, body: T): Promise<K | ErrorResponse> {
    try {
        const response = await fetch(url, {
            method: method,
            body: JSON.stringify(body)

        })

        if (response.ok) {
            return await response.json() as K
        } else {
            return await response.json() as ErrorResponse
        }
    } catch (e) {
        return { error: "Server offline?" }
    }
}

async function doRequestWithAuth<T, K>(url: string, method: string, token: string, body: T): Promise<K | ErrorResponse> {
    try {
        const response = await fetch(url, {
            method: method,
            body: JSON.stringify(body),
            headers: { Authorization: `${token}` },
        })

        if (response.ok) {
            return await response.json() as K
        } else {
            return await response.json() as ErrorResponse
        }
    } catch (e) {
        return { error: "Server offline?" }
    }
}