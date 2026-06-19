"use client";

import { Search as SearchIcon } from "lucide-react";
import { useSearch } from "./SearchProvider";

export default function Search() {
  const { searchTerm, setSearchTerm } = useSearch();
  return (
    <div className="flex items-center rounded-4xl bg-primary-dark px-6 py-3 w-80 gap-3 border border-primary-line">
      <SearchIcon className={"text-gray-400 "} />
      <input
        placeholder={"商品を検索..."}
        className={"focus:outline-none flex-1 min-w-0"}
        aria-label="search-input"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
      />
    </div>
  );
}
