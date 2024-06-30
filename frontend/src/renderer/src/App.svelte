<script lang="ts">
  import LogIn from './components/log-in.svelte'
  import ChooseList from './components/choose-list.svelte'
  import Todo from './components/to-do.svelte'

  let accessToken: string | null = null
  let todoList: string | null = null

  const onLogIn = (token: string) => {
    console.log('Received token:', token)
    accessToken = token
  }
  const onListSelected = (list: string) => {
    console.log('List selected:', list)
    todoList = list
  }
</script>

{#if todoList != null}
  <Todo {accessToken} todoListId={todoList} />
{:else if accessToken != null}
  <ChooseList {accessToken} {onListSelected} />
{:else}
  <LogIn {onLogIn} />
{/if}
