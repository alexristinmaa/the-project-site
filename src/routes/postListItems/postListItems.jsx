import PostListItem from "./postListItem/postListItem";

import style from "./postListItems.module.css"

export default function PostListItems({posts}) {

    return (
        <>
            <div id={style.gridder} style={{gridTemplateColumns: "repeat(" + posts.length + ", 1fr)"}}>
                {posts.map((post) => <PostListItem post={post} />)}
            </div>
    </>
    )
}