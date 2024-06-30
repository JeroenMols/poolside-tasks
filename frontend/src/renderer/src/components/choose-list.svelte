<script lang="ts">
  import { ensureNonEmpty } from '../utils/assertions'

  export let accessToken: string = ''
  export let onListSelected: (list: string) => void

  let todoListId = ''

  type CreateTodoListResponse = {
    todo_list_id: string
  }

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

<h1>TODO lists</h1>
<p>Create or join a todo list</p>
<input type="text" bind:value={todoListId} placeholder="Enter todolist id" />
<button on:click={() => checkListExists()}>Join list</button>
<p>Or</p>
<button
  on:click={() => {
    createList()
  }}>Create list</button
>
