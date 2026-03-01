<svelte:head>
  <title>管理 - 新規作成</title>
</svelte:head>

<div class="container">
  <div class="card">
    <h1>新規作成</h1>
    <p class="muted">Markdown を貼り付けて記事を作成します。</p>

    <form method="post" action="/admin/article/new" style="margin-top: 24px;">
      <div style="display: grid; gap: 12px;">
        <label>
          <div class="muted">タイトル</div>
          <input name="title" type="text" style="width: 100%; padding: 10px; border-radius: 10px; border: 1px solid #e2d8cc;" />
        </label>
        <label>
          <div class="muted">本文 (Markdown)</div>
          <textarea name="content_md" rows="10" style="width: 100%; padding: 10px; border-radius: 10px; border: 1px solid #e2d8cc;"></textarea>
        </label>
        <button class="cta" type="submit">保存</button>
      </div>
    </form>
  </div>
</div>
