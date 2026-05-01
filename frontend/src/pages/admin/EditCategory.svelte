<script lang="ts">
  import { Form } from '@inertiajs/svelte'
  import FlashMessage from '../../components/FlashMessage.svelte'
  import type { Flash } from '../../types/flash'

  type Category = {
    id: number
    name: string
  }

  export let categories: Category[] = []
  export let validationErrors: Record<string, string> = {}
  export let flash: Flash = {}

  let categoryDrafts: Record<number, string> = {}
  let deleteTarget: Category | null = null

  $: categoryDrafts = Object.fromEntries(
    categories.map((category) => [category.id, category.name]),
  )

  function updateCategoryDraft(categoryId: number, value: string) {
    categoryDrafts = {
      ...categoryDrafts,
      [categoryId]: value,
    }
  }

  function openDeleteModal(category: Category) {
    deleteTarget = category
  }

  function closeDeleteModal() {
    deleteTarget = null
  }
</script>

<svelte:head>
  <title>管理 - カテゴリ編集</title>
</svelte:head>

<div class="admin-page px-6 py-8 sm:px-8 lg:px-10 lg:py-10">
  <div class="mx-auto max-w-5xl">
    <FlashMessage {flash} />

    <div class="admin-panel rounded-lg border px-6 py-8 sm:px-8">
      <p class="admin-eyebrow text-sm font-semibold">Category</p>
      <h1 class="mt-3 text-3xl font-semibold tracking-tight">カテゴリ編集</h1>
      <p class="admin-copy mt-3 text-base leading-8">カテゴリの追加と名前の更新を行います。</p>

      <Form
        class="mt-6 space-y-5"
        action="/admin/category"
        method="post"
        options={{ preserveScroll: true, preserveState: false }}
        resetOnSuccess={true}
      >
        <label class="block">
          <span class="text-sm font-semibold">新しいカテゴリ名</span>
          <input
            class="mt-2 w-full rounded-lg border px-4 py-3"
            name="name"
            type="text"
            autocomplete="off"
          />
          {#if validationErrors['create.name']}
            <p class="mt-2 text-sm font-semibold text-red-600">{validationErrors['create.name']}</p>
          {/if}
        </label>

        <div class="flex flex-wrap gap-3">
          <button
            class="admin-button inline-flex items-center justify-center rounded-lg px-5 py-3 text-sm font-semibold"
            type="submit"
          >
            追加
          </button>
        </div>
      </Form>
    </div>

    <div class="admin-panel mt-6 rounded-lg border px-4 py-5 sm:px-6">
      <div class="flex items-end justify-between gap-3 border-b border-[var(--admin-border)] px-2 pb-3">
        <div>
          <h2 class="text-xl font-semibold">登録済みカテゴリ</h2>
          <p class="admin-copy mt-1 text-sm">{categories.length} 件</p>
        </div>
      </div>

      {#if categories.length === 0}
        <p class="admin-copy mt-4">カテゴリはまだありません。</p>
      {:else}
        <div class="admin-panel mt-3 overflow-x-auto rounded-lg border">
          <table class="min-w-full border-collapse text-left text-sm">
            <thead>
              <tr class="admin-table-head">
                <th class="px-4 py-3 font-semibold">カテゴリ名</th>
                <th class="whitespace-nowrap px-4 py-3 font-semibold">操作</th>
              </tr>
            </thead>
            <tbody>
              {#each categories as category}
                <tr class="admin-table-row">
                  <td class="px-4 py-3">
                    <Form
                      id={`category-form-${category.id}`}
                      action={`/admin/category/${category.id}`}
                      method="post"
                      options={{ preserveScroll: true, preserveState: false }}
                    >
                      <input
                        class="w-full rounded-lg border px-3 py-2 text-sm"
                        name="name"
                        type="text"
                        value={categoryDrafts[category.id]}
                        autocomplete="off"
                        on:input={(event) => {
                          const target = event.currentTarget
                          if (target instanceof HTMLInputElement) {
                            updateCategoryDraft(category.id, target.value)
                          }
                        }}
                      />
                    </Form>
                    {#if validationErrors['update.name']}
                      <p class="mt-1 text-xs font-semibold text-red-600">{validationErrors['update.name']}</p>
                    {/if}
                  </td>
                  <td class="whitespace-nowrap px-4 py-3">
                    <div class="flex gap-2">
                      <button
                        class={`rounded-lg border px-4 py-2 text-sm font-semibold disabled:cursor-not-allowed disabled:opacity-40 ${
                          categoryDrafts[category.id] !== category.name ? 'admin-button' : 'admin-secondary-button'
                        }`}
                        type="submit"
                        form={`category-form-${category.id}`}
                        disabled={categoryDrafts[category.id] === category.name}
                      >
                        更新
                      </button>
                      <button
                        class="admin-secondary-button rounded-lg border px-4 py-2 text-sm font-semibold"
                        type="button"
                        on:click={() => openDeleteModal(category)}
                      >
                        削除
                      </button>
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}

      <a class="admin-secondary-button mt-5 inline-flex items-center justify-center rounded-lg border px-5 py-3 text-sm font-semibold" href="/admin">
        戻る
      </a>
    </div>
    {#if deleteTarget}
      <div
        class="fixed inset-0 z-50 bg-black/50"
        aria-hidden="true"
      ></div>
      <div
        class="fixed inset-0 z-[60] flex items-center justify-center px-6"
        role="dialog"
        aria-modal="true"
        aria-labelledby="delete-category-title"
      >
        <div class="admin-panel w-full max-w-md rounded-lg border bg-[var(--admin-surface)] px-6 py-6 shadow-2xl">
          <h2 id="delete-category-title" class="text-xl font-semibold">カテゴリを削除しますか</h2>
          <p class="admin-copy mt-3 text-sm leading-6">
            {deleteTarget.name} を削除します。
          </p>

          <Form
            class="mt-6 flex flex-wrap gap-3"
            action={`/admin/category/${deleteTarget.id}/delete`}
            method="post"
            options={{ preserveScroll: true, preserveState: false }}
            onSuccess={closeDeleteModal}
          >
            <button
              class="admin-button rounded-lg px-5 py-3 text-sm font-semibold"
              type="submit"
            >
              削除
            </button>
            <button
              class="admin-secondary-button rounded-lg border px-5 py-3 text-sm font-semibold"
              type="button"
              on:click={closeDeleteModal}
            >
              キャンセル
            </button>
          </Form>
        </div>
      </div>
    {/if}
  </div>
</div>
