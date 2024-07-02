<script lang="ts">
  import { onMount } from 'svelte'
  import ErrorBanner from './error-banner.svelte'
  import type { TodoItem, TodoStatus } from '../net/models'
  import { createTodoItem, getTodoList, updateTodoItem } from '../net/requests'
  import { nextTodoStatus, previousTodoStatus } from '../utils/states'

  export let accessToken: string = ''
  export let todoListId: string = ''

  let todos: TodoItem[] = []
  let newTodoDescription = ''
  let errorMessage = ''

  onMount(async () => {
    await getTodosForList()
  })

  const getTodosForList = async () => {
    let response = await getTodoList(accessToken, todoListId)
    if ('error' in response) {
      errorMessage = response.error as string
    } else {
      // TODO validate result here
      todos = response.todos
    }
  }

  const createTodo = async () => {
    let response = await createTodoItem(accessToken, todoListId, newTodoDescription)
    if ('error' in response) {
      errorMessage = response.error as string
    } else {
      newTodoDescription = ''
      todos = [...todos, response]
    }
    await getTodosForList()
  }

  const updateTodo = async (todo: TodoItem, newStatus: TodoStatus) => {
    let response = await updateTodoItem(accessToken, todo.id, newStatus)
    if ('error' in response) {
      errorMessage = response.error as string
    } else {
      await getTodosForList()
    }
  }

  // Note: this intentionally allows invalid status transitions to demonstrate backend validation
  const prevStatus = async (todo: TodoItem) => {
    await updateTodo(todo, previousTodoStatus(todo))
  }
  const nextStatus = async (todo: TodoItem) => {
    await updateTodo(todo, nextTodoStatus(todo))
  }

  const onDismissError = () => {
    errorMessage = ''
    newTodoDescription = ''
  }
</script>

<ErrorBanner {errorMessage} {onDismissError} />

{#if todos.length == 0}
  <h1>Nothing to do</h1>
  <p>Create your first task below</p>
{:else}
  <h1>Your tasks 🚀</h1>
  <ul class="task-list">
    {#each todos as todo}
      <li class="task-item">
        <div class="task-header">
          <div>
            <h2>{todo.description}</h2>
            <small>Created by: {todo.created_by}, last updated: {todo.updated_at}</small>
          </div>
          <div>
            <span class="status {todo.status}">{todo.status}</span>
            <div>
              <button class="secondary-button" on:click={() => prevStatus(todo)}>prev</button>
              <button class="secondary-button" on:click={() => nextStatus(todo)}>next</button>
            </div>
          </div>
        </div>
      </li>
    {/each}
  </ul>
{/if}

<div class="new-task">
  <input type="text" bind:value={newTodoDescription} placeholder="New task description" />
  <button on:click={createTodo}>Add</button>
</div>

<style>
  @import '../assets/base.css';

  .container {
    width: 90%;
    margin: 0 0;
    padding-bottom: 20px;
  }
  h1 {
    font-family: 'Open Sans', sans-serif;
    font-weight: 600;
    color: var(--primary);
    margin-top: 20px;
  }
  .task-list {
    width: 90%;
    list-style-type: none;
    padding: 10px 0 0 0;
  }
  .task-item {
    background: var(--white);
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    padding: 20px;
    margin-bottom: 10px;
  }
  .task-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
    color: var(--primary);
  }
  .status {
    display: inline-block;
    padding: 5px 10px;
    border-radius: 4px;
    font-size: 1.2em;
    color: var(--white);
    text-align: center;
    width: 100%;
    margin-bottom: 10px;
  }
  .status.todo {
    background-color: var(--status-todo);
  }
  .status.ongoing {
    background-color: var(--status-ongoing);
  }
  .status.done {
    background-color: var(--status-done);
  }
  .new-task {
    display: flex;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background-color: var(--background);
    padding: 10px;
    border-top: 1px solid var(--light-gray);
  }
  .new-task input {
    flex: 1;
    padding: 10px;
    margin-right: 10px;
    border: 1px solid var(--light-gray);
    border-radius: 4px;
  }
  .new-task button {
    background-color: var(--primary);
    color: white;
    padding: 10px 15px;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    cursor: pointer;
  }
  .new-task button:hover {
    background-color: var(--primary-hover);
  }
  .secondary-button {
    background-color: var(--white);
    color: var(--acc);
    border: 2px solid var(--primary);
    padding: 3px 20px;
    font-size: 14px;
    border-radius: 5px;
    cursor: pointer;
    transition:
      background-color 0.3s,
      color 0.3s,
      border-color 0.3s;
  }

  .secondary-button:hover {
    background-color: var(--primary);
    color: var(--white);
    border-color: var(--white);
  }

  .secondary-button:active {
    background-color: #ff69b4;
    color: var(--white);
    border-color: var(--primary);
  }

  .secondary-button:focus {
    outline: none;
    box-shadow: 0 0 0 3px var(--background);
  }
</style>
