function light(fontSize?: string): string {
    return `font-light font-serif text-ink3 ${fontSize ?? "text-sm"}`
}

function bold(fontSize?: string): string {
    return `font-bold font-serif ${fontSize ?? "text-sm"}`
}

function verticalList(): string {
    return "flex flex-col gap-4"
}


export default function Footer() {
    return (
        <div className={"flex justify-between items-start px-24 bg-primary-dark min-h-60 pt-6 border border-primary-line"}>
            <div className={verticalList()}>
                <div className={bold("text-xl")}>Dumazon</div>
                <div className={light()}>暮らしと装いの、ちょうどいい道具を。</div>
            </div>

            <div className={verticalList()}>
                <div className={bold()}>ショップ</div>
                <div className={light()}>新着商品</div>
                <div className={light()}>アパレル</div>
                <div className={light()}>雑貨</div>
                <div className={light()}>ガジェット</div>
            </div>

            <div className={verticalList()}>
                <div className={bold()}>サポート</div>
                <div className={light()}>配送について</div>
                <div className={light()}>返品・交換</div>
                <div className={light()}>よくある質問</div>
                <div className={light()}>お問い合わせ</div>
            </div>

            <div className={verticalList()}>
                <div className={bold()}>会社情報</div>
                <div className={light()}>会社概要</div>
                <div className={light()}>プライバシーポリシー</div>
                <div className={light()}>特定商取引法</div>
            </div>
        </div>
    )
}