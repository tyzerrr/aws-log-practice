import Logo from "@/components/Logo";
import Search from "@/components/Search";
import Cart from "@/components/Cart";
import Likes from "@/components/Likes";
import Account from "@/components/Account";

export default function Header() {
    return (
        <div className={"flex w-full items-center justify-between px-6 py-3 bg-white border border-primary-line"}>
            <div className={"flex items-center gap-6"}>
                <Logo/>
            </div>
            <div className={"flex items-center gap-6"}>
                <Search/>
                <Likes/>
                <Account/>
                <Cart/>
            </div>
        </div>
    )
}