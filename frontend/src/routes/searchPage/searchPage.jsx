import {
    useLoaderData,
    useParams,
    Outlet,
    redirect,
    useNavigate
} from "react-router-dom";

import {
    useState,
    useMemo
} from "react";

import Dropdown from "./dropdown/dropdown";

import { getNumberOfPages, getTags } from "../../scripts/posts";

import style from "./searchPage.module.css";

export async function loader({params, request}) {
    const url = new URL(request.url);
    const hasTags = url.searchParams.has("tags");
    //const hasSearch = url.searchParams.has("s");

    const tags = await getTags();

    let activeTags = hasTags ? url.searchParams.get("tags").split(",") : [];

    activeTags = activeTags.filter((tag) => tags.map((t) => t.Name).includes(tag));

    let isRoot = activeTags.length == 0 && params.pageNumber == 1;

    return {tags, activeTags, isRoot};
}

export default function SearchPage() {
    let {tags, activeTags, isRoot} = useLoaderData();
    let {pageNumber} = useParams();
    let navigate = useNavigate();

    let [activeState, setActiveState] = useState(activeTags);
    
    let nItems = isRoot ? 3 : 4;

    pageNumber = parseInt(pageNumber);

    function search(e) {
        if(e.key == "Enter") {
            let searchValue = document.getElementById("searchInput").value;

            let tagsQuery = activeState;
            let searchQuery = encodeURIComponent(searchValue);

            let totalQuery = "?" +
                                [((tagsQuery.length != 0) ? "tags=" + tagsQuery : ""),
                                ((searchQuery != "") ? "s=" + encodeURIComponent(searchValue) : "")].join("&");

            // remove any & in the end of the query
            if(totalQuery[totalQuery.length-1] == "&") totalQuery = totalQuery.slice(0, -1);
            if(totalQuery == "?") totalQuery = "";
            

            navigate(`/posts/page/1${totalQuery}`);
        }
    }

    return (
        <>
        <div id={style.root}>
            <div id={style.searchHeader} style={{padding: `0 max(5% / ${nItems},50% / ${nItems} - 200px)`}}>
                <Dropdown tags={tags} activeState={activeState} changeActive={(e) => {setActiveState(e)}}/>
                <div id={style.searchBar}>
                    <input id="searchInput" placeholder="Search the blog..." onKeyDown={search}></input>
                    <i className="search-icon"></i>
                </div>
            </div>
            <Outlet context={pageNumber} />
        </div>
        </>
    )
}