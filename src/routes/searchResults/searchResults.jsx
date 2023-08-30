import style from "./searchResults.module.css";

import PostListSelector from "../postListSelector/postListSelector";
import PostListItems from "../postListItems/postListItems";

import { getPosts } from "../../scripts/firebase";

import {
    useLoaderData,
    useOutletContext
} from "react-router-dom";

export async function loader({params, request}) {
    const url = new URL(request.url);

    let tags = url.searchParams.has("tags") && url.searchParams.get("tags") != "" ? url.searchParams.get("tags").split(",") : [];
    let search = url.searchParams.has("s") ? url.searchParams.get("s") : "";

    const posts = await getPosts({page: parseInt(params.pageNumber), tags:tags, search: search, isRoot: url.pathname == "/"});

    return {posts};
}

export default function SearchResults() {
    let pageNumber = useOutletContext();
    let {posts} = useLoaderData();

    return(
        <>
            <div className={style.searchResults}>
                <PostListItems posts={posts.Posts}/>
                <PostListSelector length={posts.Pages} current={pageNumber}/>
            </div>       
        </>
    )
}