import {
  Activity,
  Archive,
  CircleDollarSign,
  PackagePlus,
  Search,
  SlidersHorizontal,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";

const products = [
  {
    id: "prd_001",
    name: "CloudWatch Export Bundle",
    description: "監査ログ転送の検証用商品",
    price: "JPY 12,800",
    status: "Active",
    updatedAt: "2026-05-30 18:20",
  },
  {
    id: "prd_002",
    name: "S3 Access Log Parser",
    description: "S3アクセスログの解析テンプレート",
    price: "JPY 8,400",
    status: "Active",
    updatedAt: "2026-05-29 09:12",
  },
  {
    id: "prd_003",
    name: "ALB Request Trace Kit",
    description: "ALBリクエスト調査用パッケージ",
    price: "JPY 16,200",
    status: "Active",
    updatedAt: "2026-05-28 21:45",
  },
];

export default function Home() {
  return (
    <main className="min-h-screen bg-muted/70 text-foreground">
      <div className="mx-auto flex w-full max-w-7xl flex-col gap-6 px-4 py-5 sm:px-6 lg:px-8">
        <header className="flex flex-col gap-4 border-b border-border pb-5 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <p className="text-sm font-medium text-muted-foreground">
              Product Service
            </p>
            <h1 className="mt-1 text-2xl font-semibold tracking-normal text-foreground sm:text-3xl">
              商品管理
            </h1>
          </div>
          <div className="flex flex-col gap-2 sm:flex-row">
            <Button variant="outline">
              <SlidersHorizontal />
              表示条件
            </Button>
            <Button>
              <PackagePlus />
              商品を登録
            </Button>
          </div>
        </header>

        <section className="grid gap-4 md:grid-cols-3">
          <Card>
            <CardHeader className="flex flex-row items-start justify-between gap-3">
              <div>
                <CardDescription>Active Products</CardDescription>
                <CardTitle className="mt-2 text-3xl">128</CardTitle>
              </div>
              <div className="rounded-md bg-accent p-2 text-accent-foreground">
                <Archive className="size-5" />
              </div>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-start justify-between gap-3">
              <div>
                <CardDescription>Average Price</CardDescription>
                <CardTitle className="mt-2 text-3xl">JPY 12,460</CardTitle>
              </div>
              <div className="rounded-md bg-emerald-50 p-2 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300">
                <CircleDollarSign className="size-5" />
              </div>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-start justify-between gap-3">
              <div>
                <CardDescription>Recent Updates</CardDescription>
                <CardTitle className="mt-2 text-3xl">24</CardTitle>
              </div>
              <div className="rounded-md bg-sky-50 p-2 text-sky-700 dark:bg-sky-950 dark:text-sky-300">
                <Activity className="size-5" />
              </div>
            </CardHeader>
          </Card>
        </section>

        <section className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_380px]">
          <Card>
            <CardHeader className="gap-4 sm:flex sm:flex-row sm:items-center sm:justify-between">
              <div>
                <CardTitle>商品一覧</CardTitle>
                <CardDescription className="mt-1">
                  API接続前の表示確認用プレースホルダー
                </CardDescription>
              </div>
              <div className="relative w-full sm:w-72">
                <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                <Input placeholder="商品名で検索" className="pl-9" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="overflow-hidden rounded-md border border-border">
                <div className="grid grid-cols-[1.1fr_1.5fr_0.8fr_0.7fr_1fr] bg-secondary px-4 py-3 text-xs font-medium uppercase text-muted-foreground max-md:hidden">
                  <span>Product</span>
                  <span>Description</span>
                  <span>Price</span>
                  <span>Status</span>
                  <span>Updated</span>
                </div>
                <div className="divide-y divide-border bg-card">
                  {products.map((product) => (
                    <div
                      key={product.id}
                      className="grid gap-3 px-4 py-4 text-sm md:grid-cols-[1.1fr_1.5fr_0.8fr_0.7fr_1fr] md:items-center"
                    >
                      <div>
                        <p className="font-medium">{product.name}</p>
                        <p className="mt-1 font-mono text-xs text-muted-foreground">
                          {product.id}
                        </p>
                      </div>
                      <p className="text-muted-foreground">
                        {product.description}
                      </p>
                      <p className="font-medium">{product.price}</p>
                      <Badge variant="success">{product.status}</Badge>
                      <p className="text-muted-foreground">
                        {product.updatedAt}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>登録フォーム</CardTitle>
              <CardDescription>
                ConnectRPCのmutation接続を差し込むための雛形
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form className="grid gap-4">
                <div className="grid gap-2">
                  <Label htmlFor="name">商品名</Label>
                  <Input id="name" placeholder="ALB Request Trace Kit" />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="description">説明</Label>
                  <Textarea id="description" placeholder="商品説明を入力" />
                </div>
                <div className="grid grid-cols-[1fr_120px] gap-3">
                  <div className="grid gap-2">
                    <Label htmlFor="price">価格</Label>
                    <Input id="price" inputMode="numeric" placeholder="12800" />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="currency">通貨</Label>
                    <Input id="currency" placeholder="JPY" />
                  </div>
                </div>
                <Separator />
                <div className="flex items-center justify-between gap-3">
                  <div>
                    <p className="text-sm font-medium">公開状態</p>
                    <p className="text-sm text-muted-foreground">
                      登録時はActiveとして扱う想定
                    </p>
                  </div>
                  <Badge variant="secondary">Active</Badge>
                </div>
                <Button type="button" className="w-full">
                  <PackagePlus />
                  登録処理を接続する
                </Button>
              </form>
            </CardContent>
          </Card>
        </section>
      </div>
    </main>
  );
}
