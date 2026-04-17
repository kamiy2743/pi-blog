<script lang="ts">
  import { Link } from '@inertiajs/svelte'
  import PublicSiteLink from '../../components/PublicSiteLink.svelte'
  import { formatDate } from '../../utils/date'

  type Article = {
    id: number
    title: string
    body: string
    date: string
    categoryNames: string[]
  }

  export let article: Article
</script>

<svelte:head>
  <title>{article.title}</title>
</svelte:head>

<div class="mx-auto max-w-5xl px-6 py-10 sm:px-8 lg:px-10 lg:py-14">
  <div class="blog-shell overflow-hidden rounded-[2rem] border">
    <article class="px-6 py-10 sm:px-8 lg:px-10 lg:py-12">
      <div class="max-w-3xl space-y-5">
        <PublicSiteLink />
        <div class="flex flex-wrap gap-2">
          {#each article.categoryNames as categoryName}
            <span class="blog-category-pill px-3 py-1 text-sm font-semibold">
              {categoryName}
            </span>
          {/each}
        </div>
        <h1 class="blog-title text-3xl font-semibold tracking-tight sm:text-4xl">
          {article.title}
        </h1>
        <p class="muted text-sm font-semibold">{formatDate(article.date)}</p>
      </div>

      <div class="mt-8 max-w-3xl">
        <p class="whitespace-pre-wrap text-base leading-8">{article.body}</p>
      </div>

      <div class="mt-8">
        <Link
          class="blog-button inline-flex items-center gap-2 rounded-full px-5 py-3 text-sm font-semibold transition"
          href="/article"
        >
          記事一覧へ戻る
        </Link>
      </div>
    </article>
  </div>
</div>
