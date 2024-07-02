import { expect, test } from 'vitest'
import { nextTodoStatus, previousTodoStatus } from './states'
import type { TodoItem, TodoStatus } from '../net/models'

function testItem(status: TodoStatus): TodoItem {
  return {
    id: '1',
    created_by: '1',
    description: 'Test',
    status: status,
    updated_at: '2021-01-01T00:00:00Z'
  }
}

test('next todo > ongoing', () => {
  expect(nextTodoStatus(testItem('todo'))).toBe('ongoing')
})

test('next ongoing > done', () => {
  expect(nextTodoStatus(testItem('ongoing'))).toBe('done')
})

test('next done > todo', () => {
  expect(nextTodoStatus(testItem('done'))).toBe('todo')
})

test('previous todo > ongoing', () => {
  expect(previousTodoStatus(testItem('todo'))).toBe('done')
})

test('previous ongoing > done', () => {
  expect(previousTodoStatus(testItem('ongoing'))).toBe('todo')
})

test('previous done > todo', () => {
  expect(previousTodoStatus(testItem('done'))).toBe('ongoing')
})
