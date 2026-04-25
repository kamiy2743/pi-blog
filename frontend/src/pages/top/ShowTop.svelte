<svelte:head>
  <title>パンダの開発ブログ</title>
</svelte:head>

<script lang="ts">
  import { Link } from '@inertiajs/svelte'
  import PublicSiteLink from '../../components/PublicSiteLink.svelte'
  import { formatDate } from '../../utils/date'

  type Article = {
    id: number
    title: string
    date: string
    categoryNames: string[]
  }

  type Category = {
    id: number
    name: string
  }

  export let latestArticles: Article[] = []
  export let categories: Category[] = []
</script>

<div class="mx-auto max-w-7xl px-6 py-10 sm:px-8 lg:px-10 lg:py-12">
  <div class="blog-shell overflow-hidden rounded-[2rem] border">
    <div class="blog-hero border-b px-6 py-10 sm:px-8 lg:px-10 lg:py-12">
      <div class="max-w-3xl space-y-4">
        <PublicSiteLink />
        <h1 class="blog-title text-3xl font-semibold tracking-tight sm:text-3xl">
          パンダの開発ブログ
        </h1>
        <p class="blog-copy max-w-2xl text-base leading-8 sm:text-lg">
          Raspberry Pi で自宅サーバーを開発・運用しながら学んだことを発信していきます。
        </p>
      </div>
    </div>

    <div class="grid gap-8 px-6 py-8 sm:px-8 lg:grid-cols-[minmax(100px,0.5fr)_minmax(0,1.7fr)] lg:gap-8 lg:px-10 lg:py-10">
      <aside class="space-y-5">
        <div class="blog-side-card rounded-[1.5rem] border p-6">
          <h2 class="blog-section-title text-2xl font-semibold">カテゴリ</h2>
          <div class="mt-5 space-y-3">
            {#each categories as category}
              <Link
                class="blog-category-link block rounded-2xl border px-4 py-3 font-semibold transition"
                href={`/article?categoryId=${encodeURIComponent(category.id)}`}
              >
                {category.name}
              </Link>
            {/each}
          </div>
          <Link
            class="blog-button mt-5 inline-flex items-center gap-2 rounded-full px-5 py-3 text-sm font-semibold transition"
            href="/article"
          >
            記事一覧
            <span aria-hidden="true">→</span>
          </Link>
        </div>
      </aside>

      <section class="space-y-5">
        <div class="space-y-2">
          <h2 class="blog-section-title text-2xl font-semibold">最新記事</h2>
        </div>

        <div class="space-y-4">
          {#each latestArticles as article}
            <Link
              class="blog-article-card group block rounded-[1.35rem] border px-5 py-5 transition"
              href={`/article/${article.id}`}
            >
              <p class="blog-accent text-sm font-semibold">{formatDate(article.date)}</p>
              <h3 class="blog-card-title mt-2 text-xl font-semibold">
                {article.title}
              </h3>
              <div class="mt-3 flex flex-wrap gap-2">
                {#each article.categoryNames as categoryName}
                  <span class="blog-category-pill px-3 py-1 text-sm font-semibold">{categoryName}</span>
                {/each}
              </div>
            </Link>
          {/each}
        </div>
      </section>
    </div>
  </div>
</div>
