import {
    type RegisterResponse,
    type RegisterRequest,
    type ErrorResponse,
    type LogInResponse,
    type LogInRequest
} from "./models"

export async function createAccount(username: string): Promise<RegisterResponse | ErrorResponse> {
    if (username.length < 5) {
        return { error: 'Name must be at least 5 characters long' }
    }
    return await doRequest<RegisterRequest, RegisterResponse>(
        'http://localhost:8080/users/register',
        'POST',
        { name: username }
    )
}

export async function logIn(user_id: string): Promise<LogInResponse | ErrorResponse> {
    return await doRequest<LogInRequest, LogInResponse>(
        'http://localhost:8080/users/login',
        'POST',
        { user_id: user_id }
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