<svelte:head>
  <title>Blog Top</title>
</svelte:head>

<div class="container">
  <div class="card hero">
    <div>
      <p class="muted">panda-dev.net / blog</p>
      <h1>静かな学びを記録する Blog</h1>
      <p>
        Raspberry Pi 5 で動く自作ブログ。Go + SQLite を基盤に、
        記事の検証・公開までをシンプルに運用できる構成を目指しています。
      </p>
    </div>

    <div class="badges">
      <span class="badge">Go</span>
      <span class="badge">SQLite</span>
      <span class="badge">Inertia</span>
      <span class="badge">Svelte</span>
      <span class="badge">Raspberry Pi 5</span>
    </div>

    <div class="grid">
      <div class="card">
        <h3>最新の記事</h3>
        <p class="muted">まだ記事はありません。まずは最初の投稿を準備中です。</p>
      </div>
      <div class="card">
        <h3>運用ステータス</h3>
        <p class="muted">/health で監視可。systemd による常駐運用。</p>
      </div>
      <div class="card">
        <h3>管理</h3>
        <p class="muted">/admin から記事作成・公開設定を管理します。</p>
      </div>
    </div>

    <a class="cta" href="/article">記事一覧へ →</a>
  </div>
</div>
