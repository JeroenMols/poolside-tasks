<script lang="ts">
  import { onMount } from 'svelte'

  export let accessToken: string = ''
  export let todoListId: string = ''

  let todos: TodoItem[] = []
  let newTodoDescription = ''

  onMount(async () => {
    await getTodosForList()
  })

  const getTodosForList = async () => {
    let response = await fetch(`http://localhost:8080/todolists/${todoListId}`, {
      headers: { Authorization: `${accessToken}` },
      method: 'GET'
    })

    if (response.ok) {
      const registerResponse = (await response.json()) as TodoListReponses
      // TODO validate result here
      todos = registerResponse.todos
    } else {
      alert(`Failed to create todo list - backend error (${response.status})`)
    }
  }

  const createTodo = async () => {
    let response = await fetch('http://localhost:8080/todos', {
      headers: { Authorization: `${accessToken}` },
      method: 'POST',
      body: JSON.stringify({ description: newTodoDescription, todo_list_id: todoListId })
    })

    if (response.ok) {
      const item = (await response.json()) as TodoItem
      // TODO validate result here
      newTodoDescription = ''
      todos = [...todos, item]
    } else {
      alert(`Failed to create todo list - backend error (${response.status})`)
    }
  }
</script>

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
            <small>Last updated: {todo.updated_at}</small>
          </div>
          <span class="status {todo.status}">{todo.status}</span>
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
  .container {
    width: 90%;
    margin: 0 0;
    padding-bottom: 20px;
  }
  h1 {
    font-family: 'Open Sans', sans-serif;
    font-weight: 600;
    color: #c2185b;
    margin-top: 20px;
  }
  .task-list {
    width: 90%;
    list-style-type: none;
    padding: 10px 0 0 0;
  }
  .task-item {
    background: #ffffff;
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
  }
  .status {
    display: inline-block;
    padding: 5px 10px;
    border-radius: 4px;
    font-size: 14px;
    color: white;
  }
  .status.todo {
    background-color: #f44336;
  }
  .status.ongoing {
    background-color: #ffc107;
  }
  .status.done {
    background-color: #4caf50;
  }
  .new-task {
    display: flex;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background-color: #fce4ec;
    padding: 10px;
    border-top: 1px solid #cccccc;
  }
  .new-task input {
    flex: 1;
    padding: 10px;
    margin-right: 10px;
    border: 1px solid #cccccc;
    border-radius: 4px;
  }
  .new-task button {
    background-color: var(--ev-c-gray-1);
    color: white;
    padding: 10px 15px;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    cursor: pointer;
  }
  .new-task button:hover {
    background-color: #c2185b;
  }
</style>
