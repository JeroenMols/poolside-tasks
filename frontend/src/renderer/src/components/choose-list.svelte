<script lang="ts">
  import { createTodoList, getTodoList } from '../net/requests'
  import { ensureNonEmpty } from '../utils/assertions'
  import ErrorBanner from './error-banner.svelte'

  export let accessToken: string = ''
  export let onListSelected: (list: string) => void

  let todoListId = ''
  let errorMessage = ''

  const createList = async () => {
    const response = await createTodoList(accessToken)
    if ('error' in response) {
      errorMessage = response.error as string
    } else {
      onListSelected(ensureNonEmpty(response.todo_list_id))
    }
  }

  const checkListExists = async () => {
    const response = await getTodoList(accessToken, todoListId)
    if ('error' in response) {
      todoListId = ''
      errorMessage = response.error as string
    } else {
      onListSelected(ensureNonEmpty(response.todo_list_id))
    }
  }

  const logOut = () => {
    localStorage.removeItem('access_token')
    location.reload()
  }

  const onDismissError = () => {
    errorMessage = ''
    todoListId = ''
  }
</script>

<ErrorBanner {errorMessage} {onDismissError} />
<div class="card">
  <h1>Choose a list</h1>
  <input bind:value={todoListId} type="text" placeholder="join existing list" />
  <button on:click={() => checkListExists()}>join</button>

  <p>Or</p>

  <button on:click={() => createList()}>create new</button>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <p class="logout" on:click={() => logOut()}>logout</p>
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

  .logout:hover {
    color: var(--primary);
    border-color: var(--white);
  }
</style>
