<script lang="ts">
  export let articleId: string
</script>

<svelte:head>
  <title>記事 {articleId}</title>
</svelte:head>

<div class="container">
  <div class="card">
    <h1>記事 {articleId}</h1>
    <p class="muted">本文はまだ準備中です。</p>
  </div>
</div>
