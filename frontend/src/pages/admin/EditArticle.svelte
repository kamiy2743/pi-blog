<script lang="ts">
  import { Form, Link } from '@inertiajs/svelte'
  import type { FormDataConvertible } from '@inertiajs/core'
  import ConfirmModal from '../../components/ConfirmModal.svelte'
  import { getOldInputString, getOldInputStringList, type OldInput, type ValidationErrors } from '../../types/form'
  import { fromJstDatetimeLocalValue, toJstDatetimeLocalValue } from '../../utils/date'

  type Article = {
    id: number
    title: string
    bodyMarkdown: string
    isPublished: boolean
    publishStartAt: string
    publishEndAt: string
    categoryIds: number[]
  }

  type Category = {
    id: number
    name: string
  }

  export let article: Article
  export let categories: Category[] = []
  export let oldInput: OldInput = {}
  export let validationErrors: ValidationErrors = {}

  function hasOldInput(field: string) {
    return Object.prototype.hasOwnProperty.call(oldInput, field)
  }

  let isPublished = hasOldInput('isPublished')
    ? getOldInputString(oldInput, 'isPublished') === 'true'
    : article.isPublished
  let selectedCategoryIds = hasOldInput('categoryIds')
    ? getOldInputStringList(oldInput, 'categoryIds')
    : article.categoryIds.map(String)
  let isDeleteModalOpen = false
</script>

<svelte:head>
  <title>管理 - 記事編集</title>
</svelte:head>

<div class="admin-page px-6 py-8 sm:px-8 lg:px-10 lg:py-10">
  <div class="mx-auto max-w-6xl">
    <header class="admin-hero rounded-lg border px-6 py-8 sm:px-8">
      <p class="admin-eyebrow text-sm font-semibold">Article</p>
      <h1 class="mt-3 text-3xl font-semibold tracking-tight">記事編集</h1>
      <p class="admin-copy mt-3 text-base leading-8">
        Markdown 本文と公開条件を更新します。
      </p>
    </header>

    <Form
      class="mt-6 grid gap-6 lg:grid-cols-[minmax(0,1fr)_320px]"
      action={`/admin/article/${article.id}`}
      method="post"
      options={{ preserveScroll: true, preserveState: false }}
      transform={(data: Record<string, FormDataConvertible>) => ({
        ...data,
        categoryIds: selectedCategoryIds,
        publishStartAt: fromJstDatetimeLocalValue(String(data.publishStartAt ?? '')),
        publishEndAt: fromJstDatetimeLocalValue(String(data.publishEndAt ?? '')),
      })}
    >
      <input type="hidden" name="isPublished" value={isPublished ? 'true' : 'false'} />

      <section class="admin-panel rounded-lg border px-6 py-6 sm:px-8">
        <div class="space-y-6">
          <label class="block">
            <span class="text-sm font-semibold">タイトル</span>
            <input
              class="mt-2 w-full rounded-lg border px-4 py-3"
              name="title"
              type="text"
              value={hasOldInput('title') ? getOldInputString(oldInput, 'title') : article.title}
              autocomplete="off"
            />
            {#if validationErrors.title}
              <p class="mt-2 text-sm font-semibold text-red-600">{validationErrors.title}</p>
            {/if}
          </label>

          <label class="block">
            <span class="text-sm font-semibold">本文</span>
            <textarea
              class="mt-2 min-h-[32rem] w-full resize-y rounded-lg border px-4 py-3 font-mono text-sm leading-7 lg:min-h-[42rem]"
              name="body"
              value={hasOldInput('body') ? getOldInputString(oldInput, 'body') : article.bodyMarkdown}
              placeholder="Markdown"
            ></textarea>
            {#if validationErrors.body}
              <p class="mt-2 text-sm font-semibold text-red-600">{validationErrors.body}</p>
            {/if}
          </label>
        </div>
      </section>

      <aside class="space-y-4">
        <section class="admin-panel rounded-lg border px-5 py-5">
          <h2 class="text-lg font-semibold">公開設定</h2>

          <label class="mt-5 flex items-center gap-3 text-sm font-semibold">
            <input
              class="h-4 w-4"
              type="checkbox"
              bind:checked={isPublished}
            />
            <span>公開する</span>
          </label>

          <label class="mt-5 block">
            <span class="text-sm font-semibold">公開開始時刻</span>
            <input
              class="mt-2 w-full rounded-lg border px-3 py-2 text-sm"
              name="publishStartAt"
              type="datetime-local"
              value={toJstDatetimeLocalValue(hasOldInput('publishStartAt') ? getOldInputString(oldInput, 'publishStartAt') : article.publishStartAt)}
            />
            {#if validationErrors.publishStartAt}
              <p class="mt-2 text-sm font-semibold text-red-600">{validationErrors.publishStartAt}</p>
            {/if}
          </label>

          <label class="mt-5 block">
            <span class="text-sm font-semibold">公開終了時刻</span>
            <input
              class="mt-2 w-full rounded-lg border px-3 py-2 text-sm"
              name="publishEndAt"
              type="datetime-local"
              value={toJstDatetimeLocalValue(hasOldInput('publishEndAt') ? getOldInputString(oldInput, 'publishEndAt') : article.publishEndAt)}
            />
            {#if validationErrors.publishEndAt}
              <p class="mt-2 text-sm font-semibold text-red-600">{validationErrors.publishEndAt}</p>
            {/if}
          </label>
        </section>

        <section class="admin-panel rounded-lg border px-5 py-5">
          <h2 class="text-lg font-semibold">カテゴリ</h2>
          {#if categories.length === 0}
            <p class="admin-copy mt-3 text-sm leading-6">カテゴリはまだありません。</p>
          {:else}
            <div class="mt-4 flex flex-wrap gap-2">
              {#each categories as category}
                <label
                  class={`admin-filter-chip inline-flex items-center gap-2 rounded-lg border px-3 py-2 text-sm font-semibold transition ${
                    selectedCategoryIds.includes(String(category.id)) ? 'admin-filter-chip-selected' : ''
                  }`}
                >
                  <input
                    type="checkbox"
                    value={String(category.id)}
                    bind:group={selectedCategoryIds}
                  />
                  <span>{category.name}</span>
                </label>
              {/each}
            </div>
          {/if}
          {#if validationErrors.categoryIds}
            <p class="mt-2 text-sm font-semibold text-red-600">{validationErrors.categoryIds}</p>
          {/if}
        </section>

        <section class="admin-panel rounded-lg border px-5 py-5">
          <button
            class="admin-button inline-flex w-full items-center justify-center rounded-lg px-5 py-3 text-sm font-semibold"
            type="submit"
          >
            更新
          </button>
          <Link
            class="admin-secondary-button mt-3 inline-flex w-full items-center justify-center rounded-lg border px-5 py-3 text-sm font-semibold"
            href="/admin"
          >
            戻る
          </Link>
          <button
            class="mt-3 inline-flex w-full items-center justify-center rounded-lg border border-red-300 px-5 py-3 text-sm font-semibold text-red-700 transition hover:bg-red-50"
            type="button"
            on:click={() => {
              isDeleteModalOpen = true
            }}
          >
            削除
          </button>
        </section>
      </aside>
    </Form>
    {#if isDeleteModalOpen}
      <ConfirmModal
        title="記事を削除しますか"
        message={`${article.title} を削除します。`}
        action={`/admin/article/${article.id}/delete`}
        confirmLabel="削除"
        onClose={() => {
          isDeleteModalOpen = false
        }}
      />
    {/if}
  </div>
</div>
