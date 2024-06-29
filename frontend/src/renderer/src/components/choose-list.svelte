<script lang="ts">
  export let accessToken: string = ''
  export let onListSelected: (list: string) => void
</script>

<h1>TODO lists</h1>
<p>Create or join a todo list</p>
<input type="text" placeholder="Enter todolist id" />
<button
  on:click={() => {
    onListSelected('fake_list_id')
  }}>Join list</button
>
<p>Or</p>
<button
  on:click={() => {
    onListSelected('fake_list_id')
  }}>Create list</button
>
