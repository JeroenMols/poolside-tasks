import type { TodoItem, TodoStatus } from '../net/models'

const allowedStatus: TodoStatus[] = ['todo', 'ongoing', 'done']

export function nextTodoStatus(todo: TodoItem): TodoStatus {
  return allowedStatus[(allowedStatus.indexOf(todo.status) + 1) % allowedStatus.length]
}
export function previousTodoStatus(todo: TodoItem): TodoStatus {
  return allowedStatus[
    (allowedStatus.indexOf(todo.status) + allowedStatus.length - 1) % allowedStatus.length
  ]
}
