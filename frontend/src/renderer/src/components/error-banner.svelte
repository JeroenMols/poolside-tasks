<script lang="ts">
  export let errorMessage = ''
  export let onDismissError: () => void
</script>

{#if errorMessage}
  <div class="error-container">
    <p class="error-message">{errorMessage}</p>
    <button class="close-button" on:click={onDismissError}>&times;</button>
  </div>
{/if}

<style>
  .error-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    background-color: var(--status-todo);
    color: var(--white);
    padding: 15px 20px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    display: flex;
    justify-content: space-between;
    align-items: center;
    z-index: 1000;
  }
  .error-message {
    margin: 0;
  }
  .close-button {
    background: none;
    border: none;
    color: var(--white);
    cursor: pointer;
    font-size: 20px;
  }
  .close-button:hover {
    opacity: 0.7;
  }
</style>
