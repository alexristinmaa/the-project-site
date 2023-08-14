import style from "./postListItem.module.css"

import {
    Link
} from "react-router-dom";

export default function PostListItem({post}) {
    console.log(post)
    return(
        <>
        <article className={style.post}>
            <div>
                <img src="test-thumbnail.jpeg" alt="A Test Thumbail"/>
            </div>
            <div className={style.text}>
                <p className={style.date}>{post.Date}</p>
                <h2 className={style.header}>{post.Title}</h2>
                <p className={style.description}>{post.Description}</p>
                <div className={style.footer}>
                    <div className={style.divisor}></div>
                    <div className={style.tag}>{post.Tag}</div>
                </div>
            </div>
        </article>
        </>
    )
}