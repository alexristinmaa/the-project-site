import {
    useLoaderData
} from "react-router-dom";

import { getPost } from "../../scripts/posts";

import style from "./post.module.css";

export async function loader({params}) {
    let post = await getPost({id: params.postId});

    return {post};
}

export default function Post() {
    let {post} = useLoaderData();

    return(
        <>
        <div className={style.post}>
            <header>
                <h1 className={style.title}>{post.Title}</h1>
                <p className={style.description}>{post.description}</p>
            </header>
            <main>

            </main>
            <footer>
                <span className={style.readTime}>{Math.floor(post.Body.split(" ").length / 200)} minute read.</span>
                <span className={style.author}>{post.AuthorName}</span>
            </footer>
        </div>
        </>
    )
}