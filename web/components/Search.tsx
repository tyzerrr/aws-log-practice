import {Search as SearchIcon} from "lucide-react"


export default function Search() {
    return (
        <div className="flex items-center rounded-4xl bg-primary-dark px-6 py-3 w-80 gap-3 border border-primary-line">
            <SearchIcon className={"text-gray-400 "}/>
            <input placeholder={"商品を検索..."} className={"focus:outline-none flex-1 min-w-0"}/>
        </div>
    )
}