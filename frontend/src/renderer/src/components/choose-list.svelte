<script lang="ts">
  import { ensureNonEmpty } from '../utils/assertions'

  export let accessToken: string = ''
  export let onListSelected: (list: string) => void

  let todoListId = ''

  const createList = async () => {
    let response = await fetch('http://localhost:8080/todolists', {
      headers: { Authorization: `${accessToken}` },
      method: 'POST',
      body: '{}'
    })

    if (response.ok) {
      const registerResponse = (await response.json()) as CreateTodoListResponse
      ensureNonEmpty(registerResponse.todo_list_id)
      onListSelected(registerResponse.todo_list_id)
    } else {
      alert(`Failed to create todo list - backend error (${response.status})`)
    }
  }

  const checkListExists = async () => {
    let response = await fetch(`http://localhost:8080/todolists/${todoListId}`, {
      headers: { Authorization: `${accessToken}` },
      method: 'GET'
    })

    if (response.ok) {
      onListSelected(todoListId)
    } else {
      alert(`Failed to get todo list - backend error (${response.status})`)
    }
  }
</script>

<div class="card">
  <h1>Choose a list</h1>
  <input bind:value={todoListId} type="text" placeholder="join existing list" />
  <button on:click={() => checkListExists()}>join</button>

  <p>Or</p>

  <button on:click={() => createList()}>create new</button>
</div>

<style>
  @import '../assets/base.css';

  .card {
    background: var(--white);
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    padding: 40px;
    width: 350px;
    text-align: center;
  }
  .card h1 {
    font-weight: 600;
    margin-bottom: 20px;
    color: var(--primary);
  }
  .card input[type='text'] {
    width: 100%;
    padding: 10px;
    margin: 10px 0;
    border: 1px solid var(--light-gray);
    border-radius: 4px;
    font-size: 16px;
  }
  .card button {
    background-color: var(--primary);
    color: var(--white);
    padding: 10px 15px;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    cursor: pointer;
    width: 100%;
  }
  .card button:hover {
    background-color: var(--primary-hover);
  }
  .card .footer {
    margin-top: 20px;
    font-size: 14px;
    color: var(--dark-gray);
  }

  p {
    margin: 20px;
  }
</style>
