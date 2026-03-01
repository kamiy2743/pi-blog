<script lang="ts">
  export let articleId: string
</script>

<svelte:head>
  <title>管理 - 編集 {articleId}</title>
</svelte:head>

<div class="container">
  <div class="card">
    <h1>編集: {articleId}</h1>
    <p class="muted">記事を更新します。</p>

    <form method="post" action={`/admin/article/edit/${articleId}`} style="margin-top: 24px;">
      <div style="display: grid; gap: 12px;">
        <label>
          <div class="muted">タイトル</div>
          <input name="title" type="text" style="width: 100%; padding: 10px; border-radius: 10px; border: 1px solid #e2d8cc;" />
        </label>
        <label>
          <div class="muted">本文 (Markdown)</div>
          <textarea name="content_md" rows="10" style="width: 100%; padding: 10px; border-radius: 10px; border: 1px solid #e2d8cc;"></textarea>
        </label>
        <button class="cta" type="submit">更新</button>
      </div>
    </form>

    <form method="post" action={`/admin/article/publish/${articleId}`} style="margin-top: 12px;">
      <button class="cta" type="submit">公開/非公開を切替</button>
    </form>
  </div>
</div>
