# Frontend Test Investigation

## Summary

Next.js App Router の frontend test は、1つの道具ですべてを確認するより、対象の大きさで分ける。

- Unit test: `chunk`, `hashKey`, `rankingHandlers` などの純粋関数を確認する。
- Component test: `RankingSelector`, `ProductTable` などの Client Component の振る舞いを確認する。
- E2E test: Next.js page として `/` が表示され、実ブラウザ上で主要導線が動くことを確認する。

当面は `Vitest + React Testing Library` を Unit / Component test に使い、Server Component や route 全体は `next build` と将来的な Playwright E2E で確認するのがよい。

## Next.js Specific Practice

Next.js App Router では Server Components がデフォルトで、Client Component は `"use client"` を明示する。テストもこの境界を意識する。

- Server Component は可能な限り server 側に残す。
- `useState`, `onClick` などの browser interaction は Client Component に閉じ込める。
- Server Component から Client Component に渡す props は、基本的に serializable な値にする。
- `async Server Components` は test runner 側の対応が限定的なので、細かい component test より E2E / build verification を重視する。

この project では `app/page.tsx` は Server Component として残し、sorting interaction は `ProductTable` / `RankingSelector` 側の Client Component で扱う方針にした。

## Vitest In Next.js

Vitest は Vite-powered な test runner だが、application bundler が Vite である必要はない。

この project の runtime/build は Next.js + Turbopack だが、Unit / Component test は Vitest の Vite-based transform で実行できる。

ただし、Vitest は Next.js / Turbopack と完全に同じ実行経路ではない。そのため、以下を分けて考える。

- Vitest: pure TypeScript logic, hooks, Client Component の小さな振る舞い
- `next build`: Next.js として build できること、Server / Client boundary が壊れていないこと
- Playwright: 実ブラウザ上で page と user flow が動くこと

Added devDependencies:

- `vitest`
- `@vitejs/plugin-react`
- `jsdom`
- `@testing-library/react`
- `@testing-library/dom`
- `vite-tsconfig-paths`

## React Testing Library

React Testing Library は React component を DOM に render し、user から見える振る舞いを検証するための library。

見るべきもの:

- button, heading, textbox など user が認識する要素があるか
- click や input で期待する UI 変化が起きるか
- callback が user action に応じて呼ばれるか

避けるべきもの:

- component 内部の state 変数を直接検証する
- CSS class の詳細に依存する
- implementation detail に強く依存する mock を増やす

例として `RankingSelector` では、`onSelectOrder` が正しい `order` で呼ばれるかを見る。`useState` の中身そのものは直接テストしない。

## Playwright

Playwright は本物の browser で Next.js app を開く E2E test に使う。

確認対象:

- `/` が表示される
- 商品一覧が表示される
- sort button を押すと商品順が変わる
- image optimization や routing を含めて実ブラウザで破綻しない

React Testing Library と Playwright は思想が近いが、対象の大きさが違う。

- React Testing Library: component 単体または小さい component composition
- Playwright: page 全体、routing、browser behavior、backend integration を含む user flow

## What To Test

優先順位は次の通り。

1. Pure logic
   - `chunk` が正しく配列を分割する
   - `rankingHandlers.newer` が登録日時の新しい順に並べる
   - `rankingHandlers.lowerPrice` / `higherPrice` が価格順に並べる
   - `hashKey` が同じ入力に対して同じ key を返す

2. Client Component behavior
   - `RankingSelector` の各 button が表示される
   - click によって正しい `RankingSelectorContext` が通知される
   - selected state が user に見える形で変わる

3. Component integration
   - `ProductTable` が商品を表示する
   - sort button click で表示順が変わる
   - tag や price が期待通り表示される

4. Page / E2E
   - `/` にアクセスして商品一覧を確認できる
   - 初期表示が新しい順になっている
   - sort interaction が実ブラウザで動く

5. Backend integration, later
   - ConnectRPC client で商品一覧を取得できる
   - loading / error / empty state を表示できる
   - local backend と AWS backend の接続先を切り替えられる

## Storybook Position

Storybook は test runner というより、UI component catalog / isolated development environment。

用途:

- component の見た目を props pattern ごとに確認する
- design review のために UI state を一覧化する
- visual regression test の土台にする

現時点では必須ではない。`ProductCard` や `ProductTable` の状態パターンが増え、見た目のレビューや visual regression が必要になった時点で導入を検討する。

## Recommended First Tests

最初に書くなら、この順番がよい。

1. `web/lib/lib.ts`
   - `chunk`
   - `rankingHandlers`

2. `web/lib/key.ts`
   - `hashKey`

3. `web/components/RankingSelector.tsx`
   - sort button rendering
   - click callback

4. `web/components/ProductTable.tsx`
   - product rendering
   - sort behavior

5. E2E later
   - `/` page rendering
   - sort flow

## References

- Next.js Testing Guide: https://nextjs.org/docs/app/guides/testing
- Next.js Vitest Guide: https://nextjs.org/docs/app/guides/testing/vitest
- Testing Library Guiding Principles: https://testing-library.com/docs/guiding-principles/
- Playwright Best Practices: https://playwright.dev/docs/best-practices
- Vitest Guide: https://vitest.dev/guide/
