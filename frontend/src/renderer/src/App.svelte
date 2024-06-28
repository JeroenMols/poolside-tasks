<script>
  import { onMount } from 'svelte'
  import ComingSoon from './components/ComingSoon.svelte'

  let greeting = '...loading'

  onMount(() => {
    fetch('http://localhost:8080/hello/jeroen')
      .then((response) => response.json())
      .then((data) => {
        greeting = data.greeting
      })
  })
</script>

<h1>{greeting}</h1>
<ComingSoon />
